package properties

import (
	"bytes"
	"testing"
)

func Test_Accept(t *testing.T) {
	str := `
	#comment1
	key1=1
	#comment2
	key 2 : 2
	key3 : 2
	#coment3
	`

	count := 0

	doc, _ := Load(bytes.NewBufferString(str))
	doc.Accept(func(typo byte, value string, key string) bool {
		count++
		return typo != ':'
	})

	expect(t, "Accept提前中断", 5 == count)
}

func Test_Foreach(t *testing.T) {
	str := `
	#comment1
	key1=1
	#comment2
	key 2 : 2
	key3 : 2
	#coment3
	`

	count := 0

	doc, _ := Load(bytes.NewBufferString(str))
	doc.Foreach(func(value string, key string) bool {
		count++
		return "key 2" != key
	})

	expect(t, "Foreach提前中断", 2 == count)
}
