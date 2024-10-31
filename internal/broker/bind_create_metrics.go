package broker

import (
	"context"
	"github.com/kyma-project/kyma-environment-broker/internal/event"
	"github.com/pivotal-cf/brokerapi/v8/domain"
	"time"
)

type Binder interface {
	Bind(ctx context.Context, instanceID, bindingID string, details domain.BindDetails, asyncAllowed bool) (domain.Binding, error)
}

// todo: find a name for this struct
type BindMetrics struct {
	target    Binder
	publisher event.Publisher
}

type BindRequestProcessed struct {
	ProcessingDuration time.Duration
	Error              error
}

type UnbindRequestProcessed struct {
	ProcessingDuration time.Duration
	Error              error
}

func NewBindMetrics(target Binder, publisher event.Publisher) *BindMetrics {
	return &BindMetrics{
		target:    target,
		publisher: publisher,
	}
}

func (b *BindMetrics) Bind(ctx context.Context, instanceID, bindingID string, details domain.BindDetails, asyncAllowed bool) (domain.Binding, error) {
	start := time.Now()
	response, err := b.target.Bind(ctx, instanceID, bindingID, details, asyncAllowed)
	processingDuration := time.Since(start)

	b.publisher.Publish(ctx, BindRequestProcessed{ProcessingDuration: processingDuration, Error: err})

	return response, err
}
