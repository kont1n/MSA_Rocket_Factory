package model

import (
	"testing"
)

func TestCategory_String(t *testing.T) {
	tests := []struct {
		name     string
		category Category
		expected string
	}{
		{
			name:     "ENGINE категория",
			category: ENGINE,
			expected: "ENGINE",
		},
		{
			name:     "FUEL категория",
			category: FUEL,
			expected: "FUEL",
		},
		{
			name:     "PORTHOLE категория",
			category: PORTHOLE,
			expected: "PORTHOLE",
		},
		{
			name:     "WING категория",
			category: WING,
			expected: "WING",
		},
		{
			name:     "UNKNOWN категория",
			category: UNKNOWN,
			expected: "UNKNOWN",
		},
		{
			name:     "Неизвестное значение",
			category: Category(999),
			expected: "UNKNOWN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.category.String()
			if result != tt.expected {
				t.Errorf("Category.String() = %v, expected %v", result, tt.expected)
			}
		})
	}
}

func TestCategoryName(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		expected string
	}{
		{
			name:     "ENGINE по значению 1",
			value:    1,
			expected: "ENGINE",
		},
		{
			name:     "FUEL по значению 2",
			value:    2,
			expected: "FUEL",
		},
		{
			name:     "PORTHOLE по значению 3",
			value:    3,
			expected: "PORTHOLE",
		},
		{
			name:     "WING по значению 4",
			value:    4,
			expected: "WING",
		},
		{
			name:     "UNKNOWN по значению 0",
			value:    0,
			expected: "UNKNOWN",
		},
		{
			name:     "Неизвестное значение",
			value:    999,
			expected: "UNKNOWN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CategoryName(tt.value)
			if result != tt.expected {
				t.Errorf("CategoryName(%d) = %v, expected %v", tt.value, result, tt.expected)
			}
		})
	}
}

func TestToCategory(t *testing.T) {
	tests := []struct {
		name     string
		value    int
		expected Category
	}{
		{
			name:     "Конвертация в ENGINE",
			value:    1,
			expected: ENGINE,
		},
		{
			name:     "Конвертация в FUEL",
			value:    2,
			expected: FUEL,
		},
		{
			name:     "Конвертация в PORTHOLE",
			value:    3,
			expected: PORTHOLE,
		},
		{
			name:     "Конвертация в WING",
			value:    4,
			expected: WING,
		},
		{
			name:     "Конвертация невалидного значения в UNKNOWN",
			value:    999,
			expected: UNKNOWN,
		},
		{
			name:     "Конвертация отрицательного значения в UNKNOWN",
			value:    -1,
			expected: UNKNOWN,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ToCategory(tt.value)
			if result != tt.expected {
				t.Errorf("ToCategory(%d) = %v, expected %v", tt.value, result, tt.expected)
			}
		})
	}
}

// Тест проверяющий соответствие констант их значениям
func TestCategoryConstants(t *testing.T) {
	expectedValues := map[Category]int{
		UNKNOWN:  0,
		ENGINE:   1,
		FUEL:     2,
		PORTHOLE: 3,
		WING:     4,
	}

	for category, expectedValue := range expectedValues {
		if int(category) != expectedValue {
			t.Errorf("Константа %s имеет значение %d, ожидалось %d",
				category.String(), int(category), expectedValue)
		}
	}
}
