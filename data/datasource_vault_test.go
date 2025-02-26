package data

import (
	"context"
	"net/url"
	"testing"

	"github.com/hairyhenderson/gomplate/v4/vault"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadVault(t *testing.T) {
	ctx := context.Background()

	expected := "{\"value\":\"foo\"}\n"
	server, v := vault.MockServer(200, `{"data":`+expected+`}`)
	defer server.Close()

	source := &Source{
		Alias:     "foo",
		URL:       &url.URL{Scheme: "vault", Path: "/secret/foo"},
		mediaType: textMimetype,
		vc:        v,
	}

	r, err := readVault(ctx, source)
	require.NoError(t, err)
	assert.Equal(t, []byte(expected), r)

	r, err = readVault(ctx, source, "bar")
	require.NoError(t, err)
	assert.Equal(t, []byte(expected), r)

	r, err = readVault(ctx, source, "?param=value")
	require.NoError(t, err)
	assert.Equal(t, []byte(expected), r)

	source.URL, _ = url.Parse("vault:///secret/foo?param1=value1&param2=value2")
	r, err = readVault(ctx, source)
	require.NoError(t, err)
	assert.Equal(t, []byte(expected), r)

	expected = "[\"one\",\"two\"]\n"
	server, source.vc = vault.MockServer(200, `{"data":{"keys":`+expected+`}}`)
	defer server.Close()
	source.URL, _ = url.Parse("vault:///secret/foo/")
	r, err = readVault(ctx, source)
	require.NoError(t, err)
	assert.Equal(t, []byte(expected), r)
}
