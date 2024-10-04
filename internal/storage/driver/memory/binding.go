package memory

import (
	"sync"

	"github.com/kyma-project/kyma-environment-broker/internal"
	"github.com/kyma-project/kyma-environment-broker/internal/storage/dberr"
)

type Binding struct {
	mu   sync.Mutex
	data map[string]internal.Binding
}

func NewBinding() *Binding {
	return &Binding{
		data: make(map[string]internal.Binding),
	}
}

func (s *Binding) GetByBindingID(bindingId string) (*internal.Binding, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	binding, found := s.data[bindingId]
	if !found {
		return nil, dberr.NotFound("instance with id %s not exist", bindingId)
	}
	return &binding, nil
}

func (s *Binding) Insert(binding *internal.Binding) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, found := s.data[binding.ID]; found {
		return dberr.AlreadyExists("instance with id %s already exist", binding.ID)
	}
	s.data[binding.ID] = *binding

	return nil
}

func (s *Binding) DeleteByBindingID(ID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.data, ID)
	return nil
}

func (s *Binding) ListByInstanceID(instanceID string) ([]internal.Binding, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var bindings []internal.Binding
	for _, binding := range s.data {
		if binding.InstanceID == instanceID {
			bindings = append(bindings, binding)
		}
	}

	return bindings, nil
}
