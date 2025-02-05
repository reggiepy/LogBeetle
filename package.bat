@echo off
REM 设置源代码目录和输出目录
SET SOURCE_DIR="github.com/reggiepy/LogBeetle/cmd/LogBeetle"
SET OUTPUT_DIR=.\bin
SET APP_NAME=LogBeetle

REM 创建输出目录（如果不存在）
IF NOT EXIST %OUTPUT_DIR% (
    mkdir %OUTPUT_DIR%
)

REM 启用延迟扩展
SETLOCAL EnableDelayedExpansion

REM 目标平台的操作系统和架构
SET PLATFORMS=linux/amd64 windows/amd64

REM 循环遍历所有目标平台进行编译
FOR %%P IN (%PLATFORMS%) DO (
    REM 分解目标平台
    FOR /F "tokens=1,2 delims=/" %%A IN ("%%P") DO (
        SET GOOS=%%A
        SET GOARCH=%%B

        REM 设置输出文件名，包含操作系统和架构
        SET OUTPUT_FILE=%OUTPUT_DIR%\%APP_NAME%-%%A-%%B

        REM Windows 平台需要添加 .exe 后缀
        IF %%A==windows (
            SET OUTPUT_FILE=!OUTPUT_FILE!.exe
        )

        REM 设置GOOS和GOARCH环境变量并进行编译
        echo Building for %%A/%%B...
        SET GOOS=%%A
        SET GOARCH=%%B
        go build -o !OUTPUT_FILE! %SOURCE_DIR%

        REM 检查编译是否成功
        IF %ERRORLEVEL% EQU 0 (
            echo Successfully built for %%A/%%B: !OUTPUT_FILE!
        ) ELSE (
            echo Failed to build for %%A/%%B
        )
    )
)

echo Build process completed.
pause
