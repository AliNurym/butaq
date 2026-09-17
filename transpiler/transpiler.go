package transpiler

import (
	"butaq/locales"
	"butaq/parser"
	"fmt"
	"strings"
)

// Transpiler translates an AST program into Butaq source code in a target locale.
type Transpiler struct {
	targetLocale *locales.Locale
	indentStr    string
}

// New creates a new Transpiler targeting the specified locale.
func New(targetLocale *locales.Locale) *Transpiler {
	if targetLocale == nil {
		targetLocale = locales.Default()
	}
	return &Transpiler{
		targetLocale: targetLocale,
		indentStr:    "    ",
	}
}

// Transpile converts a parsed Program AST into source code of the target human language.
func (t *Transpiler) Transpile(program *parser.Program) string {
	if program == nil {
		return ""
	}

	var sb strings.Builder
	for i, stmt := range program.Statements {
		if i > 0 {
			sb.WriteString("\n")
		}
		sb.WriteString(t.transpileStatement(stmt, 0))
	}
	sb.WriteString("\n")
	return sb.String()
}

func (t *Transpiler) transpileStatement(stmt parser.Statement, depth int) string {
	indent := strings.Repeat(t.indentStr, depth)

	switch s := stmt.(type) {
	case *parser.VarAssignStatement:
		valStr := t.transpileExpression(s.Value)
		if s.IsDeclaration {
			varKw := t.targetLocale.KeywordForToken("VAR")
			return fmt.Sprintf("%s%s %s = %s", indent, varKw, s.Name.Value, valStr)
		}
		return fmt.Sprintf("%s%s = %s", indent, s.Name.Value, valStr)

	case *parser.StructFieldAssignStatement:
		valStr := t.transpileExpression(s.Value)
		targetStr := ""
		if s.Target != nil {
			targetStr = t.transpileExpression(s.Target)
		} else {
			targetStr = s.StructName
		}
		return fmt.Sprintf("%s%s.%s = %s", indent, targetStr, s.Field, valStr)

	case *parser.IndexAssignStatement:
		idxStr := t.transpileExpression(s.Index)
		valStr := t.transpileExpression(s.Value)
		return fmt.Sprintf("%s%s[%s] = %s", indent, s.Array.Value, idxStr, valStr)

	case *parser.PrintStatement:
		printKw := t.targetLocale.KeywordForToken("PRINT")
		var args []string
		if len(s.Values) > 0 {
			for _, v := range s.Values {
				args = append(args, t.transpileExpression(v))
			}
		} else if s.Value != nil {
			args = append(args, t.transpileExpression(s.Value))
		}
		return fmt.Sprintf("%s%s(%s)", indent, printKw, strings.Join(args, ", "))

	case *parser.ReturnStatement:
		valStr := t.transpileExpression(s.Value)
		retKw := t.targetLocale.KeywordForToken("RETURN")
		if valStr == "" {
			return fmt.Sprintf("%s%s", indent, retKw)
		}
		return fmt.Sprintf("%s%s %s", indent, retKw, valStr)

	case *parser.CallStatement:
		callStr := t.transpileExpression(s.Call)
		return fmt.Sprintf("%s%s", indent, callStr)

	case *parser.FreeStatement:
		valStr := t.transpileExpression(s.Value)
		freeKw := t.targetLocale.KeywordForToken("FREE")
		return fmt.Sprintf("%s%s(%s)", indent, freeKw, valStr)

	case *parser.BreakStatement:
		return fmt.Sprintf("%s%s", indent, t.targetLocale.KeywordForToken("BREAK"))

	case *parser.ContinueStatement:
		return fmt.Sprintf("%s%s", indent, t.targetLocale.KeywordForToken("CONTINUE"))

	case *parser.ImportStatement:
		importKw := t.targetLocale.KeywordForToken("IMPORT")
		return fmt.Sprintf("%s%s %q", indent, importKw, s.Path)

	case *parser.FileWriteStatement:
		pathStr := t.transpileExpression(s.Path)
		contentStr := t.transpileExpression(s.Content)
		fwKw := t.targetLocale.KeywordForToken("FILE_WRITE")
		return fmt.Sprintf("%s%s(%s, %s)", indent, fwKw, pathStr, contentStr)

	case *parser.IfStatement:
		condStr := t.transpileExpression(s.Condition)
		ifKw := t.targetLocale.KeywordForToken("IF")
		if s.Alternative == nil && s.Consequence != nil && len(s.Consequence.Statements) == 1 {
			thenKw := t.targetLocale.KeywordForToken("THEN")
			thenStmt := strings.TrimSpace(t.transpileStatement(s.Consequence.Statements[0], 0))
			return fmt.Sprintf("%s%s %s %s %s", indent, ifKw, condStr, thenKw, thenStmt)
		}

		conseqStr := t.transpileBlock(s.Consequence, depth)
		out := fmt.Sprintf("%s%s %s\n%s", indent, ifKw, condStr, conseqStr)
		if s.Alternative != nil {
			elseKw := t.targetLocale.KeywordForToken("ELSE")
			altStr := t.transpileBlock(s.Alternative, depth)
			out += fmt.Sprintf("%s%s\n%s", indent, elseKw, altStr)
		}
		return strings.TrimRight(out, "\n")

	case *parser.WhileStatement:
		condStr := t.transpileExpression(s.Condition)
		whileKw := t.targetLocale.KeywordForToken("WHILE")
		bodyStr := t.transpileBlock(s.Body, depth)
		return strings.TrimRight(fmt.Sprintf("%s%s %s\n%s", indent, whileKw, condStr, bodyStr), "\n")

	case *parser.FunctionStatement:
		funcKw := t.targetLocale.KeywordForToken("FUNC")
		bodyStr := t.transpileBlock(s.Body, depth)
		params := strings.Join(s.Parameters, ", ")
		return strings.TrimRight(fmt.Sprintf("%s%s %s(%s)\n%s", indent, funcKw, s.Name, params, bodyStr), "\n")

	case *parser.StructStatement:
		structKw := t.targetLocale.KeywordForToken("STRUCT")
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%s%s %s {\n", indent, structKw, s.Name))
		fieldIndent := strings.Repeat(t.indentStr, depth+1)
		for i, field := range s.Fields {
			fieldType := "САН"
			if i < len(s.Types) {
				fieldType = s.Types[i]
			}
			isWeak := false
			if i < len(s.IsWeak) {
				isWeak = s.IsWeak[i]
			}

			translatedType := t.transpileType(fieldType, isWeak)
			sb.WriteString(fmt.Sprintf("%s%s %s\n", fieldIndent, field, translatedType))
		}
		sb.WriteString(fmt.Sprintf("%s}", indent))
		return sb.String()

	case *parser.InterfaceStatement:
		interfaceKw := t.targetLocale.KeywordForToken("INTERFACE")
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%s%s %s {\n", indent, interfaceKw, s.Name))
		methodIndent := strings.Repeat(t.indentStr, depth+1)
		for _, m := range s.Methods {
			paramTypes := strings.Join(m.Parameters, ", ")
			retType := ""
			if m.ReturnType != "" {
				retType = " " + t.transpileType(m.ReturnType, false)
			}
			sb.WriteString(fmt.Sprintf("%s%s(%s)%s\n", methodIndent, m.Name, paramTypes, retType))
		}
		sb.WriteString(fmt.Sprintf("%s}", indent))
		return sb.String()

	case *parser.ThreadStatement:
		threadKw := t.targetLocale.KeywordForToken("THREAD")
		bodyStr := t.transpileStatement(s.Body, 0)
		return fmt.Sprintf("%s%s %s", indent, threadKw, strings.TrimSpace(bodyStr))

	case *parser.ExpressionStatement:
		return indent + t.transpileExpression(s.Expression)

	default:
		if stmt != nil {
			return indent + stmt.String()
		}
		return ""
	}
}

func (t *Transpiler) transpileBlock(block *parser.BlockStatement, depth int) string {
	if block == nil {
		return ""
	}
	var sb strings.Builder
	for _, stmt := range block.Statements {
		sb.WriteString(t.transpileStatement(stmt, depth+1))
		sb.WriteString("\n")
	}
	return sb.String()
}

func (t *Transpiler) transpileExpression(expr parser.Expression) string {
	if expr == nil {
		return ""
	}

	switch e := expr.(type) {
	case *parser.NumberLiteral:
		return fmt.Sprintf("%v", e.Value)

	case *parser.IntLiteral:
		return fmt.Sprintf("%d", e.Value)

	case *parser.StringLiteral:
		return fmt.Sprintf("%q", e.Value)

	case *parser.BoolLiteral:
		if e.Value {
			return t.targetLocale.KeywordForToken("TRUE")
		}
		return t.targetLocale.KeywordForToken("FALSE")

	case *parser.Identifier:
		// Check if identifier corresponds to a localized builtin function
		return t.translateBuiltinOrIdent(e.Value)

	case *parser.CallExpression:
		fnName := t.translateBuiltinOrIdent(e.Function)
		var args []string
		for _, arg := range e.Arguments {
			args = append(args, t.transpileExpression(arg))
		}
		return fmt.Sprintf("%s(%s)", fnName, strings.Join(args, ", "))

	case *parser.PostfixExpression:
		// Check for unary negation: 0 - x -> -x
		if il, ok := e.Left.(*parser.IntLiteral); ok && il.Value == 0 && (e.Operator == "алу" || e.Operator == "-") {
			rightStr := t.transpileExpression(e.Right)
			if isLowerPrecedence(e.Right, e.Operator) {
				rightStr = "(" + rightStr + ")"
			}
			return "-" + rightStr
		}

		leftStr := t.transpileExpression(e.Left)
		rightStr := t.transpileExpression(e.Right)
		opStr := t.transpileInfixOperator(e.Operator)

		if isLowerPrecedence(e.Left, e.Operator) {
			leftStr = "(" + leftStr + ")"
		}
		if isLowerPrecedence(e.Right, e.Operator) {
			rightStr = "(" + rightStr + ")"
		}
		return fmt.Sprintf("%s %s %s", leftStr, opStr, rightStr)

	case *parser.UnaryExpression:
		rightStr := t.transpileExpression(e.Right)
		notKw := t.targetLocale.KeywordForToken("NOT")
		return fmt.Sprintf("%s %s", notKw, rightStr)

	case *parser.StructCreateExpression:
		newKw := t.targetLocale.KeywordForToken("NEW")
		return fmt.Sprintf("%s %s", e.StructName, newKw)

	case *parser.StructFieldAccessExpression:
		if e.Target != nil {
			return fmt.Sprintf("%s.%s", t.transpileExpression(e.Target), e.Field)
		}
		return fmt.Sprintf("%s.%s", e.StructName, e.Field)

	case *parser.ArrayLiteral:
		var elements []string
		for _, el := range e.Elements {
			elements = append(elements, t.transpileExpression(el))
		}
		return fmt.Sprintf("[%s]", strings.Join(elements, ", "))

	case *parser.IndexExpression:
		leftStr := t.transpileExpression(e.Left)
		idxStr := t.transpileExpression(e.Index)
		return fmt.Sprintf("%s[%s]", leftStr, idxStr)

	case *parser.LengthExpression:
		valStr := t.transpileExpression(e.Value)
		lenKw := t.targetLocale.LocalNameForBuiltin("ұзындық")
		if lenKw == "" {
			lenKw = "len"
		}
		return fmt.Sprintf("%s(%s)", lenKw, valStr)

	case *parser.CharAtExpression:
		strVal := t.transpileExpression(e.Str)
		idxVal := t.transpileExpression(e.Index)
		return fmt.Sprintf("%s[%s]", strVal, idxVal)

	case *parser.StrConcatExpression:
		leftStr := t.transpileExpression(e.Left)
		rightStr := t.transpileExpression(e.Right)
		return fmt.Sprintf("%s + %s", leftStr, rightStr)

	case *parser.StrLenExpression:
		valStr := t.transpileExpression(e.Value)
		strlenKw := t.targetLocale.LocalNameForBuiltin("ұзындық")
		if strlenKw == "" {
			strlenKw = "len"
		}
		return fmt.Sprintf("%s(%s)", strlenKw, valStr)

	case *parser.StrEqExpression:
		leftStr := t.transpileExpression(e.Left)
		rightStr := t.transpileExpression(e.Right)
		return fmt.Sprintf("%s == %s", leftStr, rightStr)

	case *parser.ToStrExpression:
		valStr := t.transpileExpression(e.Value)
		tostrKw := t.targetLocale.LocalNameForBuiltin("санды_мәтін")
		if tostrKw == "" {
			tostrKw = "str"
		}
		return fmt.Sprintf("%s(%s)", tostrKw, valStr)

	case *parser.CharCodeExpression:
		valStr := t.transpileExpression(e.Value)
		ccKw := t.targetLocale.KeywordForToken("CHAR_CODE")
		return fmt.Sprintf("%s(%s)", ccKw, valStr)

	case *parser.FileReadExpression:
		pathStr := t.transpileExpression(e.Path)
		frKw := t.targetLocale.KeywordForToken("FILE_READ")
		return fmt.Sprintf("%s(%s)", frKw, pathStr)

	case *parser.InputExpression:
		return fmt.Sprintf("%s()", t.targetLocale.KeywordForToken("INPUT"))

	case *parser.ErrorLiteral:
		msgStr := t.transpileExpression(e.Message)
		errKw := t.targetLocale.KeywordForToken("ERROR_LITERAL")
		return fmt.Sprintf("%s %s", msgStr, errKw)

	case *parser.TryErrorExpression:
		leftStr := t.transpileExpression(e.Left)
		tryKw := t.targetLocale.KeywordForToken("TRY_ERROR")
		blockStr := t.transpileBlock(e.Block, 0)
		return fmt.Sprintf("%s %s %s {\n%s}", leftStr, tryKw, e.VarName, blockStr)

	default:
		return expr.String()
	}
}

func (t *Transpiler) transpileOperator(op string) string {
	// Look up the canonical token name for this operator
	canonical := ""
	for _, loc := range locales.Available() {
		if tokName, ok := loc.Keywords[op]; ok {
			canonical = tokName
			break
		}
	}

	if canonical != "" {
		return t.targetLocale.KeywordForToken(canonical)
	}
	return op
}

func (t *Transpiler) transpileType(typeName string, isWeak bool) string {
	baseType := typeName
	weakPrefix := false
	if strings.HasPrefix(typeName, "әлсіз_") {
		baseType = strings.TrimPrefix(typeName, "әлсіз_")
		weakPrefix = true
	} else if strings.HasPrefix(typeName, "weak_") {
		baseType = strings.TrimPrefix(typeName, "weak_")
		weakPrefix = true
	}

	canonicalTok := ""
	switch strings.ToUpper(baseType) {
	case "БҮТІН", "INT", "INTERO", "ЦЕЛОЕ":
		canonicalTok = "TYPE_INT"
	case "САН", "FLOAT", "NUMBER", "NUMERO", "ЧИСЛО":
		canonicalTok = "TYPE_FLOAT"
	case "МӘТІН", "STRING", "TEXT", "TESTO", "СТРОКА":
		canonicalTok = "TYPE_STRING"
	case "АҚИҚАТ", "BOOL", "BOOLEANO", "БУЛЕВО":
		canonicalTok = "TYPE_BOOL"
	case "БАЙТ", "BYTE":
		canonicalTok = "TYPE_BYTE"
	default:
		canonicalTok = ""
	}

	result := baseType
	if canonicalTok != "" {
		result = t.targetLocale.KeywordForToken(canonicalTok)
	}

	if isWeak || weakPrefix {
		weakKw := t.targetLocale.KeywordForToken("WEAK")
		return weakKw + " " + result
	}
	return result
}

func (t *Transpiler) translateBuiltinOrIdent(name string) string {
	// 1. Find canonical builtin name across all locales
	canonical := ""
	for _, loc := range locales.Available() {
		if canon, ok := loc.Builtins[name]; ok {
			canonical = canon
			break
		}
	}

	// 2. If it is a known builtin, translate to target locale
	if canonical != "" {
		return t.targetLocale.LocalNameForBuiltin(canonical)
	}

	// 3. Otherwise return identifier unmodified
	return name
}

func opPrecedence(op string) int {
	switch op {
	case "немесе", "or", "o", "или":
		return 1
	case "және", "and", "e", "и":
		return 2
	case "==", "!=", "тең", "тең_емес":
		return 3
	case "<", "<=", ">", ">=", "кіші", "кіші_тең", "үлкен", "үлкен_тең":
		return 4
	case "+", "-", "қосу", "алу":
		return 5
	case "*", "/", "көбейту", "бөлу":
		return 6
	default:
		return 0
	}
}

func isLowerPrecedence(expr parser.Expression, parentOp string) bool {
	if pe, ok := expr.(*parser.PostfixExpression); ok {
		return opPrecedence(pe.Operator) < opPrecedence(parentOp)
	}
	return false
}

func (t *Transpiler) transpileInfixOperator(op string) string {
	switch op {
	case "қосу", "+":
		return "+"
	case "алу", "-":
		return "-"
	case "көбейту", "*":
		return "*"
	case "бөлу", "/":
		return "/"
	case "тең", "==":
		return "=="
	case "тең_емес", "!=":
		return "!="
	case "кіші", "<":
		return "<"
	case "кіші_тең", "<=":
		return "<="
	case "үлкен", ">":
		return ">"
	case "үлкен_тең", ">=":
		return ">="
	case "және", "and":
		return t.targetLocale.KeywordForToken("AND")
	case "немесе", "or":
		return t.targetLocale.KeywordForToken("OR")
	default:
		return t.transpileOperator(op)
	}
}
