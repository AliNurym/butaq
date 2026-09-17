@echo off
chcp 65001 > nul
title Butaq Fast Matrix Fibonacci O(log n) - Rome Congress Demo
mode con: cols=75 lines=26
cls
"%~dp0butaq.exe" run "%~dp0examples\fibonacci_matrix_kk.btq"
echo.
pause
