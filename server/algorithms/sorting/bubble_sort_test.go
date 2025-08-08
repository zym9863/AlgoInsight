package sorting

import (
	"fmt"
	"gin/models"
	"testing"
)

func TestBubbleSort_Execute(t *testing.T) {
	bs := NewBubbleSort()
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
			input:    []interface{}{3, 1, 4, 1, 5, 9, 2, 6},
			expected: []interface{}{1, 1, 2, 3, 4, 5, 6, 9},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := bs.Execute(tt.input, tracker)
			if err != nil {
				t.Errorf("Execute() error = %v", err)
				return
			}

			resultArray, ok := result.([]interface{})
			if !ok {
				t.Errorf("Execute() result is not []interface{}")
				return
			}

			if len(resultArray) != len(tt.expected) {
				t.Errorf("Execute() result length = %v, expected %v", len(resultArray), len(tt.expected))
				return
			}

			for i, v := range resultArray {
				if v != tt.expected[i] {
					t.Errorf("Execute() result[%d] = %v, expected %v", i, v, tt.expected[i])
				}
			}
		})
	}
}

func TestBubbleSort_ValidateInput(t *testing.T) {
	bs := NewBubbleSort()

	tests := []struct {
		name    string
		input   interface{}
		wantErr bool
	}{
		{
			name:    "Valid array",
			input:   []interface{}{1, 2, 3},
			wantErr: false,
		},
		{
			name:    "Nil input",
			input:   nil,
			wantErr: true,
		},
		{
			name:    "Non-array input",
			input:   "not an array",
			wantErr: true,
		},
		{
			name:    "Too large array",
			input:   make([]interface{}, 20000),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := bs.ValidateInput(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateInput() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestBubbleSort_Properties(t *testing.T) {
	bs := NewBubbleSort()

	if !bs.IsStable() {
		t.Error("BubbleSort should be stable")
	}

	if !bs.IsInPlace() {
		t.Error("BubbleSort should be in-place")
	}

	if !bs.IsAdaptive() {
		t.Error("BubbleSort should be adaptive")
	}

	info := bs.GetInfo()
	if info.ID != "bubble_sort" {
		t.Errorf("Expected ID 'bubble_sort', got '%s'", info.ID)
	}

	if info.Name != "冒泡排序" {
		t.Errorf("Expected name '冒泡排序', got '%s'", info.Name)
	}
}

func BenchmarkBubbleSort(b *testing.B) {
	bs := NewBubbleSort()
	tracker := models.NewStepTracker()

	sizes := []int{10, 100, 1000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("size_%d", size), func(b *testing.B) {
			// 生成测试数据
			data := make([]interface{}, size)
			for i := 0; i < size; i++ {
				data[i] = size - i // 逆序数据（最坏情况）
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				// 复制数据以避免影响下次测试
				testData := make([]interface{}, len(data))
				copy(testData, data)

				_, err := bs.Execute(testData, tracker)
				if err != nil {
					b.Fatalf("Execute failed: %v", err)
				}
			}
		})
	}
}
