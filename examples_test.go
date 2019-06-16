package apis_test

import (
	"encoding/json"
	"fmt"

	"go.indent.com/apis/pkg/access/condition"
	"go.indent.com/apis/pkg/access/engine"
	access "go.indent.com/apis/pkg/access/v1"
)

func Example_accessRequest() {
	req := access.Request{
		Actor:     access.Actor{ID: "vogre923f2.1033994"},
		Actions:   access.Actions{"{provider}:actions::{type}:{action}"},
		Resources: access.Resources{"{provider}:resources::{type}:{resource}"},
	}

	fmt.Printf("%v", req)
	// Output:
	// [AccessResult.v1 REDACTED]
}

func Example_blockMissingReason() {
	ae := engine.New(engine.PolicyEngineInput{
		Policies: []access.Policy{
			access.Policy{
				Rules: []access.Rule{
					access.Rule{
						Effect:    "block",
						Actions:   access.Actions{"*"},
						Resources: access.Resources{"*"},
						Conditions: []condition.RuleCondition{
							condition.RuleCondition{
								Labels: access.Labels{"reason": ""},
							},
						},
					},
				},
			},
		},
	})

	res := ae.Process(access.Request{
		Actor:     access.Actor{ID: "4D4BC592A2F37"},
		Actions:   access.Actions{"{provider}:actions::{type}:{action}"},
		Resources: access.Resources{"{provider}:resources::{type}:{resource}"},
	})
	j, _ := json.Marshal(res)

	fmt.Printf("%s\n", string(j))

	res = ae.Process(access.Request{
		Actor:     access.Actor{ID: "4D4BC592A2F37"},
		Actions:   access.Actions{"{provider}:actions::{type}:{action}"},
		Resources: access.Resources{"{provider}:resources::{type}:{resource}"},
		Metadata: &access.Metadata{
			Labels: access.Labels{
				"reason": "support: billing plan change",
			},
		},
	})
	j, _ = json.Marshal(res)

	fmt.Printf("%s", string(j))
	// Output:
	// {"effect":"block","errors":[{"code":"idt_rule_fail","message":"indent: access(core): engine(EvalCondition): req.Metadata.Labels: missing 'reason'"}]}
	// {"effect":"allow","errors":[]}
}
