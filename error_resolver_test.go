package cache

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func Test_errorResolver(t *testing.T) {
	// original := GetManager()
	// defer SetManager(original)
	manager := NewInMemory(time.Hour)
	c := manager.GetCache("theCache", "theVersion")

	resolver := NewResolverCachingErrors[string](c)

	resolveCount := 0
	resolveFunc := func() (string, error) {
		resolveCount++
		return "theValue", nil
	}

	val, err := resolver.Resolve("theKey", resolveFunc)
	require.NoError(t, err)
	require.Equal(t, 1, resolveCount)
	require.Equal(t, "theValue", val)

	val, err = resolver.Resolve("theKey", resolveFunc)
	require.NoError(t, err)
	require.Equal(t, 1, resolveCount)
	require.Equal(t, "theValue", val)

	errorCount := 0
	errorFunc := func() (string, error) {
		errorCount++
		return "", fmt.Errorf("an error")
	}

	_, err = resolver.Resolve("errorValue", errorFunc)
	require.ErrorContains(t, err, "an error")
	require.Equal(t, 1, errorCount)

	_, err = resolver.Resolve("errorValue", errorFunc)
	require.ErrorContains(t, err, "an error")
	require.Equal(t, 1, errorCount)
}
