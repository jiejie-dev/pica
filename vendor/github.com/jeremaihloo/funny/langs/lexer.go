package langs

import (
	"fmt"
	"unicode/utf8"
)

// Position of one token
type Position struct {
	Line int
	Col  int
}

// String of one token
func (p *Position) String() string {
	return fmt.Sprintf("[Position] Line: %4d, Col: %4d", p.Line, p.Col)
}

// Token a part of code
type Token struct {
	Position Position
	Kind     string
	Data     string
}

func (t *Token) String() string {
	dt := t.Data
	if t.Data == "\n" {
		dt = "\\n"
	}
	return fmt.Sprintf("[Token] Kind: %6s, %6s, Data: %6s", t.Kind, t.Position.String(), "["+dt+"]")
}

// Lexer the lexer
type Lexer struct {
	Offset     int
	CurrentPos Position

	SaveOffset int
	SavePos    Position
	Data       []byte
	Elements   []Token
}

// NewLexer create a new lexer
func NewLexer(data []byte) *Lexer {
	return &Lexer{
		Data: data,
		CurrentPos: Position{
			Line: 1,
			Col:  1,
		},
	}
}

// LA next char
func (l *Lexer) LA(n int) rune {
	offset := l.Offset
	for {
		ch, size := utf8.DecodeRune(l.Data[offset:])
		if offset+size > len(l.Data) {
			return -1
		}
		chString := string(ch)
		fmt.Sprintf(chString)
		offset += size
		n--
		if n == 0 {
			return ch
		}
	}
}

// Consume next char and move position
func (l *Lexer) Consume(n int) rune {
	for {
		ch, size := utf8.DecodeRune(l.Data[l.Offset:])
		if l.Offset+size > len(l.Data) {
			return -1
		}
		chString := string(ch)
		fmt.Sprintf(chString)
		l.Offset += size
		l.CurrentPos.Col += size
		n--
		if n == 0 {
			return ch
		}
	}
}

// CreateToken create a new token and move position
func (l *Lexer) CreateToken(kind string) Token {
	st := l.Data[l.SaveOffset:l.Offset]
	token := Token{
		Kind: kind,
		Data: string(st),
		Position: Position{
			Col:  l.CurrentPos.Col - 1,
			Line: l.CurrentPos.Line,
		},
	}
	//l.CurrentPos.Col += len(token.Data)
	return token
}

// NewLine move to next line
func (l *Lexer) NewLine() Token {
	token := l.CreateToken(NEW_LINE)
	l.CurrentPos.Col = 1
	l.CurrentPos.Line = l.CurrentPos.Line + 1

	l.Reset()
	return token
}

func isNameStart(ch rune) bool {
	chString := string(ch)
	fmt.Sprintf(chString)
	return ch == '_' || (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

// ReadInt get a into from current position
func (l *Lexer) ReadInt() Token {
	for {
		ch := l.LA(1)
		switch ch {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			l.Consume(1)
		default:
			return l.CreateToken(INT)
		}
	}
}

// Reset reset the position
func (l *Lexer) Reset() {
	l.SaveOffset = l.Offset
	l.SavePos = l.CurrentPos
}

// Next get next token
func (l *Lexer) Next() Token {
	for {
		l.Reset()
		ch := l.LA(1)
		chString := string(ch)
		fmt.Sprintf(chString)
		switch ch {
		case -1:
			l.Consume(1)
			return l.CreateToken(EOF)
		case '\n':
			l.Consume(1)
			return l.NewLine()
		case ' ':
			l.Consume(1)
			break
		case '/':
			if chNext := l.LA(2); chNext == '/' {
				l.Consume(2)
				return l.ReadComments()
			}
			return l.CreateToken(DEVIDE)
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			return l.ReadInt()
		case '=':
			if l.LA(2) == '=' {
				l.Consume(2)
				return l.CreateToken(DOUBLE_EQ)
			}
			l.Consume(1)
			return l.CreateToken(EQ)
		case '+':
			l.Consume(1)
			return l.CreateToken(PLUS)
		case '-':
			l.Consume(1)
			return l.CreateToken(MINUS)
		case '*':
			l.Consume(1)
			return l.CreateToken(TIMES)
		case '(':
			l.Consume(1)
			return l.CreateToken(LParenthese)
		case ')':
			l.Consume(1)
			return l.CreateToken(RParenthese)
		case '[':
			l.Consume(1)
			return l.CreateToken(LBracket)
		case ']':
			l.Consume(1)
			return l.CreateToken(RBracket)
		case '{':
			l.Consume(1)
			return l.CreateToken(LBrace)
		case '}':
			l.Consume(1)
			return l.CreateToken(RBrace)
		case ',':
			l.Consume(1)
			return l.CreateToken(COMMA)
		case '.':
			l.Consume(1)
			return l.CreateToken(DOT)
		case '>':
			if l.LA(2) == '=' {
				l.Consume(2)
				return l.CreateToken(GTE)
			}
			l.Consume(1)
			return l.CreateToken(GT)
		case '<':
			if l.LA(2) == '=' {
				l.Consume(2)
				return l.CreateToken(LTE)
			}
			l.Consume(1)
			return l.CreateToken(LT)
		case '!':
			if l.LA(2) == '=' {
				l.Consume(2)
				return l.CreateToken(NOTEQ)
			}
		case '\'':
			if l.LA(2) == '"' {
				l.Consume(2)
				return l.CreateToken(STRING)
			}
			return l.ReadString()
		default:

			if isNameStart(ch) {
				l.Consume(1)
				for {
					chNext := l.LA(1)
					chNS := string(chNext)
					fmt.Sprintf("%s", chNS)
					if !isNameStart(chNext) {
						return l.CreateToken(NAME)
					}
					l.Consume(1)
				}
			}
			l.Consume(1)
			return l.CreateToken(EOF)
		}
	}
}

// ReadString read next string token
func (l *Lexer) ReadString() Token {
	// TODO: Fix (using state machine)
	l.Consume(1)
	l.Reset()

	for {
		ch := l.LA(1)
		switch ch {
		case '\'':
			token := l.CreateToken(STRING)
			l.Consume(1)
			return token
		default:
			l.Consume(1)
			break
		}
	}
}

// ReadComments read comments
func (l *Lexer) ReadComments() Token {
	l.Reset()
	for {
		ch := l.LA(1)
		switch ch {
		case 65533, -1:
			return l.CreateToken(EOF)
		case '\n':
			token := l.CreateToken(COMMENT)
			l.Consume(1)
			return token
		default:
			l.Consume(1)
			break
		}
	}
}
