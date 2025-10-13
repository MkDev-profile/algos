package main

import "fmt"

func RemoveDuplicatesFromSortedArray(nums []int) []int {
    if len(nums) == 0 {
        return nil
    }
    
    slow := 0
    for fast := 1; fast < len(nums); fast++ {
        // если дубликат, то просто дальше передвигаем fast (а slow нет)
        if nums[fast] == nums[slow] {
            continue
        }

        // set-им недубликатное значение на slow++ позицию 
        slow++
        nums[slow] = nums[fast]
    }
    
    // slow это последний индекс of "недубликатной" части массива
    // возвращаем slice до позиции slow
    return nums[:(slow + 1)]
}

func Main_removedublicates() {
    nums := []int{0, 0, 1, 1, 1, 2, 2, 3, 3, 4}
    unique_nums := RemoveDuplicatesFromSortedArray(nums)
    fmt.Println("Unique elements:", unique_nums) // Output: [0 1 2 3 4]
}






