package transpiler

import (
	"butaq/interpreter"
	"butaq/lexer"
	"butaq/locales"
	"butaq/parser"
	"strings"
	"testing"
)

const fibKK = `функция фиб(n)
    егер n <= 1 сонда қайтару n
    қайтару фиб(n - 1) + фиб(n - 2)

жазу("fib(5) = ", фиб(5))
`

func runCode(t *testing.T, src string, loc *locales.Locale) string {
	l := lexer.NewWithLocale(src, loc)
	p := parser.New(l)
	prog := p.ParseProgram()
	if len(p.Errors()) > 0 {
		t.Fatalf("Parse errors in %s: %s\nCode:\n%s", loc.Code, strings.Join(p.Errors(), ", "), src)
	}

	ip := interpreter.New()
	out, err := ip.Run(prog)
	if err != nil {
		t.Fatalf("Runtime error in %s: %v", loc.Code, err)
	}
	return out
}

func TestModernSyntaxAndTranspilation(t *testing.T) {
	locKK, ok := locales.Get("kk")
	if !ok {
		t.Fatal("locales kk not found")
	}

	outKK := runCode(t, fibKK, locKK)
	if !strings.Contains(outKK, "fib(5) = 5") {
		t.Fatalf("Expected output containing 'fib(5) = 5', got: %s", outKK)
	}

	// Transpile to EN, IT, RU
	targetLangs := []string{"en", "it", "ru"}
	for _, lang := range targetLangs {
		targetLoc, ok := locales.Get(lang)
		if !ok {
			t.Fatalf("Locale %s not found", lang)
		}

		// Parse source
		l := lexer.NewWithLocale(fibKK, locKK)
		p := parser.New(l)
		prog := p.ParseProgram()
		if len(p.Errors()) > 0 {
			t.Fatalf("Parse errors: %v", p.Errors())
		}

		tr := New(targetLoc)
		transpiled := tr.Transpile(prog)

		// Run transpiled code
		out := runCode(t, transpiled, targetLoc)
		if !strings.Contains(out, "fib(5) = 5") {
			t.Fatalf("Lang %s: expected 'fib(5) = 5', got: %s\nTranspiled code:\n%s", lang, out, transpiled)
		}
	}
}
