package deprovisioning

import (
	"fmt"
	"time"

	"github.com/kyma-project/kyma-environment-broker/internal/storage/dberr"

	"github.com/kyma-project/control-plane/components/provisioner/pkg/gqlschema"
	"github.com/kyma-project/kyma-environment-broker/internal"
	"github.com/kyma-project/kyma-environment-broker/internal/process"
	"github.com/kyma-project/kyma-environment-broker/internal/provisioner"
	"github.com/kyma-project/kyma-environment-broker/internal/storage"
	"github.com/sirupsen/logrus"
)

type CheckRuntimeRemovalStep struct {
	operationManager  *process.OperationManager
	provisionerClient provisioner.Client
	instanceStorage   storage.Instances
	timeout           time.Duration
}

var _ process.Step = &CheckRuntimeRemovalStep{}

func NewCheckRuntimeRemovalStep(operations storage.Operations, instances storage.Instances,
	provisionerClient provisioner.Client, timeout time.Duration) *CheckRuntimeRemovalStep {
	return &CheckRuntimeRemovalStep{
		operationManager:  process.NewOperationManager(operations),
		provisionerClient: provisionerClient,
		instanceStorage:   instances,
		timeout:           timeout,
	}
}

func (s *CheckRuntimeRemovalStep) Name() string {
	return "Check_Runtime_Removal"
}

func (s *CheckRuntimeRemovalStep) Run(operation internal.Operation, log logrus.FieldLogger) (internal.Operation, time.Duration, error) {
	if time.Since(operation.UpdatedAt) > s.timeout {
		log.Infof("operation has reached the time limit: %s updated operation time: %s", s.timeout, operation.UpdatedAt)
		return s.operationManager.OperationFailed(operation, fmt.Sprintf("CheckRuntimeRemovalStep operation has reached the time limit: %s", s.timeout), nil, log)
	}
	if operation.ProvisionerOperationID == "" {
		log.Infof("ProvisionerOperationID is empty, skipping (there is no runtime)")
		return operation, 0, nil
	}

	instance, err := s.instanceStorage.GetByID(operation.InstanceID)
	switch {
	case err == nil:
	case dberr.IsNotFound(err):
		log.Infof("instance already deleted", err)
		return operation, 0 * time.Second, nil
	default:
		log.Errorf("unable to get instance from storage: %s", err)
		return operation, 1 * time.Second, nil
	}

	status, err := s.provisionerClient.RuntimeOperationStatus(instance.GlobalAccountID, operation.ProvisionerOperationID)
	if err != nil {
		log.Errorf("call to provisioner RuntimeOperationStatus failed: %s, GlobalAccountID=%s, Provisioner OperationID=%s", err.Error(), instance.GlobalAccountID, operation.ProvisionerOperationID)
		return operation, 1 * time.Minute, nil
	}
	var msg string
	if status.Message != nil {
		msg = *status.Message
	}
	log.Infof("call to provisioner returned %s status: %s", status.State.String(), msg)

	switch status.State {
	case gqlschema.OperationStateSucceeded:
		msg := fmt.Sprintf("Provisioner succeeded in %s.", time.Since(operation.UpdatedAt))
		log.Info(msg)
		operation.EventInfof(msg)
		return operation, 0, nil
	case gqlschema.OperationStateInProgress:
		return operation, 30 * time.Second, nil
	case gqlschema.OperationStatePending:
		return operation, 30 * time.Second, nil
	case gqlschema.OperationStateFailed:
		lastErr := provisioner.OperationStatusLastError(status.LastError)
		return s.operationManager.OperationFailed(operation, "provisioner client returns failed status", lastErr, log)
	}

	lastErr := provisioner.OperationStatusLastError(status.LastError)
	return s.operationManager.OperationFailed(operation, fmt.Sprintf("unsupported provisioner client status: %s", status.State.String()), lastErr, log)
}
