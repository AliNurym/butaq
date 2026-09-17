package interpreter

import (
	"butaq/lexer"
	"butaq/locales"
	"butaq/parser"
	"testing"
)

func TestInterpreterLoopWithExpressionStatement(t *testing.T) {
	locales.Init()
	code := `
болсын i = 0
болсын count = 0
әзірше i < 5
    ұйықтау(0.001)
    count = count + 1
    i = i + 1
жазу(count)
`
	l := lexer.New(code)
	p := parser.New(l)
	prog := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("parse errors: %v", p.Errors())
	}

	ip := New()
	out, err := ip.Run(prog)
	if err != nil {
		t.Fatalf("runtime error: %v", err)
	}

	expected := "5\n"
	if out != expected {
		t.Fatalf("expected %q, got %q", expected, out)
	}
}
