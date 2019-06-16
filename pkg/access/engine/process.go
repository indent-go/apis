package engine

import (
	"errors"
	"github.com/indent-go/apis/pkg/access/condition"
	"github.com/indent-go/apis/pkg/access/v1"
)

// Process will evaluate the access request against the policy engine
func (e PolicyEngine) Process(r v1.Request) (res v1.Result) {
	errors := []v1.Error{}
	relevantEffects := []v1.Effect{}

	for _, policy := range e.Policies {
		for _, rule := range policy.Rules {
			if !matchResources(rule.Resources, r.Resources) {
				continue
			}

			if !matchActions(rule.Actions, r.Actions) {
				continue
			}

			if !matchActor(rule.Actors, r.Actor) {
				continue
			}

			for _, cnd := range rule.Conditions {
				ok, err := EvalCondition(cnd, r)

				if err != nil {
					errors = append(errors, v1.Error{
						Error:   err,
						Code:    "idt_rule_fail",
						Message: "indent: access(core): engine(EvalCondition): " + err.Error(),
					})
					continue
				}

				var effect v1.Effect

				if ok {
					effect = rule.Effect
				} else {
					effect = getOppositeEffect(rule.Effect)
				}

				relevantEffects = append(relevantEffects, effect)
			}
		}
	}

	res.Effect = getMostLimitedEffect(relevantEffects)
	res.Errors = errors

	if res.Effect == "" {
		if len(res.Errors) == 0 {
			res.Effect = "allow"
		} else {
			res.Effect = "block"
		}
	}

	return res
}

// EvalCondition ...
func EvalCondition(c condition.RuleCondition, r v1.Request) (passed bool, finalErr error) {
	if len(c.Labels) != 0 {
		if r.Metadata != nil {
			for label, value := range c.Labels {
				if metaValue, ok := r.Metadata.Labels[label]; ok {
					if metaValue == value {
						passed = true
					} else {
						passed = false
						finalErr = errors.New(".Metadata.Labels[" + label + "]: " + r.Metadata.Labels[label] + " != " + metaValue)
						break
					}
				} else {
					passed = false
					finalErr = errors.New(".Metadata.Labels[" + label + "]: " + r.Metadata.Labels[label] + " != ''")
					break
				}
			}
		} else {
			passed = false
			finalErr = errors.New(".Metadata.Labels: not found")
		}
	}

	return
}
