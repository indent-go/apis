package apis_test

import (
	"fmt"

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
	// [AccessRequest REDACTED]
}
