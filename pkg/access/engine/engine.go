package engine

import (
	"go.indent.com/apis/pkg/access/v1"
)

// PolicyEngine is an instance of the
// Indent Access Policy Engine with all
// the policies loaded.
//
// Any data fetching happens during startup
// or for dynamic policies, out-of-band for
// performance.
type PolicyEngine struct {
	Claims   []v1.Claim
	Policies []v1.Policy
}

type PolicyEngineInput struct {
	Claims   []v1.Claim
	Policies []v1.Policy
}

// New returns a PolicyEngine instance.
func New(input PolicyEngineInput) PolicyEngine {
	return PolicyEngine{
		Policies: input.Policies,
	}
}
