package condition

type Opcode string
type Operand string
type EvalCondition struct {
	Opcode
	Operands []Operand
}

type RuleCondition struct {
	Labels       map[string]string `json:"labels"`
	DetectRules  `json:"detect"`
	ExternalEval `json:"ext_eval"`
}

type Reason string
type AppDomain string
type ExternalEval struct {
	URL string `json:"url"`
}

type DetectRules struct {
	TimeRangeCounter
}

type TimeRangeCounter struct {
	EvalCondition
	RelativeWindowMs int
}
