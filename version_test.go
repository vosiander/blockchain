package blockchain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidVersions(t *testing.T) {
	d := []struct {
		Version      string
		Challenge    string
		IsCompatible bool
	}{
		{"0.0.1", "0.0.1", true},
		{"0.0.1", "0.0.2", false},
		{"0.0.1", "no-real-versioning", false},
		{"1.0.5", "1.0.3", true},
		{"1.1.1", "1.1.0", true},
		{"2.1.1", "1.1.0", false},
		{"0.9.8", "0.9.7", true},
	}

	for _, a := range d {
		assert.Equal(t, a.IsCompatible, IsCompatible(a.Version, a.Challenge), "incompatible version", a)
	}
}
