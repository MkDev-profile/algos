package main

// MergeSort performs iterative merge sort on the input slice
func MergeSortIterative(arr []int) []int {
    if len(arr) <= 1 {
        return arr
    }
    
    n := len(arr)
    // Create a temporary array to avoid modifying the original during merging
    temp := make([]int, n)
    copy(temp, arr)
    
    // Start with subarrays of size 1 and double each time
    for width := 1; width < n; width *= 2 { // обход level-ов
        for i := 0; i < n; i += 2 * width { // обход sub-массивов of current level-a
            // Find the boundaries for the current subarrays
            left := i
            mid := min(i+width, n)
            right := min(i+2*width, n)
            
            // Merge the two subarrays
            mergeInternalForIterative(temp, left, mid, right)
        }
        
        // Copy the partially sorted array back for the next iteration
        copy(arr, temp)
    }
    
    return arr
}

// merge combines two sorted subarrays [left, mid) and [mid, right)
func mergeInternalForIterative(arr []int, left, mid, right int) {
    i, j := left, mid
    
    // Create a temporary slice for the merged result
    temp := make([]int, right-left)
    idx := 0
    
    // Merge the two subarrays
    for i < mid && j < right {
        if arr[i] <= arr[j] {
            temp[idx] = arr[i]
            i++
        } else {
            temp[idx] = arr[j]
            j++
        }
        idx++
    }
    
    // Copy remaining elements from left subarray
    for i < mid {
        temp[idx] = arr[i]
        i++
        idx++
    }
    
    // Copy remaining elements from right subarray
    for j < right {
        temp[idx] = arr[j]
        j++
        idx++
    }
    
    // Copy the merged result back to the original array
    copy(arr[left:right], temp)
}

// min returns the minimum of two integers
func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}











