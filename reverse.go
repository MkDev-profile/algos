package main

import "fmt"

func reverseString(s []byte) {
	left, right := 0, len(s)-1
	
	for left < right {
		s[left], s[right] = s[right], s[left]
		left++
		right--
	}
}

func Main_reverse() {
	s := []byte("hello")
	reverseString(s)
	fmt.Println("Reversed:", string(s)) // Output: "olleh"
}
