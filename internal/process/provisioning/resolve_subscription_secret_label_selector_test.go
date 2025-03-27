package provisioning

import (
	"golang.org/x/exp/maps"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kyma-project/kyma-environment-broker/common/gardener"
	"github.com/kyma-project/kyma-environment-broker/common/hyperscaler"
	"github.com/kyma-project/kyma-environment-broker/common/hyperscaler/rules"
	"github.com/kyma-project/kyma-environment-broker/internal"
	"github.com/kyma-project/kyma-environment-broker/internal/broker"
	"github.com/kyma-project/kyma-environment-broker/internal/fixture"
	"github.com/kyma-project/kyma-environment-broker/internal/process"
	"github.com/kyma-project/kyma-environment-broker/internal/ptr"
	"github.com/kyma-project/kyma-environment-broker/internal/storage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/sets"
	k8sTesting "k8s.io/client-go/testing"
)

func TestSecretBindingLabelSelector(t *testing.T) {
	/**
	  This test checks if the label selector used by new 3s step implementation is the same as the old one produces.
	*/

	for name, tc := range map[string]struct {
		providerType      string
		planID            string
		hyperscalerRegion string
		platformRegion    string
	}{
		"aws": {
			providerType:      "aws",
			planID:            broker.AWSPlanID,
			hyperscalerRegion: "us-east-1",
			platformRegion:    "cf-eu10",
		},
		"aws(PR=cf-eu11)": {
			providerType:      "aws",
			planID:            broker.AWSPlanID,
			hyperscalerRegion: "eu-central-1",
			platformRegion:    "cf-eu11",
		},
		"build-runtime-aws": {
			providerType:      "aws",
			planID:            broker.BuildRuntimeAWSPlanID,
			hyperscalerRegion: "ap-northeast-1",
			platformRegion:    "cf-us11",
		},
		"build-runtime-aws(PR=cf-eu11)": {
			providerType:      "aws",
			planID:            broker.BuildRuntimeAWSPlanID,
			hyperscalerRegion: "eu-central-1",
			platformRegion:    "cf-eu11",
		},
		"azure": {
			providerType:      "azure",
			planID:            broker.AzurePlanID,
			hyperscalerRegion: "eastus",
			platformRegion:    "cf-eu20",
		},
		"azure(PR=cf-ch20)": {
			providerType:      "azure",
			planID:            broker.AzurePlanID,
			hyperscalerRegion: "switzerlandnorth",
			platformRegion:    "cf-ch20",
		},
		"build-runtime-azure": {
			providerType:      "azure",
			planID:            broker.BuildRuntimeAzurePlanID,
			hyperscalerRegion: "westus2",
			platformRegion:    "cf-us20",
		},
		"build-runtime-azure(PR=cf-ch20)": {
			providerType:      "azure",
			planID:            broker.BuildRuntimeAzurePlanID,
			hyperscalerRegion: "switzerlandnorth",
			platformRegion:    "cf-ch20",
		},
		"gcp": {
			providerType:      "gcp",
			planID:            broker.GCPPlanID,
			hyperscalerRegion: "asia-southeast1",
			platformRegion:    "cf-us21",
		},
		"gcp(PR=cf-sa30)": {
			providerType:      "gcp",
			planID:            broker.GCPPlanID,
			hyperscalerRegion: "me-central2",
			platformRegion:    "cf-sa30",
		},
		"build-runtime-gcp": {
			providerType:      "gcp",
			planID:            broker.BuildRuntimeGCPPlanID,
			hyperscalerRegion: "us-central1",
			platformRegion:    "cf-us21",
		},
		"build-runtime-gcp(PR=cf-sa30)": {
			providerType:      "gcp",
			planID:            broker.BuildRuntimeGCPPlanID,
			hyperscalerRegion: "me-central2",
			platformRegion:    "cf-sa30",
		},
		"trial": {
			providerType:      "aws",
			planID:            broker.TrialPlanID,
			hyperscalerRegion: "eu-central-1",
			platformRegion:    "cf-eu10",
		},
		"trial-azure": {
			providerType:      "azure",
			planID:            broker.TrialPlanID,
			hyperscalerRegion: "eastus",
			platformRegion:    "cf-us20",
		},
		"sap-converged-cloud": {
			providerType:      "openstack",
			planID:            broker.SapConvergedCloudPlanID,
			hyperscalerRegion: "eu-de-1",
			platformRegion:    "cf-eu20-staging",
		},
		"azure_lite": {
			providerType:      "azure",
			planID:            broker.AzureLitePlanID,
			hyperscalerRegion: "westus",
			platformRegion:    "cf-us20",
		},
		"free": {
			providerType:      "aws",
			planID:            broker.FreemiumPlanID,
			hyperscalerRegion: "eu-central-1",
			platformRegion:    "cf-eu10",
		},
		"preview": {
			providerType:      "aws",
			planID:            broker.PreviewPlanID,
			hyperscalerRegion: "eu-central-1",
			platformRegion:    "cf-eu10",
		},
	} {
		t.Run(name, func(t *testing.T) {
			var referenceSelector *string = ptr.String("")
			var gotSelector *string = ptr.String("")

			// given
			operation := fixProvisioningOperation(tc.hyperscalerRegion, tc.platformRegion, tc.planID, tc.providerType)

			// todo: replace running old step by hardcoded expected label selector (a reference) when the old step is removed
			step := fixResolveCredentialStep(t, referenceSelector, operation)

			_, backoff, _ := step.Run(operation, fixLogger())
			// after the old step is run - we have referenceSelector with a reference value, which is used for an assertion

			require.Zero(t, backoff)
			newStep := fixNewResolveStep(t, gotSelector, operation)

			// when
			_, backoff, _ = newStep.Run(operation, fixLogger())
			require.Zero(t, backoff)

			// then
			assertSelectors(t, referenceSelector, gotSelector)
		})
	}

}

func assertSelectors(t *testing.T, expected *string, got *string) {
	t.Helper()
	t.Log("expectedSet: ", *expected)
	assert.NotZerof(t, len(*expected), "expected selector is empty")
	expectedParts := strings.Split(*expected, ",")
	expectedSet := sets.New(expectedParts...)

	gotParts := strings.Split(*got, ",")
	gotSet := sets.New(gotParts...)

	if expectedSet.Has("shared=true") {
		assert.True(t, gotSet.Has("!euAccess"))
		gotSet.Delete("!euAccess")
	}

	assert.Equal(t, expectedSet, gotSet)
}

func dummySecretBinding() gardener.SecretBinding {
	name := uuid.New().String()
	sb := gardener.SecretBinding{}
	sb.SetName(name)
	sb.SetNamespace(namespace)
	sb.SetSecretRefName(name)
	return sb
}

func savingLabelSelectorReactor(selector *string) func(action k8sTesting.Action) (bool, runtime.Object, error) {
	return func(action k8sTesting.Action) (bool, runtime.Object, error) {
		labelSelector := action.(k8sTesting.ListActionImpl).GetListRestrictions().Labels.String()
		*selector = labelSelector

		labels := map[string]string{}
		requirements, _ := action.(k8sTesting.ListActionImpl).GetListRestrictions().Labels.Requirements()
		for _, r := range requirements {
			if len(r.Values()) > 0 {
				labels[r.Key()] = r.Values().List()[0]
			}
		}
		sb := dummySecretBinding()
		sb.SetLabels(labels)
		listToReturn := &unstructured.UnstructuredList{
			Items: []unstructured.Unstructured{sb.Unstructured},
		}
		return true, listToReturn, nil
	}
}

func fixProvisioningOperation(region string, platformRegion string, planID, provider string) internal.Operation {
	operation := fixture.FixProvisioningOperation("op-id", "inst-id")
	operation.ProvisioningParameters.Parameters.Region = ptr.String(region)
	operation.ProvisioningParameters.PlatformRegion = platformRegion
	operation.ProvisioningParameters.PlanID = planID
	operation.ProviderValues = &internal.ProviderValues{
		ProviderType: provider,
		Region:       region,
	}
	return operation
}

func fixResolveCredentialStep(t *testing.T, selector *string, operation internal.Operation) *ResolveCredentialsStep {
	gardenerK8sClient := gardener.NewDynamicFakeClient()
	gardenerK8sClient.PrependReactor("list",
		gardener.SecretBindingResource.Resource,
		savingLabelSelectorReactor(selector))
	memoryStorage := storage.NewMemoryStorage()
	accountProvider := hyperscaler.NewAccountProvider(hyperscaler.NewAccountPool(gardenerK8sClient, namespace), hyperscaler.NewSharedGardenerAccountPool(gardenerK8sClient, namespace))
	step := NewResolveCredentialsStep(memoryStorage.Operations(), accountProvider, &rules.RulesService{})
	err := memoryStorage.Operations().InsertOperation(operation)
	require.NoError(t, err)
	return step
}

func fixNewResolveStep(t *testing.T, selector *string, operation internal.Operation) process.Step {
	memoryStorage := storage.NewMemoryStorage()
	err := memoryStorage.Operations().InsertOperation(operation)
	require.NoError(t, err)
	gardenerK8sClient := gardener.NewDynamicFakeClient()
	gardenerK8sClient.PrependReactor("list",
		gardener.SecretBindingResource.Resource,
		savingLabelSelectorReactor(selector))
	return NewResolveSubscriptionSecretStep(memoryStorage.Operations(), gardener.NewClient(gardenerK8sClient, namespace), rulesService(t),
		internal.RetryTuple{Timeout: time.Millisecond, Interval: time.Millisecond})
}

func rulesService(t *testing.T) *rules.RulesService {
	content := `rule:
      - aws                             # pool: hyperscalerType: aws
      - aws(PR=cf-eu11) -> EU           # pool: hyperscalerType: aws; euAccess: true
      - build-runtime-aws               # pool: hyperscalerType: aws
      - build-runtime-aws(PR=cf-eu11) -> EU # pool: hyperscalerType: aws; euAccess: true
      - azure                           # pool: hyperscalerType: azure
      - azure(PR=cf-ch20) -> EU         # pool: hyperscalerType: azure; euAccess: true
      - build-runtime-azure             # pool: hyperscalerType: azure
      - build-runtime-azure(PR=cf-ch20) -> EU # pool: hyperscalerType: azure; euAccess: true
      - gcp                             # pool: hyperscalerType: gcp
      - gcp(PR=cf-sa30) -> PR           # pool: hyperscalerType: gcp_cf-sa30
      - build-runtime-gcp               # pool: hyperscalerType: gcp
      - build-runtime-gcp(PR=cf-sa30) -> PR # pool: hyperscalerType: gcp_cf-sa30
      - trial -> S                      # pool: hyperscalerType: azure; shared: true - TRIAL POOL, pool: hyperscalerType: aws; shared: true - TRIAL POOL
      - sap-converged-cloud -> S, HR    # pool: hyperscalerType: openstack_<HYPERSCALER_REGION>; shared: true
      - azure_lite                      # pool: hyperscalerType: azure
      - free
      - preview`
	tmpfile, err := rules.CreateTempFile(content)
	require.NoError(t, err)
	defer os.Remove(tmpfile)

	rs, err := rules.NewRulesServiceFromFile(tmpfile, sets.New(maps.Keys(broker.PlanIDsMapping)...), sets.New("aws", "azure", "gcp", "trial", "build-runtime-aws", "build-runtime-azure", "build-runtime-gcp", "sap-converged-cloud", "azure_lite", "free", "preview"))
	require.NoError(t, err)

	return rs
}
