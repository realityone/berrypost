package proxy

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeBinHeader(t *testing.T) {
	v, err := decodeMetadataHeader("Md-Bin", "ZGFuZ2Vyb3VzZA")
	assert.NoError(t, err)
	assert.Equal(t, "dangerousd", v)
}
