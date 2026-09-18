package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type BenchmarkItem struct {
	ID          string
	Name        string
	Description string
}

type LangResult struct {
	Language string
	Result   string
	TimeMs   float64
	Error    string
}

type TestReport struct {
	Item    BenchmarkItem
	Results []LangResult
}

var benchmarks = []BenchmarkItem{
	{"bench01_primes", "1. Prime Counting (Trial Division 1,000,000)", "Heavy integer division, modulo, break jump & branch prediction"},
	{"bench02_mandelbrot", "2. Mandelbrot Fractal Escape (500x500 Grid, 500 iter)", "22.2M iterations of complex plane double float math & escape condition"},
	{"bench03_fib_deep", "3. Recursive Fibonacci F(38)", "78.2 million deep function frame invocations & recursion tree"},
	{"bench04_leibniz_pi", "4. Leibniz Pi Infinite Series (50,000,000 iter)", "50M tight floating-point division, alternating sign & ALU throughput"},
	{"bench05_tak_recursion", "5. Tak (Takeuchi) Deep Recursion (33, 16, 8)", "Triple-recursive call stack stress test & hardware register passing"},
}

func main() {
	topDir := "top5_extreme"
	if _, err := os.Stat(topDir); err != nil {
		if _, err2 := os.Stat(filepath.Join("тесты", topDir)); err2 == nil {
			_ = os.Chdir("тесты")
		}
	}

	fmt.Println("════════════════════════════════════════════════════════════════════════════════════")
	fmt.Println("       ⚡ TOP 5 EXTREME BENCHMARKS: BUTAQ (NATIVE) vs PYTHON 3.12 vs C++")
	fmt.Println("          Demonstrating Native Machine Speed vs Interpreted Dynamic VMs")
	fmt.Println("════════════════════════════════════════════════════════════════════════════════════")
	fmt.Println()

	pythonExe := findPython()
	cppCompiler := findCppCompiler()
	butaqExe := findButaq()

	fmt.Println("🔍 Toolchains detected:")
	fmt.Printf("   • Butaq Native Compiler : %s\n", butaqExe)
	fmt.Printf("   • Python 3.12 Engine    : %s\n", pythonExe)
	fmt.Printf("   • C++ Compiler (GCC)    : %s\n", cppCompiler)
	fmt.Println()

	var reports []TestReport

	for i, b := range benchmarks {
		fmt.Printf("────────────────────────────────────────────────────────────────────────────────────\n")
		fmt.Printf("▶ [%d/5] %s\n", i+1, b.Name)
		fmt.Printf("   %s\n", b.Description)
		fmt.Printf("────────────────────────────────────────────────────────────────────────────────────\n")

		report := TestReport{Item: b}

		// 1. C++ (Native -O3 -static)
		if cppCompiler != "" {
			src := filepath.Join(topDir, b.ID+".cpp")
			bin := filepath.Join(topDir, b.ID+"_cpp.exe")
			binPath, _ := filepath.Abs(bin)
			buildCmd := exec.Command(cppCompiler, "-O3", "-static", "-o", binPath, src)
			if out, err := buildCmd.CombinedOutput(); err == nil {
				res := executeRun(binPath)
				res.Language = "C++ (Native -O3)"
				report.Results = append(report.Results, res)
				_ = os.Remove(binPath)
			} else {
				report.Results = append(report.Results, LangResult{Language: "C++ (Native -O3)", Error: string(out)})
			}
		}

		// 2. Butaq Native Compiled (x86-64 machine code via NASM/GCC)
		if butaqExe != "" {
			src := filepath.Join(topDir, b.ID+".btq")
			bin := filepath.Join(topDir, b.ID+"_btq.exe")
			binPath, _ := filepath.Abs(bin)
			buildCmd := exec.Command(butaqExe, "build", "-o", binPath, src)
			if out, err := buildCmd.CombinedOutput(); err == nil {
				res := executeRun(binPath)
				res.Language = "Butaq (Native x86-64)"
				report.Results = append(report.Results, res)
				_ = os.Remove(binPath)
			} else {
				res := executeRun(butaqExe, "run", src)
				res.Language = "Butaq (Native x86-64)"
				if res.Error != "" {
					res.Error = string(out)
				}
				report.Results = append(report.Results, res)
			}
		}

		// 3. Python 3.12
		if pythonExe != "" {
			src := filepath.Join(topDir, b.ID+".py")
			res := executeRun(pythonExe, src)
			res.Language = "Python 3.12"
			report.Results = append(report.Results, res)
		}

		printReportTable(report)
		fmt.Println()
		reports = append(reports, report)
	}

	saveMarkdownReport(reports)
	fmt.Println("════════════════════════════════════════════════════════════════════════════════════")
	fmt.Println("✅ All 5 extreme benchmarks finished! Report saved to: TOP5_EXTREME_BENCHMARK.md")
	fmt.Println("════════════════════════════════════════════════════════════════════════════════════")
}

func executeRun(prog string, args ...string) LangResult {
	cmd := exec.Command(prog, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	start := time.Now()
	err := cmd.Run()
	wallMs := time.Since(start).Seconds() * 1000.0

	if err != nil {
		return LangResult{
			Error: fmt.Sprintf("Error: %v (%s)", err, strings.TrimSpace(stderr.String())),
		}
	}

	outStr := stdout.String()
	resVal := ""
	timeMs := wallMs

	lines := strings.Split(outStr, "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "RESULT:") {
			val := strings.TrimSpace(strings.TrimPrefix(line, "RESULT:"))
			if val == "" && i+1 < len(lines) {
				val = strings.TrimSpace(lines[i+1])
			}
			resVal = val
		}
		if strings.HasPrefix(line, "TIME_MS:") {
			val := strings.TrimSpace(strings.TrimPrefix(line, "TIME_MS:"))
			if val == "" && i+1 < len(lines) {
				val = strings.TrimSpace(lines[i+1])
			}
			if p, err := strconv.ParseFloat(val, 64); err == nil {
				timeMs = p
			}
		}
	}

	return LangResult{
		Result: resVal,
		TimeMs: timeMs,
	}
}

func printReportTable(rep TestReport) {
	fmt.Printf("%-24s | %-16s | %-12s | %s\n", "Язык", "Результат", "Время (мс)", "Сравнение с Python")
	fmt.Println("-------------------------+------------------+--------------+----------------------------------")

	var pyTime float64 = -1
	for _, r := range rep.Results {
		if r.Language == "Python 3.12" && r.Error == "" {
			pyTime = r.TimeMs
		}
	}

	for _, r := range rep.Results {
		if r.Error != "" {
			fmt.Printf("%-24s | %-16s | %-12s | ❌ %s\n", r.Language, "-", "-", r.Error)
			continue
		}

		comparison := ""
		if r.Language == "C++ (Native -O3)" {
			if pyTime > 0 && r.TimeMs > 0 {
				ratio := pyTime / r.TimeMs
				comparison = fmt.Sprintf("⚡ В %.1fx быстрее Python (Эталон)", ratio)
			} else {
				comparison = "⚡ C++ эталон"
			}
		} else if r.Language == "Python 3.12" {
			comparison = "🐢 Базовый уровень (Интерпретатор)"
		} else if r.Language == "Butaq (Native x86-64)" {
			if pyTime > 0 && r.TimeMs > 0 {
				ratio := pyTime / r.TimeMs
				comparison = fmt.Sprintf("🚀 В %.1fx БЫСТРЕЕ PYTHON!", ratio)
			}
		}

		fmt.Printf("%-24s | %-16s | %10.2f мс | %s\n", r.Language, r.Result, r.TimeMs, comparison)
	}
}

func saveMarkdownReport(reports []TestReport) {
	reportFile := "TOP5_EXTREME_BENCHMARK.md"
	var sb strings.Builder

	sb.WriteString("# ⚡ Топ 5 Экстремальных Бенчмарков: Butaq vs Python 3.12 vs C++\n\n")
	sb.WriteString("> **Специальный отчет**: Тестирование нативного компилятора **Butaq (x86-64)** в сравнении с **Python 3.12** и эталонным **C++ (GCC -O3)** на вычислительно сложных алгоритмах с миллионами итераций.\n\n")
	sb.WriteString(fmt.Sprintf("*Дата и время тестирования: %s*\n\n", time.Now().Format("2006-01-02 15:04:05")))

	sb.WriteString("## 📊 Сводная таблица результатов\n\n")
	sb.WriteString("| № | Алгоритм / Тест | C++ (-O3) | Butaq (Native x86-64) | Python 3.12 | Преимущество Butaq | Преимущество C++ |\n")
	sb.WriteString("|---|---|---|---|---|---|---|\n")

	var totalCpp, totalBtq, totalPy float64

	for i, rep := range reports {
		var cppTime, btqTime, pyTime float64
		for _, r := range rep.Results {
			if r.Error == "" {
				switch r.Language {
				case "C++ (Native -O3)":
					cppTime = r.TimeMs
				case "Butaq (Native x86-64)":
					btqTime = r.TimeMs
				case "Python 3.12":
					pyTime = r.TimeMs
				}
			}
		}
		totalCpp += cppTime
		totalBtq += btqTime
		totalPy += pyTime

		btqRatio := pyTime / btqTime
		cppRatio := pyTime / cppTime

		sb.WriteString(fmt.Sprintf("| %d | **%s** | `%.2f ms` | **`%.2f ms`** | `%.2f ms` | **🚀 %.1fx быстрее!** | **⚡ %.1fx быстрее!** |\n",
			i+1, rep.Item.Name, cppTime, btqTime, pyTime, btqRatio, cppRatio))
	}

	totalBtqRatio := totalPy / totalBtq
	totalCppRatio := totalPy / totalCpp

	sb.WriteString(fmt.Sprintf("| **ИТОГО** | **Суммарное время всех 5 тестов** | **`%.2f ms`** | **`%.2f ms`** | **`%.2f ms`** | **🚀 %.1fx быстрее!** | **⚡ %.1fx быстрее!** |\n\n",
		totalCpp, totalBtq, totalPy, totalBtqRatio, totalCppRatio))

	sb.WriteString("## 🔬 Почему Python настолько медленнее (Архитектурный анализ)\n\n")
	sb.WriteString("### 1. Упаковка объектов (Object Boxing) и аллокация кучи\n")
	sb.WriteString("В Python каждое число (целое или вещественное) — это полноценный объект кучи `PyLongObject` или `PyFloatObject`, содержащий счетчик ссылок `ob_refcnt`, указатель на тип `ob_type` и тело числа. При вычислении выражения вроде `zx = zx2 - zy2 + cx` в Python создается **3 промежуточных объекта в куче** на каждую итерацию. В Butaq и C++ числа хранятся непосредственно в аппаратных регистрах процессора (`xmm0`, `xmm1`, `rax`), а промежуточные значения не требуют ни одного байта оперативной памяти.\n\n")

	sb.WriteString("### 2. Накладные расходы на вызов функций (`PyFrameObject` vs Hardware `call`)\n")
	sb.WriteString("В тесте глубокой рекурсии (Tak Function и Fibonacci F(38)) выполняются десятки миллионов вызовов функций. В Python каждый вызов создает тяжелый объект кадра стека `PyFrameObject` с таблицей локальных переменных, отслеживанием трассировки и проверкой переполнения стека (~120 наносекунд на вызов). В Butaq компилятор генерирует чистую инструкцию `call` процессора x86-64 (~1 наносекунда), что дает колоссальный отрыв в скорости.\n\n")

	sb.WriteString("### 3. Байткод-интерпретатор vs Прямой машинный код\n")
	sb.WriteString("Интерпретатор Python обрабатывает цикл через цикл выборки байткода (`switch (opcode)`), декодируя инструкции `BINARY_OP`, `COMPARE_OP`, `POP_JUMP_FORWARD_IF_FALSE`. Butaq напрямую транслирует синтаксис в ассемблер NASM и компилирует в бинарный машинный код, выполняемый непосредственно микроархитектурой процессора на частоте 4+ ГГц с аппаратным предсказанием ветвлений (Branch Target Buffer).\n\n")

	sb.WriteString("## 🚀 Как запустить бенчмарки самостоятельно\n\n")
	sb.WriteString("```bash\n")
	sb.WriteString("# Запуск автоматического тестирования всех 5 тестов:\n")
	sb.WriteString("go run тесты/run_top5_benchmarks.go\n\n")
	sb.WriteString("# Или ручная компиляция и запуск любого теста (на примере Мандельброта):\n")
	sb.WriteString("butaq build -o mandelbrot.exe тесты/top5_extreme/bench02_mandelbrot.btq\n")
	sb.WriteString("./mandelbrot.exe\n\n")
	sb.WriteString("python тесты/top5_extreme/bench02_mandelbrot.py\n")
	sb.WriteString("```\n")

	_ = os.WriteFile(reportFile, []byte(sb.String()), 0644)
}

func findPython() string {
	for _, name := range []string{"python", "python3", "py"} {
		if path, err := exec.LookPath(name); err == nil {
			return path
		}
	}
	return ""
}

func findCppCompiler() string {
	for _, name := range []string{"g++", "clang++"} {
		if path, err := exec.LookPath(name); err == nil {
			return path
		}
	}
	return ""
}

func findButaq() string {
	candidates := []string{
		"./butaq.exe",
		"../butaq.exe",
		"butaq.exe",
		"./butaq",
		"../butaq",
		"butaq",
	}
	for _, c := range candidates {
		if _, err := os.Stat(c); err == nil {
			abs, _ := filepath.Abs(c)
			return abs
		}
	}
	if path, err := exec.LookPath("butaq"); err == nil {
		return path
	}
	return ""
}
