package strategies

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRoundRobin(t *testing.T) {
	t.Run("with no backends", func(t *testing.T) {
		t.Parallel()

		rr := NewRoundRobin([]string{})

		assert.Equal(t, "", rr.ChooseBackend())
		assert.Equal(t, "", rr.ChooseBackend())

		rr.UpdateBackendHealth("A", true)

		assert.Equal(t, "", rr.ChooseBackend())
	})
	t.Run("with 1 backend", func(t *testing.T) {
		t.Parallel()

		backends := []string{"A"}

		rr := NewRoundRobin(backends)

		assert.Equal(t, "", rr.ChooseBackend())
		assert.Equal(t, "", rr.ChooseBackend())

		rr.UpdateBackendHealth("A", true)

		assert.Equal(t, "A", rr.ChooseBackend())
		assert.Equal(t, "A", rr.ChooseBackend())

		rr.UpdateBackendHealth("B", true)

		assert.Equal(t, "A", rr.ChooseBackend())

		rr.UpdateBackendHealth("A", false)

		assert.Equal(t, "", rr.ChooseBackend())
	})
	t.Run("with more backends", func(t *testing.T) {
		t.Parallel()

		backends := []string{"A", "B", "C"}

		rr := NewRoundRobin(backends)

		assert.Equal(t, "", rr.ChooseBackend())

		rr.UpdateBackendHealth("A", true)

		assert.Equal(t, "A", rr.ChooseBackend())
		assert.Equal(t, "A", rr.ChooseBackend())

		rr.UpdateBackendHealth("B", true)

		assert.Equal(t, "A", rr.ChooseBackend())
		assert.Equal(t, "B", rr.ChooseBackend())

		rr.UpdateBackendHealth("C", true)

		assert.Equal(t, "C", rr.ChooseBackend())
		assert.Equal(t, "A", rr.ChooseBackend())
		assert.Equal(t, "B", rr.ChooseBackend())
	})
}
