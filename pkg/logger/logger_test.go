package logger

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetLogger(t *testing.T) {
	got := GetLogger()
	require.NotNil(t, got)
}
