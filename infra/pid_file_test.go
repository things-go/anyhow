package infra

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPidFile(t *testing.T) {
	dir := "/tmp/pidfile"
	require.NoError(t, WritePidFile(dir))
	require.NoError(t, RemovePidFile(dir))
}
