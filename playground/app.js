// Predefined Code Examples
const examples = {
    fib: `функция фиб(n)
    егер n <= 1 сонда қайтару n
    қайтару фиб(n - 1) + фиб(n - 2)

жазу("=== Фибоначчи қатары ===")

болсын i = 0
әзірше i <= 10
    жазу("fib(", i, ") = ", фиб(i))
    i = i + 1`,

    donut: `# 3D Donut (Айналмалы 3D Тор / Бублик)
болсын A = 1.0
болсын B = 1.0

болсын b = жаңа_тізім(880, " ")
болсын z = жаңа_тізім(880, 0.0)

болсын sin_A = синус(A)
болсын cos_A = косинус(A)
болсын sin_B = синус(B)
болсын cos_B = косинус(B)

болсын j = 0.0
әзірше j < 6.28
    болсын sin_j = синус(j)
    болсын cos_j = косинус(j)
    
    болсын i = 0.0
    әзірше i < 6.28
        болсын sin_i = синус(i)
        болсын cos_i = косинус(i)
        
        болсын h = cos_j + 2.0
        болсын D = 1.0 / (sin_i * h * sin_A + sin_j * cos_A + 5.0)
        болсын t = sin_i * h * cos_A - sin_j * sin_A
        
        болсын x = бүтін(20.0 + 18.0 * D * (cos_i * h * cos_B - t * sin_B))
        болсын y = бүтін(11.0 + 9.0 * D * (cos_i * h * sin_B + t * cos_B))
        
        болсын N = бүтін(8.0 * ((sin_j * sin_A - sin_i * cos_j * cos_A) * cos_B - sin_i * cos_j * sin_A - sin_j * cos_A - cos_i * cos_j * sin_B))
        
        егер y >= 0 және y < 22 және x >= 0 және x < 40
            болсын idx = y * 40 + x
            егер D > z[idx]
                z[idx] = D
                егер N > 0 және N < 12
                    b[idx] = ".,-~:;=!*#$@"[N]
                әйтпесе
                    b[idx] = "."
        
        i = i + 0.08
    j = j + 0.04

болсын жол_y = 0
әзірше жол_y < 22
    болсын жол = ""
    болсын баған_x = 0
    әзірше баған_x < 40
        жол = жол + b[жол_y * 40 + баған_x]
        баған_x = баған_x + 1
    жазу(жол)
    жол_y = жол_y + 1`,
    
    hello: `жазу("Сәлем, Әлем!")`,
    
    math: `# Математикалық инфикс өрнектер
болсын нәтиже = (5 + 10) * 3
# Нәтиже: (5 + 10) * 3 = 45
жазу("Нәтиже: ", нәтиже)`,
    
    loop: `# әзірше циклімен 1-ден 5-ке дейін санау
болсын санауыш = 1

әзірше санауыш <= 5
    жазу("Санауыш: ", санауыш)
    санауыш = санауыш + 1`,
    
    func: `# Екі санды қосатын функция
функция қосу_екі(сан1, сан2)
    қайтару сан1 + сан2

болсын нәтиже = қосу_екі(25, 35)
жазу("Нәтиже = ", нәтиже)`,
    
    struct: `# Құрылымды жариялау және қолдану
құрылым Адам {
    аты МӘТІН
    жасы БҮТІН
}

болсын әли = Адам жасау
әли.аты = "Әлихан"
әли.жасы = 21

жазу("Аты: ", әли.аты)
жазу("Жасы: ", әли.жасы)`
};

// UI Elements
const editor = document.getElementById('code-editor');
const lineNumbers = document.getElementById('line-numbers');
const runBtn = document.getElementById('run-btn');
const transpileBtn = document.getElementById('transpile-btn');
const examplesDropdown = document.getElementById('examples-dropdown');
const targetLocaleSelect = document.getElementById('target-locale-select');
const compilerStatus = document.getElementById('compiler-status');
const detectedLocaleBadge = document.getElementById('detected-locale-badge');
const transpileToast = document.getElementById('transpile-toast');

const outputConsole = document.getElementById('output-console');
const llvmConsole = document.getElementById('llvm-console');
const nasmConsole = document.getElementById('nasm-console');

let toastTimeout = null;

function showToast(message, isError = false) {
    if (toastTimeout) clearTimeout(toastTimeout);
    transpileToast.textContent = message;
    transpileToast.className = 'transpile-toast' + (isError ? ' error-toast' : '');
    transpileToast.classList.remove('hidden');
    toastTimeout = setTimeout(() => {
        transpileToast.classList.add('hidden');
    }, 6000);
}

function updateLocaleBadge(code) {
    const map = {
        kk: 'Қазақша (KK)',
        en: 'English (EN)',
        it: 'Italiano (IT)',
        ru: 'Русский (RU)'
    };
    detectedLocaleBadge.textContent = 'Тіл: ' + (map[code.toLowerCase()] || code.toUpperCase());
}

// Simple heuristic detector for editor input
function detectEditorLocale(code) {
    if (/(\b(funzione|ritorna|se|allora|altrimenti|mentre|stampa|sia)\b)/i.test(code)) return 'it';
    if (/(\b(func|return|if|then|else|while|print|let)\b)/i.test(code)) return 'en';
    if (/(\b(функция|вернуть|если|тогда|иначе|пока|печать|пусть)\b)/i.test(code)) return 'ru';
    if (/(\b(болсын|егер|сонда|әйтпесе|әзірше|қайтару|жазу)\b)/i.test(code)) return 'kk';
    return 'kk';
}

// Load selected example
examplesDropdown.addEventListener('change', (e) => {
    editor.value = examples[e.target.value] || '';
    updateLineNumbers();
    updateLocaleBadge(detectEditorLocale(editor.value));
});

// Update Line Numbers & auto-detect language
function updateLineNumbers() {
    const lines = editor.value.split('\n');
    lineNumbers.innerHTML = Array(lines.length).fill(0).map((_, i) => i + 1).join('<br>');
}

editor.addEventListener('input', () => {
    updateLineNumbers();
    updateLocaleBadge(detectEditorLocale(editor.value));
});

editor.addEventListener('scroll', () => {
    lineNumbers.scrollTop = editor.scrollTop;
});

// Setup Tabs
document.querySelectorAll('.tab-btn').forEach(btn => {
    btn.addEventListener('click', () => {
        document.querySelectorAll('.tab-btn').forEach(b => b.classList.remove('active'));
        document.querySelectorAll('.tab-pane').forEach(p => p.classList.remove('active'));
        
        btn.classList.add('active');
        const tabId = btn.getAttribute('data-tab') + '-pane';
        document.getElementById(tabId).classList.add('active');
    });
});

// WASM Integration
const go = new Go();
WebAssembly.instantiateStreaming(fetch('butaq.wasm'), go.importObject).then((result) => {
    go.run(result.instance);
    compilerStatus.textContent = 'Дайын';
    compilerStatus.classList.add('ready');
    
    // Set default value
    editor.value = examples.fib;
    updateLineNumbers();
    updateLocaleBadge('kk');

    // Wire up benchmark button now that WASM is ready
    const benchRunBtn = document.getElementById('bench-run-btn');
    if (benchRunBtn && !benchRunBtn._benchWired) {
        benchRunBtn.addEventListener('click', runBenchmark);
        benchRunBtn._benchWired = true;
    }
}).catch(err => {
    console.error('WASM loading error:', err);
    compilerStatus.textContent = 'Қате';
    compilerStatus.style.color = '#ef4444';
});

// AST Transpile Action
transpileBtn.addEventListener('click', () => {
    if (compilerStatus.textContent !== 'Дайын') {
        alert('Компилятор әлі жүктелуде, сәл күте тұрыңыз...');
        return;
    }

    const code = editor.value;
    const targetLang = targetLocaleSelect.value;
    const fromLang = detectEditorLocale(code);

    if (fromLang === targetLang) {
        showToast(`ℹ️ Код қазірдің өзінде [${targetLang.toUpperCase()}] тілінде жазылған.`);
        return;
    }

    const res = window.transpileButaqCode(code, targetLang, fromLang);
    if (res.error) {
        showToast('❌ Транспиляция қатесі: ' + res.error, true);
        return;
    }

    editor.value = res.code;
    updateLineNumbers();
    updateLocaleBadge(res.toLocale);
    showToast(`✨ AST-транспиляция сәтті аяқталды [${res.fromLocale.toUpperCase()} → ${res.toLocale.toUpperCase()}]! Код логикасы толық сақталды.`);
});

// Run Code
runBtn.addEventListener('click', () => {
    if (compilerStatus.textContent !== 'Дайын') {
        alert('Компилятор әлі жүктелуде, сәл күте тұрыңыз...');
        return;
    }

    const code = editor.value;
    const detectedLang = detectEditorLocale(code);

    // 1. Run Interpreter
    outputConsole.classList.remove('error');
    outputConsole.textContent = 'Орындалуда...';
    
    const interpreterResult = window.runButaqCode(code, detectedLang);
    if (interpreterResult.error) {
        outputConsole.classList.add('error');
        outputConsole.textContent = interpreterResult.error;
    } else {
        outputConsole.textContent = interpreterResult.output || 'Бағдарлама сәтті аяқталды (шығыс мәліметтер жоқ)';
        if (interpreterResult.locale) {
            updateLocaleBadge(interpreterResult.locale);
        }
    }

    // 2. Compile to LLVM IR
    const llvmResult = window.compileToLlvm(code, detectedLang);
    if (llvmResult.error) {
        llvmConsole.textContent = llvmResult.error;
    } else {
        llvmConsole.textContent = llvmResult.code;
    }

    // 3. Compile to NASM
    const nasmResult = window.compileToNasm(code, detectedLang);
    if (nasmResult.error) {
        nasmConsole.textContent = nasmResult.error;
    } else {
        nasmConsole.textContent = nasmResult.code;
    }
});

// ============================================
// BENCHMARK RUNNER
// ============================================

// Fib code templates per locale
function buildFibCode(locale, n) {
    const templates = {
        kk: `функция фиб(n)
    егер n <= 1 сонда қайтару n
    қайтару фиб(n - 1) + фиб(n - 2)
жазу(фиб(${n}))`,

        en: `func fib(n)
    if n <= 1 then return n
    return fib(n - 1) + fib(n - 2)
print(fib(${n}))`,

        it: `funzione fib(n)
    se n <= 1 allora ritorna n
    ritorna fib(n - 1) + fib(n - 2)
stampa(fib(${n}))`,

        ru: `функция фиб(n)
    если n <= 1 тогда вернуть n
    вернуть фиб(n - 1) + фиб(n - 2)
печать(фиб(${n}))`
    };
    return templates[locale] || templates.kk;
}

const BENCH_LOCALES = [
    { key: 'kk', label: 'Қазақша (KK)', flag: '🇰🇿' },
    { key: 'en', label: 'English (EN)',  flag: '🇬🇧' },
    { key: 'it', label: 'Italiano (IT)', flag: '🇮🇹' },
    { key: 'ru', label: 'Русский (RU)',  flag: '🇷🇺' },
];

function resetBenchUI() {
    BENCH_LOCALES.forEach(({ key }) => {
        const card = document.getElementById('bench-' + key);
        const bar  = document.getElementById('bench-bar-' + key);
        const time = document.getElementById('bench-time-' + key);
        const out  = document.getElementById('bench-out-' + key);
        card.className = 'bench-locale';
        bar.style.width = '0%';
        time.textContent = '—';
        out.textContent  = '';
    });
    const verdict = document.getElementById('bench-verdict');
    verdict.className = 'bench-verdict hidden';
    verdict.textContent = '';
}

async function runBenchmark() {
    if (!window.runButaqCode) {
        alert('Компилятор әлі жүктелуде...');
        return;
    }

    const benchBtn = document.getElementById('bench-run-btn');
    const nInput   = document.getElementById('bench-n');
    const n = Math.max(5, Math.min(35, parseInt(nInput.value, 10) || 28));
    nInput.value = n;

    benchBtn.disabled = true;
    benchBtn.innerHTML = '<span class="bench-btn-icon">⏳</span> Орындалуда...';

    // Update subtitle
    document.querySelector('.bench-subtitle').textContent =
        `Fibonacci(${n}) · Interpreter · All Locales`;

    resetBenchUI();

    const results = {};

    // Run each locale sequentially with small delay for visual effect
    for (const { key } of BENCH_LOCALES) {
        const card = document.getElementById('bench-' + key);
        const timeEl = document.getElementById('bench-time-' + key);
        const outEl  = document.getElementById('bench-out-' + key);

        card.classList.add('running');
        timeEl.textContent = '⏳';

        // Yield to browser for repaint
        await new Promise(r => setTimeout(r, 60));

        const code = buildFibCode(key, n);
        const t0 = performance.now();
        const res = window.runButaqCode(code, key);
        const elapsed = performance.now() - t0;

        card.classList.remove('running');

        if (res.error) {
            timeEl.textContent = 'Қате';
            outEl.textContent  = '✗ ' + res.error.slice(0, 80);
            results[key] = { ms: Infinity, output: '', error: true };
        } else {
            const ms = elapsed;
            results[key] = { ms, output: (res.output || '').trim() };
            timeEl.textContent = ms < 10 ? ms.toFixed(2) + ' ms'
                                : ms < 1000 ? ms.toFixed(1) + ' ms'
                                : (ms / 1000).toFixed(2) + ' s';
            outEl.textContent  = '✓ ' + (res.output || '').trim().split('\n')[0];
        }

        card.classList.add('revealed');
        // Small gap between locales
        await new Promise(r => setTimeout(r, 80));
    }

    // Find best (min ms, ignoring errors)
    const valid = BENCH_LOCALES.filter(({ key }) => !results[key]?.error);
    if (valid.length === 0) {
        benchBtn.disabled = false;
        benchBtn.innerHTML = '<span class="bench-btn-icon">▶</span> Запустить тест';
        return;
    }

    const maxMs  = Math.max(...valid.map(({ key }) => results[key].ms));
    const minMs  = Math.min(...valid.map(({ key }) => results[key].ms));
    const winner = valid.reduce((a, b) => results[a.key].ms < results[b.key].ms ? a : b);

    // Animate bars proportionally (fastest = 100%, slowest = relative %)
    // All bars scale relative to slowest so fastest is always 100%
    BENCH_LOCALES.forEach(({ key }) => {
        if (results[key]?.error) return;
        const pct = maxMs > 0 ? (results[key].ms / maxMs) * 100 : 100;
        // Invert so fastest bar is longest (100%)
        const barPct = maxMs > 0 ? ((maxMs - results[key].ms + minMs) / maxMs) * 100 : 100;
        document.getElementById('bench-bar-' + key).style.width = Math.max(barPct, 4) + '%';
    });

    // Highlight winner
    document.getElementById('bench-' + winner.key).classList.add('winner');

    // Show verdict
    const spread = maxMs - minMs;
    const verdict = document.getElementById('bench-verdict');
    verdict.className = 'bench-verdict';
    verdict.innerHTML = `🏆 Победитель: <strong>${winner.flag} ${winner.label}</strong> — ${results[winner.key].ms.toFixed(1)} ms<br>
        <span style="font-weight:400;opacity:0.8;">Разброс: ${spread.toFixed(1)} ms · fib(${n}) = ${results[winner.key].output}</span>`;

    benchBtn.disabled = false;
    benchBtn.innerHTML = '<span class="bench-btn-icon">▶</span> Запустить тест';
}

// Fallback wire-up for DOMContentLoaded (if WASM was already loaded)
document.addEventListener('DOMContentLoaded', () => {
    const benchRunBtn = document.getElementById('bench-run-btn');
    if (benchRunBtn && !benchRunBtn._benchWired) {
        benchRunBtn.addEventListener('click', runBenchmark);
        benchRunBtn._benchWired = true;
    }
});
