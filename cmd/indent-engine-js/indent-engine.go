package main

//go:generate gopherjs build --minify -o out/indent-engine.js

import (
	"github.com/gopherjs/gopherjs/js"
	"go.indent.com/apis/pkg/access/context"
	access "go.indent.com/apis/pkg/access/v1"
)

func main() {
	js.Module.Get("exports").Set("canAccess", canAccess)
}

type AccessRequestInput struct {
	*js.Object
	ActorID   string           `js:"actorId"`
	Actions   access.Actions   `js:"actions"`
	Resources access.Resources `js:"resources"`
}

// type AccessRequestInput struct {
// 	*js.Object
// 	Actor access.Actor `js:"actor"`
// 	Actions access.Actions `js:"actions"`
// 	Resources access.Resources `js:"resources"`
// }

func canAccess(oreq *js.Object) *js.Object {
	req := &AccessRequestInput{Object: oreq}

	lc := context.Local()
	res := lc.Process(access.Request{
		Actor:     access.Actor{ID: oreq.Get("").String()},
		Actions:   req.Actions,
		Resources: req.Resources,
	})

	o := js.Global.Get("Object").New()

	o.Set("Effect", res.Effect)
	o.Set("Errors", res.Errors)

	return o
}
