package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersParse(t *testing.T) {

	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Valid double header with done
	headers = NewHeaders()
	data = []byte(" Host: localhost:42069  \r\n  gUest: freakonaleash69\r\n\r\n")
	n, done, err = headers.Parse(data)
	total := n
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 26, n)
	assert.False(t, done)
	n, done, err = headers.Parse(data[total:])
	total += n
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, "freakonaleash69", headers["gUest"])
	assert.Equal(t, 26, n)
	assert.False(t, done)
	n, done, err = headers.Parse(data[total:])
	assert.Equal(t, 2, n)
	assert.True(t, done)
	require.NoError(t, err)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Valid double header with existing headers with done
	headers = NewHeaders()
	headers["billy"] = "bob"
	headers["frank"] = "joke"
	data = []byte(" Host: localhost:42069  \r\n  gUest: freakonaleash69\r\n\r\n")
	n, done, err = headers.Parse(data)
	total = n
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 26, n)
	assert.False(t, done)
	n, done, err = headers.Parse(data[total:])
	total += n
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, 26, n)
	assert.False(t, done)
	assert.Equal(t, "bob", headers["billy"])
	assert.Equal(t, "joke", headers["frank"])
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, "freakonaleash69", headers["gUest"])
	n, done, err = headers.Parse(data[total:])
	assert.Equal(t, 2, n)
	assert.True(t, done)
	require.NoError(t, err)

}
