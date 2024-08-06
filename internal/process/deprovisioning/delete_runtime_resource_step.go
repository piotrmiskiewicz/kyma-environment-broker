package deprovisioning

import (
	"context"
	"github.com/kyma-project/kyma-environment-broker/internal"
	"github.com/kyma-project/kyma-environment-broker/internal/process"
	"github.com/kyma-project/kyma-environment-broker/internal/storage"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"

	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
)

const (
	timeoutForRuntimeDeletion = 10 * time.Minute
)

type DeleteRuntimeResourceStep struct {
	operationManager *process.OperationManager
	kcpClient        client.Client
}

func NewDeleteRuntimeResourceStep(operations storage.Operations, kcpClient client.Client) *DeleteRuntimeResourceStep {
	return &DeleteRuntimeResourceStep{
		operationManager: process.NewOperationManager(operations),
		kcpClient:        kcpClient,
	}
}

func (step *DeleteRuntimeResourceStep) Name() string {
	return "Delete_Runtime_Resource"
}

func (step *DeleteRuntimeResourceStep) Run(operation internal.Operation, logger logrus.FieldLogger) (internal.Operation, time.Duration, error) {
	resourceName := operation.RuntimeResourceName
	resourceNamespace := operation.KymaResourceNamespace

	if resourceName == "" {
		logger.Infof("Runtime resource name is empty, skipping")
		return operation, 0, nil
	}
	if resourceNamespace == "" {
		logger.Warnf("Namespace for Runtime resource not specified")
		return operation, 0, nil
	}

	runtime := &imv1.Runtime{
		ObjectMeta: v1.ObjectMeta{
			Name:      resourceName,
			Namespace: resourceNamespace,
		},
	}

	err := step.kcpClient.Delete(context.Background(), runtime)

	// check the error
	if err != nil {
		if meta.IsNoMatchError(err) {
			logger.Info("No CRD installed, skipping")
			return operation, 0, nil
		}

		// if the resource is not found, log it and return (it is not a problem)
		if errors.IsNotFound(err) {
			logger.Info("Runtime resource deleted")
			return operation, 0, nil
		} else {
			logger.Warnf("unable to delete the Runtime resource: %s", err)
			return step.operationManager.RetryOperationWithoutFail(operation, step.Name(), "unable to delete the Runtime resource", backoffForK8SOperation, timeoutForK8sOperation, logger, err)
		}
	}

	return operation, 0, nil
}