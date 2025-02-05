#!/bin/bash

# 设置源代码目录和输出目录
SOURCE_DIR="github.com/reggiepy/LogBeetle/cmd/LogBeetle"    # Go 源代码所在目录
OUTPUT_DIR="./bin" # 编译输出目录
APP_NAME="LogBeetle" # 编译输出目录

# 创建输出目录（如果不存在）
mkdir -p $OUTPUT_DIR

# 目标平台的操作系统和架构
#PLATFORMS=("linux/amd64" "linux/arm64" "darwin/amd64" "windows/amd64" "freebsd/amd64")
PLATFORMS=("linux/amd64" "windows/amd64")

# 循环遍历所有目标平台进行编译
for PLATFORM in "${PLATFORMS[@]}"; do
    # 分解目标平台
    OS=$(echo $PLATFORM | cut -d'/' -f1)
    ARCH=$(echo $PLATFORM | cut -d'/' -f2)

    # 设置GOOS和GOARCH环境变量
    export GOOS=$OS
    export GOARCH=$ARCH

    # 设置输出文件名，包含操作系统和架构
    OUTPUT_FILE="$OUTPUT_DIR/$APP_NAME-$OS-$ARCH"

    if [ "$OS" == "windows" ]; then
      OUTPUT_FILE="$OUTPUT_FILE.exe"
    fi

    # 编译
    echo "Building for $OS/$ARCH..."
    go build -o $OUTPUT_FILE $SOURCE_DIR

    if [ $? -eq 0 ]; then
        echo "Successfully built for $OS/$ARCH: $OUTPUT_FILE"
    else
        echo "Failed to build for $OS/$ARCH"
    fi
done

echo "Build process completed."
