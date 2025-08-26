package sorting

import (
	"gin/models"
	"testing"
)

func TestMergeSort_Execute(t *testing.T) {
	mergeSort := NewMergeSort()
	tracker := models.NewStepTracker()

	tests := []struct {
		name     string
		input    []interface{}
		expected []interface{}
	}{
		{
			name:     "Empty array",
			input:    []interface{}{},
			expected: []interface{}{},
		},
		{
			name:     "Single element",
			input:    []interface{}{5},
			expected: []interface{}{5},
		},
		{
			name:     "Already sorted",
			input:    []interface{}{1, 2, 3, 4, 5},
			expected: []interface{}{1, 2, 3, 4, 5},
		},
		{
			name:     "Reverse sorted",
			input:    []interface{}{5, 4, 3, 2, 1},
			expected: []interface{}{1, 2, 3, 4, 5},
		},
		{
			name:     "Random order",
			input:    []interface{}{3, 1, 4, 1, 5, 9},
			expected: []interface{}{1, 1, 3, 4, 5, 9},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := mergeSort.Execute(tt.input, tracker)
			if err != nil {
				t.Errorf("MergeSort.Execute() error = %v", err)
				return
			}

			if tt.name == "Empty array" && len(result.([]interface{})) == 0 {
				return
			}

			resultArr := result.([]interface{})
			if len(resultArr) != len(tt.expected) {
				t.Errorf("MergeSort.Execute() result length = %v, expected %v", len(resultArr), len(tt.expected))
				return
			}

			for i, v := range resultArr {
				if v != tt.expected[i] {
					t.Errorf("MergeSort.Execute() result[%d] = %v, expected %v", i, v, tt.expected[i])
				}
			}
		})
	}
}

func TestInsertionSort_Execute(t *testing.T) {
	insertionSort := NewInsertionSort()
	tracker := models.NewStepTracker()

	tests := []struct {
		name     string
		input    []interface{}
		expected []interface{}
	}{
		{
			name:     "Empty array",
			input:    []interface{}{},
			expected: []interface{}{},
		},
		{
			name:     "Single element",
			input:    []interface{}{5},
			expected: []interface{}{5},
		},
		{
			name:     "Already sorted",
			input:    []interface{}{1, 2, 3, 4, 5},
			expected: []interface{}{1, 2, 3, 4, 5},
		},
		{
			name:     "Reverse sorted",
			input:    []interface{}{5, 4, 3, 2, 1},
			expected: []interface{}{1, 2, 3, 4, 5},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := insertionSort.Execute(tt.input, tracker)
			if err != nil {
				t.Errorf("InsertionSort.Execute() error = %v", err)
				return
			}

			if tt.name == "Empty array" && len(result.([]interface{})) == 0 {
				return
			}

			resultArr := result.([]interface{})
			if len(resultArr) != len(tt.expected) {
				t.Errorf("InsertionSort.Execute() result length = %v, expected %v", len(resultArr), len(tt.expected))
				return
			}

			for i, v := range resultArr {
				if v != tt.expected[i] {
					t.Errorf("InsertionSort.Execute() result[%d] = %v, expected %v", i, v, tt.expected[i])
				}
			}
		})
	}
}

func TestHeapSort_Execute(t *testing.T) {
	heapSort := NewHeapSort()
	tracker := models.NewStepTracker()

	tests := []struct {
		name     string
		input    []interface{}
		expected []interface{}
	}{
		{
			name:     "Empty array",
			input:    []interface{}{},
			expected: []interface{}{},
		},
		{
			name:     "Single element",
			input:    []interface{}{5},
			expected: []interface{}{5},
		},
		{
			name:     "Basic test",
			input:    []interface{}{4, 1, 3, 2, 16, 9, 10, 14, 8, 7},
			expected: []interface{}{1, 2, 3, 4, 7, 8, 9, 10, 14, 16},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := heapSort.Execute(tt.input, tracker)
			if err != nil {
				t.Errorf("HeapSort.Execute() error = %v", err)
				return
			}

			if tt.name == "Empty array" && len(result.([]interface{})) == 0 {
				return
			}

			resultArr := result.([]interface{})
			if len(resultArr) != len(tt.expected) {
				t.Errorf("HeapSort.Execute() result length = %v, expected %v", len(resultArr), len(tt.expected))
				return
			}

			for i, v := range resultArr {
				if v != tt.expected[i] {
					t.Errorf("HeapSort.Execute() result[%d] = %v, expected %v", i, v, tt.expected[i])
				}
			}
		})
	}
}