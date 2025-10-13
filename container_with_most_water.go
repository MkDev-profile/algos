// Container With Most Water.
// using 2 pointers as "Opposite Ends"(left and right)

package main

func maxArea(height []int) int {
	// сначала площадь самых крайних столбцов, т.е. width of контейнера максимальная
    left, right := 0, len(height)-1
    maxArea := 0
    
    for left < right {
        // Calculate current area
        width := right - left
        currentHeight := min(height[left], height[right])
        currentArea := width * currentHeight
        
        // Update max area if current is larger
        if currentArea > maxArea {
            maxArea = currentArea
        }
        
        // Move the pointer with smaller height
		/*
		p.s. "greater pointer" (например справа) не имеет смысла смещать, т.к. тогда "smaller pointer" (например слева) остается, area выравнивается по smaller pointer-y, area = (prev_width - 1) * "smaller pointer"), т.е. area < prev_area. 
		*/
        if height[left] < height[right] {
            left++
        } else {
            right--
        }
    }
    
    return maxArea
}










