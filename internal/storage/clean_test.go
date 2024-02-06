package storage

import (
    "fmt"
    "github.com/kyma-project/kyma-environment-broker/internal/events"
    "github.com/sirupsen/logrus"
    "github.com/stretchr/testify/require"
    "testing"
    "time"
)

func TestClean(t *testing.T) {

    cfg := Config{
        User:            "broker",
        Password:        "ckXajBpeRyBh",
        Host:            "localhost",
        Port:            "5432",
        Name:            "broker",
        SSLMode:         "disable",
        SecretKey:       "################################",
        MaxOpenConns:    1,
        MaxIdleConns:    1,
        ConnMaxLifetime: time.Minute,
    }

    s, _, err := NewFromConfig(cfg, events.Config{}, NewEncrypter(cfg.SecretKey), logrus.StandardLogger())

    require.NoError(t, err)

    ids, err := s.Operations().FindDeletedInstanceIDs()
    require.NoError(t, err)

    fmt.Println(ids)
}
