package strategies

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoundRobin_SendOrder(t *testing.T) {
	t.Run("with no backends", func(t *testing.T) {
		t.Parallel()

		rr := NewRoundRobin([]string{})

		assert.Empty(t, rr.SendOrder())
		assert.Empty(t, rr.SendOrder())
	})
	t.Run("with 1 backend", func(t *testing.T) {
		t.Parallel()

		backends := []string{"A"}

		rr := NewRoundRobin(backends)

		assert.Equal(t, backends, rr.SendOrder())
		assert.Equal(t, backends, rr.SendOrder())
	})
	t.Run("with more backends", func(t *testing.T) {
		t.Parallel()

		backends := []string{"A", "B", "C"}

		rr := NewRoundRobin(backends)

		assert.Equal(t, []string{"A", "B", "C"}, rr.SendOrder())
		assert.Equal(t, []string{"B", "C", "A"}, rr.SendOrder())
		assert.Equal(t, []string{"C", "A", "B"}, rr.SendOrder())
		assert.Equal(t, []string{"A", "B", "C"}, rr.SendOrder())
	})
}
