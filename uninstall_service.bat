@echo off
echo Uninstalling CompressImg Service...

REM 停止服务
echo Stopping service...
sc stop "CompressImg"

REM 等待服务停止
timeout /t 3 /nobreak > nul

REM 删除服务
echo Removing service...
sc delete "CompressImg"

if %errorlevel% == 0 (
    echo Service uninstalled successfully!
) else (
    echo Failed to uninstall service. Error code: %errorlevel%
)

pause
