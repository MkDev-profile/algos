package main

import "fmt"

func threeSum_forSortedArray(nums []int) [][]int {
	result := [][]int{}
    n := len(nums)
    
	// p.s. До (n-2) т.к. для "последнего item-а" в цикле нужно найти 2 правых item-a (которые будут с ним в тройке) 
    for i := 0; i < n-2; i++ { 

        // Skip duplicates for "current" item
        // (иначе для того же current-a будут находиться те же left и right значения)
        if i > 0 && nums[i] == nums[i-1] {
            continue
        }

		// p.s. left это второе слагаемое (в тройке), 
		// right это третье слагаемое (в тройке),  
		// nums[i] это первое слагаемое (в тройке)
        left, right := i+1, n-1

		// ищем таких left, right которые: 
		// left + right= -(nums[i]) (т.е. равны отрицательному nums[i])
        target := -nums[i]
        
        for left < right {
            sum := nums[left] + nums[right]
            
            if sum == target {
                result = append(result, []int{nums[i], nums[left], nums[right]})
                
                // Skip duplicates for left pointer
                // (иначе для того же current-a и left-a найдется тот же right value)
                for left < right && nums[left] == nums[left+1] {
                    left++
                }

                // Skip duplicates for right pointer
                // (аналогично)
                for left < right && nums[right] == nums[right-1] {
                    right--
                }

                // перемещаем дальше, чтобы искать следующую тройку для current num-a
                left++
                right--
            } else if sum < target {
                left++
            } else {
                right--
            }
        }
    }
    
    return result
}

func Main_three_sum() {
	testCases := [][]int{
		{-3, -2, 1, 1, 1, 2},
 		{-4, -1, -1, 0, 1, 2},
		{0, 1, 1},
		{0, 0, 0},
		{-3, -3, 2, 2, 1, 1, 0, 0},
		{-2, 0, 1, 1, 2}, 
	}

	for _, nums := range testCases {
		fmt.Printf("Input: %v\n", nums)
		result := threeSum_forSortedArray(nums)
		fmt.Printf("Output: %v\n\n", result)
	}

    /*
    
    output:
    
Input: [-3 -2 1 1 1 2]
Output: [[-3 1 2] [-2 1 1]]

Input: [-4 -1 -1 0 1 2]
Output: [[-1 -1 2] [-1 0 1]]

Input: [0 1 1]
Output: []

Input: [0 0 0]
Output: [[0 0 0]]

Input: [-3 -3 2 2 1 1 0 0]
Output: []

Input: [-2 0 1 1 2]
Output: [[-2 0 2] [-2 1 1]]
    
    */
}

/*

Формулировка задачи:
найти все уникальные тройки [nums[i], nums[j], nums[k]] такие, что:

nums[i] + nums[j] + nums[k] == 0

И:
i != j, i != k, j != k
(то есть индексы не повторяются, то есть каждый элемент массива содержится только в одной тройке)

И:
все тройки должны быть уникальные
(то есть три значения в одной тройке не должны совпадать с тремя значениями в другой тройке (независимо от расположения значений внутри троек),
то есть нельзя: [1,2,3],[2,1,3],[3,1,2], т.д.)

*/





