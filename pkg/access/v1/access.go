package v1

import (
	"github.com/dgrijalva/jwt-go"
	"go.indent.com/apis/pkg/access/condition"
)

type Request struct {
	Actor     `json:"actor"`
	Actions   `json:"actions"`
	Resources `json:"resources"`
	*Metadata `json:"metadata"`
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
	ID      string        `json:"id,omitempty"`
	Groups  []string      `json:"groups,omitempty"`
	Context *ActorContext `json:"context,omitempty"`
}

type ActorContext struct {
	IPAddr string `json:"ip,omitempty"`
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

type ClaimTokens []string
type Labels map[string]string
type Metadata struct {
	Labels      `json:"labels"`
	ClaimTokens `json:"claim_tokens"`
}

type AccessClaims struct {
	jwt.StandardClaims
}

// SetClaims will set key to a value in the metadata store.
type SetClaims struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
	jwt.StandardClaims
}
