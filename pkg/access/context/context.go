package context

import (
	"go.indent.com/apis/pkg/access/engine"
	"go.indent.com/apis/pkg/access/v1"
)

type Context interface {
	Process(v1.Request) v1.Result
}

type LocalContext struct {
	Config v1.Config `json:"config"`
}

// Local returns a pointer to a LocalContext
// for evaluating access requests
func Local() *LocalContext {
	lc := LocalContext{
		Config: v1.Config{
			Env: v1.Env{
				RuntimeURN: "",
				LibraryURN: "",
			},
		},
	}

	return &lc
}

func (c *LocalContext) Process(r v1.Request) v1.Result {
	ie := engine.New(engine.PolicyEngineInput{
		Policies: []v1.Policy{
			v1.Policy{
				Rules: []v1.Rule{
					v1.Rule{
						Actors:    []v1.Actor{{ID: "*"}},
						Actions:   v1.Actions{"*"},
						Resources: v1.Resources{"*"},
						Effect:    "block",
					},
				},
			},
		},
	})

	return ie.Process(r)
}
