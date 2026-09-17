@echo off
chcp 65001 > nul
title Butaq vs C++ vs Python - Fibonacci Benchmark
mode con: cols=90 lines=32
cls
cd /d "%~dp0"
"C:\Program Files\Go\bin\go.exe" run run_benchmarks.go
echo.
pause
