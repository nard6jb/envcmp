package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubtract_NoOptions(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2", "C": "3"}
	right := map[string]string{"B": "2", "C": "3"}

	result := Subtract(left, right, SubtractOptions{})

	assert.Len(t, result, 1)
	assert.Equal(t, "A", result[0].Key)
	assert.Equal(t, "1", result[0].Value)
}

func TestSubtract_AllKeysPresent(t *testing.T) {
	left := map[string]string{"X": "10", "Y": "20"}
	right := map[string]string{"X": "10", "Y": "20"}

	result := Subtract(left, right, SubtractOptions{})

	assert.Empty(t, result)
}

func TestSubtract_EmptyRight(t *testing.T) {
	left := map[string]string{"A": "1", "B": "2"}
	right := map[string]string{}

	result := Subtract(left, right, SubtractOptions{})

	assert.Len(t, result, 2)
	assert.Equal(t, "A", result[0].Key)
	assert.Equal(t, "B", result[1].Key)
}

func TestSubtract_EmptyLeft(t *testing.T) {
	left := map[string]string{}
	right := map[string]string{"A": "1"}

	result := Subtract(left, right, SubtractOptions{})

	assert.Empty(t, result)
}

func TestSubtract_MaskSecrets(t *testing.T) {
	left := map[string]string{
		"API_KEY":  "super-secret",
		"APP_NAME": "myapp",
		"PASSWORD": "hunter2",
	}
	right := map[string]string{}

	result := Subtract(left, right, SubtractOptions{MaskSecrets: true})

	byKey := make(map[string]string)
	for _, e := range result {
		byKey[e.Key] = e.Value
	}

	assert.Equal(t, "***", byKey["API_KEY"])
	assert.Equal(t, "***", byKey["PASSWORD"])
	assert.Equal(t, "myapp", byKey["APP_NAME"])
}

func TestSubtract_SortedOutput(t *testing.T) {
	left := map[string]string{"Z": "z", "A": "a", "M": "m"}
	right := map[string]string{}

	result := Subtract(left, right, SubtractOptions{})

	keys := make([]string, len(result))
	for i, e := range result {
		keys[i] = e.Key
	}
	assert.Equal(t, []string{"A", "M", "Z"}, keys)
}
