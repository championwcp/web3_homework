// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract IntegerToRoman {
    /**
     * @dev 将整数转换为罗马数字
     * @param num 输入的整数 (1 <= num <= 3999)
     * @return 对应的罗马数字字符串
     */
    function intToRoman(uint256 num) public pure returns (string memory) {
        require(num >= 1 && num <= 3999, "Number must be between 1 and 3999");
        
        // 定义数值和对应的罗马符号
        uint256[] memory values = new uint256[](13);
        string[] memory symbols = new string[](13);
        
        values[0] = 1000; symbols[0] = "M";
        values[1] = 900;  symbols[1] = "CM";
        values[2] = 500;  symbols[2] = "D";
        values[3] = 400;  symbols[3] = "CD";
        values[4] = 100;  symbols[4] = "C";
        values[5] = 90;   symbols[5] = "XC";
        values[6] = 50;   symbols[6] = "L";
        values[7] = 40;   symbols[7] = "XL";
        values[8] = 10;   symbols[8] = "X";
        values[9] = 9;    symbols[9] = "IX";
        values[10] = 5;   symbols[10] = "V";
        values[11] = 4;   symbols[11] = "IV";
        values[12] = 1;   symbols[12] = "I";
        
        bytes memory result;
        
        for (uint256 i = 0; i < values.length; i++) {
            while (num >= values[i]) {
                // 将字符串转换为bytes并追加到结果中
                bytes memory symbolBytes = bytes(symbols[i]);
                for (uint256 j = 0; j < symbolBytes.length; j++) {
                    // 动态扩展bytes数组
                    bytes memory temp = new bytes(result.length + 1);
                    for (uint256 k = 0; k < result.length; k++) {
                        temp[k] = result[k];
                    }
                    temp[result.length] = symbolBytes[j];
                    result = temp;
                }
                num -= values[i];
            }
        }
        
        return string(result);
    }
}