package langs

import (
	"fmt"
	"strings"
)

func collectBlock(block Block) []string {
	flag := 0
	var s []string
	for _, item := range block {
		if item == nil {
			break
		}
		switch item.(type) {
		case *NewLine:
			flag++
			if flag < 1 {
				continue
			}
			break
		}
		flag = 0
		s = append(s, item.String())
	}
	return s
}

func intent(s string) string {
	ss := strings.Split(s, "\n")
	for index, item := range ss {
		if item == "" {
			continue
		}
		ss[index] = fmt.Sprintf("  %s", strings.TrimRight(item, " \n"))
	}
	return strings.Join(ss, "\n")
}

// Statement abstract
type Statement interface {
	Position() Position
	String() string
}

// NewLine \n
type NewLine struct {
	pos Position
}

// Position of NewLine
func (n *NewLine) Position() Position {
	return n.pos
}

func (n *NewLine) String() string {
	return "\n"
}

// Variable means var
type Variable struct {
	pos  Position
	Name string
}

func (v *Variable) String() string {
	if strings.Index(v.Name, "-") > -1 {
		return fmt.Sprintf("'%s'", v.Name)
	}
	return fmt.Sprintf("%s", v.Name)
}

// Position of Variable
func (v *Variable) Position() Position {
	return v.pos
}

// Literal like 1
type Literal struct {
	pos   Position
	Value interface{}
}

// Position of Literal
func (l *Literal) Position() Position {
	return l.pos
}

func (l *Literal) String() string {
	if Typing(l.Value) == "string" {
		return fmt.Sprintf("'%v'", l.Value)
	}
	return fmt.Sprintf("%v", l.Value)
}

// Expression abstract
type Expression interface {
	Position() Position
	String() string
}

// BinaryExpression like a > 10
type BinaryExpression struct {
	pos      Position
	Left     Expression
	Operator Token
	Right    Expression
}

// Position of BinaryExpression
func (b *BinaryExpression) Position() Position {
	return b.pos
}

func (b *BinaryExpression) String() string {
	return fmt.Sprintf("%s %s %s", b.Left.String(), b.Operator.Data, b.Right.String())
}

// Assign like a = 2
type Assign struct {
	pos    Position
	Target Expression
	Value  Expression
}

// Position of Assign
func (a *Assign) Position() Position {
	return a.pos
}

func (a *Assign) String() string {
	switch a.Value.(type) {
	case *Block:
		return fmt.Sprintf("%s = {%s}", a.Target.String(), intent(a.Value.String()))
	case *List:
		return fmt.Sprintf("%s = [%s]", a.Target.String(), intent(a.Value.String()))
	}
	return fmt.Sprintf("%s = %s", a.Target.String(), a.Value.String())
}

// List like [1, 2, 3]
type List struct {
	pos    Position
	Values []Expression
}

// Position of List
func (l *List) Position() Position {
	return l.pos
}

func (l *List) String() string {
	var s []string
	for _, item := range l.Values {
		switch item.(type) {
		case *Block:
			s = append(s, fmt.Sprintf("\n{%s}\n", intent(item.String())))
			break
		default:
			s = append(s, item.String())
		}
	}
	return fmt.Sprintf("%s", strings.Join(s, ", "))
}

// ListAccess like a[0]
type ListAccess struct {
	pos   Position
	Index int
	List  Variable
}

// Position of ListAccess
func (l *ListAccess) Position() Position {
	return l.pos
}

func (l *ListAccess) String() string {
	return fmt.Sprintf("%s[%d]", l.List.String(), l.Index)
}

// Block contains many statments
type Block []Statement

// Position of Block
func (b *Block) Position() Position {
	return Position{}
}

func (b *Block) String() string {
	var s []string
	for _, item := range *b {
		s = append(s, item.String())
	}
	return strings.Join(s, "")
}

// Function like test(a, b){}
type Function struct {
	pos        Position
	Name       string
	Parameters []Expression
	Body       Block
}

// Position of Function
func (f *Function) Position() Position {
	return f.pos
}

func (f *Function) String() string {
	var args []string
	for _, item := range f.Parameters {
		args = append(args, item.String())
	}
	s := block(f.Body)
	return fmt.Sprintf("%s(%s) {%s}", f.Name, strings.Join(args, ", "), s)
}

// FunctionCall like test(a, b)
type FunctionCall struct {
	pos        Position
	Name       string
	Parameters []Expression
}

// Position of FunctionCall
func (c *FunctionCall) Position() Position {
	return c.pos
}

func (c *FunctionCall) String() string {
	var args []string
	for _, item := range c.Parameters {
		args = append(args, item.String())
	}
	return fmt.Sprintf("%s(%s)", c.Name, strings.Join(args, ", "))
}

func block(b Block) string {
	s := collectBlock(b)
	var ss []string
	for _, item := range s {
		ss = append(ss, intent(item))
	}
	return strings.Join(ss, "")
}

// Program means the whole application
type Program struct {
	Statements Block
}

func (p *Program) String() string {
	return p.Statements.String()
}

// IFStatement like if
type IFStatement struct {
	pos       Position
	Condition Expression
	Body      Block
	Else      Block
}

// Position of IFStatement
func (i *IFStatement) Position() Position {
	return i.pos
}

func (i *IFStatement) String() string {
	if i.Else != nil && len(i.Else) != 0 {
		return fmt.Sprintf("if %s {%s} else {%s}", i.Condition.String(), block(i.Body), block(i.Else))
	} else {
		return fmt.Sprintf("if %s {%s}", i.Condition.String(), block(i.Body))
	}
}

// FORStatement like for
type FORStatement struct {
	pos      Position
	Iterable IterableExpression
	Block    Block

	CurrentIndex Variable
	CurrentItem  Expression
}

// Position of FORStatement
func (f *FORStatement) Position() Position {
	return f.pos
}

func (f *FORStatement) String() string {
	return fmt.Sprintf("for %s, %s in %s {\n%s\n}",
		f.CurrentIndex.String(),
		f.CurrentItem.String(),
		f.Iterable.Name.String(),
		intent(f.Block.String()))
}

// IterableExpression like for in
type IterableExpression struct {
	pos   Position
	Name  Variable
	Index int
	Items []Expression
}

// Position of IterableExpression
func (i *IterableExpression) Position() Position {
	return i.pos
}

func (i *IterableExpression) String() string {
	return fmt.Sprintf("")
}

// Next part of IterableExpression
func (i *IterableExpression) Next() (int, Expression) {
	if i.Index+1 >= len(i.Items) {
		return -1, nil
	}
	item := i.Items[i.Index]
	i.Index++
	return i.Index, item
}

// Break like break in for
type Break struct {
	pos Position
}

// Position of Break
func (b *Break) Position() Position {
	return b.pos
}

func (b *Break) String() string {
	return "break"
}

// Continue like continue in for
type Continue struct {
	pos Position
}

// Position of Continue
func (b *Continue) Position() Position {
	return b.pos
}

func (b *Continue) String() string {
	return "continue"
}

// Return like return varA
type Return struct {
	pos   Position
	Value Expression
}

// Position of Return
func (r *Return) Position() Position {
	return r.pos
}

func (r *Return) String() string {
	switch r.Value.(type) {
	case *Block:
		return fmt.Sprintf("return {%s}", intent(r.Value.String()))
	}
	return fmt.Sprintf("return %s", r.Value.String())
}

// Field like obj.age
type Field struct {
	pos      Position
	Variable Variable
	Value    Expression
}

// Position of Field
func (f *Field) Position() Position {
	return f.pos
}

func (f *Field) String() string {
	if v, ok := f.Value.(*Variable); ok && strings.Index(v.Name, "-") > -1 {
		return fmt.Sprintf("%s[%s]", f.Variable.String(), f.Value.String())
	}
	return fmt.Sprintf("%s.%s", f.Variable.String(), f.Value.String())
}

// Boolen like true, false
type Boolen struct {
	pos   Position
	Value bool
}

// Position of Boolen
func (b *Boolen) Position() Position {
	return b.pos
}

func (b *Boolen) String() string {
	if b.Value {
		return "true"
	}
	return "false"
}

// StringExpression like 'hello world !'
type StringExpression struct {
	pos   Position
	Value string
}

// Position of StringExpression
func (s *StringExpression) Position() Position {
	return s.pos
}

func (s *StringExpression) String() string {
	return s.Value
}

// Comment line for sth
type Comment struct {
	pos   Position
	Value string
}

// Position of comment
func (c *Comment) Position() Position {
	return c.pos
}

func (c *Comment) String() string {
	return fmt.Sprintf("//%s\n", c.Value)
}
