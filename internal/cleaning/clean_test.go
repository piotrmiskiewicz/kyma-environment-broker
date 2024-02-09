package cleaning

import (
    "fmt"
    "github.com/kyma-project/kyma-environment-broker/internal/events"
    "github.com/kyma-project/kyma-environment-broker/internal/storage"
    "github.com/kyma-project/kyma-environment-broker/internal/storage/dberr"
    "github.com/sirupsen/logrus"
    "github.com/stretchr/testify/require"
    "os"
    "testing"
    "time"
)

func TestClean(t *testing.T) {

    cfg := storage.Config{
        User:            "broker",
        Password:        os.Getenv("DB_PASSWORD"),
        Host:            "localhost",
        Port:            "5432",
        Name:            "broker",
        SSLMode:         "disable",
        SecretKey:       os.Getenv("DB_SECRET"),
        MaxOpenConns:    1,
        MaxIdleConns:    1,
        ConnMaxLifetime: time.Minute,
    }

    s, _, err := storage.NewFromConfig(cfg, events.Config{}, storage.NewEncrypter(cfg.SecretKey), logrus.StandardLogger())

    require.NoError(t, err)

    // warning: think about paging etc.
    ids, err := s.Operations().FindDeletedInstanceIDs()
    require.NoError(t, err)

    max := 10
    for i, instanceId := range ids {
        fmt.Println()
        fmt.Printf("InstanceID=%s\n", instanceId)
        _, errInstance := s.Instances().GetByID(instanceId)
        if !dberr.IsNotFound(errInstance) {
            panic("Instance exists!!!")
        }



        operations, _ := s.Operations().ListOperationsByInstanceID(instanceId)
        // check if deprovisioning exists ???


        for _, operation := range operations {
            fmt.Printf("\t\toperation=%s\t%s\t%s\n", operation.ID, operation.Type, operation.CreatedAt)


            runtimeStates, _ := s.RuntimeStates().ListByOperationID(operation.ID)
            for _, rs := range runtimeStates {
                fmt.Printf("\t\t\t\t\tRuntimeState: %s %s\n", rs.ID, rs.CreatedAt)
            }
        }

        if i > max {
            break
        }
    }
}

func TestR(t *testing.T) {
    fmt.Println((70 * time.Hour).Round(24 * time.Hour))

    time.
}
