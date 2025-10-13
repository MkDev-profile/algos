package main

// MergeSortRecursive performs the merge sort algorithm
func MergeSortRecursive(arr []int) []int {
    if len(arr) <= 1 {
        return arr
    }
    
    mid := len(arr) / 2
    left := MergeSortRecursive(arr[:mid])
    right := MergeSortRecursive(arr[mid:])
    
    return mergeInternalForRecursive(left, right)
}

// merge combines two sorted slices into one sorted slice
func mergeInternalForRecursive(left, right []int) []int {
    result := make([]int, 0, len(left)+len(right))
    i, j := 0, 0
    
    for i < len(left) && j < len(right) {
        if left[i] <= right[j] {
            result = append(result, left[i])
            i++
        } else {
            result = append(result, right[j])
            j++
        }
    }
    
    // Append remaining elements
	if i < len(left) {
    	result = append(result, left[i:]...)
	}
	if j < len(right) {
    	result = append(result, right[j:]...)
	}

    return result
}

/*


Time Complexity: n*(logN)
p.s. (number of levels) * (work per level).
explaining:
 - сколько раз делим массив пополам: logN
example: n=8, divideCounts=logN=log2(8)=~ 3
p.s. 
1st: 2 arrays по 4 items  ([****][****]), 
2nd: 4 arrays по 2 items  ([**][**][**][**]), 
3rd: 8 arrays по 1 item-y ([*][*][*][*][*][*][*][*]) 
 - сколько раз сравниваем 2 элемента (1 item from left sub-array and 1 item from right sub-array): n
p.s. 
1st level: compare 4 times (4 items of left array with 4 items of right array) =~ n
2nd level: compare 4 times (2 times per len(leftArr)) =~ n
3rd level: compare 4 times (4 times per len(leftArr)) =~ n
=> 3*4 (or more robust: 3*n) = 12 times (or more robust: 24) =~ n*2 (or more robust: n*3) =~ n


Space Complexity: O(n)
p.s.:
(Input space) + (AllocatedSpace inside func merge) + (Recursion Stack Space)
= O(n) + O(n) + O(log n)
= O(n)



*/





