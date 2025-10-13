// Minimum Path Sum
// p.s. uses dynamic programming

// Time Complexity: 
// 	O(m × n) - must visit every cell once

// Space Complexity:
//	O(1) (modifies input)

package main

func minPathSum(grid [][]int) int {
    if len(grid) == 0 || len(grid[0]) == 0 {
        return 0
    }
    
    m, n := len(grid), len(grid[0])
    
    // Calculate sums along first row (т.к. starts From Left) (to -> right)
    for j := 1; j < n; j++ {
        grid[0][j] += grid[0][j-1]
    }
    
    // Calculate sums along first column (т.к. starts From Top) (to -> bottom)
    for i := 1; i < m; i++ {
        grid[i][0] += grid[i-1][0]
    }
    
    // for each remaining cell: choose min sum из adjacent cells
    for i := 1; i < m; i++ {
        for j := 1; j < n; j++ {
            grid[i][j] += min(grid[i-1][j], grid[i][j-1])
        }
    }
    
	// return last bottom right cell (contains min sum (of path: from top left cell to bottom right cell))
    return grid[m-1][n-1]
}










