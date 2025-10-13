package main

func isValidParentheses(s string) bool {
    stack := []rune{}

    pairs := map[rune]rune{
        ')': '(',
        '}': '{',
        ']': '[',
    }
    
    for _, char := range s {
        // если opened-скобка то кладем ее в стек
        if char == '(' || char == '{' || char == '[' {
            stack = append(stack, char)
        } else {
            // иначе если closed-скобка, то находим в map-e соответствующую ей opened-скобку

            // если в стеке (в верхушке) не находится эта opened-скобка то return false
            if len(stack) == 0 || pairs[char] != stack[len(stack)-1] {
                return false
            }

            // иначе если в стеке (в верхушке) эта opened-скобка, то выполняем "pop" (то есть убираем top-item из стека)  
            stack = stack[:len(stack)-1]
        }
    }
    
    // return: для всех ли open-скобок были closed-скобки
    return len(stack) == 0
}

// Пример использования
func Main_isValidParentheses() {
    println(isValidParentheses("()[]{}"))  // true
    println(isValidParentheses("([)]"))    // false
}
