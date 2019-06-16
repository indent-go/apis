package pgproto

// MsgType is the type of a Postgresql wire message.
type MsgType byte

// String returns the message type as defined in the Postgres manual, frontend should be true if a frontend message.
func (t MsgType) String(frontend bool) string {
	if str, ok := bothMsgs[t]; ok {
		return str
	} else if frontend {
		if str, ok = frontendMsgs[t]; ok {
			return str
		}
	} else if str, ok = backendMsgs[t]; ok {
		return str
	}
	return "Unknown"
}

// Postgres protocol message types
//
// Defined in https://www.postgresql.org/docs/11/protocol-message-formats.html

// Frontend message types
const (
	// DescribeMsg identifies the message as a Describe command.
	DescribeMsg MsgType = 'D'
	// PasswordMessageMsg identifies the message as a password response. Note that this is also used for
	// GSSAPI, SSPI and SASL response messages. The exact message type can be deduced from the context.
	PasswordMessageMsg MsgType = 'p'
	// QueryMsg identifies the message as a simple query.
	QueryMsg MsgType = 'Q'
	// TerminateMsg identifies the message as a termination.
	TerminateMsg MsgType = 'X'
)

var frontendMsgs = map[MsgType]string{
	DescribeMsg:        "Describe",
	PasswordMessageMsg: "PasswordMessage",
	QueryMsg:           "Query",
	TerminateMsg:       "Terminate",
}

// Backend message types
const (
	// AuthenticationMsg identifies the message as an authentication request.
	AuthenticationMsg MsgType = 'R'
	// CommandCompleteMsg identifies the message as a command-completed response.
	CommandCompleteMsg MsgType = 'C'
	// CopyInResponseMsg identifies the message as a Start Copy In response. The frontend must now send copy-in data
	// (if not prepared to do so, send a CopyFail message).
	CopyInResponseMsg MsgType = 'G'
	// DataRowMsg identifies the message as a data row.
	DataRowMsg MsgType = 'D'
	// EmptyQueryResponseMsg identifies the message as a response to an empty query string.
	// (This substitutes for CommandComplete.)
	EmptyQueryResponseMsg MsgType = 'I'
	// ErrorResponseMsg identifies the message as an error.
	ErrorResponseMsg MsgType = 'E'
	// NoticeResponseMsg identifies the message as a notice.
	NoticeResponseMsg MsgType = 'N'
	// ReadyForQueryMsg is sent whenever the backend is ready for a new query cycle.
	ReadyForQueryMsg MsgType = 'Z'
	// RowDescriptionMsg identifies the message as a row description.
	RowDescriptionMsg MsgType = 'T'
)

var backendMsgs = map[MsgType]string{
	AuthenticationMsg:     "Authentication",
	CommandCompleteMsg:    "CommandComplete",
	CopyInResponseMsg:     "CopyInResponse",
	DataRowMsg:            "DataRow",
	EmptyQueryResponseMsg: "EmptyQueryResponse",
	ErrorResponseMsg:      "ErrorResponse",
	NoticeResponseMsg:     "NoticeResponse",
	ReadyForQueryMsg:      "ReadyForQuery",
	RowDescriptionMsg:     "RowDescription",
}

// Backend / Frontend message types
const (
	// CopyDataMsg identifies the message as COPY data.
	CopyDataMsg MsgType = 'd'
	// CopyDoneMsg identifies the message as a COPY-complete indicator.
	CopyDoneMsg MsgType = 'c'
	// CopyFailMsg identifies the message as a COPY-failure indicator.
	CopyFailMsg MsgType = 'f'
)

var bothMsgs = map[MsgType]string{
	CopyDataMsg: "CopyData",
	CopyDoneMsg: "CopyDone",
	CopyFailMsg: "CopyFail",
}
