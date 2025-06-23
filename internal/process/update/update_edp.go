package update

import (
	"fmt"
	"github.com/kyma-project/kyma-environment-broker/internal"
	"github.com/kyma-project/kyma-environment-broker/internal/edp"
	"github.com/kyma-project/kyma-environment-broker/internal/process"
	"github.com/kyma-project/kyma-environment-broker/internal/process/provisioning"
	"log/slog"
	"strings"
	"time"
)

const (
	edpRetryInterval = 30 * time.Second
	edpRetryTimeout  = 20 * time.Minute
)

type EDPUpdater interface {
	UpdateMetadataTenant(name, env string, metadataKey, metadataValue string, log *slog.Logger) error
}

type EDPUpdateStep struct {
	client           EDPUpdater
	config           edp.Config
	operationManager *process.OperationManager
}

func (s *EDPUpdateStep) Name() string {
	return "EDP_Update"
}

func (s *EDPUpdateStep) Run(operation internal.Operation, log *slog.Logger) (internal.Operation, time.Duration, error) {
	if operation.UpdatedPlanID == "" {
		log.Info("No plan update detected, skipping EDP update step")
		return operation, 0, nil
	}
	maasConsumerServicePlan := provisioning.SelectServicePlan(operation.UpdatedPlanID)
	subAccountID := strings.ToLower(operation.ProvisioningParameters.ErsContext.SubAccountID)

	err := s.client.UpdateMetadataTenant(subAccountID, s.config.Environment, edp.MaasConsumerServicePlan, maasConsumerServicePlan, log)

	if err != nil {
		log.Warn(fmt.Sprintf("request to EDP failed: %s. Retry...", err))
		if s.config.Required {
			return s.operationManager.RetryOperation(operation, "request to EDP failed", err, edpRetryInterval, edpRetryTimeout, log)
		} else {
			return s.operationManager.RetryOperationWithoutFail(operation, s.Name(), "request to EDP failed", edpRetryInterval, edpRetryTimeout, log, err)
		}
	}
	log.Info(fmt.Sprintf("EDP metadata updated for subaccount %s with plan %s", subAccountID, maasConsumerServicePlan))
	return operation, 0, nil
}
