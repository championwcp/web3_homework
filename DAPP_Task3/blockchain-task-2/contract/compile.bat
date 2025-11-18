@echo off
echo 正在编译 Counter.sol...

:: 使用 solcjs 编译
solcjs --bin --abi Counter.sol

if %errorlevel% equ 0 (
    echo 编译成功！
    echo 重命名文件...
    ren Counter_sol_Counter.abi Counter.abi
    ren Counter_sol_Counter.bin Counter.bin
    echo 生成的文件:
    dir *.bin *.abi
) else (
    echo 编译失败！
)

pause