package engine

import (
	"go.indent.com/apis/pkg/access/v1"
)

func getMostLimitedEffect(effects []v1.Effect) v1.Effect {
	var mostLimitedEffect v1.Effect

	for _, effect := range effects {
		if mostLimitedEffect == "" {
			mostLimitedEffect = effect
			continue
		}

		if effect == "block" && mostLimitedEffect != "block" {
			mostLimitedEffect = effect
		}

		if effect == "escalate" && mostLimitedEffect != "escalate" && mostLimitedEffect != "block" {
			mostLimitedEffect = effect
		}
	}

	return mostLimitedEffect
}

func getOppositeEffect(eff v1.Effect) (opp v1.Effect) {
	switch eff {
	case "block":
		opp = "allow"
	case "escalate":
		opp = "allow"
	case "allow":
	default:
		opp = "block"
	}

	return
}
