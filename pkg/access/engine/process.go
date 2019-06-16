package engine

import (
	"errors"

	"go.indent.com/apis/pkg/access/condition"
	"go.indent.com/apis/pkg/access/v1"
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

			conditionsMatched := true
			conditionErrors := []v1.Error{}

			for _, cnd := range rule.Conditions {
				ok, err := EvalCondition(cnd, r)

				if !ok {
					conditionsMatched = false
				}

				if err != nil {
					conditionErrors = append(conditionErrors, v1.Error{
						Error:   err,
						Code:    "idt_rule_fail",
						Message: "indent: access(core): engine(EvalCondition): " + err.Error(),
					})
				}
			}

			eff := rule.Effect

			if conditionsMatched {
				errors = append(errors, conditionErrors...)
			} else {
				eff = getOppositeEffect(eff)
			}

			relevantEffects = append(relevantEffects, eff)
		}
	}

	res.Effect = getMostLimitedEffect(relevantEffects)
	res.Errors = errors

	if len(res.Errors) != 0 {
		res.Effect = "block"
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
						finalErr = errors.New("req.Metadata.Labels[" + label + "]: '" + metaValue + "' != '" + value + "'")
						break
					}
				} else {
					if value == "" {
						passed = true
						break
					}

					passed = false
					finalErr = errors.New("req.Metadata.Labels[" + label + "]: nil != '" + value + "'")
					break
				}
			}
		} else {
			// If req.Metadata.Labels is not found, check if conditions
			// are looking for empty labels
			for label, value := range c.Labels {
				if value == "" {
					passed = true
					finalErr = errors.New("req.Metadata.Labels: missing '" + label + "'")
				}
			}

			if !passed {
				finalErr = errors.New("req.Metadata.Labels: not found")
			}
		}
	}

	return
}
