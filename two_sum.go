// Two Sum.
// using 2 pointers as "Opposite Ends"(left and right)

package main

import "fmt"

func twoSum_forSortedArray(numbers []int, target int) []int {
	left, right := 0, len(numbers)-1

	for left < right {
		sum := numbers[left] + numbers[right]

		if sum == target {
			return []int{left, right}
		} else if sum < target {
			left++
		} else {
			right--
		}
	}

	return nil
}

func twoSum_forUnsortedArray(nums []int, target int) []int {
    numMap := make(map[int]int)
    
    for i, num := range nums {
		// для num удовлетворяющее число = target - num
        dif := target - num

		// если такое число существует в map-e то возвращаем результат
        if idx, exists := numMap[dif]; exists {
            return []int{idx, i}
        }

		// иначе добавляем в map-у: [num]=[индекс of num в массиве]
        numMap[num] = i
    }
    return nil
}

func Main_twosum() {
	nums := []int{2, 7, 11, 15}
	target := 9
	result := twoSum_forSortedArray(nums, target)
	fmt.Println("Two Sum indices:", result) // Output: [0 1]

	nums = []int{2, 15, 11, 7}
	target = 9
	result = twoSum_forUnsortedArray(nums, target)
	fmt.Println("Two Sum indices:", result) // Output: [0 3]
}








