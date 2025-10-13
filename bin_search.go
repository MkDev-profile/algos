// binary search

package main

import (
    "cmp"
    "fmt"
)

// BinarySearch returns the index of target in sorted slice arr
// Returns -1 if target is not found
func BinarySearch[T cmp.Ordered](arr []T, target T) int {
    if len(arr) == 0 {
        return -1
    }
    
    low, high := 0, len(arr)-1
    
    for low <= high {
        mid := low + (high-low)/2
        
        switch {
        case arr[mid] == target:
            return mid
        case arr[mid] < target:
            low = mid + 1
        default:
            high = mid - 1
        }
    }
    
    return -1
}

func Main_BinarySearch() {
    // Example usage
 	words := []string{"apple", "banana", "cherry", "date", "elderberry"}
    wordIndex := BinarySearch(words, "cherry")
    fmt.Printf("Found 'cherry' at index %d\n", wordIndex)
}















