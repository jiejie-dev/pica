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
			flag += 1
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

type Statement interface {
	Position() Position
	String() string
}

type NewLine struct {
	pos Position
}

func (n *NewLine) Position() Position {
	return n.pos
}

func (n *NewLine) String() string {
	return "\n"
}

// Variable
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

func (v *Variable) Position() Position {
	return v.pos
}

// Literal
type Literal struct {
	pos   Position
	Value interface{}
}

func (l *Literal) Position() Position {
	return l.pos
}

func (l *Literal) String() string {
	if Typing(l.Value) == "string" {
		return fmt.Sprintf("'%v'", l.Value)
	}
	return fmt.Sprintf("%v", l.Value)
}

// Expression
type Expression interface {
	Position() Position
	String() string
}

// BinaryExpression
type BinaryExpression struct {
	pos      Position
	Left     Expression
	Operator Token
	Right    Expression
}

func (b *BinaryExpression) Position() Position {
	return b.pos
}

func (b *BinaryExpression) String() string {
	return fmt.Sprintf("%s %s %s", b.Left.String(), b.Operator.Data, b.Right.String())
}

// Assign
type Assign struct {
	pos    Position
	Target Expression
	Value  Expression
}

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

// List
type List struct {
	pos    Position
	Values []Expression
}

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

type ListAccess struct {
	pos   Position
	Index int
	List  Variable
}

func (l *ListAccess) Position() Position {
	return l.pos
}

func (l *ListAccess) String() string {
	return fmt.Sprintf("%s[%d]", l.List.String(), l.Index)
}

// Block
type Block []Statement

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

// Function
type Function struct {
	pos        Position
	Name       string
	Parameters []Expression
	Body       Block
}

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

type FunctionCall struct {
	pos        Position
	Name       string
	Parameters []Expression
}

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

// Program
type Program struct {
	Statements Block
}

func (p *Program) String() string {
	return p.Statements.String()
}

// IFStatement
type IFStatement struct {
	pos       Position
	Condition Expression
	Body      Block
	Else      Block
}

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

type FORStatement struct {
	pos      Position
	Iterable IterableExpression
	Block    Block

	CurrentIndex Variable
	CurrentItem  Expression
}

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

// IterableExpression
type IterableExpression struct {
	pos   Position
	Name  Variable
	Index int
	Items []Expression
}

func (i *IterableExpression) Position() Position {
	return i.pos
}

func (i *IterableExpression) String() string {
	return fmt.Sprintf("")
}

func (i *IterableExpression) Next() (int, Expression) {
	if i.Index+1 >= len(i.Items) {
		return -1, nil
	}
	item := i.Items[i.Index]
	i.Index++
	return i.Index, item
}

type Break struct {
	pos Position
}

func (b *Break) Position() Position {
	return b.pos
}

func (b *Break) String() string {
	return "break"
}

type Continue struct {
	pos Position
}

func (b *Continue) Position() Position {
	return b.pos
}

func (b *Continue) String() string {
	return "continue"
}

type Return struct {
	pos   Position
	Value Expression
}

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

type Field struct {
	pos      Position
	Variable Variable
	Value    Expression
}

func (f *Field) Position() Position {
	return f.pos
}

func (f *Field) String() string {
	if v, ok := f.Value.(*Variable); ok && strings.Index(v.Name, "-") > -1 {
		return fmt.Sprintf("%s[%s]", f.Variable.String(), f.Value.String())
	}
	return fmt.Sprintf("%s.%s", f.Variable.String(), f.Value.String())
}

type Boolen struct {
	pos   Position
	Value bool
}

func (b *Boolen) Position() Position {
	return b.pos
}

func (b *Boolen) String() string {
	if b.Value {
		return "true"
	}
	return "false"
}

type StringExpression struct {
	pos   Position
	Value string
}

func (s *StringExpression) Position() Position {
	return s.pos
}

func (s *StringExpression) String() string {
	return s.Value
}

type Comment struct {
	pos   Position
	Value string
}

func (c *Comment) Position() Position {
	return c.pos
}

func (c *Comment) String() string {
	return fmt.Sprintf("//%s\n", c.Value)
}
