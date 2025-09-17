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
	//负数、结尾为0的非零数不是
    if x < 0 || (x % 10 == 0 && x != 0) {
        return false
    }
    
    reversedHalf := 0
    original := x
    // 计算回文数
    for x > 0 {
        reversedHalf = reversedHalf * 10 + x % 10
        x /= 10
    }
    
    return original == reversedHalf
}



/*
*

	给定一个只包括 '('，')'，'{'，'}'，'['，']' 的字符串，判断字符串是否有效

*
*/
func isValid(s string) bool {
    // 使用栈来跟踪开括号
    stack := make([]rune, 0)
    
    // 创建映射表，用于快速查找匹配的括号
    mapping := map[rune]rune{
        ')': '(',
        '}': '{',
        ']': '[',
    }
    
    // 遍历字符串中的每个字符
    for _, char := range s {
        // 如果是闭括号
        if matchingOpen, isClose := mapping[char]; isClose {
            // 检查栈是否为空或栈顶元素是否匹配
            if len(stack) == 0 || stack[len(stack)-1] != matchingOpen {
                return false
            }
            // 弹出栈顶元素
            stack = stack[:len(stack)-1]
        } else {
            // 是开括号，压入栈中
            stack = append(stack, char)
        }
    }
    
    // 如果栈为空，说明所有括号都正确匹配
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
    
    prefix := strs[0]
    
    for i := 1; i < len(strs); i++ {
        // 逐步缩短前缀，直到匹配当前字符串
        for len(prefix) > 0 {
            // 检查当前字符串长度是否足够
            if len(strs[i]) < len(prefix) {
                prefix = prefix[:len(prefix)-1]
                continue
            }
            // 比较前缀部分
            if strs[i][:len(prefix)] != prefix {
                prefix = prefix[:len(prefix)-1]
            } else {
                break
            }
        }
        
        // 如果前缀为空，提前返回
        if prefix == "" {
            return ""
        }
    }
    
    return prefix
}

/*
*

	给定一个由整数组成的非空数组所表示的非负整数，在该数的基础上加一

*
*/
func plusOne(digits []int) []int {
    for i := len(digits) - 1; i >= 0; i-- {
        if digits[i]==9{
            digits[i] = 0
        }else{
            digits[i] += 1
            return digits
        }
    }
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
    slow := 0
    for  fast:=1;fast<len(nums);fast++{
        if nums[fast]>nums[slow]{
            slow++
            nums[slow]=nums[fast]
        }
    }
    slow+=1
    return slow
}

/*
*

	合并区间：以数组 intervals 表示若干个区间的集合，其中单个区间为 intervals[i] = [starti, endi] 。请你合并所有重叠的区间，并返回一个不重叠的
	区间数组，该数组需恰好覆盖输入中的所有区间。可以先对区间数组按照区间的起始位置进行排序，然后使用一个切片来存储合并后的区间，遍历排序后的区间数组，
	将当前区间与切片中最后一个区间进行比较，如果有重叠，则合并区间；如果没有重叠，则将当前区间添加到切片中

*
*/

func merge(intervals [][]int) [][]int {
    if len(intervals) == 0 {
        return [][]int{}
    }
    
    // 自己实现快速排序
    quickSort(intervals, 0, len(intervals)-1)
    
    merged := [][]int{}
    for _, interval := range intervals {
        if len(merged) == 0 || interval[0] > merged[len(merged)-1][1] {
            merged = append(merged, interval)
        } else {
            if interval[1] > merged[len(merged)-1][1] {
                merged[len(merged)-1][1] = interval[1]
            }
        }
    }
    
    return merged
}

func quickSort(arr [][]int, left, right int) {
    if left >= right {
        return
    }
    
    pivotIndex := partition(arr, left, right)
    quickSort(arr, left, pivotIndex-1)
    quickSort(arr, pivotIndex+1, right)
}

func partition(arr [][]int, left, right int) int {
    pivot := arr[right]
    i := left - 1
    
    for j := left; j < right; j++ {
        if arr[j][0] < pivot[0] || (arr[j][0] == pivot[0] && arr[j][1] < pivot[1]) {
            i++
            arr[i], arr[j] = arr[j], arr[i]
        }
    }
    
    arr[i+1], arr[right] = arr[right], arr[i+1]
    return i + 1
}

/**

	给定一个整数数组 nums 和一个目标值 target，请你在该数组中找出和为目标值的那两个整数

**/

func twoSum(nums []int, target int) []int {
    hashmap := make(map[int]int)
    
    for index, num := range nums {
        anotherNum := target - num
        if anotherIndex, exists := hashmap[anotherNum]; exists {
            return []int{anotherIndex, index}
        }
        hashmap[num] = index
    }
    return nil
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


