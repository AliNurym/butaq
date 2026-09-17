package parser

import (
	"butaq/lexer"
	"fmt"
	"strconv"
	"strings"
)

type ParseError struct {
	Message string
	Line    int
	Col     int
}

func (e ParseError) Error() string {
	return fmt.Sprintf("[жол %d, баған %d] %s", e.Line, e.Col, e.Message)
}

type Parser struct {
	l *lexer.Lexer

	curToken  lexer.Token
	peekToken lexer.Token
	errors    []ParseError
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []ParseError{}}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Errors() []string {
	var errs []string
	for _, e := range p.errors {
		errs = append(errs, e.Error())
	}
	return errs
}

func (p *Parser) ErrorsStructured() []ParseError {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) errorf(format string, args ...interface{}) {
	p.errors = append(p.errors, ParseError{
		Message: fmt.Sprintf(format, args...),
		Line:    p.curToken.Line,
		Col:     p.curToken.Col,
	})
}

func (p *Parser) setPos(node Node) Node {
	if node != nil {
		node.SetPosition(p.curToken.Line, p.curToken.Col)
	}
	return node
}

func (p *Parser) setPosAt(node Node, line, col int) Node {
	if node != nil {
		node.SetPosition(line, col)
	}
	return node
}

// ---------------------------------------------------------------------------
// Top-level
// ---------------------------------------------------------------------------

func (p *Parser) ParseProgram() *Program {
	program := &Program{Statements: []Statement{}}
	p.setPos(program)
	for p.curToken.Type != lexer.EOF {
		stmt := p.parseTopLevelStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}

func (p *Parser) parseTopLevelStatement() Statement {
	if p.curToken.Type == lexer.FUNC {
		return p.parseFunctionStatement()
	}
	if p.curToken.Type == lexer.STRUCT {
		return p.parseStructStatement()
	}
	if p.curToken.Type == lexer.INTERFACE {
		return p.parseInterfaceStatement()
	}
	return p.parseStatement()
}

// құрылым Адам { аты жасы }
func (p *Parser) parseStructStatement() *StructStatement {
	line, col := p.curToken.Line, p.curToken.Col
	p.nextToken() // cur = name
	if p.curToken.Type != lexer.IDENTIFIER {
		p.errorf("құрылым атауы күтілді")
		return nil
	}
	name := p.curToken.Literal
	p.nextToken() // cur = {

	if p.curToken.Type != lexer.LBRACE {
		p.errorf("'{' күтілді")
		return nil
	}
	p.nextToken() // cur = first field
 
	var fields []string
	var types []string
	var isWeakList []bool
	for p.curToken.Type != lexer.RBRACE && p.curToken.Type != lexer.EOF {
		if p.curToken.Type == lexer.IDENTIFIER {
			fieldName := p.curToken.Literal
			fieldType := "САН" // default type is float/number
			isWeak := false

			if p.peekToken.Type == lexer.WEAK {
				p.nextToken() // cur = әлсіз
				isWeak = true
				if p.peekToken.Type == lexer.TYPE_INT ||
					p.peekToken.Type == lexer.TYPE_FLOAT ||
					p.peekToken.Type == lexer.TYPE_STRING ||
					p.peekToken.Type == lexer.TYPE_BOOL ||
					p.peekToken.Type == lexer.TYPE_BYTE ||
					p.peekToken.Type == lexer.IDENTIFIER {
					p.nextToken() // cur = type
					fieldType = "әлсіз_" + p.curToken.Literal
				} else {
					p.errorf("әлсіз сілтеме үшін тип күтілді")
					return nil
				}
			} else if p.peekToken.Type == lexer.TYPE_INT ||
				p.peekToken.Type == lexer.TYPE_FLOAT ||
				p.peekToken.Type == lexer.TYPE_STRING ||
				p.peekToken.Type == lexer.TYPE_BOOL ||
				p.peekToken.Type == lexer.TYPE_BYTE ||
				p.peekToken.Type == lexer.IDENTIFIER {
				p.nextToken()
				fieldType = p.curToken.Literal
			}

			fields = append(fields, fieldName)
			types = append(types, fieldType)
			isWeakList = append(isWeakList, isWeak)
		}
		p.nextToken()
	}
	// cur = }

	res := &StructStatement{Name: name, Fields: fields, Types: types, IsWeak: isWeakList}
	p.setPosAt(res, line, col)
	return res
}

// интерфейс Көлік { жүру(САН) САН тоқтау() }
func (p *Parser) parseInterfaceStatement() *InterfaceStatement {
	line, col := p.curToken.Line, p.curToken.Col
	p.nextToken() // cur = name
	if p.curToken.Type != lexer.IDENTIFIER {
		p.errorf("интерфейс атауы күтілді")
		return nil
	}
	name := p.curToken.Literal
	p.nextToken() // cur = {

	if p.curToken.Type != lexer.LBRACE {
		p.errorf("'{' күтілді")
		return nil
	}
	p.nextToken() // cur = first method or }

	var methods []MethodSignature
	for p.curToken.Type != lexer.RBRACE && p.curToken.Type != lexer.EOF {
		if p.curToken.Type == lexer.IDENTIFIER {
			methodName := p.curToken.Literal
			p.nextToken() // cur = (
			if p.curToken.Type != lexer.LPAREN {
				p.errorf("'(' күтілді")
				return nil
			}
			p.nextToken() // cur = first param type or )

			var params []string
			for p.curToken.Type != lexer.RPAREN && p.curToken.Type != lexer.EOF {
				if p.curToken.Type == lexer.TYPE_INT ||
					p.curToken.Type == lexer.TYPE_FLOAT ||
					p.curToken.Type == lexer.TYPE_STRING ||
					p.curToken.Type == lexer.TYPE_BOOL ||
					p.curToken.Type == lexer.TYPE_BYTE ||
					p.curToken.Type == lexer.IDENTIFIER {
					params = append(params, p.curToken.Literal)
				}
				if p.peekToken.Type == lexer.COMMA {
					p.nextToken()
				}
				p.nextToken()
			}
			// cur = )
			
			returnType := ""
			// Проверяем, есть ли после ) возвращаемый тип
			if p.peekToken.Type == lexer.TYPE_INT ||
				p.peekToken.Type == lexer.TYPE_FLOAT ||
				p.peekToken.Type == lexer.TYPE_STRING ||
				p.peekToken.Type == lexer.TYPE_BOOL ||
				p.peekToken.Type == lexer.TYPE_BYTE ||
				p.peekToken.Type == lexer.IDENTIFIER {
				p.nextToken() // cur = return type
				returnType = p.curToken.Literal
			}

			methods = append(methods, MethodSignature{
				Name:       methodName,
				Parameters: params,
				ReturnType: returnType,
			})
		}
		p.nextToken()
	}

	res := &InterfaceStatement{Name: name, Methods: methods}
	p.setPosAt(res, line, col)
	return res
}

// ---------------------------------------------------------------------------
// Statement dispatch
// ---------------------------------------------------------------------------

const (
	PREC_LOWEST  = 0
	PREC_OR      = 1
	PREC_AND     = 2
	PREC_EQUALS  = 3
	PREC_COMPARE = 4
	PREC_SUM     = 5
	PREC_PRODUCT = 6
	PREC_PREFIX  = 7
	PREC_CALL    = 8
)

func (p *Parser) tokenPrecedence(tt lexer.TokenType) int {
	switch tt {
	case lexer.OR:
		return PREC_OR
	case lexer.AND:
		return PREC_AND
	case lexer.EQ, lexer.NEQ:
		return PREC_EQUALS
	case lexer.LT, lexer.LTE, lexer.GT, lexer.GTE:
		return PREC_COMPARE
	case lexer.PLUS, lexer.MINUS:
		return PREC_SUM
	case lexer.MUL, lexer.DIV:
		return PREC_PRODUCT
	default:
		return PREC_LOWEST
	}
}

func (p *Parser) isStartOfInfixExpr() bool {
	switch p.curToken.Type {
	case lexer.INT_LITERAL, lexer.NUMBER, lexer.STRING, lexer.TRUE, lexer.FALSE,
		lexer.IDENTIFIER, lexer.NOT, lexer.MINUS, lexer.LPAREN, lexer.LBRACKET:
		return true
	default:
		return false
	}
}

func (p *Parser) parseInfixExpression(precedence int) Expression {
	line, col := p.curToken.Line, p.curToken.Col
	var left Expression

	switch p.curToken.Type {
	case lexer.INT_LITERAL:
		val, _ := strconv.ParseInt(p.curToken.Literal, 10, 64)
		left = p.setPosAt(&IntLiteral{Value: val}, line, col).(Expression)
	case lexer.NUMBER:
		val, _ := strconv.ParseFloat(p.curToken.Literal, 64)
		left = p.setPosAt(&NumberLiteral{Value: val}, line, col).(Expression)
	case lexer.STRING:
		left = p.setPosAt(&StringLiteral{Value: p.curToken.Literal}, line, col).(Expression)
	case lexer.TRUE:
		left = p.setPosAt(&BoolLiteral{Value: true}, line, col).(Expression)
	case lexer.FALSE:
		left = p.setPosAt(&BoolLiteral{Value: false}, line, col).(Expression)
	case lexer.NOT:
		p.nextToken()
		right := p.parseInfixExpression(PREC_PREFIX)
		left = p.setPosAt(&UnaryExpression{Operator: "емес", Right: right}, line, col).(Expression)
	case lexer.MINUS:
		p.nextToken()
		right := p.parseInfixExpression(PREC_PREFIX)
		left = p.setPosAt(&PostfixExpression{
			Left:     &IntLiteral{Value: 0},
			Right:    right,
			Operator: "алу",
		}, line, col).(Expression)
	case lexer.LPAREN:
		p.nextToken()
		left = p.parseInfixExpression(PREC_LOWEST)
		if p.peekToken.Type == lexer.RPAREN {
			p.nextToken()
		}
	case lexer.LBRACKET:
		var elements []Expression
		if p.peekToken.Type != lexer.RBRACKET {
			p.nextToken() // cur = first element
			for p.curToken.Type != lexer.RBRACKET && p.curToken.Type != lexer.EOF {
				elem := p.parseInfixExpression(PREC_LOWEST)
				if elem != nil {
					elements = append(elements, elem)
				}
				if p.peekToken.Type == lexer.COMMA {
					p.nextToken() // cur = ,
					p.nextToken() // cur = next elem
				} else if p.peekToken.Type == lexer.RBRACKET {
					p.nextToken() // cur = ]
					break
				} else {
					p.nextToken()
				}
			}
		} else {
			p.nextToken() // cur = ]
		}
		left = p.setPosAt(&ArrayLiteral{Elements: elements}, line, col).(Expression)
	case lexer.IDENTIFIER:
		id := p.curToken.Literal
		if p.peekToken.Type == lexer.LPAREN {
			left = p.parseCallExpression(id, line, col)
		} else if strings.Contains(id, ".") {
			parts := strings.Split(id, ".")
			left = p.setPosAt(&StructFieldAccessExpression{
				StructName: parts[0],
				Field:      parts[1],
			}, line, col).(Expression)
		} else {
			left = p.setPosAt(&Identifier{Value: id}, line, col).(Expression)
		}
	default:
		return nil
	}

	for p.peekToken.Type == lexer.LBRACKET {
		p.nextToken() // cur = [
		p.nextToken() // cur = start of index
		idx := p.parseInfixExpression(PREC_LOWEST)
		if p.peekToken.Type == lexer.RBRACKET {
			p.nextToken() // cur = ]
		}
		left = p.setPosAt(&IndexExpression{
			Left:  left,
			Index: idx,
		}, line, col).(Expression)
	}

	for p.peekToken.Type != lexer.EOF &&
		p.peekToken.Type != lexer.THEN &&
		p.peekToken.Type != lexer.COMMA &&
		p.peekToken.Type != lexer.RPAREN &&
		p.peekToken.Type != lexer.RBRACKET &&
		p.peekToken.Type != lexer.COLON &&
		precedence < p.tokenPrecedence(p.peekToken.Type) {

		p.nextToken() // cur = operator
		opTok := p.curToken
		prec := p.tokenPrecedence(opTok.Type)
		p.nextToken() // cur = start of right expression
		right := p.parseInfixExpression(prec)
		left = p.setPosAt(&PostfixExpression{
			Left:     left,
			Right:    right,
			Operator: canonicalOp(opTok),
		}, line, col).(Expression)
	}

	return left
}

func (p *Parser) parseIndentedBlock(parentCol int) *BlockStatement {
	block := &BlockStatement{Statements: []Statement{}}
	p.setPosAt(block, p.curToken.Line, p.curToken.Col)

	startLine := p.curToken.Line
	for p.peekToken.Type != lexer.EOF {
		if p.peekToken.Line > startLine && p.peekToken.Col <= parentCol {
			break
		}
		if p.peekToken.Type == lexer.RBRACE || p.peekToken.Type == lexer.ELSE || p.peekToken.Type == lexer.FUNC {
			break
		}
		p.nextToken()
		if p.curToken.Type == lexer.COLON {
			continue
		}
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
	}
	return block
}

func (p *Parser) parseStatement() Statement {
	line, col := p.curToken.Line, p.curToken.Col

	// 1. Modern VAR: болсын i = 0 / let i = 0 / sia i = 0 / пусть i = 0
	if p.curToken.Type == lexer.VAR && p.peekToken.Type == lexer.IDENTIFIER {
		p.nextToken() // cur = identifier
		varName := p.curToken.Literal
		idNode := p.setPos(&Identifier{Value: varName}).(*Identifier)
		if p.peekToken.Type == lexer.ASSIGN {
			p.nextToken() // cur = =
			p.nextToken() // cur = start of value expr
			val := p.parseInfixExpression(PREC_LOWEST)
			stmt := &VarAssignStatement{Name: idNode, Value: val, IsDeclaration: true}
			p.setPosAt(stmt, line, col)
			return stmt
		}
	}

	// 2. Modern Assignment: i = expr
	if p.curToken.Type == lexer.IDENTIFIER && p.peekToken.Type == lexer.ASSIGN {
		varName := p.curToken.Literal
		idNode := p.setPos(&Identifier{Value: varName}).(*Identifier)
		p.nextToken() // cur = =
		p.nextToken() // cur = start of value expr
		val := p.parseInfixExpression(PREC_LOWEST)
		stmt := &VarAssignStatement{Name: idNode, Value: val, IsDeclaration: false}
		p.setPosAt(stmt, line, col)
		return stmt
	}

	// 2b. Modern Array Index Assignment: arr[idx] = expr
	if p.curToken.Type == lexer.IDENTIFIER && p.peekToken.Type == lexer.LBRACKET {
		arrName := p.curToken.Literal
		idNode := p.setPos(&Identifier{Value: arrName}).(*Identifier)
		p.nextToken() // cur = [
		p.nextToken() // cur = start of index expr
		idxExpr := p.parseInfixExpression(PREC_LOWEST)
		if p.peekToken.Type == lexer.RBRACKET {
			p.nextToken() // cur = ]
		}
		if p.peekToken.Type == lexer.ASSIGN {
			p.nextToken() // cur = =
			p.nextToken() // cur = start of value expr
			valExpr := p.parseInfixExpression(PREC_LOWEST)
			stmt := &IndexAssignStatement{Array: idNode, Index: idxExpr, Value: valExpr}
			p.setPosAt(stmt, line, col)
			return stmt
		}
	}

	// 3. Modern Print: жазу(a, b, c) / print(...) / stampa(...) / печать(...)
	if p.curToken.Type == lexer.PRINT && p.peekToken.Type == lexer.LPAREN {
		p.nextToken() // cur = (
		p.nextToken() // cur = first arg or )
		var vals []Expression
		for p.curToken.Type != lexer.RPAREN && p.curToken.Type != lexer.EOF {
			expr := p.parseInfixExpression(PREC_LOWEST)
			if expr != nil {
				vals = append(vals, expr)
			}
			if p.peekToken.Type == lexer.COMMA {
				p.nextToken() // cur = ,
				p.nextToken() // cur = next arg
			} else if p.peekToken.Type == lexer.RPAREN {
				p.nextToken() // cur = )
				break
			} else {
				p.nextToken()
			}
		}
		stmt := &PrintStatement{Values: vals}
		p.setPosAt(stmt, line, col)
		return stmt
	}

	// 4. Modern Return: қайтару expr / return expr / ritorna expr / вернуть expr
	if p.curToken.Type == lexer.RETURN {
		p.nextToken() // cur = start of expr
		val := p.parseInfixExpression(PREC_LOWEST)
		stmt := &ReturnStatement{Value: val}
		p.setPosAt(stmt, line, col)
		return stmt
	}

	// 5. Modern IF: егер cond сонда stmt OR егер cond block
	if p.curToken.Type == lexer.IF {
		p.nextToken() // cur = start of condition
		cond := p.parseInfixExpression(PREC_LOWEST)
		stmt := &IfStatement{Condition: cond}
		p.setPosAt(stmt, line, col)

		// Check for THEN: сонда / then / allora / тогда
		if p.peekToken.Type == lexer.THEN {
			p.nextToken() // cur = THEN
			p.nextToken() // cur = start of then statement
			thenStmt := p.parseStatement()
			stmt.Consequence = &BlockStatement{Statements: []Statement{thenStmt}}
			return stmt
		}

		if p.peekToken.Type == lexer.LBRACE {
			p.nextToken() // cur = {
			stmt.Consequence = p.parseBlockStatement()
		} else {
			if p.peekToken.Type == lexer.COLON {
				p.nextToken()
			}
			stmt.Consequence = p.parseIndentedBlock(col)
		}

		if p.peekToken.Type == lexer.ELSE {
			p.nextToken() // cur = ELSE
			if p.peekToken.Type == lexer.LBRACE {
				p.nextToken() // cur = {
				stmt.Alternative = p.parseBlockStatement()
			} else {
				if p.peekToken.Type == lexer.COLON {
					p.nextToken()
				}
				stmt.Alternative = p.parseIndentedBlock(col)
			}
		}

		return stmt
	}

	// 6. Modern WHILE: әзірше cond block
	if p.curToken.Type == lexer.WHILE {
		p.nextToken() // cur = start of condition
		cond := p.parseInfixExpression(PREC_LOWEST)
		stmt := &WhileStatement{Condition: cond}
		p.setPosAt(stmt, line, col)

		if p.peekToken.Type == lexer.LBRACE {
			p.nextToken() // cur = {
			stmt.Body = p.parseBlockStatement()
		} else {
			if p.peekToken.Type == lexer.COLON {
				p.nextToken()
			}
			stmt.Body = p.parseIndentedBlock(col)
		}
		return stmt
	}

	// 7. Fallback to legacy SOV parser
	node := p.parseExpression()
	if node == nil {
		return nil
	}
	if s, ok := node.(Statement); ok {
		return s
	}
	if expr, ok := node.(Expression); ok {
		el, ec := expr.Position()
		res := &ExpressionStatement{Expression: expr}
		p.setPosAt(res, el, ec)
		return res
	}
	return nil
}

// ---------------------------------------------------------------------------
// Function definition: функция атауы(x, y) [block]
// ---------------------------------------------------------------------------

func (p *Parser) parseFunctionStatement() *FunctionStatement {
	line, col := p.curToken.Line, p.curToken.Col
	// cur = функция
	p.nextToken() // cur = name

	if p.curToken.Type != lexer.IDENTIFIER {
		p.errorf("функция атауы күтілді, бірақ '%s' табылды", p.curToken.Literal)
		return nil
	}
	name := p.curToken.Literal
	p.nextToken() // cur = (

	if p.curToken.Type != lexer.LPAREN {
		p.errorf("'(' күтілді функция параметрлері алдында, бірақ '%s' табылды", p.curToken.Literal)
		return nil
	}
	p.nextToken() // cur = first param or )

	var params []string
	for p.curToken.Type != lexer.RPAREN && p.curToken.Type != lexer.EOF {
		if p.curToken.Type == lexer.IDENTIFIER {
			params = append(params, p.curToken.Literal)
		}
		p.nextToken()
		if p.curToken.Type == lexer.COMMA {
			p.nextToken()
		}
	}
	// cur = )
	if p.peekToken.Type == lexer.COLON {
		p.nextToken() // consume optional colon
	}

	var body *BlockStatement
	if p.peekToken.Type == lexer.LBRACE {
		p.nextToken() // cur = {
		body = p.parseBlockStatement()
	} else {
		body = p.parseIndentedBlock(col)
	}

	res := &FunctionStatement{Name: name, Parameters: params, Body: body}
	p.setPosAt(res, line, col)
	return res
}


// ---------------------------------------------------------------------------
// Main expression parser (stack-based SOV)
// ---------------------------------------------------------------------------

func (p *Parser) parseExpression() Node {
	var stack []Node

	for p.curToken.Type != lexer.EOF &&
		p.curToken.Type != lexer.LBRACE &&
		p.curToken.Type != lexer.RBRACE {

		switch p.curToken.Type {

		// --- Literals ---
		case lexer.NUMBER:
			val, err := strconv.ParseFloat(p.curToken.Literal, 64)
			if err != nil {
				p.errorf("жарамсыз сан: '%s'", p.curToken.Literal)
				return nil
			}
			stack = append(stack, p.setPos(&NumberLiteral{Value: val}))

		case lexer.INT_LITERAL:
			val, err := strconv.ParseInt(p.curToken.Literal, 10, 64)
			if err != nil {
				p.errorf("жарамсыз бүтін сан: '%s'", p.curToken.Literal)
				return nil
			}
			stack = append(stack, p.setPos(&IntLiteral{Value: val}))

		case lexer.STRING:
			stack = append(stack, p.setPos(&StringLiteral{Value: p.curToken.Literal}))

		case lexer.IDENTIFIER:
			// Peek ahead — if next is '(' this is a call expression
			if p.peekToken.Type == lexer.LPAREN {
				callExpr := p.parseCallExpression(p.curToken.Literal, p.curToken.Line, p.curToken.Col)
				if callExpr == nil {
					return nil
				}
				stack = append(stack, callExpr)
			} else {
				id := p.curToken.Literal
				if len(id) > 0 && id[0] == '.' {
					if len(stack) < 1 {
						p.errorf("өріске кіру үшін объект қажет")
						return nil
					}
					target := stack[len(stack)-1].(Expression)
					stack = stack[:len(stack)-1]
					stack = append(stack, p.setPos(&StructFieldAccessExpression{
						StructName: "",
						Field:      id[1:],
						Target:     target,
					}))
				} else {
					isStructField := false
					for i := 0; i < len(id); i++ {
						if id[i] == '.' {
							stack = append(stack, p.setPos(&StructFieldAccessExpression{
								StructName: id[:i],
								Field:      id[i+1:],
							}))
							isStructField = true
							break
						}
					}
					if !isStructField {
						stack = append(stack, p.setPos(&Identifier{Value: id}))
					}
				}
			}

		case lexer.TRUE:
			stack = append(stack, p.setPos(&BoolLiteral{Value: true}))

		case lexer.FALSE:
			stack = append(stack, p.setPos(&BoolLiteral{Value: false}))

		// --- Array literal: [ ... ] ---
		case lexer.LBRACKET:
			arr := p.parseArrayLiteral()
			if arr == nil {
				return nil
			}
			stack = append(stack, arr)

		// --- Binary operators: need 2 operands ---
		case lexer.PLUS, lexer.MINUS, lexer.MUL, lexer.DIV,
			lexer.GT, lexer.LT, lexer.EQ, lexer.NEQ, lexer.GTE, lexer.LTE,
			lexer.AND, lexer.OR, lexer.LSHIFT, lexer.RSHIFT:
			if len(stack) < 2 {
				p.errorf("'%s' операторы 2 операнд талап етеді", p.curToken.Literal)
				return nil
			}
			right := stack[len(stack)-1].(Expression)
			left := stack[len(stack)-2].(Expression)
			stack = stack[:len(stack)-2]
			stack = append(stack, p.setPos(&PostfixExpression{Left: left, Right: right, Operator: canonicalOp(p.curToken)}))


		// --- Unary NOT ---
		case lexer.NOT:
			if len(stack) < 1 {
				p.errorf("'емес' операторы 1 операнд талап етеді")
				return nil
			}
			val := stack[len(stack)-1].(Expression)
			stack = stack[:len(stack)-1]
			stack = append(stack, p.setPos(&UnaryExpression{Operator: "емес", Right: val}))

		// --- Array length: arr ұзындық ---
		case lexer.ARRAY_LEN:
			if len(stack) < 1 {
				p.errorf("'ұзындық' 1 операнд талап етеді")
				return nil
			}
			val := stack[len(stack)-1].(Expression)
			stack = stack[:len(stack)-1]
			stack = append(stack, p.setPos(&LengthExpression{Value: val}))

		// --- Char at: str idx символ ---
		case lexer.CHAR_AT:
			if len(stack) < 2 {
				p.errorf("'символ' 2 операнд талап етеді (мәтін, индекс)")
				return nil
			}
			idx := stack[len(stack)-1].(Expression)
			str := stack[len(stack)-2].(Expression)
			stack = stack[:len(stack)-2]
			stack = append(stack, p.setPos(&CharAtExpression{Str: str, Index: idx}))

		// --- String concat: str1 str2 біріктіру ---
		case lexer.STR_CONCAT:
			if len(stack) < 2 {
				p.errorf("'біріктіру' 2 операнд талап етеді (мәтін1, мәтін2)")
				return nil
			}
			right := stack[len(stack)-1].(Expression)
			left := stack[len(stack)-2].(Expression)
			stack = stack[:len(stack)-2]
			stack = append(stack, p.setPos(&StrConcatExpression{Left: left, Right: right}))

		// --- String length: str ұзындық_жол ---
		case lexer.STR_LEN:
			if len(stack) < 1 {
				p.errorf("'ұзындық_жол' 1 операнд талап етеді")
				return nil
			}
			val := stack[len(stack)-1].(Expression)
			stack = stack[:len(stack)-1]
			stack = append(stack, p.setPos(&StrLenExpression{Value: val}))

		// --- String equals: str1 str2 мәтін_тең ---
		case lexer.STR_EQ:
			if len(stack) < 2 {
				p.errorf("'мәтін_тең' 2 операнд талап етеді")
				return nil
			}
			right := stack[len(stack)-1].(Expression)
			left := stack[len(stack)-2].(Expression)
			stack = stack[:len(stack)-2]
			stack = append(stack, p.setPos(&StrEqExpression{Left: left, Right: right}))

		// --- Number to string: num санды_мәтін ---
		case lexer.TO_STR:
			if len(stack) < 1 {
				p.errorf("'санды_мәтін' 1 операнд талап етеді")
				return nil
			}
			val := stack[len(stack)-1].(Expression)
			stack = stack[:len(stack)-1]
			stack = append(stack, p.setPos(&ToStrExpression{Value: val}))

		// --- Char code: str таңба_коды → byte value ---
		case lexer.CHAR_CODE:
			if len(stack) < 1 {
				p.errorf("'таңба_коды' 1 операнд талап етеді")
				return nil
			}
			val := stack[len(stack)-1].(Expression)
			stack = stack[:len(stack)-1]
			stack = append(stack, p.setPos(&CharCodeExpression{Value: val}))

		// --- PRINT: <value> жазу ---
		case lexer.PRINT:
			if len(stack) < 1 {
				p.errorf("'жазу' 1 операнд талап етеді")
				return nil
			}
			val := stack[len(stack)-1].(Expression)
			stack = stack[:len(stack)-1]
			stmt := &PrintStatement{Value: val}
			p.setPos(stmt)
			stack = append(stack, stmt)
			return stack[0]

		// --- RETURN: қайтару <value> ---
		case lexer.RETURN:
			if len(stack) < 1 {
				p.errorf("'қайтару' 1 операнд талап етеді")
				return nil
			}
			val := stack[len(stack)-1].(Expression)
			stack = stack[:len(stack)-1]
			stmt := &ReturnStatement{Value: val}
			p.setPos(stmt)
			stack = append(stack, stmt)
			return stack[0]

		// --- CALL (шақыру) as statement: funcName(args) шақыру ---
		case lexer.CALL:
			if len(stack) < 1 {
				p.errorf("'шақыру' алдында функция шақыруы болуы керек")
				return nil
			}
			top := stack[len(stack)-1]
			callExpr, ok := top.(*CallExpression)
			if !ok {
				p.errorf("'шақыру' функция шақыруымен ғана жұмыс істейді")
				return nil
			}
			stack = stack[:len(stack)-1]
			stmt := &CallStatement{Call: callExpr}
			p.setPos(stmt)
			stack = append(stack, stmt)
			return stack[0]

		// --- FILE READ: "path" файл_оқу → FileReadExpression ---
		case lexer.FILE_READ:
			if len(stack) < 1 {
				p.errorf("'файл_оқу' 1 операнд талап етеді (жол атауы)")
				return nil
			}
			path := stack[len(stack)-1].(Expression)
			stack = stack[:len(stack)-1]
			stack = append(stack, p.setPos(&FileReadExpression{Path: path}))

		// --- INPUT: кіру → InputExpression (reads a line from stdin) ---
		case lexer.INPUT:
			stack = append(stack, p.setPos(&InputExpression{}))

		// --- ERROR: msg қате → ErrorLiteral ---
		case lexer.ERROR_LITERAL:
			if len(stack) < 1 {
				p.errorf("'қате' 1 операнд талап етеді (хабарлама)")
				return nil
			}
			msg := stack[len(stack)-1].(Expression)
			stack = stack[:len(stack)-1]
			stack = append(stack, p.setPos(&ErrorLiteral{Token: p.curToken, Message: msg}))

		// --- TRY_ERROR: expr қатемен var { block } → TryErrorExpression ---
		case lexer.TRY_ERROR:
			if len(stack) < 1 {
				p.errorf("'қатемен' алдында өрнек болуы керек")
				return nil
			}
			left := stack[len(stack)-1].(Expression)
			stack = stack[:len(stack)-1]

			tryToken := p.curToken

			// Считываем имя переменной
			p.nextToken()
			if p.curToken.Type != lexer.IDENTIFIER {
				p.errorf("'қатемен' кейін айнымалы атауы күтілді, бірақ '%s' табылды", p.curToken.Literal)
				return nil
			}
			varName := p.curToken.Literal

			// Считываем '{'
			p.nextToken()
			if p.curToken.Type != lexer.LBRACE {
				p.errorf("'қатемен' айнымалысынан кейін '{' күтілді, бірақ '%s' табылды", p.curToken.Literal)
				return nil
			}
			block := p.parseBlockStatement()
			if block == nil {
				return nil
			}

			stmt := &TryErrorExpression{
				Token:   tryToken,
				Left:    left,
				VarName: varName,
				Block:   block,
			}
			p.setPosAt(stmt, tryToken.Line, tryToken.Col)
			stack = append(stack, stmt)

		// --- FILE WRITE: path content файл_жазу → FileWriteStatement ---
		case lexer.FILE_WRITE:
			if len(stack) < 2 {
				p.errorf("'файл_жазу' 2 операнд талап етеді (жол, мазмұн)")
				return nil
			}
			content := stack[len(stack)-1].(Expression)
			path := stack[len(stack)-2].(Expression)
			stack = stack[:len(stack)-2]
			stmt := &FileWriteStatement{Path: path, Content: content}
			p.setPos(stmt)
			stack = append(stack, stmt)
			return stack[0]

		// --- IF: <cond> егер { ... } [әйтпесе { ... }] ---
		case lexer.IF:
			if len(stack) < 1 {
				p.errorf("'егер' шарты жоқ")
				return nil
			}
			cond := stack[len(stack)-1].(Expression)
			stack = stack[:len(stack)-1]

			stmt := &IfStatement{Condition: cond}
			p.setPos(stmt)
			p.nextToken()
			if p.curToken.Type != lexer.LBRACE {
				p.errorf("'егер' кейін '{' күтілді, бірақ '%s' табылды", p.curToken.Literal)
				return nil
			}
			stmt.Consequence = p.parseBlockStatement()

			if p.peekToken.Type == lexer.ELSE {
				p.nextToken() // cur = әйтпесе
				p.nextToken() // cur = {
				if p.curToken.Type != lexer.LBRACE {
					p.errorf("'әйтпесе' кейін '{' күтілді, бірақ '%s' табылды", p.curToken.Literal)
					return nil
				}
				stmt.Alternative = p.parseBlockStatement()
			}

			stack = append(stack, stmt)
			return stack[0]

		// --- THREAD: ағын { ... } немесе ағын функция() ---
		case lexer.THREAD:
			threadToken := p.curToken
			p.nextToken()
			var body Statement
			if p.curToken.Type == lexer.LBRACE {
				body = p.parseBlockStatement()
			} else if p.curToken.Type == lexer.IDENTIFIER && p.peekToken.Type == lexer.LPAREN {
				name := p.curToken.Literal
				callExpr := p.parseCallExpression(name, p.curToken.Line, p.curToken.Col)
				if callExpr == nil {
					return nil
				}
				body = &ExpressionStatement{Expression: callExpr}
				el, ec := callExpr.Position()
				p.setPosAt(body, el, ec)
			} else {
				p.errorf("ағыннан кейін '{' немесе функцияны шақыру күтілді, бірақ '%s' табылды", p.curToken.Literal)
				return nil
			}
			stmt := &ThreadStatement{Token: threadToken, Body: body}
			p.setPosAt(stmt, threadToken.Line, threadToken.Col)
			stack = append(stack, stmt)

		// --- WHILE: <cond> әзірше { ... } ---
		case lexer.WHILE:
			if len(stack) < 1 {
				p.errorf("'әзірше' шарты жоқ")
				return nil
			}
			cond := stack[len(stack)-1].(Expression)
			stack = stack[:len(stack)-1]

			stmt := &WhileStatement{Condition: cond}
			p.setPos(stmt)
			p.nextToken()
			if p.curToken.Type != lexer.LBRACE {
				p.errorf("'әзірше' кейін '{' күтілді, бірақ '%s' табылды", p.curToken.Literal)
				return nil
			}
			stmt.Body = p.parseBlockStatement()

			stack = append(stack, stmt)
			return stack[0]

		// --- ILLEGAL token ---
		case lexer.ILLEGAL:
			p.errorf("жарамсыз таңба: '%s'", p.curToken.Literal)
			return nil

		// --- INDEX GET: arr idx тізім_алу ---
		case lexer.INDEX_GET:
			if len(stack) < 2 {
				p.errorf("'тізім_алу' 2 операнд талап етеді (тізім, индекс)")
				return nil
			}
			idx := stack[len(stack)-1].(Expression)
			arr := stack[len(stack)-2].(Expression)
			stack = stack[:len(stack)-2]
			stack = append(stack, p.setPos(&IndexExpression{Left: arr, Index: idx}))

		// --- INDEX SET: arr idx val тізім_қой ---
		case lexer.INDEX_SET:
			if len(stack) < 3 {
				p.errorf("'тізім_қой' 3 операнд талап етеді (тізім, индекс, мән)")
				return nil
			}
			val := stack[len(stack)-1].(Expression)
			idxSet := stack[len(stack)-2].(Expression)
			arrNode, ok := stack[len(stack)-3].(*Identifier)
			if !ok {
				p.errorf("'тізім_қой' бірінші операнды идентификатор болуы керек")
				return nil
			}
			stack = stack[:len(stack)-3]
			stmt := &IndexAssignStatement{Array: arrNode, Index: idxSet, Value: val}
			p.setPos(stmt)
			stack = append(stack, stmt)
			return stack[0]

		// --- Struct Create: Адам жасау ---
		case lexer.NEW:
			if len(stack) < 1 {
				p.errorf("'жасау' 1 операнд талап етеді (құрылым атауы)")
				return nil
			}
			val := stack[len(stack)-1].(Expression)
			stack = stack[:len(stack)-1]

			id, ok := val.(*Identifier)
			if !ok {
				p.errorf("'жасау' алдында құрылым атауы болуы керек")
				return nil
			}
			stack = append(stack, p.setPos(&StructCreateExpression{StructName: id.Value}))

		// --- Variable Assign: x 10 болсын, OR adam.name "Ali" болсын ---
		case lexer.VAR:
			if len(stack) < 2 {
				p.errorf("айнымалыға меншіктеу ('болсын') 2 операнд талап етеді: мән және айнымалы")
				return nil
			}
			val := stack[len(stack)-1].(Expression)
			idExpr := stack[len(stack)-2].(Expression)
			stack = stack[:len(stack)-2]

			idNode, ok := idExpr.(*Identifier)
			if !ok {
				if sfa, isSfa := idExpr.(*StructFieldAccessExpression); isSfa {
					stmt := &StructFieldAssignStatement{
						StructName: sfa.StructName,
						Field:      sfa.Field,
						Value:      val,
						Target:     sfa.Target,
					}
					p.setPos(stmt)
					stack = append(stack, stmt)
					return stack[0]
				}

				p.errorf("айнымалы атауы жарамсыз: %v", idExpr)
				return nil
			}

			stmt := &VarAssignStatement{Name: idNode, Value: val, IsDeclaration: true}
			p.setPos(stmt)
			stack = append(stack, stmt)
			return stack[0]

		// --- BREAK: үзу ---
		case lexer.BREAK:
			stmt := &BreakStatement{}
			p.setPos(stmt)
			stack = append(stack, stmt)
			return stack[0]

		// --- CONTINUE: жалғастыру ---
		case lexer.CONTINUE:
			stmt := &ContinueStatement{}
			p.setPos(stmt)
			stack = append(stack, stmt)
			return stack[0]

		// --- FREE: arr бос ---
		case lexer.FREE:
			if len(stack) < 1 {
				p.errorf("'бос' 1 операнд талап етеді")
				return nil
			}
			val := stack[len(stack)-1].(Expression)
			stack = stack[:len(stack)-1]
			stmt := &FreeStatement{Value: val}
			p.setPos(stmt)
			stack = append(stack, stmt)
			return stack[0]

		// --- IMPORT: "path" енгізу ---
		case lexer.IMPORT:
			if len(stack) < 1 {
				p.errorf("'енгізу' 1 операнд талап етеді (жол атауы)")
				return nil
			}
			pathNode, ok := stack[len(stack)-1].(*StringLiteral)
			if !ok {
				p.errorf("'енгізу' жол атауын талап етеді")
				return nil
			}
			stack = stack[:len(stack)-1]
			stmt := &ImportStatement{Path: pathNode.Value}
			p.setPos(stmt)
			stack = append(stack, stmt)
			return stack[0]
		}

		// болсын: assignment — triggered when VAR token literal is "болсын"
		if p.curToken.Type == lexer.VAR && p.curToken.Literal == "болсын" {
			if len(stack) < 2 {
				p.errorf("'болсын' идентификатор мен мән талап етеді")
				return nil
			}
			value := stack[len(stack)-1].(Expression)
			ident, ok := stack[len(stack)-2].(*Identifier)
			if !ok {
				p.errorf("'болсын' сол жағы идентификатор болуы керек")
				return nil
			}
			stack = stack[:len(stack)-2]
			stmt := &VarAssignStatement{Name: ident, Value: value, IsDeclaration: true}
			p.setPos(stmt)
			stack = append(stack, stmt)
			return stack[0]
		}

		// Check for early termination hints
		if p.peekToken.Type == lexer.LBRACE || p.peekToken.Type == lexer.EOF || p.peekToken.Line > p.curToken.Line {
			break
		}

		p.nextToken()
	}

	if len(stack) > 0 {
		return stack[len(stack)-1]
	}
	return nil
}

// ---------------------------------------------------------------------------
// Call expression: funcName(arg1, arg2)
// ---------------------------------------------------------------------------

func (p *Parser) parseCallExpression(name string, line, col int) *CallExpression {
	// cur = funcName, peek = (
	p.nextToken() // cur = (
	p.nextToken() // cur = first arg or )

	var args []Expression
	for p.curToken.Type != lexer.RPAREN && p.curToken.Type != lexer.EOF {
		var expr Expression
		if p.isStartOfInfixExpr() {
			expr = p.parseInfixExpression(PREC_LOWEST)
		}
		if expr == nil {
			expr = p.parseArgExpression()
		}
		if expr != nil {
			args = append(args, expr)
		}
		if p.peekToken.Type == lexer.COMMA {
			p.nextToken() // cur = ,
			p.nextToken() // cur = next arg
		} else if p.peekToken.Type == lexer.RPAREN {
			p.nextToken() // cur = )
			break
		} else if p.curToken.Type == lexer.COMMA {
			p.nextToken()
		} else {
			break
		}
	}
	// cur = )
	res := &CallExpression{Function: name, Arguments: args}
	p.setPosAt(res, line, col)
	return res
}

// parseArgExpression parses one argument expression which may be a postfix
// (SOV) binary expression, e.g. "x y қосу" or just a literal/identifier.
func (p *Parser) parseArgExpression() Expression {
	var stack []Expression
	for p.curToken.Type != lexer.RPAREN &&
		p.curToken.Type != lexer.COMMA &&
		p.curToken.Type != lexer.RBRACKET &&
		p.curToken.Type != lexer.EOF {

		switch p.curToken.Type {
		case lexer.NUMBER:
			val, err := strconv.ParseFloat(p.curToken.Literal, 64)
			if err == nil {
				stack = append(stack, p.setPos(&NumberLiteral{Value: val}).(Expression))
			}
		case lexer.INT_LITERAL:
			val, err := strconv.ParseInt(p.curToken.Literal, 10, 64)
			if err == nil {
				stack = append(stack, p.setPos(&IntLiteral{Value: val}).(Expression))
			}
		case lexer.STRING:
			stack = append(stack, p.setPos(&StringLiteral{Value: p.curToken.Literal}).(Expression))
		case lexer.TRUE:
			stack = append(stack, p.setPos(&BoolLiteral{Value: true}).(Expression))
		case lexer.FALSE:
			stack = append(stack, p.setPos(&BoolLiteral{Value: false}).(Expression))
		case lexer.LBRACKET:
			arr := p.parseArrayLiteral()
			if arr != nil {
				stack = append(stack, arr)
			}
		case lexer.NEW:
			if len(stack) >= 1 {
				val := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				id, ok := val.(*Identifier)
				if ok {
					stack = append(stack, p.setPos(&StructCreateExpression{StructName: id.Value}).(Expression))
				}
			}
		case lexer.FILE_READ:
			if len(stack) >= 1 {
				path := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				stack = append(stack, p.setPos(&FileReadExpression{Path: path}).(Expression))
			}
		case lexer.INPUT:
			stack = append(stack, p.setPos(&InputExpression{}).(Expression))
		case lexer.IDENTIFIER:
			if p.peekToken.Type == lexer.LPAREN {
				callExpr := p.parseCallExpression(p.curToken.Literal, p.curToken.Line, p.curToken.Col)
				if callExpr != nil {
					stack = append(stack, callExpr)
					continue
				}
			} else {
				id := p.curToken.Literal
				if len(id) > 0 && id[0] == '.' {
					if len(stack) < 1 {
						p.errorf("өріске кіру үшін объект қажет")
						return nil
					}
					target := stack[len(stack)-1]
					stack = stack[:len(stack)-1]
					stack = append(stack, p.setPos(&StructFieldAccessExpression{
						StructName: "",
						Field:      id[1:],
						Target:     target,
					}).(Expression))
				} else {
					isStructField := false
					for i := 0; i < len(id); i++ {
						if id[i] == '.' {
							stack = append(stack, p.setPos(&StructFieldAccessExpression{
								StructName: id[:i],
								Field:      id[i+1:],
							}).(Expression))
							isStructField = true
							break
						}
					}
					if !isStructField {
						stack = append(stack, p.setPos(&Identifier{Value: id}).(Expression))
					}
				}
			}
		case lexer.PLUS, lexer.MINUS, lexer.MUL, lexer.DIV,
			lexer.GT, lexer.LT, lexer.EQ, lexer.NEQ, lexer.GTE, lexer.LTE,
			lexer.AND, lexer.OR, lexer.LSHIFT, lexer.RSHIFT:
			if len(stack) >= 2 {
				right := stack[len(stack)-1]
				left := stack[len(stack)-2]
				stack = stack[:len(stack)-2]
				stack = append(stack, p.setPos(&PostfixExpression{Left: left, Right: right, Operator: canonicalOp(p.curToken)}).(Expression))
			}
		case lexer.NOT:
			if len(stack) >= 1 {
				val := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				stack = append(stack, p.setPos(&UnaryExpression{Operator: "емес", Right: val}).(Expression))
			}
		case lexer.STR_CONCAT:
			if len(stack) >= 2 {
				right := stack[len(stack)-1]
				left := stack[len(stack)-2]
				stack = stack[:len(stack)-2]
				stack = append(stack, p.setPos(&StrConcatExpression{Left: left, Right: right}).(Expression))
			}
		case lexer.STR_LEN:
			if len(stack) >= 1 {
				val := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				stack = append(stack, p.setPos(&StrLenExpression{Value: val}).(Expression))
			}
		case lexer.STR_EQ:
			if len(stack) >= 2 {
				right := stack[len(stack)-1]
				left := stack[len(stack)-2]
				stack = stack[:len(stack)-2]
				stack = append(stack, p.setPos(&StrEqExpression{Left: left, Right: right}).(Expression))
			}
		case lexer.TO_STR:
			if len(stack) >= 1 {
				val := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				stack = append(stack, p.setPos(&ToStrExpression{Value: val}).(Expression))
			}
		case lexer.CHAR_CODE:
			if len(stack) >= 1 {
				val := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				stack = append(stack, p.setPos(&CharCodeExpression{Value: val}).(Expression))
			}
		case lexer.CHAR_AT:
			if len(stack) >= 2 {
				idx := stack[len(stack)-1]
				str := stack[len(stack)-2]
				stack = stack[:len(stack)-2]
				stack = append(stack, p.setPos(&CharAtExpression{Str: str, Index: idx}).(Expression))
			}
		case lexer.ARRAY_LEN:
			if len(stack) >= 1 {
				val := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				stack = append(stack, p.setPos(&LengthExpression{Value: val}).(Expression))
			}
		case lexer.INDEX_GET:
			if len(stack) >= 2 {
				idx := stack[len(stack)-1]
				arr := stack[len(stack)-2]
				stack = stack[:len(stack)-2]
				stack = append(stack, p.setPos(&IndexExpression{Left: arr, Index: idx}).(Expression))
			}
		}
		p.nextToken()
	}

	if len(stack) > 0 {
		return stack[len(stack)-1]
	}
	return nil
}

// parseSingleExpression parses one atomic expression (literal or identifier).
func (p *Parser) parseSingleExpression() Expression {
	switch p.curToken.Type {
	case lexer.NUMBER:
		val, err := strconv.ParseFloat(p.curToken.Literal, 64)
		if err != nil {
			p.errorf("жарамсыз сан: '%s'", p.curToken.Literal)
			return nil
		}
		return p.setPos(&NumberLiteral{Value: val}).(Expression)
	case lexer.INT_LITERAL:
		val, err := strconv.ParseInt(p.curToken.Literal, 10, 64)
		if err != nil {
			p.errorf("жарамсыз бүтін сан: '%s'", p.curToken.Literal)
			return nil
		}
		return p.setPos(&IntLiteral{Value: val}).(Expression)
	case lexer.STRING:
		return p.setPos(&StringLiteral{Value: p.curToken.Literal}).(Expression)
	case lexer.IDENTIFIER:
		id := p.curToken.Literal
		// Check if it is a struct field access (e.g. adam.аты)
		for i := 0; i < len(id); i++ {
			if id[i] == '.' {
				return p.setPos(&StructFieldAccessExpression{
					StructName: id[:i],
					Field:      id[i+1:],
				}).(Expression)
			}
		}
		return p.setPos(&Identifier{Value: id}).(Expression)
	case lexer.TRUE:
		return p.setPos(&BoolLiteral{Value: true}).(Expression)
	case lexer.FALSE:
		return p.setPos(&BoolLiteral{Value: false}).(Expression)
	}
	return nil
}

// ---------------------------------------------------------------------------
// Array literal: [ expr expr expr ]
// ---------------------------------------------------------------------------

func (p *Parser) parseArrayLiteral() *ArrayLiteral {
	line, col := p.curToken.Line, p.curToken.Col
	// cur = [
	p.nextToken() // cur = first element or ]
	var elements []Expression
	for p.curToken.Type != lexer.RBRACKET && p.curToken.Type != lexer.EOF {
		expr := p.parseArgExpression()
		if expr != nil {
			elements = append(elements, expr)
		}
		if p.curToken.Type == lexer.COMMA {
			p.nextToken()
		}
	}
	// cur = ]
	res := &ArrayLiteral{Elements: elements}
	p.setPosAt(res, line, col)
	return res
}

// ---------------------------------------------------------------------------
// Block: { stmt stmt ... }
// ---------------------------------------------------------------------------

func (p *Parser) parseBlockStatement() *BlockStatement {
	line, col := p.curToken.Line, p.curToken.Col
	block := &BlockStatement{Statements: []Statement{}}
	p.nextToken() // skip {

	for p.curToken.Type != lexer.RBRACE && p.curToken.Type != lexer.EOF {
		// Handle nested function definitions inside blocks
		var stmt Statement
		if p.curToken.Type == lexer.FUNC {
			stmt = p.parseFunctionStatement()
		} else {
			stmt = p.parseStatement()
		}
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	p.setPosAt(block, line, col)
	return block
}

func canonicalOp(tok lexer.Token) string {
	switch tok.Type {
	case lexer.PLUS:
		return "қосу"
	case lexer.MINUS:
		return "алу"
	case lexer.MUL:
		return "көбейту"
	case lexer.DIV:
		return "бөлу"
	case lexer.GT:
		return "үлкен"
	case lexer.LT:
		return "кіші"
	case lexer.EQ:
		return "тең"
	case lexer.NEQ:
		return "тең_емес"
	case lexer.GTE:
		return "үлкен_тең"
	case lexer.LTE:
		return "кіші_тең"
	case lexer.AND:
		return "және"
	case lexer.OR:
		return "немесе"
	case lexer.LSHIFT:
		return "жылжыту_сол"
	case lexer.RSHIFT:
		return "жылжыту_оң"
	default:
		return tok.Literal
	}
}

