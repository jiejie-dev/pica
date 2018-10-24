package langs

const (
	LBrace      = "{"
	RBrace      = "}"
	LBracket    = "["
	RBracket    = "]"
	LParenthese = "("
	RParenthese = ")"
	EQ          = "="
	DOUBLE_EQ   = "=="
	PLUS        = "+"
	MINUS       = "-"
	TIMES       = "*"
	DEVIDE      = "/"
	Quote       = "\""
	GT          = ">"
	LT          = "<"
	GTE         = ">="
	LTE         = "<="
	NOTEQ       = "!="
	COMMA       = ","
	DOT         = "."
	EOF         = "EOF"
	INT         = "INT"
	NAME        = "NAME"
	STRING      = "STRING"

	IF       = "if"
	ELSE     = "else"
	TRUE     = "true"
	FALSE    = "false"
	FOR      = "for"
	AND      = "and"
	IN       = "in"
	NIL      = "nil"
	NOT      = "not"
	OR       = "or"
	RETURN   = "return"
	BREAK    = "break"
	CONTINUE = "continue"

	NEW_LINE = "\\n"
	COMMENT  = "comment"
)

var Keywords = map[string]string{
	"and":      "and",
	"else":     "else",
	"false":    "false",
	"for":      "for",
	"if":       "if",
	"in":       "in",
	"nil":      "nil",
	"not":      "not",
	"or":       "or",
	"return":   "return",
	"true":     "true",
	"break":    "break",
	"continue": "continue",
}
