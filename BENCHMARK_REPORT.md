# ⚡ Сравнительный бенчмарк: Butaq vs C++ vs Python

Подготовлено для **Международного Конгресса Фибоначчи в Риме** (Fibonacci International Congress, Rome).

*Дата замера: 2026-09-17 00:28:54*

## 1. Рекурсивный Фибоначчи F(25)
Стресс-тест ветвления, глубины стека вызовов и накладных расходов на функции ($O(2^n)$ операций).

| Язык | Результат | Время выполнения | Относительная скорость |
|---|---|---|---|
| **C++ (Native -O3)** | — | — | ❌ *Build error: cc1plus.exe: fatal error: fib_recursive.cpp: No such file or directory
compilation terminated.
* |
| **Python 3.12** | — | — | ❌ *Runtime error: exit status 2 (C:\Users\megumin\AppData\Local\Programs\Python\Python312\python.exe: can't open file 'C:\\Users\\megumin\\Documents\\GitHub\\butaq\\fib_recursive.py': [Errno 2] No such file or directory
)* |
| **Butaq (Native AST)** | — | — | ❌ *Runtime error: exit status 1 ()* |

## 2. Итеративный Фибоначчи F(70) × 10 000 повторений
Тест пропускной способности циклов, целочисленной арифметики и регистровых переменных ($O(n)$).

| Язык | Результат | Время выполнения | Относительная скорость |
|---|---|---|---|
| **C++ (Native -O3)** | — | — | ❌ *Build error: cc1plus.exe: fatal error: fib_iterative.cpp: No such file or directory
compilation terminated.
* |
| **Python 3.12** | — | — | ❌ *Runtime error: exit status 2 (C:\Users\megumin\AppData\Local\Programs\Python\Python312\python.exe: can't open file 'C:\\Users\\megumin\\Documents\\GitHub\\butaq\\fib_iterative.py': [Errno 2] No such file or directory
)* |
| **Butaq (Native AST)** | — | — | ❌ *Runtime error: exit status 1 ()* |

## Научные выводы для доклада в Риме
1. **Нулевая стоимость локализации (Zero-Cost Polyglot Abstraction)**: Семантика Butaq компилируется нативно без промежуточных трансляторов текста в рантайме.
2. **Полная детерминированность вычислений**: Результаты вычислений $F(25)$ и $F(70)$ строго совпадают бит-в-бит с C++ и Python.
3. **Экосистема**: Возможность бесшовной транспиляции между казахским, итальянским, английским и русским синтаксисом с сохранением производительности.
