#!/bin/bash

contractName=$1       # 合约目录名，例如 erc721
buildOption=$2        # 编译方式: tinygo 或 go
targetARCH=$3         # 可选架构参数
crypto=""

if [[ -z $contractName ]]; then
  echo "Usage: ./build.sh <contractName> [tinygo|go] [targetARCH]"
  exit 1
fi

contractPath="./$contractName"
contractMain="$contractPath/$contractName.go"

if [[ ! -d $contractPath ]]; then
  echo "Error: $contractName not found in ./demo/"
  exit 1
fi

if [[ ! -f $contractMain ]]; then
  echo "Error: $contractMain not found"
  exit 1
fi

# 默认编译器
if [[ -z $buildOption ]]; then
  buildOption="tinygo"
fi

if [[ "$(uname)" == "Linux" ]]; then
  crypto="-tags crypto"
fi

# 进入合约目录
cd "$contractPath" || exit 1

# 编译输出文件名
outputName=""
if [[ $buildOption == "tinygo" ]]; then


  outputName="$contractName-tinygo.wasm"
  echo "[*] Using TinyGo to compile $contractName..."
  tinygo build -no-debug -opt=s -o "$outputName" -target=wasip1
else


  outputName="$contractName-go.wasm"
  echo "[*] Using Go to compile $contractName..."
  GOOS=wasip1 GOARCH=wasm go build -o "$outputName"
fi

# 返回 demo 目录
cd - >/dev/null

# 拷贝 wasm 到 demo 目录
cp "$contractPath/$outputName" "./$outputName"
echo "[✓] Done compiling $contractName → $outputName"
