# CompressImg Windows Service

这是一个自动压缩图片的Windows服务程序。

## 功能
- 监控指定文件夹中的图片文件
- 自动使用TinyPNG API压缩图片
- 支持PNG、JPG、JPEG格式
- 以Windows服务方式运行

## 安装和使用

### 方法1：使用PowerShell脚本（推荐）

1. 以**管理员身份**打开PowerShell
2. 导航到程序目录
3. 执行以下命令：

```powershell
# 安装服务
.\manage_service.ps1 -Action install

# 查看服务状态
.\manage_service.ps1 -Action status

# 停止服务
.\manage_service.ps1 -Action stop

# 启动服务
.\manage_service.ps1 -Action start

# 卸载服务
.\manage_service.ps1 -Action uninstall
```

### 方法2：使用批处理文件

1. 以**管理员身份**运行 `install_service.bat` 安装服务
2. 以**管理员身份**运行 `uninstall_service.bat` 卸载服务

### 方法3：手动命令

```cmd
# 安装服务
sc create "CompressImg" binPath= "D:\path\to\compressImg.exe" DisplayName= "Image Compression Service" start= auto

# 启动服务
sc start "CompressImg"

# 停止服务
sc stop "CompressImg"

# 删除服务
sc delete "CompressImg"
```

## 配置

默认监控路径：`C:\Users\83795\Downloads`

如果需要监控其他路径，可以通过命令行参数指定：
```cmd
compressImg.exe "C:\Your\Custom\Path"
```

## 注意事项

1. **必须以管理员权限运行**安装/卸载命令
2. 服务会自动删除原始图片文件，请确保已备份重要文件
3. 需要网络连接访问TinyPNG API
4. 确保防火墙允许程序访问网络

## 服务管理

可以通过以下方式管理服务：
- Windows服务管理器 (services.msc)
- 任务管理器的服务选项卡
- PowerShell的Get-Service命令
- 命令行的sc命令

## 日志

服务运行时的日志信息会输出到Windows事件日志中。
