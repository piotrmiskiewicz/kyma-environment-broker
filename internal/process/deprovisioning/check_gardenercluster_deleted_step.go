package deprovisioning

import (
	"context"
	"time"

	"github.com/kyma-project/kyma-environment-broker/internal/process/steps"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/sirupsen/logrus"

	"github.com/kyma-project/kyma-environment-broker/internal"
	"github.com/kyma-project/kyma-environment-broker/internal/process"
	"github.com/kyma-project/kyma-environment-broker/internal/storage"
)

type CheckGardenerClusterDeletedStep struct {
	operationManager *process.OperationManager
	kcpClient        client.Client
}

func NewCheckGardenerClusterDeletedStep(operations storage.Operations, kcpClient client.Client) *CheckGardenerClusterDeletedStep {
	return &CheckGardenerClusterDeletedStep{
		operationManager: process.NewOperationManager(operations),
		kcpClient:        kcpClient,
	}
}

func (step *CheckGardenerClusterDeletedStep) Name() string {
	return "Check_Kyma_Resource_Deleted"
}

func (step *CheckGardenerClusterDeletedStep) Run(operation internal.Operation, logger logrus.FieldLogger) (internal.Operation, time.Duration, error) {
	if operation.KymaResourceNamespace == "" {
		logger.Warnf("namespace for Kyma resource not specified")
		return operation, 0, nil
	}
	resourceName := operation.GardenerClusterName
	if resourceName == "" {
		logger.Infof("GardenerCluster resource name is empty, skipping")
		return operation, 0, nil
	}

	gardenerClusterUnstructured := &unstructured.Unstructured{}
	gardenerClusterUnstructured.SetGroupVersionKind(steps.GardenerClusterGVK())
	err := step.kcpClient.Get(context.Background(), client.ObjectKey{
		Namespace: operation.KymaResourceNamespace,
		Name:      resourceName,
	}, gardenerClusterUnstructured)

	if err == nil {
		logger.Infof("GardenerCluster resource still exists")
		return step.operationManager.RetryOperationWithoutFail(operation, step.Name(), "GardenerCluster resource still exists", 5*time.Second, 15*time.Minute, logger)
	}

	if !errors.IsNotFound(err) {
		logger.Errorf("unable to check GardenerCluster resource existence: %s", err)
		return step.operationManager.RetryOperationWithoutFail(operation, step.Name(), "unable to check GardenerCluster resource existence", backoffForK8SOperation, timeoutForK8sOperation, logger)
	}

	return operation, 0, nil
}
