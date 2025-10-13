package main

import "fmt"

func isPalindrome(s string) bool {
	left, right := 0, len(s)-1

	for left < right {
		if s[left] != s[right] {
			return false
		}

		left++
		right--
	}

	return true
}

func isPalindromeAlphaNum(s string) bool {
	left, right := 0, len(s)-1

	for left < right {

		// Skip non-alphanumeric characters слева
		for left < right && !isAlphanumeric(s[left]) {
			left++
		}

		// Skip non-alphanumeric characters справа
		for left < right && !isAlphanumeric(s[right]) {
			right--
		}

		if toLower(s[left]) != toLower(s[right]) {
			return false
		}

		left++
		right--
	}

	return true
}

func isAlphanumeric(c byte) bool {
	return (c >= 'a' && c <= 'z') || 
		(c >= 'A' && c <= 'Z') || 
		(c >= '0' && c <= '9')
}

func toLower(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c + 32
	}

	return c
}

func Main_palindrome() {
	fmt.Println("Palindrome:", isPalindrome("racecar"))                                      // true
	
	fmt.Println("Valid Palindrome:", isPalindromeAlphaNum("A man, a plan, a canal: Panama")) // true
}
