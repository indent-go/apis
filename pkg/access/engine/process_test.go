package engine_test

import (
	"fmt"
	"github.com/indent-go/apis/pkg/access"
	"github.com/indent-go/apis/pkg/access/condition"
	"github.com/indent-go/apis/pkg/access/engine"
	"testing"
)

func newEngine() engine.PolicyEngine {
	return engine.New(engine.PolicyEngineInput{
		Policies: []access.Policy{
			access.Policy{
				Rules: []access.Rule{
					access.Rule{
						Effect:    "",
						Actions:   access.Actions{"*"},
						Resources: access.Resources{"*"},
						Conditions: []condition.RuleCondition{
							condition.RuleCondition{
								Labels: access.Labels{
									"reason": "",
								},
							},
						},
					},
				},
			},
		},
	})
}

func newEngineWithClaims() engine.PolicyEngine {
	return engine.New(engine.PolicyEngineInput{
		Policies: []access.Policy{
			access.Policy{
				Rules: []access.Rule{
					access.Rule{
						Effect:    "",
						Actions:   access.Actions{"*"},
						Resources: access.Resources{"*"},
						Conditions: []condition.RuleCondition{
							condition.RuleCondition{
								Labels: access.Labels{
									"reason": "",
								},
							},
						},
					},
				},
			},
		},
	})
}

func TestProcessBasic(t *testing.T) {
	ie := newEngine()
	res := ie.Process(access.Request{
		Actor:     access.Actor{ID: "github:actors::user:human"},
		Actions:   access.Actions{"github:actions::repo.clone"},
		Resources: access.Resources{"github:resources::repo:123"},
	})

	if res.Effect != "block" {
		fmt.Println("TestProcessBasic: res.Effect != block")
		t.Fail()
	}

	res = ie.Process(access.Request{
		Actor:     access.Actor{ID: "github:actors::user:human"},
		Actions:   access.Actions{"github:actions::repo.clone"},
		Resources: access.Resources{"github:resources::repo:123"},
		Metadata: &access.Metadata{
			Labels: access.Labels{
				"reason": "",
			},
		},
	})

	if res.Effect != "allow" {
		fmt.Println("TestProcessBasic: res.Effect != allow")
		t.Fail()
	}
}

func TestProcessClaim(t *testing.T) {
	ie := newEngineWithClaims()
	res := ie.Process(access.Request{
		Actor: access.Actor{ID: "lattice:resources::user:123"},
		Actions: access.Actions{
			"lattice:actions::goal:view",
			"http:actions::post:/v1/api/graphql",
		},
		Resources: access.Resources{"github:resources::repo:123"},
	})

	if res.Effect != "block" {
		fmt.Println("TestProcessBasic: res.Effect != block")
		t.Fail()
	}

	res = ie.Process(access.Request{
		Actor:     access.Actor{ID: "github:actors::user:human"},
		Actions:   access.Actions{"github:actions::repo.clone"},
		Resources: access.Resources{"github:resources::repo:123"},
		Metadata: &access.Metadata{
			Labels: access.Labels{
				"reason": "",
			},
		},
	})

	if res.Effect != "allow" {
		fmt.Println("TestProcessBasic: res.Effect != allow")
		t.Fail()
	}
}

func _testProcessCounter(t *testing.T) {
	ie := engine.New(engine.PolicyEngineInput{
		Policies: []access.Policy{
			{
				Rules: []access.Rule{
					{
						Effect:    "block",
						Actions:   access.Actions{"*"},
						Resources: access.Resources{"*"},
						Conditions: []condition.RuleCondition{
							{
								Labels: map[string]string{
									"reason": "",
								},
							},
						},
					},
				},
			},
			{
				Rules: []access.Rule{
					access.Rule{
						Effect:    "alert",
						Actions:   access.Actions{"github:actions:repo.*"},
						Resources: access.Resources{"github:*"},
						Conditions: []condition.RuleCondition{
							condition.RuleCondition{
								DetectRules: condition.DetectRules{
									TimeRangeCounter: condition.TimeRangeCounter{
										EvalCondition: condition.EvalCondition{
											Opcode: ">",
											Operands: []condition.Operand{
												"counter:count(@action)",
												"counter:median_count(@action)",
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	})
	res := ie.Process(access.Request{
		Actor:     access.Actor{ID: "github:actors::user:human"},
		Actions:   access.Actions{"github:actions::repo.clone"},
		Resources: access.Resources{"github:resources::repo:123"},
	})

	if res.Effect != "allow" {
		t.Fail()
	}
}
