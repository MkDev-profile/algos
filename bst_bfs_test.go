package main

import (
    "reflect"
    "testing"
)

func TestBST_BFS(t *testing.T) {
    bst := NewBST()
    values := []int{5, 3, 7, 2, 4, 6, 8}
    for _, value := range values {
        bst.Insert(value)
    }

    expected := []int{5, 3, 7, 2, 4, 6, 8}
    result := bst.BFS()

    if !reflect.DeepEqual(result, expected) {
        t.Errorf("BFS() = %v, expected %v", result, expected)
    }
}

func TestBST_BFSWithLevels(t *testing.T) {
    bst := NewBST()
    values := []int{5, 3, 7, 2, 4, 6, 8}
    for _, value := range values {
        bst.Insert(value)
    }

    expected := [][]int{
        {5},
        {3, 7},
        {2, 4, 6, 8},
    }
    result := bst.BFSWithLevels()

    if !reflect.DeepEqual(result, expected) {
        t.Errorf("BFSWithLevels() = %v, expected %v", result, expected)
    }
}






















