package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func NewHeaders() Headers {
	return Headers{}
}

func TestHeaderParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Capital letters in field name
	headers = NewHeaders()
	data = []byte("  HOST: localhost:42069  \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 27, n)
	assert.False(t, done)

	// Test: Invalid character in field name
	headers = NewHeaders()
	data = []byte("  H©ST: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Valid single header
	headers = NewHeaders()
	data = []byte("Set-Person: lane-loves-go;\r\n")
	_, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 0, n)
	require.False(t, done)

	data = []byte("Set-Person: prime-loves-zig;\r\n")
	_, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 0, n)
	require.False(t, done)

	data = []byte("Set-Person: tj-loves-ocaml;\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	assert.Equal(t, 29, n)
	require.False(t, done)

	assert.Equal(t, "lane-loves-go, prime-loves-zig, tj-loves-ocaml", headers["set-person"])
}
