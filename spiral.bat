@echo off
chcp 65001 > nul
title Butaq Fibonacci Phyllotaxis Spiral - Rome Congress Demo
mode con: cols=68 lines=28
cls
"%~dp0butaq.exe" run -i "%~dp0examples\fibonacci_spiral_kk.btq"
pause
