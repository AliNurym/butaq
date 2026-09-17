@echo off
chcp 65001 > nul
title Butaq 3D Fibonacci Golden Sphere - Rome Congress Demo
mode con: cols=68 lines=28
cls
"%~dp0butaq.exe" run -i "%~dp0examples\fibonacci_sphere_kk.btq"
pause
