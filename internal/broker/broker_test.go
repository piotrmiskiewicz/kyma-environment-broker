package broker_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/kyma-incubator/compass/components/kyma-environment-broker/internal/broker"
	"github.com/kyma-incubator/compass/components/kyma-environment-broker/internal/provisioner"
	schema "github.com/kyma-incubator/compass/components/provisioner/pkg/gqlschema"
	"github.com/pivotal-cf/brokerapi/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	serviceID = "47c9dcbf-ff30-448e-ab36-d3bad66ba281"
	planID    = "4deee563-e5ec-4731-b9b1-53b42d855f0c"
)

func TestBroker_Services(t *testing.T) {
	// given
	broker, err := broker.NewBroker(nil, broker.ProvisioningConfig{})
	require.NoError(t, err)

	// when
	services, err := broker.Services(context.TODO())

	// then
	require.NoError(t, err)
	assert.Len(t, services, 1)
	assert.Len(t, services[0].Plans, 1)
}

func TestBroker_ProvisioningScenario(t *testing.T) {
	// given
	const instID = "inst-id"
	const clusterName = "cluster-testing"
	fCli := provisioner.NewFakeClient()
	broker, err := broker.NewBroker(fCli, broker.ProvisioningConfig{})
	require.NoError(t, err)

	// when
	res, err := broker.Provision(context.TODO(), instID, domain.ProvisionDetails{
		ServiceID:     serviceID,
		PlanID:        planID,
		RawParameters: json.RawMessage(fmt.Sprintf(`{"name": "%s"}`, clusterName)),
		RawContext:    json.RawMessage(`{}`),
	}, true)
	require.NoError(t, err)

	// then
	assert.Equal(t, clusterName, fCli.GetProvisionRuntimeInput(0).ClusterConfig.GardenerConfig.Name)

	// when
	op, err := broker.LastOperation(context.TODO(), instID, domain.PollDetails{
		ServiceID:     serviceID,
		PlanID:        planID,
		OperationData: res.OperationData,
	})

	// then
	require.NoError(t, err)
	assert.Equal(t, domain.InProgress, op.State)

	// when
	fCli.FinishProvisionerOperation(res.OperationData, schema.OperationStateSucceeded)
	op, err = broker.LastOperation(context.TODO(), instID, domain.PollDetails{
		ServiceID:     serviceID,
		PlanID:        planID,
		OperationData: res.OperationData,
	})

	// then
	require.NoError(t, err)
	assert.Equal(t, domain.Succeeded, op.State)
}
