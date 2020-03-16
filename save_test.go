// nolint gomnd
package properties

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/bingoohuang/gou/ran"
	"github.com/stretchr/testify/assert"
)

func TestSaveFile(t *testing.T) {
	ioutil.WriteFile("save_test.properties", []byte("key=value"), 0644)

	doc, err := LoadFile("save_test.properties")
	assert.Nil(t, err)

	os.Remove("save_test.properties")

	val := ran.String(10)
	doc.Set("key", val)

	err = doc.ExportFile("save_test.properties")
	assert.Nil(t, err)

	nv, _ := ioutil.ReadFile("save_test.properties")
	assert.Equal(t, string(nv), "key="+val+"\n")
	os.Remove("save_test.properties")
}
