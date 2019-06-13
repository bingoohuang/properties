package properties

import (
	"testing"

	"github.com/bingoohuang/gou"
	"github.com/stretchr/testify/assert"
)

func TestSaveFile(t *testing.T) {
	doc, err := LoadFile("save_test.properties")
	assert.Nil(t, err)

	doc.Set("key", gou.RandomString(10))

	err = doc.ExportFile("save_test.properties")
	assert.Nil(t, err)
}
