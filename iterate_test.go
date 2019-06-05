package properties

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const str1 = `
	#comment1
	key1=1
	#comment2
	key 2 : 2
	key3 : 2
	#coment3
	`

func TestDoc_Map(t *testing.T) {
	doc, _ := LoadString(str1)
	m := doc.Map()
	assert.Equal(t, map[string]string{
		"key1":  "1",
		"key 2": "2",
		"key3":  "2",
	}, m)
}

func Test_Accept(t *testing.T) {
	count := 0

	doc, _ := LoadString(str1)
	doc.Accept(func(typo byte, value string, key string) bool {
		count++
		return typo != ':'
	})

	expect(t, "Accept提前中断", count == 5)
}

func Test_Foreach(t *testing.T) {
	count := 0

	doc, _ := LoadString(str1)
	doc.Foreach(func(value string, key string) bool {
		count++
		return key != "key 2"
	})

	expect(t, "Foreach提前中断", count == 2)
}
