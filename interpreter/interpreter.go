package interpreter

import (
	"butaq/locales"
	"butaq/parser"
	"fmt"
	"math"
	"os"
	"strings"
	"time"
)

// Environment holds variables and functions
type Environment struct {
	vars   map[string]interface{}
	parent *Environment
}

func NewEnvironment(parent *Environment) *Environment {
	return &Environment{
		vars:   make(map[string]interface{}),
		parent: parent,
	}
}

func (e *Environment) Get(name string) (interface{}, bool) {
	val, ok := e.vars[name]
	if !ok && e.parent != nil {
		return e.parent.Get(name)
	}
	return val, ok
}

func (e *Environment) Set(name string, val interface{}) {
	if e.parent != nil {
		if _, ok := e.parent.Get(name); ok {
			e.parent.Set(name, val)
			return
		}
	}
	e.vars[name] = val
}

func (e *Environment) SetLocal(name string, val interface{}) {
	e.vars[name] = val
}

// User-defined function representation
type UserFunction struct {
	Params []string
	Body   *parser.BlockStatement
}

// Struct instance representation
type StructInstance struct {
	Name   string
	Fields map[string]interface{}
}

// Interpreter executes the program AST
type Interpreter struct {
	Output       strings.Builder
	StreamOutput bool // if true, print directly to os.Stdout instead of buffering
	env          *Environment
	funcs        map[string]*UserFunction
	structDefs   map[string]*parser.StructStatement
}

func New() *Interpreter {
	return &Interpreter{
		env:        NewEnvironment(nil),
		funcs:      make(map[string]*UserFunction),
		structDefs: make(map[string]*parser.StructStatement),
	}
}

// Return/Break/Continue signals
type ReturnValue struct{ Value interface{} }
type BreakSignal struct{}
type ContinueSignal struct{}

func (ip *Interpreter) Run(program *parser.Program) (string, error) {
	ip.Output.Reset()
	for _, stmt := range program.Statements {
		res, err := ip.evalStatement(stmt, ip.env)
		if err != nil {
			return "", err
		}
		if _, ok := res.(ReturnValue); ok {
			break
		}
	}
	return ip.Output.String(), nil
}

func (ip *Interpreter) evalStatement(stmt parser.Statement, env *Environment) (interface{}, error) {
	switch s := stmt.(type) {
	case *parser.ExpressionStatement:
		_, err := ip.evalExpression(s.Expression, env)
		return nil, err

	case *parser.VarAssignStatement:
		val, err := ip.evalExpression(s.Value, env)
		if err != nil {
			return nil, err
		}
		env.Set(s.Name.Value, val)
		return nil, nil

	case *parser.PrintStatement:
		if len(s.Values) > 0 {
			var sb strings.Builder
			for _, v := range s.Values {
				val, err := ip.evalExpression(v, env)
				if err != nil {
					return nil, err
				}
				sb.WriteString(fmt.Sprintf("%v", val))
			}
			sb.WriteString("\n")
			if ip.StreamOutput {
				fmt.Fprint(os.Stdout, sb.String())
			} else {
				ip.Output.WriteString(sb.String())
			}
			return nil, nil
		}
		val, err := ip.evalExpression(s.Value, env)
		if err != nil {
			return nil, err
		}
		if ip.StreamOutput {
			fmt.Fprintf(os.Stdout, "%v\n", val)
		} else {
			ip.Output.WriteString(fmt.Sprintf("%v\n", val))
		}
		return nil, nil

	case *parser.BlockStatement:
		blockEnv := NewEnvironment(env)
		for _, subStmt := range s.Statements {
			res, err := ip.evalStatement(subStmt, blockEnv)
			if err != nil {
				return nil, err
			}
			if res != nil {
				switch res.(type) {
				case ReturnValue, BreakSignal, ContinueSignal:
					return res, nil
				}
			}
		}
		return nil, nil

	case *parser.IfStatement:
		cond, err := ip.evalExpression(s.Condition, env)
		if err != nil {
			return nil, err
		}
		var res interface{}
		if isTruthy(cond) {
			res, err = ip.evalStatement(s.Consequence, env)
		} else if s.Alternative != nil {
			res, err = ip.evalStatement(s.Alternative, env)
		}
		if err != nil {
			return nil, err
		}
		if res != nil {
			switch res.(type) {
			case ReturnValue, BreakSignal, ContinueSignal:
				return res, nil
			}
		}
		return nil, nil

	case *parser.WhileStatement:
		for {
			cond, err := ip.evalExpression(s.Condition, env)
			if err != nil {
				return nil, err
			}
			if !isTruthy(cond) {
				break
			}
			res, err := ip.evalStatement(s.Body, env)
			if err != nil {
				return nil, err
			}
			if res != nil {
				if _, ok := res.(BreakSignal); ok {
					break
				}
				if _, ok := res.(ContinueSignal); ok {
					continue
				}
				if _, ok := res.(ReturnValue); ok {
					return res, nil
				}
			}
		}
		return nil, nil

	case *parser.BreakStatement:
		return BreakSignal{}, nil

	case *parser.ContinueStatement:
		return ContinueSignal{}, nil

	case *parser.ReturnStatement:
		val, err := ip.evalExpression(s.Value, env)
		if err != nil {
			return nil, err
		}
		return ReturnValue{Value: val}, nil

	case *parser.FunctionStatement:
		ip.funcs[s.Name] = &UserFunction{
			Params: s.Parameters,
			Body:   s.Body,
		}
		return nil, nil

	case *parser.StructStatement:
		ip.structDefs[s.Name] = s
		return nil, nil

	case *parser.StructFieldAssignStatement:
		targetObj, err := ip.evalExpression(s.Target, env)
		if err != nil {
			return nil, err
		}
		instance, ok := targetObj.(*StructInstance)
		if !ok {
			return nil, fmt.Errorf("field assignment on non-struct instance")
		}
		val, err := ip.evalExpression(s.Value, env)
		if err != nil {
			return nil, err
		}
		instance.Fields[s.Field] = val
		return nil, nil

	case *parser.IndexAssignStatement:
		arrObj, ok := env.Get(s.Array.Value)
		if !ok {
			return nil, fmt.Errorf("array not found: %s", s.Array.Value)
		}
		arr, ok := arrObj.([]interface{})
		if !ok {
			return nil, fmt.Errorf("object is not an array")
		}
		idxVal, err := ip.evalExpression(s.Index, env)
		if err != nil {
			return nil, err
		}
		idx, ok := toInt64(idxVal)
		if !ok || idx < 0 || idx >= int64(len(arr)) {
			return nil, fmt.Errorf("invalid index: %v", idxVal)
		}
		val, err := ip.evalExpression(s.Value, env)
		if err != nil {
			return nil, err
		}
		arr[idx] = val
		return nil, nil
	}

	return nil, nil
}

func (ip *Interpreter) evalExpression(expr parser.Expression, env *Environment) (interface{}, error) {
	if expr == nil {
		return nil, nil
	}

	switch e := expr.(type) {
	case *parser.IntLiteral:
		return e.Value, nil
	case *parser.NumberLiteral:
		return e.Value, nil
	case *parser.StringLiteral:
		return e.Value, nil
	case *parser.BoolLiteral:
		return e.Value, nil

	case *parser.Identifier:
		val, ok := env.Get(e.Value)
		if !ok {
			return nil, fmt.Errorf("identifier not found: %s", e.Value)
		}
		return val, nil

	case *parser.PostfixExpression:
		left, err := ip.evalExpression(e.Left, env)
		if err != nil {
			return nil, err
		}
		right, err := ip.evalExpression(e.Right, env)
		if err != nil {
			return nil, err
		}
		return evalBinaryOp(left, right, e.Operator)

	case *parser.UnaryExpression:
		right, err := ip.evalExpression(e.Right, env)
		if err != nil {
			return nil, err
		}
		if e.Operator == "емес" || resolveCanonicalOp(e.Operator) == "NOT" {
			return !isTruthy(right), nil
		}
		return nil, fmt.Errorf("unknown unary operator: %s", e.Operator)

	case *parser.ArrayLiteral:
		var elems []interface{}
		for _, el := range e.Elements {
			val, err := ip.evalExpression(el, env)
			if err != nil {
				return nil, err
			}
			elems = append(elems, val)
		}
		return elems, nil

	case *parser.IndexExpression:
		left, err := ip.evalExpression(e.Left, env)
		if err != nil {
			return nil, err
		}
		idxVal, err := ip.evalExpression(e.Index, env)
		if err != nil {
			return nil, err
		}
		if str, ok := left.(string); ok {
			idx, ok := toInt64(idxVal)
			if !ok || idx < 0 || idx >= int64(len(str)) {
				return "", nil
			}
			return string(str[idx]), nil
		}
		arr, ok := left.([]interface{})
		if !ok {
			return nil, fmt.Errorf("index access on non-array")
		}
		idx, ok := toInt64(idxVal)
		if !ok || idx < 0 || idx >= int64(len(arr)) {
			return nil, fmt.Errorf("invalid index: %v", idxVal)
		}
		return arr[idx], nil

	case *parser.LengthExpression:
		val, err := ip.evalExpression(e.Value, env)
		if err != nil {
			return nil, err
		}
		if arr, ok := val.([]interface{}); ok {
			return int64(len(arr)), nil
		}
		return nil, fmt.Errorf("length of non-array")
	
	case *parser.StrLenExpression:
		val, err := ip.evalExpression(e.Value, env)
		if err != nil {
			return nil, err
		}
		if str, ok := val.(string); ok {
			return int64(len(str)), nil
		}
		return nil, fmt.Errorf("strlen of non-string")

	case *parser.StrConcatExpression:
		left, err := ip.evalExpression(e.Left, env)
		if err != nil {
			return nil, err
		}
		right, err := ip.evalExpression(e.Right, env)
		if err != nil {
			return nil, err
		}
		return fmt.Sprintf("%v%v", left, right), nil

	case *parser.ToStrExpression:
		val, err := ip.evalExpression(e.Value, env)
		if err != nil {
			return nil, err
		}
		return fmt.Sprintf("%v", val), nil

	case *parser.StrEqExpression:
		left, err := ip.evalExpression(e.Left, env)
		if err != nil {
			return nil, err
		}
		right, err := ip.evalExpression(e.Right, env)
		if err != nil {
			return nil, err
		}
		return fmt.Sprintf("%v", left) == fmt.Sprintf("%v", right), nil

	case *parser.CharAtExpression:
		strVal, err := ip.evalExpression(e.Str, env)
		if err != nil {
			return nil, err
		}
		idxVal, err := ip.evalExpression(e.Index, env)
		if err != nil {
			return nil, err
		}
		str := fmt.Sprintf("%v", strVal)
		idx, ok := toInt64(idxVal)
		if !ok || idx < 0 || int(idx) >= len(str) {
			return "", fmt.Errorf("string index out of range")
		}
		return string(str[idx]), nil

	case *parser.CharCodeExpression:
		val, err := ip.evalExpression(e.Value, env)
		if err != nil {
			return nil, err
		}
		str := fmt.Sprintf("%v", val)
		if len(str) == 0 {
			return float64(0), nil
		}
		return float64(str[0]), nil

	case *parser.StructCreateExpression:
		def, ok := ip.structDefs[e.StructName]
		if !ok {
			return nil, fmt.Errorf("struct definition not found: %s", e.StructName)
		}
		fields := make(map[string]interface{})
		for _, fName := range def.Fields {
			fields[fName] = nil
		}
		return &StructInstance{
			Name:   e.StructName,
			Fields: fields,
		}, nil

	case *parser.StructFieldAccessExpression:
		targetObj, err := ip.evalExpression(e.Target, env)
		if err != nil {
			return nil, err
		}
		instance, ok := targetObj.(*StructInstance)
		if !ok {
			return nil, fmt.Errorf("field access on non-struct")
		}
		val, ok := instance.Fields[e.Field]
		if !ok {
			return nil, fmt.Errorf("field %s not found on struct %s", e.Field, instance.Name)
		}
		return val, nil

	case *parser.CallExpression:
		// Execute user function or built-in
		fn, ok := ip.funcs[e.Function]
		if !ok {
			// Check builtins
			return ip.callBuiltin(e.Function, e.Arguments, env)
		}

		// Evaluate args
		var args []interface{}
		for _, arg := range e.Arguments {
			val, err := ip.evalExpression(arg, env)
			if err != nil {
				return nil, err
			}
			args = append(args, val)
		}

		// Bind parameters
		fnEnv := NewEnvironment(ip.env)
		for idx, param := range fn.Params {
			if idx < len(args) {
				fnEnv.SetLocal(param, args[idx])
			} else {
				fnEnv.SetLocal(param, nil)
			}
		}

		// Execute function body
		res, err := ip.evalStatement(fn.Body, fnEnv)
		if err != nil {
			return nil, err
		}

		if ret, ok := res.(ReturnValue); ok {
			return ret.Value, nil
		}
		return nil, nil
	}

	return nil, nil
}

func (ip *Interpreter) callBuiltin(name string, args []parser.Expression, env *Environment) (interface{}, error) {
	canonical := locales.CanonicalizeBuiltin(name)
	switch canonical {
	case "синус":
		if len(args) < 1 {
			return 0.0, nil
		}
		val, err := ip.evalExpression(args[0], env)
		if err != nil {
			return nil, err
		}
		if v, ok := toFloat64(val); ok {
			return math.Sin(v), nil
		}
		return 0.0, nil

	case "косинус":
		if len(args) < 1 {
			return 0.0, nil
		}
		val, err := ip.evalExpression(args[0], env)
		if err != nil {
			return nil, err
		}
		if v, ok := toFloat64(val); ok {
			return math.Cos(v), nil
		}
		return 0.0, nil

	case "түбір":
		if len(args) < 1 {
			return 0.0, nil
		}
		val, err := ip.evalExpression(args[0], env)
		if err != nil {
			return nil, err
		}
		if v, ok := toFloat64(val); ok {
			return math.Sqrt(v), nil
		}
		return 0.0, nil

	case "дәреже":
		if len(args) < 2 {
			return 0.0, nil
		}
		baseVal, _ := ip.evalExpression(args[0], env)
		expVal, _ := ip.evalExpression(args[1], env)
		b, _ := toFloat64(baseVal)
		e, _ := toFloat64(expVal)
		return math.Pow(b, e), nil

	case "ұйықтау":
		if len(args) > 0 {
			val, _ := ip.evalExpression(args[0], env)
			if v, ok := toFloat64(val); ok {
				time.Sleep(time.Duration(v * float64(time.Second)))
			}
		}
		return 0.0, nil

	case "бүтін":
		if len(args) < 1 {
			return int64(0), nil
		}
		val, _ := ip.evalExpression(args[0], env)
		if v, ok := toInt64(val); ok {
			return v, nil
		}
		if v, ok := toFloat64(val); ok {
			return int64(v), nil
		}
		return int64(0), nil

	case "сан":
		if len(args) < 1 {
			return 0.0, nil
		}
		val, _ := ip.evalExpression(args[0], env)
		if v, ok := toFloat64(val); ok {
			return v, nil
		}
		return 0.0, nil

	case "мәтін":
		if len(args) < 1 {
			return "", nil
		}
		val, _ := ip.evalExpression(args[0], env)
		return fmt.Sprintf("%v", val), nil

	case "таңба":
		if len(args) < 1 {
			return "", nil
		}
		val, _ := ip.evalExpression(args[0], env)
		if code, ok := toInt64(val); ok {
			return string(rune(code)), nil
		}
		if code, ok := toFloat64(val); ok {
			return string(rune(int64(code))), nil
		}
		return "", nil

	case "мд5":
		val, _ := ip.evalExpression(args[0], env)
		return fmt.Sprintf("[MD5 placeholder for %v]", val), nil
	case "ша256":
		val, _ := ip.evalExpression(args[0], env)
		return fmt.Sprintf("[SHA256 placeholder for %v]", val), nil
	case "уақыт":
		return float64(time.Now().UnixNano()) / 1e9, nil

	case "уақыт_мәтіні":
		layout := "2006-01-02 15:04:05"
		if len(args) > 0 {
			if v, err := ip.evalExpression(args[0], env); err == nil {
				s := fmt.Sprintf("%v", v)
				// Basic conversion from %Y-%m-%d %H:%M:%S format
				s = strings.ReplaceAll(s, "%Y", "2006")
				s = strings.ReplaceAll(s, "%m", "01")
				s = strings.ReplaceAll(s, "%d", "02")
				s = strings.ReplaceAll(s, "%H", "15")
				s = strings.ReplaceAll(s, "%M", "04")
				s = strings.ReplaceAll(s, "%S", "05")
				layout = s
			}
		}
		return time.Now().Format(layout), nil

	case "жаңа_тізім":
		if len(args) < 2 {
			return nil, fmt.Errorf("жаңа_тізім өлшем мен мәнді талап етеді")
		}
		sizeVal, _ := ip.evalExpression(args[0], env)
		defVal, _ := ip.evalExpression(args[1], env)
		sz, _ := toInt64(sizeVal)
		res := make([]interface{}, sz)
		for i := range res {
			res[i] = defVal
		}
		return res, nil
	}
	return nil, fmt.Errorf("unknown function: %s", name)
}

func resolveCanonicalOp(op string) string {
	for _, loc := range locales.Available() {
		if tokName, ok := loc.Keywords[op]; ok {
			return tokName
		}
	}
	return op
}

func evalBinaryOp(left, right interface{}, op string) (interface{}, error) {
	canonical := resolveCanonicalOp(op)
	switch canonical {
	case "PLUS", "қосу":
		if l, ok := toFloat64(left); ok {
			if r, ok := toFloat64(right); ok {
				return l + r, nil
			}
		}
		if ls, ok := left.(string); ok {
			return fmt.Sprintf("%s%v", ls, right), nil
		}
		if rs, ok := right.(string); ok {
			return fmt.Sprintf("%v%s", left, rs), nil
		}
	case "MINUS", "алу":
		if l, ok := toFloat64(left); ok {
			if r, ok := toFloat64(right); ok {
				return l - r, nil
			}
		}
	case "MUL", "көбейту":
		if l, ok := toFloat64(left); ok {
			if r, ok := toFloat64(right); ok {
				return l * r, nil
			}
		}
	case "DIV", "бөлу":
		if l, ok := toFloat64(left); ok {
			if r, ok := toFloat64(right); ok {
				if r == 0 {
					return nil, fmt.Errorf("division by zero")
				}
				return l / r, nil
			}
		}
	case "EQ", "тең":
		if l, ok := toFloat64(left); ok {
			if r, ok := toFloat64(right); ok {
				return l == r, nil
			}
		}
		return left == right, nil
	case "NEQ", "тең_емес":
		if l, ok := toFloat64(left); ok {
			if r, ok := toFloat64(right); ok {
				return l != r, nil
			}
		}
		return left != right, nil
	case "GT", "үлкен":
		if l, ok := toFloat64(left); ok {
			if r, ok := toFloat64(right); ok {
				return l > r, nil
			}
		}
	case "LT", "кіші":
		if l, ok := toFloat64(left); ok {
			if r, ok := toFloat64(right); ok {
				return l < r, nil
			}
		}
	case "GTE", "үлкен_тең":
		if l, ok := toFloat64(left); ok {
			if r, ok := toFloat64(right); ok {
				return l >= r, nil
			}
		}
	case "LTE", "кіші_тең":
		if l, ok := toFloat64(left); ok {
			if r, ok := toFloat64(right); ok {
				return l <= r, nil
			}
		}
	case "AND", "және":
		return isTruthy(left) && isTruthy(right), nil
	case "OR", "немесе":
		return isTruthy(left) || isTruthy(right), nil
	}
	return nil, fmt.Errorf("invalid operation: %v %s %v", left, op, right)
}

func isTruthy(val interface{}) bool {
	if val == nil {
		return false
	}
	if b, ok := val.(bool); ok {
		return b
	}
	return true
}

func toFloat64(val interface{}) (float64, bool) {
	switch v := val.(type) {
	case int64:
		return float64(v), true
	case float64:
		return v, true
	}
	return 0, false
}

func toInt64(val interface{}) (int64, bool) {
	switch v := val.(type) {
	case int64:
		return v, true
	case float64:
		return int64(v), true
	}
	return 0, false
}
