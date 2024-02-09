package cleaning

import (
    "fmt"
    "github.com/kyma-project/kyma-environment-broker/internal"
    "github.com/kyma-project/kyma-environment-broker/internal/storage"
    "github.com/kyma-project/kyma-environment-broker/internal/storage/dberr"
    "github.com/sirupsen/logrus"
    "time"
)

type Service struct {
    instances     storage.Instances
    operations    storage.Operations
    runtimeStates storage.RuntimeStates

    dryRun              bool
    performDeletion     bool
    performSanitization bool

    sanitizationDuration time.Duration
    retentionDuration    time.Duration

    log logrus.FieldLogger
}

func (s *Service) Run() error {
    s.log.Info("Starting Operation and RuntimeStates Cleaning process")
    instanceIDs, err := s.operations.FindDeletedInstanceIDs()
    if err != nil {
        s.log.Errorf("Unable to get instance IDs: %s", err.Error())
        return err
    }
    s.log.Infof("Got %d instance IDs to process", len(instanceIDs))

    for _, instanceId := range instanceIDs {
        logger := s.log.WithField("instanceID", instanceId)
        // check if the instance really does not exists
        instance, errInstance := s.instances.GetByID(instanceId)
        if err == nil {
            logger.Errorf("the instance (createdAt: %s, planName: %s) still exists, aborting the process",
                instance.InstanceID, instance.CreatedAt, instance.ServicePlanName)
            return fmt.Errorf("instance exists")
        }
        if !dberr.IsNotFound(errInstance) {
            return errInstance
        }


        operations, err := s.operations.ListOperationsByInstanceID(instanceId)
        if err != nil {
            logger.Errorf("unable to get operations: %s", err.Error())
        }
        if len(operations) == 0 {
            logger.Warnf("operations not found")
        }
        // operations are sorted by date, so the first one should be "Deprovision"
        if operations[0].Type != internal.OperationTypeDeprovision {
            logger.Errorf("the last operation is not deprovision")
        }
        lastDeprovisioningFinishedAt := operations[0].UpdatedAt

        if time.Since(lastDeprovisioningFinishedAt) < s.sanitizationDuration {
            time.Since(lastDeprovisioningFinishedAt).Round(24 * time.Hour)
            logger.Infof("instance was deprovisioned %s")
        }

        //for _, operation := range operations {
        //
        //}



    }

}
