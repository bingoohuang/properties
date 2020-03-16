package properties

import (
	"container/list"
)

type line struct {
	//  #   注释行
	//  !   注释行
	//  ' ' 空白行或者空行
	//  =   等号分隔的属性行
	//  :   冒号分隔的属性行
	typo  byte   //  行类型
	value string //  值,如果是注释注释引导符也包含在内。
	key   string //  如果是属性行这里表示属性的key
}

// Doc The properties document in memory.
type Doc struct {
	lines *list.List
	props map[string]*list.Element
}

func isComment(typo byte) bool {
	return typo == '#' || typo == '!'
}

func isProperty(typo byte) bool {
	return typo == '=' || typo == ':'
}

func (i line) isProperty() bool {
	return isProperty(i.typo)
}
