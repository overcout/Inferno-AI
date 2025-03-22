@echo off
cls
echo ================================
echo        Inferno-AI Launcher      
echo ================================
echo [1] Run Telegram Bot
echo [2] Run OAuth Server
echo [3] Run Both
echo [0] Exit
echo.

set /p choice=Select option: 

if "%choice%"=="1" (
    cd cmd\bot
    go run .
    cd ../..
    pause
)

if "%choice%"=="2" (
    cd cmd\oauth
    go run .
    cd ../..
    pause
)

if "%choice%"=="3" (
    start cmd /k "cd cmd\bot && go run ."
    start cmd /k "cd cmd\oauth && go run ."
    pause
)

if "%choice%"=="0" (
    exit
)

