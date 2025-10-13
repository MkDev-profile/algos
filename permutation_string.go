package main

import "fmt"

func checkInclusion(s1 string, s2 string) bool {
    if len(s1) > len(s2) {
        return false
    }
    
    // Frequency maps for s1 and current window in s2
    s1Count := [26]int{}
    windowCount := [26]int{}
    
    // Initialize frequency maps for first window
    for i := 0; i < len(s1); i++ {
		fmt.Println(s1[i])
		fmt.Println(s1[i]-'a')
        s1Count[s1[i]-'a']++
        windowCount[s2[i]-'a']++
    }

	fmt.Println(s1Count)
	fmt.Println(windowCount)
    
    // Check if first window is a permutation
    if s1Count == windowCount {
		fmt.Println("equals")
        return true
    }
    
    // Slide the window through s2
    for i := len(s1); i < len(s2); i++ {
        // Add new character to window
        windowCount[s2[i]-'a']++
        // Remove leftmost character from window
        windowCount[s2[i-len(s1)]-'a']--
        
		fmt.Println(windowCount)

        // Check if current window is a permutation
        if s1Count == windowCount {
			fmt.Println("equals on i=", i)
            return true
        }
    }
    
    return false
}

func Main_permutation_string() {
	checkInclusion("ab", "eidbaooo")
}

/*

go run .
97
0
98
1
[1 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
[0 0 0 0 1 0 0 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
[0 0 0 1 0 0 0 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
[0 1 0 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
[1 1 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0 0]
equals on i= 4

*/










