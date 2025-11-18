@echo off
echo 生成 Go 绑定代码...

:: 使用 abigen 生成 Go 代码
abigen --bin=Counter.bin --abi=Counter.abi --pkg=bindings --type=Counter --out=../bindings/counter.go

if %errorlevel% equ 0 (
    echo Go 绑定代码生成成功！
    echo 文件位置: ../bindings/counter.go
) else (
    echo 代码生成失败！
)

pause