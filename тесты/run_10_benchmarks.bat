@echo off
chcp 65001 > nul
title Butaq vs Python vs C++ - 10 Extreme Speed Benchmarks
mode con: cols=105 lines=38
cls
cd /d "%~dp0"
"C:\Program Files\Go\bin\go.exe" run run_10_benchmarks.go
echo.
pause
