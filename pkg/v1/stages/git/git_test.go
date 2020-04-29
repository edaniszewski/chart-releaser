package git

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStage_Name(t *testing.T) {
	assert.Equal(t, "git", Stage{}.Name())
}

func TestStage_String(t *testing.T) {
	assert.Equal(t, "parsing git information", Stage{}.String())
}

func TestStage_Run(t *testing.T) {
	// todo
}
