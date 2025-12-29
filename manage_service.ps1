# CompressImg Service Management Script
# 需要以管理员权限运行

param(
    [Parameter(Mandatory=$true)]
    [ValidateSet("install", "uninstall", "start", "stop", "status")]
    [string]$Action
)

$ServiceName = "CompressImg"
$ServiceDisplayName = "Image Compression Service"
$ServiceDescription = "Automatically compresses images in specified folders"
$ExePath = Join-Path $PSScriptRoot "compressImg.exe"

function Install-Service {
    Write-Host "Installing $ServiceDisplayName..." -ForegroundColor Green
    
    if (-not (Test-Path $ExePath)) {
        Write-Error "Executable not found: $ExePath"
        return
    }
    
    try {
        New-Service -Name $ServiceName -BinaryPathName $ExePath -DisplayName $ServiceDisplayName -Description $ServiceDescription -StartupType Automatic
        Write-Host "Service installed successfully!" -ForegroundColor Green
        
        Write-Host "Starting service..." -ForegroundColor Yellow
        Start-Service -Name $ServiceName
        Write-Host "Service started successfully!" -ForegroundColor Green
    }
    catch {
        Write-Error "Failed to install service: $($_.Exception.Message)"
    }
}

function Uninstall-Service {
    Write-Host "Uninstalling $ServiceDisplayName..." -ForegroundColor Yellow
    
    try {
        $service = Get-Service -Name $ServiceName -ErrorAction SilentlyContinue
        if ($service) {
            if ($service.Status -eq "Running") {
                Write-Host "Stopping service..." -ForegroundColor Yellow
                Stop-Service -Name $ServiceName -Force
                Start-Sleep -Seconds 3
            }
            
            Remove-Service -Name $ServiceName
            Write-Host "Service uninstalled successfully!" -ForegroundColor Green
        }
        else {
            Write-Host "Service not found." -ForegroundColor Yellow
        }
    }
    catch {
        Write-Error "Failed to uninstall service: $($_.Exception.Message)"
    }
}

function Start-CompressService {
    try {
        Start-Service -Name $ServiceName
        Write-Host "Service started successfully!" -ForegroundColor Green
    }
    catch {
        Write-Error "Failed to start service: $($_.Exception.Message)"
    }
}

function Stop-CompressService {
    try {
        Stop-Service -Name $ServiceName -Force
        Write-Host "Service stopped successfully!" -ForegroundColor Green
    }
    catch {
        Write-Error "Failed to stop service: $($_.Exception.Message)"
    }
}

function Get-ServiceStatus {
    try {
        $service = Get-Service -Name $ServiceName -ErrorAction SilentlyContinue
        if ($service) {
            Write-Host "Service Status: $($service.Status)" -ForegroundColor Cyan
            Write-Host "Service Name: $($service.Name)" -ForegroundColor Cyan
            Write-Host "Display Name: $($service.DisplayName)" -ForegroundColor Cyan
        }
        else {
            Write-Host "Service not installed." -ForegroundColor Yellow
        }
    }
    catch {
        Write-Error "Failed to get service status: $($_.Exception.Message)"
    }
}

# 检查是否以管理员权限运行
if (-NOT ([Security.Principal.WindowsPrincipal] [Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole] "Administrator")) {
    Write-Error "This script requires Administrator privileges. Please run PowerShell as Administrator."
    exit 1
}

switch ($Action) {
    "install" { Install-Service }
    "uninstall" { Uninstall-Service }
    "start" { Start-CompressService }
    "stop" { Stop-CompressService }
    "status" { Get-ServiceStatus }
}
