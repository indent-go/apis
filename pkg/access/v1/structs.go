package v1

import (
	"github.com/indent-go/apis/pkg/access/condition"
)

type Request struct {
	Actor     `json:"actor"`
	Actions   `json:"actions"`
	Resources `json:"resources"`
	*Metadata `json:"metadata"`
}

type Labels map[string]string
type Metadata struct {
	Labels `json:"labels"`
}

type Result struct {
	Effect `json:"effect"`
	Errors `json:"errors"`
}

type Errors []Error
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Error   error
}

type Policy struct {
	Rules []Rule `json:"rules"`
}

type Rule struct {
	Effect        `json:"effect" required:"true"`
	Actions       `json:"actions" required:"true"`
	Resources     `json:"resources" required:"true"`
	ErrorTemplate `json:"error_tpl"`
	Actors        []Actor                   `json:"actor"`
	Conditions    []condition.RuleCondition `json:"conditions"`
}

type Actor struct {
	ID     string
	Groups []string
}

type Effect string
type Actions []string
type Resources []string
type ErrorTemplate string

type Config struct {
	Env `json:"env"`
}

type Env struct {
	RuntimeURN string // URN for current runtime
	LibraryURN string // URN for Indent binary or API host
}

type Claim struct{}

// type SetClaims struct {
// 	Issuer  string `json:"iss,omitempty"`
// 	ClaimID string `json:"id,omitempty"`
// 	jwt.Config
// }

// // SetClaim will set key to a value in the metadata store.
// type SetClaim struct {
// 	Key   string `json:"key,omitempty"`
// 	TTL   uint64 `json:"ttl,omitempty"`
// 	Value string `json:"value,omitempty"`
// }
