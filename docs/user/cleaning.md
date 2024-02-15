# Soliution 1: No additional entity

## Entities

### Instances 

removed at the end of deprovisioning, no action needed\

### Operations

currently only the last one has removed `userID`

Plan:
1. Remove  if the last deprovisioning was created before 3 months ago (`retention period`)
2. Sanitize if the last deprovisioning was created before 4 weeks ago and such operation was not cleaned

Technical detail: operations will have a new column, a flag which indicates if it was sanitized

### Runtime States

Plan:
1. Remove all data saved in RuntimeStates for deprovisioned instances

## Pros & Cons

Pros:
1. No changes in KCP CLI

Cons:
1. there is a need of a job which periodically runs the cleaning, such job must be maintained etc.
2. additional column in the DB to mark cleaned operations

# Solution 2: New entity "instance archived"

New entity is created at the end of a full succeeded deprovisioning based on data stored in regular entities.
We can store the following data (which is not sensitive):

1. InstanceID
2. GlobalAccountID
2. SubaccountID
3. SubscriptionGlobalAccountID
4. PlanID (and/or PlanName)
5. SAPRegion, SubAccountRegion (for example `cf-eu10`)
6. HyperscalerRegion (for example `eastus`)
7. CloudProvider (for example `azure`, `aws` etc.)
7. ProvisioningStartedAt
8. ProvisioningFinishedAt
9. FirstDeprovisioningStartedAt
10. FirstDeprovisioningFinishedAt
11. LastDeprovisioningFinishedAt
12. LastRuntimeID (for Trials we can have several runtime IDs, for other plans only one RuntimeID)
13. InternalUser (true, if user ID is in the `sap.com` domain)
14. ShootName

additionally we can add some summary data like:

1. NumberOfUpdates
2. NumberOfUpgradeCluster
3. NumberOfUpgradeKyma
4. NumberOfUpdates
4. NumberOfDeprovisionings
5. NumberOfSuspensions
6. NumberOfUnsuspensions
7. LastErrorMessage (this filed I'd like to introduce to the operation JSON to have better tracking errors)

## Pros & Cons

Pros:
1. No risk of accidential deletion of operations in the future (there is only one of batch deletion - migration)
2. We have clear definition which fields are archived. No risk saving some data in provisioning parameters or ERS context introduced in the future.
3. less data in a regular entities - faster SQL queries

Cons:
1. runtimes endpoint must be changed to fetch data from new entity


# This is the current output of kcp cli

Existing instance:
```json
{
      "instanceID": "a4ddcc32-6925-477c-a708-19ebbca11604",
      "runtimeID": "cbd99494-a81c-4875-aa32-2bfed0db34bc",
      "globalAccountID": "d9994f8f-7e46-42a8-b2c1-1bfff8d2fe05",
      "subscriptionGlobalAccountID": "",
      "subAccountID": "39ba9a66-2c1a-4fe4-a28e-6e5db434084e",
      "region": "europe-west3",
      "subAccountRegion": "cf-eu10",
      "shootName": "c-57b9388",
      "serviceClassID": "47c9dcbf-ff30-448e-ab36-d3bad66ba281",
      "serviceClassName": "kymaruntime",
      "servicePlanID": "ca6e5357-707f-4565-bbbd-b3ab732597c6",
      "servicePlanName": "gcp",
      "provider": "GCP",
      "status": {
        "createdAt": "2024-02-14T10:07:43.339394Z",
        "modifiedAt": "2024-02-14T10:14:52.553881Z",
        "state": "provisioning",
        "provisioning": {
          "state": "in progress",
          "description": "Operation created",
          "createdAt": "2024-02-14T10:07:43.334317Z",
          "updatedAt": "2024-02-14T10:14:52.553881Z",
          "operationID": "a66550c8-72f3-4b39-adb0-6b55e2f1bb44",
          "finishedStages": [
            "start"
          ],
          "runtimeVersion": "2.20.0"
        }
      },
      "userID": "test@test.com",
      "avsInternalEvaluationID": 385209573
    }
```

Deprovisioned:
```json
{
      "instanceID": "98b4f92f-748e-4c6b-a5c9-599a3a062abc",
      "runtimeID": "bd744044-7ceb-487b-9dc2-d1c0361d09d4",
      "globalAccountID": "e449f875-b5b2-4485-b7c0-98725c0571bf",
      "subscriptionGlobalAccountID": "",
      "subAccountID": "github-actions-keb-integration",
      "region": "",
      "subAccountRegion": "cf-eu10",
      "shootName": "c-1982ee1",
      "serviceClassID": "",
      "serviceClassName": "",
      "servicePlanID": "4deee563-e5ec-4731-b9b1-53b42d855f0c",
      "servicePlanName": "",
      "provider": "",
      "status": {
        "createdAt": "2024-02-13T00:32:25.932915Z",
        "modifiedAt": "2024-02-13T01:36:05.291172Z",
        "deletedAt": "2024-02-13T01:36:05.291172Z",
        "state": "deprovisioned",
        "deprovisioning": {
          "state": "succeeded",
          "description": "Processing finished",
          "createdAt": "2024-02-13T01:13:33.393309Z",
          "updatedAt": "2024-02-13T01:36:05.291172Z",
          "operationID": "b087644b-42b5-4658-b0b7-b1bea0f530fb",
          "finishedStages": [
            "Initialisation",
            "BTPOperator_Cleanup",
            "De-provision_AVS_Evaluations",
            "EDP_Deregistration",
            "Delete_Kyma_Resource",
            "Check_Kyma_Resource_Deleted",
            "Deregister_Cluster",
            "Check_Cluster_Deregistration",
            "Delete_GardenerCluster",
            "Check_GardenerCluster_Deleted",
            "Remove_Runtime",
            "Check_Runtime_Removal",
            "Release_Subscription",
            "Delete_Kubeconfig",
            "Remove_Instance"
          ],
          "runtimeVersion": ""
        }
      },
      "userID": "",
      "avsInternalEvaluationID": 384740555
    }
```
