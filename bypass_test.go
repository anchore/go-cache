package cache

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_bypassedCache(t *testing.T) {
	m := NewBypassed()
	c := m.GetCache("name", "version")
	err := c.Write("test", strings.NewReader("value"))
	require.NoError(t, err)
	rdr, err := c.Read("test")
	require.Nil(t, rdr)
	require.ErrorContains(t, err, "not found")
}
