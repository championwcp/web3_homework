// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract RomanToInteger {
    /**
     * @dev 将罗马数字转换为整数
     * @param s 罗马数字字符串
     * @return 对应的整数值
     */
    function romanToInt(string memory s) public pure returns (uint256) {
        bytes memory roman = bytes(s);
        uint256 length = roman.length;
        uint256 result = 0;
        
        for (uint256 i = 0; i < length; i++) {
            uint256 current = getValue(roman[i]);
            
            // 如果当前字符比下一个字符小，说明是特殊情况（如IV, IX等）
            if (i < length - 1 && current < getValue(roman[i + 1])) {
                result -= current;
            } else {
                result += current;
            }
        }
        
        return result;
    }
    
    /**
     * @dev 获取罗马字符对应的数值
     * @param c 罗马字符
     * @return 对应的数值
     */
    function getValue(bytes1 c) private pure returns (uint256) {
        if (c == 'I') return 1;
        if (c == 'V') return 5;
        if (c == 'X') return 10;
        if (c == 'L') return 50;
        if (c == 'C') return 100;
        if (c == 'D') return 500;
        if (c == 'M') return 1000;
        revert("Invalid Roman numeral");
    }
}