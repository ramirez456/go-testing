package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddSuccess(t *testing.T) {
	result := Add(1, 2)
	if result != 3 {
		t.Errorf("Add(1, 2) = %d; want 3", result)
	}
}

func TestAddSuccessWithTestify(t *testing.T) {
	c := require.New(t)
	result := Add(1, 2)
	c.Equal(3, result)
}
