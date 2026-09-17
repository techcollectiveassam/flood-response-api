package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeometryToText(t *testing.T) {
	assert.Equal(t, "value", GeometryToText("value"))
	assert.Equal(t, "bytes", GeometryToText([]byte("bytes")))
	assert.Equal(t, "", GeometryToText(nil))
	assert.Equal(t, "", GeometryToText(42))
	assert.Equal(t, "", GeometryToText(""))
}
