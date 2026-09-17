package lexer

import (
	"butaq/locales"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type TokenType int

const (
	EOF TokenType = iota
	ILLEGAL

	// Data Types
	NUMBER      // float
	INT_LITERAL // int64
	STRING
	IDENTIFIER

	// Keywords — Butaq SOV syntax
	VAR      // болсын   (assign/declare)
	IF       // егер     (if)
	ELSE     // әйтпесе  (else)
	WHILE    // әзірше   (while)
	BREAK    // үзу
	CONTINUE // жалғастыру
	PRINT    // жазу     (print)

	// Math Operations
	PLUS  // қосу
	MINUS // алу
	MUL   // көбейту
	DIV   // бөлу

	// Comparisons
	GT  // үлкен    (>)
	LT  // кіші     (<)
	EQ  // тең      (==)
	NEQ // тең_емес (!=)
	GTE // үлкен_тең (>=)
	LTE // кіші_тең  (<=)

	// Logic
	AND    // және
	OR     // немесе
	NOT    // емес
	LSHIFT // жылжыту_сол
	RSHIFT // жылжыту_оң

	// Blocks
	LBRACE // {
	RBRACE // }
	LPAREN // (
	RPAREN // )
	COMMA  // ,
	COLON  // :
	ASSIGN // =
	THEN   // сонда, then, allora, тогда
	NEWLINE
	INDENT
	DEDENT

	// Functions
	FUNC   // функция
	RETURN // қайтару
	CALL   // шақыру
	IMPORT // енгізу

	// Booleans
	TRUE  // ақиқат
	FALSE // жалған

	// Arrays
	ARRAY     // тізім
	LBRACKET  // [
	RBRACKET  // ]
	INDEX_GET // тізім_алу (array index get)
	INDEX_SET // тізім_қой (array index set)
	ARRAY_LEN // ұзындық
	FREE      // бос (free heap memory)

	// Types
	TYPE_INT    // БҮТІН
	TYPE_FLOAT  // САН
	TYPE_STRING // МӘТІН
	TYPE_BOOL   // АҚИҚАТ
	TYPE_BYTE   // БАЙТ

	// Structs
	STRUCT // құрылым
	NEW    // жасау
	WEAK   // әлсіз
	THREAD // ағын

	// String / char ops
	CHAR_AT    // символ
	STR_CONCAT // біріктіру (string concat)
	STR_LEN    // ұзындық_жол (string length)
	STR_EQ     // мәтін_тең (string equals compare)
	TO_STR     // санды_мәтін (number to string)
	CHAR_CODE  // таңба_коды (byte value of first char)

	// File I/O
	FILE_READ  // файл_оқу
	FILE_WRITE // файл_жазу

	// User Input
	INPUT // кіру

	// Error Handling
	TRY_ERROR     // қатемен
	ERROR_LITERAL // қате

	// Interfaces
	INTERFACE // интерфейс
)

// Token name to TokenType mapping (canonical internal names)
var TokenNameToType = map[string]TokenType{
	"VAR":           VAR,
	"IF":            IF,
	"ELSE":          ELSE,
	"WHILE":         WHILE,
	"BREAK":         BREAK,
	"CONTINUE":      CONTINUE,
	"PRINT":         PRINT,
	"PLUS":          PLUS,
	"MINUS":         MINUS,
	"MUL":           MUL,
	"DIV":           DIV,
	"GT":            GT,
	"LT":            LT,
	"EQ":            EQ,
	"NEQ":           NEQ,
	"GTE":           GTE,
	"LTE":           LTE,
	"AND":           AND,
	"OR":            OR,
	"NOT":           NOT,
	"LSHIFT":        LSHIFT,
	"RSHIFT":        RSHIFT,
	"ASSIGN":        ASSIGN,
	"THEN":          THEN,
	"COLON":         COLON,
	"NEWLINE":       NEWLINE,
	"INDENT":        INDENT,
	"DEDENT":        DEDENT,
	"FUNC":          FUNC,
	"RETURN":        RETURN,
	"CALL":          CALL,
	"IMPORT":        IMPORT,
	"TRUE":          TRUE,
	"FALSE":         FALSE,
	"ARRAY":         ARRAY,
	"ARRAY_LEN":     ARRAY_LEN,
	"INDEX_GET":     INDEX_GET,
	"INDEX_SET":     INDEX_SET,
	"FREE":          FREE,
	"TYPE_INT":      TYPE_INT,
	"TYPE_FLOAT":    TYPE_FLOAT,
	"TYPE_STRING":   TYPE_STRING,
	"TYPE_BOOL":     TYPE_BOOL,
	"TYPE_BYTE":     TYPE_BYTE,
	"STRUCT":        STRUCT,
	"NEW":           NEW,
	"WEAK":          WEAK,
	"THREAD":        THREAD,
	"CHAR_AT":       CHAR_AT,
	"STR_CONCAT":    STR_CONCAT,
	"STR_LEN":       STR_LEN,
	"STR_EQ":        STR_EQ,
	"TO_STR":        TO_STR,
	"CHAR_CODE":     CHAR_CODE,
	"FILE_READ":     FILE_READ,
	"FILE_WRITE":    FILE_WRITE,
	"INPUT":         INPUT,
	"TRY_ERROR":     TRY_ERROR,
	"ERROR_LITERAL": ERROR_LITERAL,
	"INTERFACE":     INTERFACE,
}

var TypeToTokenName = map[TokenType]string{
	VAR:           "VAR",
	IF:            "IF",
	ELSE:          "ELSE",
	WHILE:         "WHILE",
	BREAK:         "BREAK",
	CONTINUE:      "CONTINUE",
	PRINT:         "PRINT",
	PLUS:          "PLUS",
	MINUS:         "MINUS",
	MUL:           "MUL",
	DIV:           "DIV",
	GT:            "GT",
	LT:            "LT",
	EQ:            "EQ",
	NEQ:           "NEQ",
	GTE:           "GTE",
	LTE:           "LTE",
	AND:           "AND",
	OR:            "OR",
	NOT:           "NOT",
	LSHIFT:        "LSHIFT",
	RSHIFT:        "RSHIFT",
	FUNC:          "FUNC",
	RETURN:        "RETURN",
	CALL:          "CALL",
	IMPORT:        "IMPORT",
	TRUE:          "TRUE",
	FALSE:         "FALSE",
	ARRAY:         "ARRAY",
	ARRAY_LEN:     "ARRAY_LEN",
	INDEX_GET:     "INDEX_GET",
	INDEX_SET:     "INDEX_SET",
	FREE:          "FREE",
	TYPE_INT:      "TYPE_INT",
	TYPE_FLOAT:    "TYPE_FLOAT",
	TYPE_STRING:   "TYPE_STRING",
	TYPE_BOOL:     "TYPE_BOOL",
	TYPE_BYTE:     "TYPE_BYTE",
	STRUCT:        "STRUCT",
	NEW:           "NEW",
	WEAK:          "WEAK",
	THREAD:        "THREAD",
	CHAR_AT:       "CHAR_AT",
	STR_CONCAT:    "STR_CONCAT",
	STR_LEN:       "STR_LEN",
	STR_EQ:        "STR_EQ",
	TO_STR:        "TO_STR",
	CHAR_CODE:     "CHAR_CODE",
	FILE_READ:     "FILE_READ",
	FILE_WRITE:    "FILE_WRITE",
	INPUT:         "INPUT",
	TRY_ERROR:     "TRY_ERROR",
	ERROR_LITERAL: "ERROR_LITERAL",
	INTERFACE:     "INTERFACE",
	ASSIGN:        "ASSIGN",
	THEN:          "THEN",
	COLON:         "COLON",
	NEWLINE:       "NEWLINE",
	INDENT:        "INDENT",
	DEDENT:        "DEDENT",
}

// StringToTokenType maps token name to TokenType
func StringToTokenType(name string) (TokenType, bool) {
	tt, ok := TokenNameToType[name]
	return tt, ok
}

// TokenTypeToString maps TokenType to token name
func TokenTypeToString(tt TokenType) string {
	if s, ok := TypeToTokenName[tt]; ok {
		return s
	}
	return "UNKNOWN"
}

// Backward-compatible keywords map
var keywords = map[string]TokenType{
	"болсын":        VAR,
	"егер":          IF,
	"әйтпесе":       ELSE,
	"әзірше":        WHILE,
	"үзу":           BREAK,
	"жалғастыру":    CONTINUE,
	"жазу":          PRINT,
	"қосу":          PLUS,
	"алу":           MINUS,
	"көбейту":       MUL,
	"бөлу":          DIV,
	"үлкен":         GT,
	"кіші":          LT,
	"тең":           EQ,
	"тең_емес":      NEQ,
	"үлкен_тең":     GTE,
	"кіші_тең":      LTE,
	"және":          AND,
	"немесе":        OR,
	"емес":          NOT,
	"жылжыту_сол":   LSHIFT,
	"жылжыту_оң":    RSHIFT,
	"функция":       FUNC,
	"қайтару":       RETURN,
	"шақыру":        CALL,
	"енгізу":        IMPORT,
	"ақиқат":        TRUE,
	"жалған":        FALSE,
	"тізім":         ARRAY,
	"ұзындық":       ARRAY_LEN,
	"тізім_алу":     INDEX_GET,
	"тізім_қой":     INDEX_SET,
	"бос":           FREE,
	"БҮТІН":         TYPE_INT,
	"САН":           TYPE_FLOAT,
	"МӘТІН":         TYPE_STRING,
	"АҚИҚАТ":        TYPE_BOOL,
	"БАЙТ":          TYPE_BYTE,
	"құрылым":       STRUCT,
	"жасау":         NEW,
	"әлсіз":         WEAK,
	"ағын":          THREAD,
	"символ":        CHAR_AT,
	"біріктіру":     STR_CONCAT,
	"ұзындық_жол":   STR_LEN,
	"мәтін_тең":     STR_EQ,
	"санды_мәтін":   TO_STR,
	"таңба_коды":    CHAR_CODE,
	"файл_оқу":      FILE_READ,
	"файл_жазу":     FILE_WRITE,
	"кіру":          INPUT,
	"қатемен":       TRY_ERROR,
	"қате":          ERROR_LITERAL,
	"интерфейс":     INTERFACE,
}

type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Col     int
}

type Lexer struct {
	input        string
	position     int
	readPosition int
	ch           rune
	line         int
	col          int
	locale       *locales.Locale
}

// New creates a new Lexer with auto-detected or default locale.
func New(input string) *Lexer {
	loc := locales.AutoDetect(input)
	return NewWithLocale(input, loc)
}

// NewWithLocale creates a new Lexer with a specific human language locale.
func NewWithLocale(input string, loc *locales.Locale) *Lexer {
	if loc == nil {
		loc = locales.Default()
	}
	l := &Lexer{input: input, line: 1, col: 0, locale: loc}
	l.readChar()
	return l
}

// Locale returns the active locale of the lexer.
func (l *Lexer) Locale() *locales.Locale {
	return l.locale
}

// SetLocale changes the active locale of the lexer.
func (l *Lexer) SetLocale(loc *locales.Locale) {
	if loc != nil {
		l.locale = loc
	}
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		r, size := utf8.DecodeRuneInString(l.input[l.readPosition:])
		l.ch = r
		l.position = l.readPosition
		l.readPosition += size
		l.col++
	}
}

func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.readPosition:])
	return r
}

func (l *Lexer) lookupKeyword(ident string) (TokenType, bool) {
	// 1. Search in active locale
	if l.locale != nil {
		if tokName, ok := l.locale.Keywords[ident]; ok {
			if tt, ok := TokenNameToType[tokName]; ok {
				return tt, true
			}
		}
		return IDENTIFIER, false
	}

	// 2. If no active locale, search across all registered locales
	for _, loc := range locales.Available() {
		if tokName, ok := loc.Keywords[ident]; ok {
			if tt, ok := TokenNameToType[tokName]; ok {
				return tt, true
			}
		}
	}

	// 3. Fallback to hardcoded keywords map
	if tt, ok := keywords[ident]; ok {
		return tt, true
	}

	return IDENTIFIER, false
}

func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case '{':
		tok = l.newToken(LBRACE, string(l.ch))
	case '}':
		tok = l.newToken(RBRACE, string(l.ch))
	case '(':
		tok = l.newToken(LPAREN, string(l.ch))
	case ')':
		tok = l.newToken(RPAREN, string(l.ch))
	case '[':
		tok = l.newToken(LBRACKET, string(l.ch))
	case ']':
		tok = l.newToken(RBRACKET, string(l.ch))
	case ',':
		tok = l.newToken(COMMA, string(l.ch))
	case '+':
		tok = l.newToken(PLUS, "+")
	case '-':
		tok = l.newToken(MINUS, "-")
	case '*':
		tok = l.newToken(MUL, "*")
	case '/':
		tok = l.newToken(DIV, "/")
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(LTE, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(LT, string(l.ch))
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(GTE, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(GT, string(l.ch))
		}
	case ':':
		tok = l.newToken(COLON, string(l.ch))
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(EQ, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(ASSIGN, string(l.ch))
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(NEQ, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(ILLEGAL, string(l.ch))
		}
	case '"':
		strVal, ok := l.readString()
		if !ok {
			tok = l.newToken(ILLEGAL, "Unterminated string")
		} else {
			tok = l.newToken(STRING, strVal)
		}
		l.readChar()
		return tok
	case '#':
		l.skipComment()
		return l.NextToken()
	case 0:
		tok.Type = EOF
		tok.Literal = ""
		tok.Line = l.line
		tok.Col = l.col
		return tok
	default:
		if isDigit(l.ch) {
			tok.Line = l.line
			tok.Col = l.col
			num, ok, isFloat := l.readNumber()
			if !ok {
				tok.Type = ILLEGAL
				tok.Literal = num
				return tok
			}
			if isFloat {
				tok.Type = NUMBER
			} else {
				tok.Type = INT_LITERAL
			}
			tok.Literal = num
			return tok
		} else if isLetter(l.ch) || l.ch == '.' {
			tok.Line = l.line
			tok.Col = l.col
			tok.Literal = l.readIdentifier()
			if kw, ok := l.lookupKeyword(tok.Literal); ok {
				tok.Type = kw
			} else {
				tok.Type = IDENTIFIER
			}
			return tok
		} else {
			tok = l.newToken(ILLEGAL, string(l.ch))
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		if l.ch == '\n' {
			l.line++
			l.col = 0
		}
		l.readChar()
	}
}

func (l *Lexer) skipComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	startPos := l.position
	for isLetter(l.ch) || isDigit(l.ch) || l.ch == '_' || l.ch == '.' {
		l.readChar()
	}
	return l.input[startPos:l.position]
}

// readNumber returns the number literal, a boolean indicating validity, and a boolean indicating if it's a float.
func (l *Lexer) readNumber() (string, bool, bool) {
	startPos := l.position
	dotCount := 0
	for isDigit(l.ch) || l.ch == '.' {
		if l.ch == '.' {
			dotCount++
			if dotCount > 1 {
				// consume rest and return invalid
				for isDigit(l.ch) || l.ch == '.' {
					l.readChar()
				}
				return l.input[startPos:l.position], false, false
			}
		}
		l.readChar()
	}
	return l.input[startPos:l.position], true, dotCount == 1
}

// readString returns the string contents (without quotes) and a bool indicating if it was properly closed.
func (l *Lexer) readString() (string, bool) {
	var sb strings.Builder
	for {
		l.readChar()
		if l.ch == '"' {
			break
		}
		if l.ch == 0 {
			// EOF without closing quote
			return "", false
		}
		if l.ch == '\\' {
			l.readChar()
			switch l.ch {
			case 'n':
				sb.WriteByte('\n')
			case 't':
				sb.WriteByte('\t')
			case 'r':
				sb.WriteByte('\r')
			case '"':
				sb.WriteByte('"')
			case '\\':
				sb.WriteByte('\\')
			case 'e':
				sb.WriteByte(0x1b)
			case 'x':
				// Hex escape e.g. \x1b
				h1 := l.peekChar()
				l.readChar()
				h2 := l.peekChar()
				l.readChar()
				var b byte
				fmt.Sscanf(string([]rune{h1, h2}), "%02x", &b)
				sb.WriteByte(b)
			case '0':
				// Octal e.g. \033
				o1 := l.peekChar()
				l.readChar()
				o2 := l.peekChar()
				l.readChar()
				var b byte
				fmt.Sscanf(string([]rune{'0', o1, o2}), "%03o", &b)
				sb.WriteByte(b)
			default:
				sb.WriteRune('\\')
				sb.WriteRune(l.ch)
			}
		} else {
			sb.WriteRune(l.ch)
		}
	}
	return sb.String(), true
}

func isLetter(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}

func (l *Lexer) newToken(tokenType TokenType, ch string) Token {
	return Token{Type: tokenType, Literal: ch, Line: l.line, Col: l.col}
}

func (t Token) String() string {
	return fmt.Sprintf("Token(%d, %q, line %d, col %d)", t.Type, t.Literal, t.Line, t.Col)
}
