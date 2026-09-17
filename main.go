package main

import (
	"butaq/codegen"
	"butaq/interpreter"
	"butaq/lexer"
	"butaq/locales"
	"butaq/parser"
	"butaq/transpiler"
	"butaq/typechecker"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"unsafe"
)

// enableWindowsANSI enables VT100/ANSI escape codes in Windows console.
// Required for animations, colors, and cursor movement (\x1b[H, \x1b[2J, etc.)
func enableWindowsANSI() {
	if runtime.GOOS != "windows" {
		return
	}
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	getConsoleMode := kernel32.NewProc("GetConsoleMode")
	setConsoleMode := kernel32.NewProc("SetConsoleMode")
	handle, err := syscall.GetStdHandle(syscall.STD_OUTPUT_HANDLE)
	if err != nil {
		return
	}
	var mode uint32
	getConsoleMode.Call(uintptr(handle), uintptr(unsafe.Pointer(&mode)))
	// ENABLE_VIRTUAL_TERMINAL_PROCESSING = 0x0004
	setConsoleMode.Call(uintptr(handle), uintptr(mode|0x0004))
}

// initToolchainPath auto-discovers MinGW/LLVM/NASM compiler toolchains on Windows.
func initToolchainPath() {
	extraPaths := []string{
		"C:\\Users\\megumin\\AppData\\Local\\Microsoft\\WinGet\\Packages\\BrechtSanders.WinLibs.POSIX.UCRT_Microsoft.Winget.Source_8wekyb3d8bbwe\\mingw64\\bin",
		"C:\\Program Files\\LLVM\\bin",
		"C:\\Program Files\\NASM",
		"C:\\MinGW\\bin",
	}
	currentPath := os.Getenv("PATH")
	for _, p := range extraPaths {
		if _, err := os.Stat(p); err == nil && !strings.Contains(currentPath, p) {
			currentPath = p + ";" + currentPath
		}
	}
	os.Setenv("PATH", currentPath)
}

func main() {
	initToolchainPath()
	enableWindowsANSI()
	if len(os.Args) < 2 {
		printGeneralUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "build":
		handleBuild(os.Args[2:])
	case "run":
		handleRun(os.Args[2:])
	case "transpile":
		handleTranspile(os.Args[2:])
	case "locales":
		handleLocales(os.Args[2:])
	case "fmt":
		handleFmt(os.Args[2:])
	case "lsp":
		handleLsp(os.Args[2:])
	case "bap":
		handleBap(os.Args[2:])
	case "-h", "--help", "help":
		printGeneralUsage()
	default:
		// Backward compatibility: if first arg is a file or flag (e.g. -r, -v, or file.btq)
		handleBuild(os.Args[1:])
	}
}

func printGeneralUsage() {
	fmt.Println("═══════════════════════════════════════════════════════════════════════")
	fmt.Println("           Butaq — Zero-Cost Localized System-Level Compiler")
	fmt.Println("       Модульный полиглот-компилятор с нулевой стоимостью локализации")
	fmt.Println("═══════════════════════════════════════════════════════════════════════")
	fmt.Println("Қолдану / Usage: butaq <команда / command> [параметрлер] <файл.btq>")
	fmt.Println()
	fmt.Println("Командалар / Commands:")
	fmt.Println("  build       Компиляция в нативный бинарник (LLVM / NASM)")
	fmt.Println("  run         Компиляция и немедленный запуск программы")
	fmt.Println("  transpile   Транспиляция AST между языками (KZ, EN, IT, RU)")
	fmt.Println("  locales     Список поддерживаемых языковых локалей")
	fmt.Println("  fmt         Форматирование исходного кода (.btq)")
	fmt.Println("  lsp         Language Server Protocol для VS Code / Neovim")
	fmt.Println("  bap         Пакетный менеджер Butaq")
	fmt.Println()
	fmt.Println("Мысалдар / Examples:")
	fmt.Println("  butaq build examples/fib_kk.btq")
	fmt.Println("  butaq run examples/fib_en.btq")
	fmt.Println("  butaq transpile examples/fib_kk.btq --target=it -o examples/fib_it.btq")
	fmt.Println("  butaq transpile examples/fib_en.btq --target=ru")
	fmt.Println("  butaq locales")
	fmt.Println()
}

func handleLocales(args []string) {
	fmt.Println("═══════════════════════════════════════════════════════════════════════")
	fmt.Println("               Butaq — Қолдау көрсетілетін локальдер / Locales")
	fmt.Println("═══════════════════════════════════════════════════════════════════════")
	available := locales.Available()
	for _, loc := range available {
		fmt.Printf("  • [%-2s] %-12s — %d сөздік кілт сөздері, %d кірістірілген атаулар\n",
			loc.Code, loc.Name, len(loc.Keywords), len(loc.Builtins))
	}
	fmt.Println()
	fmt.Println("💡 Жаңа тілді қосу үшін locales/<код>.json файлын қосыңыз (қайта жинақтау қажетсіз!).")
	fmt.Println()
}

func reorderFlags(args []string) []string {
	var flags []string
	var nonFlags []string
	skipNext := false

	for i := 0; i < len(args); i++ {
		arg := args[i]
		if skipNext {
			flags = append(flags, arg)
			skipNext = false
			continue
		}
		if strings.HasPrefix(arg, "-") {
			flags = append(flags, arg)
			if !strings.Contains(arg, "=") && (arg == "-o" || arg == "--target" || arg == "-target" || arg == "--to" || arg == "-to" || arg == "--locale" || arg == "-locale" || arg == "--platform" || arg == "-platform" || arg == "--asm" || arg == "-asm") {
				skipNext = true
			}
		} else {
			nonFlags = append(nonFlags, arg)
		}
	}
	return append(flags, nonFlags...)
}

func handleTranspile(args []string) {
	args = reorderFlags(args)
	fs := flag.NewFlagSet("transpile", flag.ExitOnError)
	var (
		targetFlag string
		toFlag     string
		outputFlag string
		localeFlag string
		stdoutFlag bool
	)
	fs.StringVar(&targetFlag, "target", "en", "Мақсатты тіл / Target human language (en, kk, it, ru)")
	fs.StringVar(&toFlag, "to", "", "Мақсатты тіл (--target баламасы)")
	fs.StringVar(&outputFlag, "o", "", "Шығыс файл / Output file")
	fs.StringVar(&localeFlag, "locale", "", "Бастапқы тіл / Source locale (авто-анықтауды ауыстыру)")
	fs.BoolVar(&stdoutFlag, "stdout", false, "Тікелей терминалға шығару")

	fs.Parse(args)

	if fs.NArg() < 1 {
		fmt.Println("Қолдану: butaq transpile [жалаушалар] <файл.btq>")
		fs.PrintDefaults()
		return
	}

	inputFile := fs.Arg(0)
	sourceCode, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("❌ Қате: файлды оқу мүмкін болмады: %v\n", err)
		os.Exit(1)
	}

	// 1. Resolve source locale
	var srcLoc *locales.Locale
	if localeFlag != "" {
		if loc, ok := locales.Get(localeFlag); ok {
			srcLoc = loc
		} else {
			fmt.Printf("⚠️ Белгісіз бастапқы локаль '%s', авто-анықтау қолданылады.\n", localeFlag)
		}
	}
	if srcLoc == nil {
		srcLoc = locales.AutoDetect(string(sourceCode))
	}

	// 2. Resolve target locale
	targetCode := targetFlag
	if toFlag != "" {
		targetCode = toFlag
	}
	targetLoc, ok := locales.Get(targetCode)
	if !ok {
		fmt.Printf("❌ Қате: белгісіз мақсатты локаль '%s'. Қолжетімді: kk, en, it, ru\n", targetCode)
		os.Exit(1)
	}

	// 3. Parse AST
	l := lexer.NewWithLocale(string(sourceCode), srcLoc)
	p := parser.New(l)
	prog := p.ParseProgram()

	if len(p.ErrorsStructured()) > 0 {
		fmt.Printf("❌ %s:\n", srcLoc.Message("syntax_errors", "Синтаксистік қателер"))
		for _, e := range p.ErrorsStructured() {
			renderError(inputFile, string(sourceCode), e.Line, e.Col, e.Message, srcLoc.Message("syntax", "Синтаксис"), srcLoc)
		}
		os.Exit(1)
	}

	// 4. Transpile
	tr := transpiler.New(targetLoc)
	transpiled := tr.Transpile(prog)

	// 5. Output
	if outputFlag != "" {
		if err := os.WriteFile(outputFlag, []byte(transpiled), 0644); err != nil {
			fmt.Printf("❌ Қате: шығыс файлын жазу мүмкін болмады: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("✨ Транспиляция сәтті аяқталды [%s → %s]: %s\n", srcLoc.Code, targetLoc.Code, outputFlag)
	} else {
		fmt.Print(transpiled)
	}
}

func handleRun(args []string) {
	// Prepend -r to args and run build
	buildArgs := append([]string{"-r"}, args...)
	handleBuild(buildArgs)
}

func handleBuild(args []string) {
	args = reorderFlags(args)
	fs := flag.NewFlagSet("build", flag.ExitOnError)

	var (
		helpFlag     bool
		runFlag      bool
		outputFlag   string
		platformFlag string
		asmFileFlag  string
		verboseFlag  bool
		llvmFlag     bool
		nasmFlag     bool
		interpFlag   bool
		localeFlag   string
	)

	fs.BoolVar(&helpFlag, "h", false, "Көмек көрсету")
	fs.BoolVar(&helpFlag, "help", false, "Көмек көрсету")
	fs.BoolVar(&runFlag, "r", false, "Бағдарламаны компиляциядан кейін бірден іске қосу")
	fs.BoolVar(&runFlag, "run", false, "Бағдарламаны компиляциядан кейін бірден іске қосу")
	fs.BoolVar(&interpFlag, "i", false, "AST-интерпретатор режимінде орындау")
	fs.BoolVar(&interpFlag, "interp", false, "AST-интерпретатор режимінде орындау")
	fs.StringVar(&outputFlag, "o", "", "Шығыс екілік (binary) файлдың атауы")
	fs.StringVar(&platformFlag, "platform", "", "Мақсатты платформа (windows немесе linux)")
	fs.StringVar(&asmFileFlag, "asm", "", "Генерацияланатын ассемблер немесе LLVM файлының атауы")
	fs.BoolVar(&verboseFlag, "v", false, "Толық компиляция журналдарын көрсету")
	fs.BoolVar(&verboseFlag, "verbose", false, "Толық компиляция журналдарын көрсету")
	fs.BoolVar(&llvmFlag, "llvm", false, "LLVM IR генерациясын және clang компиляциясын қолдану")
	fs.BoolVar(&nasmFlag, "nasm", false, "NASM x86-64 backend қолдану")
	fs.StringVar(&localeFlag, "locale", "", "Бастапқы синтаксис локалі (kk, en, it, ru)")

	fs.Parse(args)

	if helpFlag || fs.NArg() < 1 {
		fmt.Println("═══════════════════════════════════════════════════════════")
		fmt.Println("             Butaq — Polyglot Compiler (CLI)")
		fmt.Println("═══════════════════════════════════════════════════════════")
		fmt.Println("Қолдану: butaq build [жалаушалар] <файл.btq>")
		fmt.Println()
		fs.PrintDefaults()
		return
	}

	inputFile := fs.Arg(0)
	sourceCode, err := os.ReadFile(inputFile)
	if err != nil {
		fmt.Printf("❌ Қате: файлды оқу мүмкін болмады: %v\n", err)
		os.Exit(1)
	}

	if verboseFlag {
		fmt.Printf("🔤 Оқылды: %s (%d байт)\n", inputFile, len(sourceCode))
	}

	// ── 0. Локальді анықтау ───────────────────────────────────────────────
	var activeLocale *locales.Locale
	if localeFlag != "" {
		if loc, ok := locales.Get(localeFlag); ok {
			activeLocale = loc
		}
	}
	if activeLocale == nil {
		activeLocale = locales.AutoDetect(string(sourceCode))
	}
	if verboseFlag {
		fmt.Printf("🌐 Анықталған локаль: [%s] %s\n", activeLocale.Code, activeLocale.Name)
	}

	// ── 1. Лексер + Парсер ──────────────────────────────────────────────────
	l := lexer.NewWithLocale(string(sourceCode), activeLocale)
	p := parser.New(l)
	p.SetLocale(activeLocale)
	program := p.ParseProgram()

	if len(p.ErrorsStructured()) != 0 {
		fmt.Printf("❌ %s:\n", activeLocale.Message("syntax_errors", "Синтаксистік қателер"))
		for _, e := range p.ErrorsStructured() {
			renderError(inputFile, string(sourceCode), e.Line, e.Col, e.Message, activeLocale.Message("syntax", "Синтаксис"), activeLocale)
		}
		os.Exit(1)
	}
	if verboseFlag {
		fmt.Printf("✅ Лексер + Парсер: %d оператор талданды\n", len(program.Statements))
	}

	// Resolve imports
	currentDir := filepath.Dir(inputFile)
	if err := parser.ResolveImports(program, currentDir); err != nil {
		fmt.Printf("❌ Енгізу қатесі: %v\n", err)
		os.Exit(1)
	}

	// ── 2. Статикалық типтер тексеру ───────────────────────────────────────
	tcEnv := typechecker.NewTypeEnv()
	tc := typechecker.New()
	tc.SetLocale(activeLocale)
	tc.Check(program, tcEnv)

	if len(tc.Errors) != 0 {
		fmt.Printf("❌ %s:\n", activeLocale.Message("type_errors", "Тип қателері"))
		for _, e := range tc.Errors {
			renderError(inputFile, string(sourceCode), e.Line, e.Col, e.Message, activeLocale.Message("type", "Тип"), activeLocale)
		}
		os.Exit(1)
	}
	if verboseFlag {
		fmt.Println("✅ Тип тексеру: сәтті аяқталды (қателер жоқ)")
	}

	// ── Linter ─────────────────────────────────────────────────────────────
	linter := typechecker.NewLinter()
	warnings := linter.Lint(program)
	if len(warnings) > 0 {
		fmt.Println("⚠️  Ескертулер:")
		for _, w := range warnings {
			renderWarning(inputFile, string(sourceCode), w.Line, w.Col, w.Message)
		}
	}

	// Default to high-performance NASM x86-64 native backend
	useLlvm := false
	if llvmFlag {
		useLlvm = true
	} else if nasmFlag {
		useLlvm = false
	}

	// ── 3. Кодогенерация ──────────────────────────────
	platform := codegen.PlatformLinux
	if platformFlag != "" {
		switch strings.ToLower(platformFlag) {
		case "windows", "win":
			platform = codegen.PlatformWindows
		case "linux":
			platform = codegen.PlatformLinux
		default:
			if runtime.GOOS == "windows" {
				platform = codegen.PlatformWindows
			}
		}
	} else {
		if runtime.GOOS == "windows" {
			platform = codegen.PlatformWindows
		}
	}

	asmFile := asmFileFlag
	if asmFile == "" {
		if useLlvm {
			asmFile = "out.ll"
		} else {
			asmFile = "out.asm"
		}
	}

	var generatedCode string
	if useLlvm {
		lg := codegen.NewLlvm(tcEnv, tc, platform)
		generatedCode = lg.Generate(program)
		if verboseFlag {
			fmt.Println("✅ LLVM IR коды сәтті генерацияланды")
		}
	} else {
		cg := codegen.NewWithTC(tcEnv, tc, platform)
		generatedCode = cg.Generate(program)
		if verboseFlag {
			fmt.Println("✅ NASM ассемблер коды сәтті генерацияланды")
		}
	}

	// ── 4. Сақтау ─────────────────────────────────────────────
	if err := os.WriteFile(asmFile, []byte(generatedCode), 0644); err != nil {
		fmt.Printf("❌ Қате: шығыс файлын жазу мүмкін болмады: %v\n", err)
		os.Exit(1)
	}

	baseName := strings.TrimSuffix(filepath.Base(inputFile), ".btq")
	outputBinary := outputFlag
	if outputBinary == "" {
		outputBinary = baseName
		if platform == codegen.PlatformWindows {
			outputBinary += ".exe"
		}
	}

	runtimeC := "runtime/runtime.c"
	if exePath, err := os.Executable(); err == nil {
		candidate := filepath.Join(filepath.Dir(exePath), "runtime", "runtime.c")
		if _, err := os.Stat(candidate); err == nil {
			runtimeC = candidate
		}
	}

	hasCompiler := false
	if useLlvm {
		_, err := exec.LookPath("clang")
		hasCompiler = (err == nil)
	} else {
		_, errNasm := exec.LookPath("nasm")
		_, errGcc := exec.LookPath("gcc")
		hasCompiler = (errNasm == nil && errGcc == nil)
	}

	// If interpFlag OR no C compiler/assembler found and running directly, fall back to pure AST Interpreter!
	if (interpFlag || !hasCompiler) && runFlag {
		if verboseFlag {
			if interpFlag {
				fmt.Println("ℹ️ AST-интерпретатор режимінде орындалуда...")
			} else {
				fmt.Println("ℹ️ Нативный компилятор не найден, запуск через встроенный AST-интерпретатор...")
			}
		}
		ip := interpreter.New()
		ip.StreamOutput = true // write directly to stdout for real-time animation
		_, err := ip.Run(program)
		if err != nil {
			fmt.Printf("❌ Орындалу қатесі: %v\n", err)
			os.Exit(1)
		}
		if asmFileFlag == "" {
			os.Remove(asmFile)
		}
		return
	}

	if useLlvm {
		// ── 5. LLVM IR компиляциясы (Clang) ──────────────
		if verboseFlag {
			fmt.Printf("🔧 Clang компиляциясы: %s + %s → %s\n", asmFile, runtimeC, outputBinary)
		}
		clangArgs := []string{"-O3", "-o", outputBinary, asmFile, runtimeC}
		if platform == codegen.PlatformLinux {
			clangArgs = append(clangArgs, "-lpthread", "-lm")
		}
		clangCmd := exec.Command("clang", clangArgs...)
		clangOut, err := clangCmd.CombinedOutput()
		if err != nil {
			fmt.Println("⚠️ Clang компиляциясы сәтсіз аяқталды немесе 'clang' табылмады.")
			fmt.Println("LLVM IR коды келесі файлға сақталды:", asmFile)
			if verboseFlag {
				fmt.Println(string(clangOut))
			}
			if runFlag {
				// Fallback to interpreter
				ip := interpreter.New()
				out, err := ip.Run(program)
				if err == nil {
					fmt.Print(out)
					return
				}
				os.Exit(1)
			}
		} else {
			if verboseFlag {
				fmt.Println("✅ Clang: сәтті орындалды")
			}
			os.Chmod(outputBinary, 0755)
			if asmFileFlag == "" {
				os.Remove("out.ll")
			}
		}
	} else {
		// ── 5. NASM арқылы объектілік файл жасау ────────────────────────────────
		objFile := baseName + ".o"
		if verboseFlag {
			fmt.Printf("🔧 NASM компиляциясы: %s → %s\n", filepath.Base(asmFile), filepath.Base(objFile))
		}

		var nasmArgs []string
		nasmExe := "nasm"
		if runtime.GOOS == "windows" {
			nasmArgs = []string{"-f", "win64", "-Ox", "-o", objFile, asmFile}
			if _, err := exec.LookPath("nasm"); err != nil {
				if _, err := os.Stat("C:\\Program Files\\NASM\\nasm.exe"); err == nil {
					nasmExe = "C:\\Program Files\\NASM\\nasm.exe"
				}
			}
		} else {
			nasmArgs = []string{"-f", "elf64", "-Ox", "-o", objFile, asmFile}
		}

		nasmOut, err := exec.Command(nasmExe, nasmArgs...).CombinedOutput()
		if err != nil {
			fmt.Println("❌ NASM қатесі:")
			fmt.Println(string(nasmOut))
			fmt.Println("\n--- Генерацияланған ASM коды ---")
			printNumbered(generatedCode)
			os.Exit(1)
		}
		if verboseFlag {
			fmt.Println("✅ NASM: объект файл жасалды")
		}

		// ── 6. Линковка ─────────────────────────────────────────────────────────
		var linkCmd *exec.Cmd
		if platform == codegen.PlatformWindows {
			linkCmd = exec.Command("gcc", "-O3", "-o", outputBinary, objFile, runtimeC, "-lkernel32", "-lmsvcrt")
		} else {
			linkCmd = exec.Command("gcc", "-O3", "-o", outputBinary, objFile, runtimeC, "-no-pie", "-lpthread", "-lc")
		}

		linkOut, err := linkCmd.CombinedOutput()
		if err != nil {
			fmt.Println("❌ Линковка қатесі:")
			fmt.Println(string(linkOut))
			os.Exit(1)
		}

		os.Chmod(outputBinary, 0755)
		os.Remove(objFile)
		if asmFileFlag == "" {
			os.Remove("out.asm")
		}
	}

	if !runFlag {
		fmt.Println()
		fmt.Println("═══════════════════════════════════════════════════════════════════════")
		fmt.Printf("  ✅ Сәтті! Дербес машиналық бағдарлама жасалды: ./%s\n", outputBinary)
		if useLlvm {
			fmt.Println("  (Таза машина коды — LLVM IR арқылы компиляцияланды!)")
		} else {
			fmt.Println("  (Таза x86-64 машина коды — NASM бэкенді!)")
		}
		fmt.Println("═══════════════════════════════════════════════════════════════════════")
	} else {
		// Run binary
		var cmd *exec.Cmd
		programArgs := fs.Args()[1:]
		if filepath.IsAbs(outputBinary) || strings.Contains(outputBinary, string(filepath.Separator)) {
			cmd = exec.Command(outputBinary, programArgs...)
		} else {
			cmd = exec.Command("."+string(filepath.Separator)+outputBinary, programArgs...)
		}
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin

		if verboseFlag {
			fmt.Printf("🚀 Бағдарлама іске қосылуда: %s\n\n", outputBinary)
		}
		runErr := cmd.Run()

		if outputFlag == "" {
			os.Remove(outputBinary)
		}

		if runErr != nil {
			os.Exit(1)
		}
	}
}

func printNumbered(code string) {
	lines := strings.Split(code, "\n")
	for i, line := range lines {
		fmt.Printf("%4d | %s\n", i+1, line)
	}
}

func renderError(filename string, sourceCode string, line int, col int, message string, errorType string, loc *locales.Locale) {
	code := "kk"
	if loc != nil {
		code = loc.Code
	}
	switch code {
	case "en":
		fmt.Printf("\n\033[1;31m%s error (line %d, col %d):\033[0m %s\n", errorType, line, col, message)
	case "ru":
		fmt.Printf("\n\033[1;31mОшибка %s (строка %d, колонка %d):\033[0m %s\n", strings.ToLower(errorType), line, col, message)
	case "it":
		fmt.Printf("\n\033[1;31mErrore di %s (linea %d, colonna %d):\033[0m %s\n", strings.ToLower(errorType), line, col, message)
	default:
		fmt.Printf("\n\033[1;31m%s қатесі (жол %d, баған %d):\033[0m %s\n", errorType, line, col, message)
	}

	lines := strings.Split(sourceCode, "\n")
	if line > 0 && line <= len(lines) {
		if line > 1 {
			fmt.Printf(" \033[34m%4d |\033[0m %s\n", line-1, lines[line-2])
		}

		errorLine := lines[line-1]
		fmt.Printf(" \033[34m%4d |\033[0m %s\n", line, errorLine)

		caretLine := ""
		for i, char := range errorLine {
			if i >= col-1 {
				break
			}
			if char == '\t' {
				caretLine += "\t"
			} else {
				caretLine += " "
			}
		}
		fmt.Printf("      \033[34m|\033[0m %s\033[1;31m^\033[0m\n", caretLine)
	}

	hint := getHint(message, loc)
	if hint != "" {
		hintLabel := "Кеңес"
		switch code {
		case "en":
			hintLabel = "Hint"
		case "ru":
			hintLabel = "Подсказка"
		case "it":
			hintLabel = "Suggerimento"
		}
		fmt.Printf(" \033[1;36m%s:\033[0m %s\n", hintLabel, hint)
	}
	fmt.Println()
}

func getHint(msg string, loc *locales.Locale) string {
	code := "kk"
	if loc != nil {
		code = loc.Code
	}
	if strings.Contains(msg, "айнымалысы жарияланбаған") || strings.Contains(msg, "undeclared") || strings.Contains(msg, "необъявленная") || strings.Contains(msg, "non dichiarata") {
		switch code {
		case "en":
			return "Declare variable before use with 'val' (or 'var' for mutable) e.g.: val x = 5"
		case "ru":
			return "Объявите переменную перед использованием с помощью 'пусть' (или 'перем') например: пусть x = 5"
		case "it":
			return "Dichiara la variabile prima dell'uso con 'val' (o 'var' per mutabile) es.: val x = 5"
		default:
			return "Айнымалыны қолданбас бұрын оны 'болсын' (немесе 'айнымалы') арқылы жариялаңыз: болсын x = 5"
		}
	}
	if strings.Contains(msg, "тұрақты") || strings.Contains(msg, "immutable") || strings.Contains(msg, "неизменяемой") || strings.Contains(msg, "immutabile") {
		switch code {
		case "en":
			return "Use 'var' or 'let mut' to declare a mutable variable"
		case "ru":
			return "Используйте 'перем' для объявления изменяемой переменной"
		case "it":
			return "Usa 'var' per dichiarare una variabile mutabile"
		default:
			return "Айнымалының мәнін өзгерту үшін 'айнымалы' (var) қолданыңыз"
		}
	}
	if strings.Contains(msg, "өзгертуге болмайды") || strings.Contains(msg, "cannot change type") || strings.Contains(msg, "нельзя изменить тип") || strings.Contains(msg, "cambiare il tipo") {
		switch code {
		case "en":
			return "Changing variable type is not allowed. Use a new variable instead"
		case "ru":
			return "Изменение типа переменной не допускается. Используйте новую переменную"
		case "it":
			return "Non è consentito modificare il tipo della variabile. Usa una nuova variabile"
		default:
			return "Айнымалының типін өзгертуге болмайды. Жаңа айнымалыны қолданыңыз"
		}
	}
	if strings.Contains(msg, "табылмады") || strings.Contains(msg, "not found") || strings.Contains(msg, "не найдена") || strings.Contains(msg, "non trovata") {
		switch code {
		case "en":
			return "Check the function name or declare it before calling"
		case "ru":
			return "Проверьте имя функции или объявите её перед вызовом"
		case "it":
			return "Controlla il nome della funzione o dichiarala prima della chiamata"
		default:
			return "Функцияның атауын тексеріңіз немесе оны шақырмас бұрын жариялаңыз"
		}
	}
	return ""
}

func renderWarning(filename string, sourceCode string, line int, col int, message string) {
	fmt.Printf("\n\033[1;33mЕскерту (жол %d, баған %d):\033[0m %s\n", line, col, message)

	lines := strings.Split(sourceCode, "\n")
	if line > 0 && line <= len(lines) {
		if line > 1 {
			fmt.Printf(" \033[34m%4d |\033[0m %s\n", line-1, lines[line-2])
		}

		errorLine := lines[line-1]
		fmt.Printf(" \033[34m%4d |\033[0m %s\n", line, errorLine)

		caretLine := ""
		for i, char := range errorLine {
			if i >= col-1 {
				break
			}
			if char == '\t' {
				caretLine += "\t"
			} else {
				caretLine += " "
			}
		}
		fmt.Printf("      \033[34m|\033[0m %s\033[1;33m^\033[0m\n", caretLine)
	}
	fmt.Println()
}
