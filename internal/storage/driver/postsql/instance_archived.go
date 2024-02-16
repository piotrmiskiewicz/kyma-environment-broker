package postsql

import (
	"github.com/kyma-project/kyma-environment-broker/internal"
	"github.com/kyma-project/kyma-environment-broker/internal/storage/dbmodel"
	"github.com/kyma-project/kyma-environment-broker/internal/storage/postsql"
)

type instanceArchivedPostgreSQLStorage struct {
	factory postsql.Factory
}

func NewInstanceArchivedPostgreSQLStorage(sess postsql.Factory) *instanceArchivedPostgreSQLStorage {
	return &instanceArchivedPostgreSQLStorage{
		factory: sess,
	}
}

func (s *instanceArchivedPostgreSQLStorage) GetByInstanceID(instanceId string) (internal.InstanceArchived, error) {
	dto, err := s.factory.NewReadSession().GetInstanceArchivedByID(instanceId)
	return dbmodel.NewInstanceArchivedFromDTO(dto), err
}

func (s *instanceArchivedPostgreSQLStorage) Insert(instance internal.InstanceArchived) error {
	return s.factory.NewWriteSession().InsertInstanceArchived(dbmodel.NewInstanceArchivedDTO(instance))
}
