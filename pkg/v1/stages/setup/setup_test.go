package setup

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStage_Name(t *testing.T) {
	assert.Equal(t, "setup", Stage{}.Name())
}

func TestStage_String(t *testing.T) {
	assert.Equal(t, "performing pre-flight setup and checks", Stage{}.String())
}

func TestStage_Run(t *testing.T) {
	// todo
}
