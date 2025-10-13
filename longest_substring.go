// Longest Substring Without Repeating Characters.
// using "2 pointers" as "sliding window"

package main

import (
	"fmt"
	"log"
)

func lengthOfLongestSubstring(s string) int {

	// Stores last occurrence index of each character
	// чтобы проверять был ли уже currentChar в current substring-e (и заодно чтобы эта проверка была быстрой поэтому используется map (т.е. hash таблица))
    charIndex := make(map[byte]int) 

	// длина наиболее длинной подстроки (в которой нет дубликатов символов, т.е. каждый символ в подстроке содержится только один раз)
    maxLength := 0

	// начало (индекс) of текущей подстроки (в которой нету дубликатов) 
    left := 0
    
	// просто итерирование символ-за-символом по всему массиву
	// (т.е. заодно: right указывает на "последний символ of текущей подстроки") 
    for right := 0; right < len(s); right++ {

		// текущий символ (current char)
        char := s[right]

        // If character exists in map and its index is within current window
		// (т.е. если символ встретился повторно "в текущей подстроке") 
        if lastIndex, exists := charIndex[char]; exists && lastIndex >= left {
			// Move left pointer to right of last occurrence
			// т.е. тогда начинается "отсчет"(начало) новой подстроки (предыдущая подстрока закончилась, т.к. встретили дубликат of символа)
            left = lastIndex + 1 
        }
		log.Println("lastIndex = ", charIndex[char], ", left = ", left)
        
		// Update character's last occurrence
		// записали (или перезаписали=обновили) ("first") индекс of current char-a в текущей substring-e
        charIndex[char] = right 
		
		// вычисление текущей длины of текущей подстроки
        currentLength := right - left + 1

		// если (текущая) длина of текущей подстроки больше чем длина of (...пред...)предыдущей подстроки
        if currentLength > maxLength {
			// обновляем return-значение: of наиболее длинной найденной подстроки 
            maxLength = currentLength
        }
    }
    
    return maxLength
}

func Main_sliding_window() {
	r := lengthOfLongestSubstring("pwwkew")
    fmt.Printf("lengthOfLongestSubstring: %d\n", r)
}






