# Service Plan Updates

Kyma Environment Broker supports updating service plan. This feature allows to change the plan of an existing Kyma runtime. Only some plan changes are possible because the new plan must use the same provider (it is not possible to change AWS to Azure).

# Configuration

To enable plan changes you must do the following steps:

1. Enable the feature by setting the value: `enablePlanUpgrades: true`
2. Define allowed plan changes in the plans configuration, for example:
```yaml
plansConfiguration:
  gcp:
    upgradableToPlans:
      - build-runtime-gcp
```

# Plan update request

The plan update request is similar to regular update request. The important part is the field `plan_id` which have the target plan ID, for example:

```http
PATCH /oauth/v2/service_instances/"{INSTANCE_ID}"?accepts_incomplete=true
{
    “service_id”: “47c9dcbf-ff30-448e-ab36-d3bad66ba281", //Kyma ID
    “plan_id": “{TARGET_PLAN_ID}”
}
```

When the plan update is not allowed, the response is HTTP 400 Bad Request.