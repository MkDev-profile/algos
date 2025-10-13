package main

import "testing"

func TestBinarySearch(t *testing.T) {
	tests := []struct {
		name   string
		arr    []int
		target int
		want   int
	}{
		{
			name:   "empty array",
			arr:    []int{},
			target: 5,
			want:   -1,
		},
		{
			name:   "Single element. Value exists",
			arr:    []int{5},
			target: 5,
			want:   0,
		},
		{
			name:   "Single element. Value doesn't exist",
			arr:    []int{3},
			target: 5,
			want:   -1,
		},
		{
			name:   "Multiple elements. Value exists at beginning",
			arr:    []int{1, 3, 5, 7, 9},
			target: 1,
			want:   0,
		},
		{
			name:   "Multiple elements. Value exists at end",
			arr:    []int{1, 3, 5, 7, 9},
			target: 9,
			want:   4,
		},
		{
			name:   "Multiple elements. Value exists in middle",
			arr:    []int{1, 3, 5, 7, 9},
			target: 5,
			want:   2,
		},
		{
			name:   "Multiple elements. Value doesn't exist",
			arr:    []int{1, 3, 5, 7, 9},
			target: 4,
			want:   -1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BinarySearch(tt.arr, tt.target); got != tt.want {
				t.Errorf("BinarySearch() = %v, want %v", got, tt.want)
			}
		})
	}
}
