package data

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mustParseURL(in string) *url.URL {
	u, _ := url.Parse(in)
	return u
}

func TestReadEnv(t *testing.T) {
	content := []byte(`hello world`)

	source, _ := NewSource("foo", mustParseURL("env:///tmp/foo"))

	actual, err := readEnv(source)
	assert.NoError(t, err)
	assert.Equal(t, content, actual)

	source, _ = NewSource("bogus", mustParseURL("env:///bogus"))
	_, err = readEnv(source)
	assert.Error(t, err)

	source, _ = NewSource("partial", mustParseURL("env:///tmp/partial"))
	actual, err = readEnv(source, "foo.txt")
	assert.NoError(t, err)
	assert.Equal(t, content, actual)

	source, _ = NewSource("dir", mustParseURL("env:///tmp/partial/"))
	actual, err = readEnv(source)
	assert.NoError(t, err)
	assert.Equal(t, []byte(`["bar.txt","baz.txt","foo.txt"]`), actual)
}
