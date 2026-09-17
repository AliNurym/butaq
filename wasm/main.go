package main

import (
	"butaq/codegen"
	"butaq/interpreter"
	"butaq/lexer"
	"butaq/locales"
	"butaq/parser"
	"butaq/transpiler"
	"butaq/typechecker"
	"strings"
	"syscall/js"
)

func parseAndTypecheck(source string, localeCode string) (*parser.Program, *typechecker.TypeEnv, *typechecker.TypeChecker, *locales.Locale, map[string]interface{}) {
	var loc *locales.Locale
	if localeCode != "" {
		if l, ok := locales.Get(localeCode); ok {
			loc = l
		}
	}
	if loc == nil {
		loc = locales.AutoDetect(source)
	}

	l := lexer.NewWithLocale(source, loc)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		return nil, nil, nil, loc, map[string]interface{}{
			"error": "Синтаксистік қателер:\n" + strings.Join(p.Errors(), "\n"),
		}
	}

	tcEnv := typechecker.NewTypeEnv()
	tc := typechecker.New()
	tc.Check(prog, tcEnv)

	if len(tc.Errors) > 0 {
		return nil, nil, nil, loc, map[string]interface{}{
			"error": "Тип қателері:\n" + strings.Join(tc.ErrorStrings(), "\n"),
		}
	}

	return prog, tcEnv, tc, loc, nil
}

func runButaqCode(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return map[string]interface{}{"error": "кіріс файлы бос"}
	}
	source := args[0].String()
	localeCode := ""
	if len(args) > 1 {
		localeCode = args[1].String()
	}

	prog, _, _, loc, errMap := parseAndTypecheck(source, localeCode)
	if errMap != nil {
		return errMap
	}

	ip := interpreter.New()
	out, err := ip.Run(prog)
	if err != nil {
		return map[string]interface{}{
			"error": "Орындалу қатесі: " + err.Error(),
		}
	}

	return map[string]interface{}{
		"output": out,
		"locale": loc.Code,
	}
}

func transpileButaqCode(this js.Value, args []js.Value) interface{} {
	if len(args) < 2 {
		return map[string]interface{}{"error": "транспиляция үшін бастапқы код және мақсатты тіл қажет"}
	}
	source := args[0].String()
	targetLang := args[1].String()
	sourceLang := ""
	if len(args) > 2 {
		sourceLang = args[2].String()
	}

	targetLoc, ok := locales.Get(targetLang)
	if !ok {
		return map[string]interface{}{"error": "белгісіз мақсатты тіл: " + targetLang}
	}

	var srcLoc *locales.Locale
	if sourceLang != "" {
		if l, ok := locales.Get(sourceLang); ok {
			srcLoc = l
		}
	}
	if srcLoc == nil {
		srcLoc = locales.AutoDetect(source)
	}

	l := lexer.NewWithLocale(source, srcLoc)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.Errors()) > 0 {
		return map[string]interface{}{
			"error": "Синтаксистік қателер:\n" + strings.Join(p.Errors(), "\n"),
		}
	}

	tr := transpiler.New(targetLoc)
	transpiledCode := tr.Transpile(prog)

	return map[string]interface{}{
		"code":       transpiledCode,
		"fromLocale": srcLoc.Code,
		"toLocale":   targetLoc.Code,
		"localeName": targetLoc.Name,
	}
}

func compileToNasm(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return map[string]interface{}{"error": "кіріс файлы бос"}
	}
	source := args[0].String()
	localeCode := ""
	if len(args) > 1 {
		localeCode = args[1].String()
	}

	prog, tcEnv, tc, _, errMap := parseAndTypecheck(source, localeCode)
	if errMap != nil {
		return errMap
	}

	cg := codegen.NewWithTC(tcEnv, tc, codegen.PlatformWindows)
	code := cg.Generate(prog)
	return map[string]interface{}{
		"code": code,
	}
}

func compileToLlvm(this js.Value, args []js.Value) interface{} {
	if len(args) < 1 {
		return map[string]interface{}{"error": "кіріс файлы бос"}
	}
	source := args[0].String()
	localeCode := ""
	if len(args) > 1 {
		localeCode = args[1].String()
	}

	prog, tcEnv, tc, _, errMap := parseAndTypecheck(source, localeCode)
	if errMap != nil {
		return errMap
	}

	lg := codegen.NewLlvm(tcEnv, tc, codegen.PlatformWindows)
	code := lg.Generate(prog)
	return map[string]interface{}{
		"code": code,
	}
}

func getAvailableLocales(this js.Value, args []js.Value) interface{} {
	locs := locales.Available()
	res := make([]interface{}, len(locs))
	for i, l := range locs {
		res[i] = map[string]interface{}{
			"code": l.Code,
			"name": l.Name,
		}
	}
	return res
}

func main() {
	c := make(chan struct{}, 0)
	js.Global().Set("runButaqCode", js.FuncOf(runButaqCode))
	js.Global().Set("transpileButaqCode", js.FuncOf(transpileButaqCode))
	js.Global().Set("compileToNasm", js.FuncOf(compileToNasm))
	js.Global().Set("compileToLlvm", js.FuncOf(compileToLlvm))
	js.Global().Set("getAvailableLocales", js.FuncOf(getAvailableLocales))
	<-c
}
