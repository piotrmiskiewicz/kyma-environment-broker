package update

import (
	"context"
	imv1 "github.com/kyma-project/infrastructure-manager/api/v1"
	"github.com/kyma-project/kyma-environment-broker/internal"
	"github.com/kyma-project/kyma-environment-broker/internal/process"
	"github.com/kyma-project/kyma-environment-broker/internal/storage"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

type UpdateRuntimeStep struct {
	operationManager *process.OperationManager
	k8sClient        client.Client
}

func NewUpdateRuntimeStep(os storage.Operations, k8sClient client.Client) *UpdateRuntimeStep {
	return &UpdateRuntimeStep{
		operationManager: process.NewOperationManager(os),
		k8sClient:        k8sClient,
	}
}

func (s *UpdateRuntimeStep) Name() string {
	return "Update_Runtime_Resource"
}

func (s *UpdateRuntimeStep) Run(operation internal.Operation, log logrus.FieldLogger) (internal.Operation, time.Duration, error) {
	// Check if the runtime exists

	var runtime = imv1.Runtime{}
	err := s.k8sClient.Get(context.Background(), client.ObjectKey{Name: operation.GetRuntimeResourceName(), Namespace: operation.GetRuntimeResourceNamespace()}, &runtime)
	if errors.IsNotFound(err) {
		// todo: after the switch to KIM, this should throw an error
		log.Infof("Runtime not found, skipping")
		return operation, 0, nil
	}

	// Update the runtime
	if operation.UpdatingParameters.MachineType != nil {
		runtime.Spec.Shoot.Provider.Workers[0].Machine.Type = *operation.UpdatingParameters.MachineType
	}
	if operation.UpdatingParameters.AutoScalerMin != nil {
		runtime.Spec.Shoot.Provider.Workers[0].Minimum = int32(*operation.UpdatingParameters.AutoScalerMin)
	}
	if operation.UpdatingParameters.AutoScalerMax != nil {
		runtime.Spec.Shoot.Provider.Workers[0].Maximum = int32(*operation.UpdatingParameters.AutoScalerMax)
	}
	if operation.UpdatingParameters.MaxSurge != nil {
		runtime.Spec.Shoot.Provider.Workers[0].Maximum = int32(*operation.UpdatingParameters.MaxSurge)
	}
	if operation.UpdatingParameters.MaxUnavailable != nil {
		v := intstr.FromInt32(int32(*operation.UpdatingParameters.MaxUnavailable))
		runtime.Spec.Shoot.Provider.Workers[0].MaxUnavailable = &v
	}
	if operation.UpdatingParameters.OIDC != nil {
		input := operation.UpdatingParameters.OIDC
		if len(input.SigningAlgs) > 0 {
			runtime.Spec.Shoot.Kubernetes.KubeAPIServer.OidcConfig.SigningAlgs = input.SigningAlgs
		}
		if input.ClientID != "" {
			runtime.Spec.Shoot.Kubernetes.KubeAPIServer.OidcConfig.ClientID = &input.ClientID
		}
		if input.IssuerURL != "" {
			runtime.Spec.Shoot.Kubernetes.KubeAPIServer.OidcConfig.IssuerURL = &input.IssuerURL
		}
		if input.GroupsClaim != "" {
			runtime.Spec.Shoot.Kubernetes.KubeAPIServer.OidcConfig.GroupsClaim = &input.GroupsClaim
		}
		if input.UsernamePrefix != "" {
			runtime.Spec.Shoot.Kubernetes.KubeAPIServer.OidcConfig.UsernamePrefix = &input.UsernamePrefix
		}
		if input.UsernameClaim != "" {
			runtime.Spec.Shoot.Kubernetes.KubeAPIServer.OidcConfig.UsernameClaim = &input.UsernameClaim
		}

	}

	err = s.k8sClient.Update(context.Background(), &runtime)
	if err != nil {
		return s.operationManager.RetryOperation(operation, "unable to update runtime", err, 10*time.Second, 1*time.Minute, log)
	}

	return operation, 0, nil
}
