package engine

import (
	"github.com/indent-go/apis/pkg/access/v1"
	// "regexp"
	// "strings"
)

func matchResources(r1, r2 v1.Resources) bool {
	if r1[0] == "*" {
		return true
	}

	matched := v1.Resources{}

	for _, rr2 := range r2 {
		for _, rr1 := range r1 {
			if matchString(rr1, rr2) {
				matched = append(matched, rr2)
			}
		}
	}

	return len(matched) > 0
}

func matchActions(a1, a2 v1.Actions) bool {
	if a1[0] == "*" {
		return true
	}

	matched := v1.Actions{}

	for _, aa2 := range a2 {
		for _, aa1 := range a1 {
			if matchString(aa1, aa2) {
				matched = append(matched, aa2)
			}
		}
	}

	return len(matched) > 0
}

func matchActor(arr1 []v1.Actor, a2 v1.Actor) bool {
	// If a rule doesn't specify any actors, it applies to everyone
	if len(arr1) == 0 {
		return true
	}

	for _, a1 := range arr1 {
		if a1.ID == "" {
			return false
		}

		if a2.ID == "" {
			return false
		}

		if matchString(a1.ID, a2.ID) {
			return true
		}
		// TODO: match `(a1,a2).Groups`
	}

	return false
}

func matchString(s1, s2 string) bool {
	if s1 == "*" {
		return true
	}

	if s1 == s2 {
		return true
	}

	// r1 := strings.Replace(s1, "*", "(.*)", -1)
	// matched, _ := regexp.MatchString(r1, s2)
	// TODO: report error
	return true //matched
}
