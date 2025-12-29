@echo off
echo Installing CompressImg Service...

REM 获取当前目录的完整路径
set "CURRENT_DIR=%~dp0"
set "EXE_PATH=%CURRENT_DIR%compressImg.exe"

echo Service executable path: %EXE_PATH%

REM 创建服务
sc create "CompressImg" binPath= "%EXE_PATH%" DisplayName= "Image Compression Service" start= auto

if %errorlevel% == 0 (
    echo Service installed successfully!
    echo Starting service...
    sc start "CompressImg"
    if %errorlevel% == 0 (
        echo Service started successfully!
    ) else (
        echo Failed to start service. Error code: %errorlevel%
    )
) else (
    echo Failed to install service. Error code: %errorlevel%
)

pause
