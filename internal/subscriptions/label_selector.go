package subscriptions

import (
	"fmt"
	"github.com/kyma-project/kyma-environment-broker/common/gardener"
	"github.com/kyma-project/kyma-environment-broker/internal/process/provisioning"
	"strings"
)

// SecretBinding selector requirements
const (
	hyperscalerTypeReqFmt = gardener.HyperscalerTypeLabelKey + "=%s"
	tenantNameReqFmt      = gardener.TenantNameLabelKey + "=%s"

	dirtyReq    = gardener.DirtyLabelKey + "=true"
	internalReq = gardener.InternalLabelKey + "=true"
	sharedReq   = gardener.SharedLabelKey + "=true"
	euAccessReq = gardener.EUAccessLabelKey + "=true"

	notSharedReq = gardener.SharedLabelKey + `!=true`

	notDirtyReq       = `!` + gardener.DirtyLabelKey
	notInternalReq    = `!` + gardener.InternalLabelKey
	notEUAccessReq    = `!` + gardener.EUAccessLabelKey
	notTenantNamedReq = `!` + gardener.TenantNameLabelKey
)

type LabelSelectorBuilder struct {
	strings.Builder
	shared bool
}

func NewLabelSelectorFromRuleset(rule provisioning.ParsedRule) *LabelSelectorBuilder {
	selector := &LabelSelectorBuilder{}
	selector.shared = rule.IsShared()
	selector.WriteString(fmt.Sprintf(hyperscalerTypeReqFmt, rule.Hyperscaler()))
	if rule.IsEUAccess() {
		selector.With(euAccessReq)
	} else {
		selector.With(notEUAccessReq)
	}
	if rule.IsShared() {
		selector.With(sharedReq)
		return selector
	}
	selector.With(notDirtyReq)

	return selector
}

func (l *LabelSelectorBuilder) With(s string) {
	if l.Len() == 0 {
		l.WriteString(s)
		return
	}
	l.WriteString("," + s)
}

func (l *LabelSelectorBuilder) BuildForTenantMatching(tenant string) string {
	if l.shared {
		return l.String()
	}
	base := l.String()
	tenantSelector := fmt.Sprintf(tenantNameReqFmt, tenant)
	return base + "," + tenantSelector
}

func (l *LabelSelectorBuilder) BuildAnySubscription() string {
	return l.String()
}

func (l *LabelSelectorBuilder) BuildForSecretBindingClaim() string {
	base := l.String()
	return fmt.Sprintf("%s,%s,%s", base, notTenantNamedReq, notSharedReq)
}
