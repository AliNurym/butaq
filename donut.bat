@echo off
chcp 65001 > nul
title Butaq 3D Donut - Fibonacci Rome Congress Demo
mode con: cols=65 lines=26
cls
"%~dp0butaq.exe" run -i "%~dp0examples\donut_kk.btq"
pause
