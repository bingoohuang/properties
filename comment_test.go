package properties

import (
	"bytes"
	"testing"
)

func Test_Comment_Uncomment(t *testing.T) {
	str := "key1=1\nkey 2 = 2"

	doc, _ := LoadString(str)
	exist := doc.Comment("NOT-EXIST", "Some comment")
	expect(t, "对一个不存在的项执行注释操作,返回false", false == exist)

	doc.Comment("key1", "This is a \ncomment \nfor a")

	buf := bytes.NewBufferString("")
	err := doc.Save(buf)
	expect(t, "格式化成功", nil == err)

	exp1 := "#This is a \n#comment \n#for a\nkey1=1\nkey 2=2\n"
	expect(t, "对已经存在的项进行注释", exp1 == buf.String())

	doc.Comment("key 2", "")

	buf = bytes.NewBufferString("")
	err = doc.Save(buf)
	expect(t, "格式化成功", nil == err)

	exp2 := "#This is a \n#comment \n#for a\nkey1=1\n#\nkey 2=2\n"
	expect(t, "对已经存在的项进行注释", exp2 == buf.String())

	exist = doc.Uncomment("key1")
	expect(t, "对已经存在的key进行注释,返回true", true == exist)

	buf = bytes.NewBufferString("")
	err = doc.Save(buf)
	expect(t, "格式化成功", nil == err)

	exp3 := "key1=1\n#\nkey 2=2\n"
	expect(t, "对已经存在的项进行注释", exp3 == buf.String())

	exist = doc.Uncomment("key 2")
	expect(t, "对已经存在的key进行注释,返回true", true == exist)
	expect(t, "对已经存在的项进行注释", exp3 == buf.String())

	buf = bytes.NewBufferString("")
	err = doc.Save(buf)
	expect(t, "格式化成功", nil == err)

	exp4 := "key1=1\nkey 2=2\n"
	expect(t, "对已经存在的项进行注释", exp4 == buf.String())

	exist = doc.Uncomment("NOT-EXIST")
	expect(t, "对不已经存在的key进行注释,返回false", false == exist)
}
