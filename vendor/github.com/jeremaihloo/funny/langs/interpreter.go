package langs

import "fmt"

type Value interface {
}

type Scope map[string]Value

type Interpreter struct {
	Vars []Scope
	Fns  map[string]BuiltinFunction
}

func NewInterpreterWithScope(vars Scope) *Interpreter {
	i := new(Interpreter)
	i.Vars = []Scope{
		vars,
	}
	i.Fns = make(map[string]BuiltinFunction)
	i.Fns = builtinFunctions

	return i
}

func (i *Interpreter) Debug() bool {
	v := i.LookupDefault("debug", Value(false))
	if v == nil {
		return false
	}
	if v, ok := v.(bool); ok {
		return v
	}
	return false
}

func (i *Interpreter) Run(v interface{}) (Value, bool) {
	if !i.Debug() {
		// defer func() {
		// 	if err := recover(); err != nil {
		// 		fmt.Printf("\nfunny runtime error: %s\n", err)
		// 	}
		// }()
	} else {
		fmt.Sprintln("Debug Mode on.")
	}
	switch v := v.(type) {
	case Statement:
		return i.EvalStatement(v)
	case Program:
		return i.Run(&v)
	case *Program:
		return i.EvalBlock(v.Statements)
	case string:
		return i.Run([]byte(v))
	case []byte:
		parser := NewParser(v)
		program := Program{
			Statements: parser.Parse(),
		}
		return i.Run(program)
	default:
		panic(fmt.Sprintf("unknow type of running value: [%v]", v))
	}
	return Value(nil), false
}

func (i *Interpreter) EvalBlock(block Block) (Value, bool) {
	for _, item := range block {
		r, has := i.EvalStatement(item)
		if has {
			return r, has
		}
	}
	return Value(nil), false
}

func (i *Interpreter) RegisterFunction(name string, fn BuiltinFunction) error {
	if _, exists := i.Fns[name]; exists {
		return fmt.Errorf("function [%s] already exists", name)
	}
	i.Fns[name] = fn
	return nil
}

func (i *Interpreter) EvalIfStatement(item IFStatement) (Value, bool) {
	exp := i.EvalExpression(item.Condition)
	if exp, ok := exp.(bool); ok {
		if exp {
			r, has := i.EvalBlock(item.Body)
			if has {
				return r, true
			}
		} else {
			r, has := i.EvalBlock(item.Else)
			if has {
				return r, true
			}
		}
	} else {
		panic(P("if statement condition must be boolen value", item.Position()))
	}
	return Value(nil), false
}

func (i *Interpreter) EvalForStatement(item FORStatement) (Value, bool) {
	panic("NOT IMPLEMENT")
}

func (i *Interpreter) EvalStatement(item Statement) (Value, bool) {
	switch item := item.(type) {
	case *Assign:
		switch a := item.Target.(type) {
		case *Variable:
			i.Assign(a.Name, i.EvalExpression(item.Value))
			break
		case *Field:
			i.AssignField(a, i.EvalExpression(item.Value))
			break
		default:
			panic(P("invalid assignment", item.Position()))
		}
	case *IFStatement:
		val, has := i.EvalIfStatement(*item)
		if has {
			return val, true
		}
	case *FORStatement:
		val, has := i.EvalForStatement(*item)
		if has {
			return val, true
		}
	case *FunctionCall:
		i.EvalFunctionCall(item)
	case *Return:
		return i.EvalExpression(item.Value), true
	case *Function:
		i.Assign(item.Name, item)
		break
	case *Field:
		i.EvalField(item)
		break
	case *NewLine:
		break
	case *Comment:
		break
	default:
		panic(P(fmt.Sprintf("invalid statement [%s]", item.String()), item.Position()))
	}
	return Value(nil), false
}

func (i *Interpreter) EvalFunctionCall(item *FunctionCall) (Value, bool) {
	var params []Value
	for _, p := range item.Parameters {
		params = append(params, i.EvalExpression(p))
	}
	if fn, ok := i.Fns[item.Name]; ok {
		return fn(i, params), true
	}
	this := i.LookupDefault("this", nil)
	var look Value
	if this != nil {
		look = this.(map[string]Value)[item.Name]
	}
	if look == nil {
		look := i.LookupDefault(item.Name, nil)
		if look == nil {
			panic(fmt.Sprintf("function [%s] not defined", item.Name))
		}
		fun := i.Lookup(item.Name).(*Function)
		return i.EvalFunction(*fun, params)

	} else {
		fun := look.(*Function)
		return i.EvalFunction(*fun, params)
	}
}

func (i *Interpreter) EvalFunction(item Function, params []Value) (Value, bool) {
	scope := Scope{}
	i.PushScope(scope)
	for index, p := range item.Parameters {
		i.Assign(p.(*Variable).Name, params[index])
	}
	r, has := i.EvalBlock(item.Body)
	i.PopScope()
	return r, has
}

func (i *Interpreter) AssignField(field *Field, val Value) {
	scope := make(map[string]Value)

	find := i.Lookup(field.Variable.Name)
	if find != nil {
		scope = find.(map[string]Value)
	}
	scope[field.Value.(*Variable).Name] = val
	i.Assign(field.Variable.Name, Value(scope))
}

func (i *Interpreter) Assign(name string, val Value) {
	i.Vars[len(i.Vars)-1][name] = val
}

func (i *Interpreter) LookupDefault(name string, defaultVal Value) Value {
	for index := len(i.Vars) - 1; index >= 0; index-- {
		item := i.Vars[index]
		for k, v := range item {
			if k == name {
				return v
			}
		}
	}
	return defaultVal
}

func (i *Interpreter) Lookup(name string) Value {
	r := i.LookupDefault(name, Value(nil))
	if r != nil {
		return r
	}
	panic(fmt.Sprintf("variable [%s] not found", name))
}

func (i *Interpreter) PopScope() {
	i.Vars = i.Vars[:len(i.Vars)-1]
}

func (i *Interpreter) PushScope(scope Scope) {
	i.Vars = append(i.Vars, scope)
}

func (i *Interpreter) EvalExpression(expression Expression) Value {
	switch item := expression.(type) {
	case *BinaryExpression:
		// TODO: string minus
		switch item.Operator.Kind {
		case PLUS:
			return i.EvalPlus(i.EvalExpression(item.Left), i.EvalExpression(item.Right))
		case MINUS:
			return i.EvalMinus(i.EvalExpression(item.Left), i.EvalExpression(item.Right))
		case TIMES:
			return i.EvalTimes(i.EvalExpression(item.Left), i.EvalExpression(item.Right))
		case DEVIDE:
			return i.EvalDevide(i.EvalExpression(item.Left), i.EvalExpression(item.Right))
		case GT:
			return i.EvalGt(i.EvalExpression(item.Left), i.EvalExpression(item.Right))
		case GTE:
			return i.EvalGte(i.EvalExpression(item.Left), i.EvalExpression(item.Right))
		case LT:
			return i.EvalLt(i.EvalExpression(item.Left), i.EvalExpression(item.Right))
		case LTE:
			return i.EvalLte(i.EvalExpression(item.Left), i.EvalExpression(item.Right))
		case DOUBLE_EQ:
			return i.EvalDoubleEq(i.EvalExpression(item.Left), i.EvalExpression(item.Right))
		default:
			panic(P(fmt.Sprintf("only support [+] [-] [*] [/] [>] [>=] [==] [<=] [<] given [%s]", expression.(*BinaryExpression).Operator.Data), expression.Position()))
		}
	case *List:
		var ls []interface{}
		for _, item := range item.Values {
			ls = append(ls, i.EvalExpression(item))
		}
		return Value(ls)
	case *Block: // dict => map[string]Value{}
		scope := make(map[string]Value)

		for _, d := range *item {
			switch d := d.(type) {
			case *Assign:
				if t, ok := d.Target.(*Variable); ok {
					scope[t.Name] = i.EvalExpression(d.Value)
				} else {
					panic(P("block assignments must be variable", item.Position()))
				}
			case *NewLine:
				break
			case *Comment:
				break
			case *Function:
				scope[d.Name] = d
				break
			default:
				panic(P("dict struct must only contains assignment and func", item.Position()))
			}
		}
		return scope
	case *Boolen:
		return Value(item.Value)
	case *Variable:
		return i.Lookup(item.Name)
	case *Literal:
		return Value(item.Value)
	case *FunctionCall:
		r, _ := i.EvalFunctionCall(item)
		return r
	case *Field:
		return i.EvalField(item)
	case *ListAccess:
		ls := i.Lookup(item.List.Name)
		lsEntry := ls.([]interface{})
		val := lsEntry[item.Index]
		return Value(val)
	}
	panic(P(fmt.Sprintf("eval expression error: [%s]", expression.String()), expression.Position()))
}

func (i *Interpreter) EvalField(item *Field) Value {
	root := i.Lookup(item.Variable.Name)
	switch v := item.Value.(type) {
	case *FunctionCall:
		this := root.(map[string]Value)
		scope := Scope{
			"this": this,
		}
		i.PushScope(scope)
		r, _ := i.EvalFunctionCall(v)
		i.PopScope()
		return r
	case *Variable:
		iii := root.(map[string]Value)
		return Value(iii[v.Name])
	}
	return Value(nil)
}

func (i *Interpreter) EvalPlus(left, right Value) Value {
	switch left := left.(type) {
	case string:
		if right, ok := right.(string); ok {
			return Value(left + right)
		}
	case int:
		if right, ok := right.(int); ok {
			return Value(left + right)
		}
	case *[]Value:
		if right, ok := right.(*[]Value); ok {
			s := make([]Value, 0, len(*left)+len(*right))
			s = append(s, *left...)
			s = append(s, *right...)
			return Value(&s)
		}
	case *Scope:
		var s []Value
		if right, ok := right.(*Scope); ok {
			for _, l := range *left {
				flag := false

				for _, r := range s {
					if !i.EvalEqual(l, r).(bool) {
						flag = true
					} else {
						flag = false
					}

				}
				if !flag {
					s = append(s, l)
				}
			}
			for _, r := range *right {
				flag := false
				for _, c := range s {
					if !i.EvalEqual(r, c).(bool) {
						flag = true
					} else {
						flag = false
					}
				}
				if !flag {
					s = append(s, r)
				}
			}
		}
		return s
	}
	panic(fmt.Sprintf("eval plus only support types: [int, list, dict] given [%s]", Typing(left)))
}

func (i *Interpreter) EvalMinus(left, right Value) Value {
	switch left := left.(type) {
	case int:
		if right, ok := right.(int); ok {
			return Value(left - right)
		}
	case *[]Value:
		var s []Value
		if right, ok := right.(*Scope); ok {
			for _, l := range *left {
				for _, r := range *right {
					if i.EvalEqual(l, r).(bool) {
						s = append(s, l)
					}
				}
			}
		}
		return s
	case *Scope:
		var s []Value
		if right, ok := right.(*Scope); ok {
			for _, l := range *left {
				for _, r := range *right {
					if i.EvalEqual(l, r).(bool) {
						s = append(s, l)
					}
				}
			}
		}
		return s
	}
	panic("eval plus only support types: [int, list, dict]")
}

func (i *Interpreter) EvalTimes(left, right Value) Value {
	if l, ok := left.(int); ok {
		if r, o := right.(int); o {
			return Value(l * r)
		}
	}
	panic("eval plus times only support types: [int]")
}

func (i *Interpreter) EvalDevide(left, right Value) Value {
	if l, o := left.(int); o {
		if r, k := right.(int); k {
			return Value(l / r)
		}
	}
	panic("eval plus devide only support types: [int]")
}

func (i *Interpreter) EvalEqual(left, right Value) Value {
	switch l := left.(type) {
	case nil:
		return Value(right == nil)
	case int:
		if r, ok := right.(int); ok {
			return Value(l == r)
		}
	case *[]Value:
		if r, ok := right.(*[]Value); ok {
			if len(*l) != len(*r) {
				return Value(false)
			}
			for _, itemL := range *l {
				for _, itemR := range *r {
					if !i.EvalEqual(itemL, itemR).(bool) {
						return Value(false)
					}
				}
			}
			return Value(true)
		}
	case *Scope:
		if r, ok := right.(*Block); ok {
			if len(*l) != len(*r) {
				return Value(false)
			}
			for _, itemL := range *l {
				for _, itemR := range *r {
					if !i.EvalEqual(itemL, itemR).(bool) {
						return Value(false)
					}
				}
			}
			return Value(true)
		}
	}
	return Value(false)
}

func (i *Interpreter) EvalGt(left, right Value) Value {
	switch left := left.(type) {
	case int:
		if right, ok := right.(int); ok {
			return Value(left > right)
		}
	}
	panic("eval gt only support: [int]")
}

func (i *Interpreter) EvalGte(left, right Value) Value {
	switch left := left.(type) {
	case int:
		if right, ok := right.(int); ok {
			return Value(left >= right)
		}
	}
	panic("eval lte only support: [int]")
}

func (i *Interpreter) EvalLt(left, right Value) Value {
	switch left := left.(type) {
	case int:
		if right, ok := right.(int); ok {
			return Value(left < right)
		}
	}
	panic("eval lt only support: [int]")
}

func (i *Interpreter) EvalLte(left, right Value) Value {
	switch left := left.(type) {
	case int:
		if right, ok := right.(int); ok {
			return Value(left <= right)
		}
	}
	panic("eval lte only support: [int]")
}

func (i *Interpreter) EvalDoubleEq(left, right Value) Value {
	return left == right
	switch left := left.(type) {
	case int:
		if right, ok := right.(int); ok {
			return Value(left == right)
		}
	case nil:
		if left == nil && right == nil {
			return Value(true)
		}
	default:
		return Value(left == right)
	}
	panic("eval double eq only support: [int]")
}
