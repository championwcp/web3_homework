package main

import (
	"fmt"
	"sort"
	"strconv"
)

/*
*

	136只出现一次的数字：给定一个非空整数数组，除了某个元素只出现一次以外，其余每个元素均出现两次。找出那个只出现了一次的元素。 可以使用 for 循环遍历数组，
	结合 if 条件判断和 map 数据结构来解决，例如通过 map 记录每个元素出现的次数，然后再遍历 map 找到出现次数为1的元素

*
*/
func singleNumber(nums []int) int {
    result := 0
    for _, num := range nums {
        result ^= num
    }
    return result
}

/*
*

	判断一个整数是否是回文数 （正序（从左向右）和倒序（从右向左）读都是一样的整数）

*
*/
func isPalindrome(x int) bool {
	if x < 0 {
		return false
	}
	// 转换为字符串
	s := strconv.Itoa(x)
	begin, end := 0, len(s)-1
	for begin < end {
		if s[begin] != s[end] {
			return false
		}
		begin++
		end--
	}
	return true
}

/*
*

	给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串，判断字符串是否有效

*
*/
func isValid(s string) bool {
	stack := []rune{} // 使用 rune 切片作为栈
	mapping := map[rune]rune{
		')': '(',
		'}': '{',
		']': '[',
	}

	for _, char := range s {
		if char == '(' || char == '{' || char == '[' {
			stack = append(stack, char) // 左括号入栈
		} else {
			// 栈为空但遇到右括号
			if len(stack) == 0 {
				return false
			}
			top := stack[len(stack)-1]
			if mapping[char] != top {
				return false // 栈顶括号与当前右括号不匹配
			}
			// 匹配成功，删除第一个元素
			stack = stack[:len(stack)-1]
		}
	}
	// 栈为空则所有括号匹配
	return len(stack) == 0
}

/*
*

	查找字符串数组中的最长公共前缀

*
*/
func longestCommonPrefix(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	// 数组仅一个元素，直接返回该元素
	if len(strs) == 1 {
		return strs[0]
	}
	// 以第一个元素为基准
	prefix := strs[0]

	for i := 1; i < len(strs); i++ {
		// 比较当前前缀与第i个字符串，找出共同前缀
		j := 0
		for j < len(prefix) && j < len(strs[i]) && prefix[j] == strs[i][j] {
			j++
		}
		// 更新前缀为共同部分
		prefix = prefix[:j]
		// 如果前缀为空，提前退出
		if prefix == "" {
			break
		}
	}

	return prefix
}

/*
*

	给定一个由整数组成的非空数组所表示的非负整数，在该数的基础上加一

*
*/
func addOne(digits []int) []int {
	// 从右到左遍历数组
	for i := len(digits) - 1; i >= 0; i-- {
		// 如果当前位不是9，直接加1并返回
		if digits[i] < 9 {
			digits[i]++
			return digits
		}
		// 如果当前位是9，进位，当前位变为0，继续遍历数组
		digits[i] = 0
	}

	// 如果所有位都是9，需要在最前面添加一个1
	// 创建一个新的数组，长度为原数组长度+1
	result := make([]int, len(digits)+1)
	result[0] = 1
	return result
}

/*
*

	删除有序数组中的重复项：给你一个有序数组 nums ，请你原地删除重复出现的元素，使每个元素只出现一次，返回删除后数组的新长度。
	不要使用额外的数组空间，你必须在原地修改输入数组并在使用 O(1) 额外空间的条件下完成。可以使用双指针法，一个慢指针 i 用于记录
	不重复元素的位置，一个快指针 j 用于遍历数组，当 nums[i] 与 nums[j] 不相等时，将 nums[j] 赋值给 nums[i + 1]，并将 i 后移一位。

*
*/
func removeDuplicates(nums []int) int {
	if len(nums) == 0 {
		return 0
	}
	// 慢指针，记录不重复元素的位置
	i := 0
	// 快指针遍历数组
	for j := 1; j < len(nums); j++ {
		if nums[j] != nums[i] {
			i++
			// 将不重复元素前移
			nums[i] = nums[j]
		}
	}
	// 新数组长度为i+1
	return i + 1
}

/*
*

	合并区间：以数组 intervals 表示若干个区间的集合，其中单个区间为 intervals[i] = [starti, endi] 。请你合并所有重叠的区间，并返回一个不重叠的
	区间数组，该数组需恰好覆盖输入中的所有区间。可以先对区间数组按照区间的起始位置进行排序，然后使用一个切片来存储合并后的区间，遍历排序后的区间数组，
	将当前区间与切片中最后一个区间进行比较，如果有重叠，则合并区间；如果没有重叠，则将当前区间添加到切片中

*
*/
func merge(intervals [][]int) [][]int {
	// 如果区间为空，直接返回
	if len(intervals) == 0 {
		return nil
	}

	// 按照区间的起始位置进行排序
	sort.Slice(intervals, func(i, j int) bool {
		return intervals[i][0] < intervals[j][0]
	})

	// 用于存储合并后的区间
	result := [][]int{intervals[0]}

	// 遍历排序后的区间
	for i := 1; i < len(intervals); i++ {
		// 获取结果集中的最后一个区间
		last := result[len(result)-1]
		// 当前区间
		current := intervals[i]

		// 如果当前区间的起始位置小于等于结果集中最后一个区间的结束位置，说明有重叠
		if current[0] <= last[1] {
			// 合并区间，取两个区间结束位置的最大值
			if current[1] > last[1] {
				last[1] = current[1]
			}
		} else {
			// 没有重叠，直接添加到结果集
			result = append(result, current)
		}
	}

	return result
}

/**

	给定一个整数数组 nums 和一个目标值 target，请你在该数组中找出和为目标值的那两个整数

**/

func twoSum(nums []int, target int) []int {
	numMap := make(map[int]int)
	for i, num := range nums {
		complement := target - num
		if j, exists := numMap[complement]; exists {
			return []int{nums[j], nums[i]}
		}
		numMap[num] = i
	}
	return []int{}
}

func main() {
	fmt.Println("=== 不重复的元素 ===")
	nums1 := [10]int{1, 1, 1, 2, 2, 3, 4, 4, 5, 5}
	result1 := oneNumbercheck(nums1)
	fmt.Printf("不重复的元素：%v\n", result1)

	// 回文数判断
	fmt.Println()
	fmt.Println("=== 回文数判断 ===")
	fmt.Println(isPalindrome(121))  // true
	fmt.Println(isPalindrome(-121)) // false
	fmt.Println(isPalindrome(1231)) // false

	// 括号判断
	fmt.Println()
	fmt.Println("=== 括号判断 ===")
	s1 := "()[]{}"
	result := isValid(s1)
	fmt.Printf("输入：s = \"%s\"\n", s1)
	fmt.Printf("输出：%t\n", result)

	s2 := "(]"
	result2 := isValid(s2)
	fmt.Printf("输入：s = \"%s\"\n", s2)
	fmt.Printf("输出：%t\n", result2)

	// 最长前缀
	fmt.Println()
	fmt.Println("=== 最长前缀 ===")
	strs1 := []string{"abcdef", "abbgdc", "abpokm"}
	result3 := longestCommonPrefix(strs1)
	fmt.Printf("输入：strs = %v\n", strs1)
	fmt.Printf("输出：\"%s\"\n", result3)

	strs2 := []string{"1", "2", "23"}
	result4 := longestCommonPrefix(strs2)
	fmt.Printf("输入：strs = %v\n", strs2)
	fmt.Printf("输出：\"%s\"\n", result4)

	// 加一
	fmt.Println()
	fmt.Println("=== 加一 ===")
	digits1 := []int{1, 2, 3, 4}
	fmt.Printf("输入：digits = %v\n", digits1)
	result5 := addOne(digits1)
	fmt.Printf("输出：%v\n", result5)

	digits2 := []int{1, 1, 0, 9}
	fmt.Printf("输入：digits = %v\n", digits2)
	result6 := addOne(digits2)
	fmt.Printf("输出：%v\n", result6)

	// removeDuplicates
	fmt.Println()
	fmt.Println("=== removeDuplicates ===")
	nums2 := []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}
	fmt.Printf("输入：nums = %v\n", nums2)
	k1 := removeDuplicates(nums2)
	fmt.Printf("输出：%d, nums = %v\n", k1, nums2[:k1])

	// merge
	fmt.Println()
	fmt.Println("=== merge ===")
	intervals1 := [][]int{{1, 3}, {2, 6}, {8, 10}, {15, 18}}
	fmt.Printf("输入：intervals = %v\n", intervals1)
	result7 := merge(intervals1)
	fmt.Printf("输出：%v\n", result7)

	// twoSum
	fmt.Println()
	fmt.Println("=== twoSum两数之和 ===")
	nums3 := []int{1, 3, 7, 9}
	target1 := 8
	fmt.Printf("输入：数组 = %v\n", nums3)
	fmt.Printf("输入：target = %v\n", target1)
	result8 := twoSum(nums3, target1)
	fmt.Printf("输出：%v\n", result8)

}
