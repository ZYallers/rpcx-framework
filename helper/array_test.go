package helper

import (
	"testing"
)

func TestRemoveDuplicateWithInt(t *testing.T) {
	ss := []int{1, 2, 2, 3, 4}
	t.Log(RemoveDuplicateWithInt(ss))
}

func TestRemoveDuplicateWithString(t *testing.T) {
	ss := []string{"a", "b", "b", "c", "bb"}
	t.Log(RemoveDuplicateWithString(ss))
}

func TestRemoveWithString(t *testing.T) {
	ss := []string{"a", "b", "c"}
	t.Log(RemoveWithString(ss, "bb"))
	t.Log(RemoveWithString(ss, "b"))
}
