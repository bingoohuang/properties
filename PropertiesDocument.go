package properties

import (
    "io"
    "bufio"
    "bytes"
    "unicode"
    "strconv"
    "container/list"
    "fmt"
    "strings"
)

type element struct {
    //  #   注释行
    //  !   注释行
    //  ' ' 空白行或者空行
    //  =   等号分隔的属性行
    //  :   冒号分隔的属性行
    typo  byte   //  行类型
    value string //  行的内容,如果是注释注释引导符也包含在内
    key   string //  如果是属性行这里表示属性的key
}

type PropertiesDocument struct {
    elems *list.List
    props map[string]*list.Element
}

func New() *PropertiesDocument {
    doc := new(PropertiesDocument)
    doc.elems = list.New()
    doc.props = make(map[string]*list.Element)
    return doc
}

func Save(doc *PropertiesDocument, writer io.Writer) {
    doc.Accept(func(typo byte, value string, key string) bool {
        switch typo {
        case '#', '!', ' ':
            {
                fmt.Fprintln(writer, value)
            }
        case '=', ':':
            {
                fmt.Fprintf(writer, "%s%c%s\n", key, typo, value)
            }
        }
        return true
    })
}

func Load(reader io.Reader) (p *PropertiesDocument, err error) {
    
    //  创建一个Properties对象
    p = New()
    
    //  创建一个扫描器
    scanner := bufio.NewScanner(reader)
    for scanner.Scan() {
        //  逐行读取
        line := scanner.Bytes()
        
        //  遇到空行
        if 0 == len(line) {
            p.elems.PushBack(&element{typo: ' ', value: string(line)})
            continue
        }
        
        //  找到第一个非空白字符
        pos := bytes.IndexFunc(line, func(r rune) bool {
            return !unicode.IsSpace(r)
        })
        
        //  遇到空白行
        if -1 == pos {
            p.elems.PushBack(&element{typo: ' ', value: string(line)})
            continue
        }
        
        //  遇到注释行
        if '#' == line[pos] {
            p.elems.PushBack(&element{typo: '#', value: string(line)})
            continue
        }
        
        if '!' == line[pos] {
            p.elems.PushBack(&element{typo: '!', value: string(line)})
            continue
        }
        
        //  找到第一个等号的位置
        end := bytes.IndexFunc(line[pos+1:], func(r rune) bool {
            return ('=' == r) || (':' == r)
        })
        
        //  没有=，说明该配置项只有key
        key := ""
        value := ""
        if -1 == end {
            key = string(bytes.TrimRightFunc(line[pos:], func(r rune) bool {
                return unicode.IsSpace(r)
            }))
        } else {
            key = string(bytes.TrimRightFunc(line[pos:pos+1+end], func(r rune) bool {
                return unicode.IsSpace(r)
            }))
            
            value = string(bytes.TrimSpace(line[pos+1+end+1:]))
        }
        
        var typo byte = '='
        if end > 0 {
            typo = line[end]
        }
        elem := &element{typo: typo, key: key, value: value}
        listelem := p.elems.PushBack(elem)
        p.props[key] = listelem
    }
    
    if err = scanner.Err(); nil != err {
        return nil, err
    }
    
    return p, nil
}

func (p PropertiesDocument) Get(key string) (value string, exist bool) {
    e, ok := p.props[key]
    return e.Value.(*element).value, ok
}

func (p*PropertiesDocument) Set(key string, value string) {
    e, ok := p.props[key]
    if ok {
        e.Value.(*element).value = value
        return
    }
    
    p.props[key] = p.elems.PushBack(&element{typo: '=', key: key, value: value})
}

func (p*PropertiesDocument) Del(key string) bool {
    e, ok := p.props[key]
    if !ok {
        return false
    }
    
    p.Uncomment(key)
    p.elems.Remove(e)
    delete(p.props, key)
    return true
}

func (p*PropertiesDocument) Comment(key string, comments string) bool {
    e, ok := p.props[key]
    if !ok {
        return false
    }
    
    //  如果所有注释为空
    if len(comments) <= 0 {
        p.elems.InsertBefore(&element{typo: '#', value: "#"}, e)
        return true
    }
    
    //  创建一个新的Scanner
    scanner := bufio.NewScanner(strings.NewReader(comments))
    for scanner.Scan() {
        
        line := scanner.Text()
        
        if len(line) <= 0 {
            p.elems.InsertBefore(&element{typo: '#', value: "#"}, e)
            continue
        }
        
        if ('#' != line[0]) && ('!' != line[0]) {
            p.elems.InsertBefore(&element{typo: '#', value: "#" + line}, e)
        }
    }
    
    return true
}

func (p*PropertiesDocument) Uncomment(key string) bool {
    e, ok := p.props[key]
    if !ok {
        return false
    }
    
    for item := e.Prev(); nil != item; {
        del := item
        item = item.Prev()
        
        if ('=' == del.Value.(*element).typo) ||
            (':' == del.Value.(*element).typo) ||
            (' ' == del.Value.(*element).typo) {
            break
        }
        
        p.elems.Remove(del)
    }
    
    return true
}

func (p PropertiesDocument) Accept(f func(typo byte, value string, key string) bool) {
    for e := p.elems.Front(); e != nil; e = e.Next() {
        elem := e.Value.(*element)
        continues := f(elem.typo, elem.value, elem.key)
        if !continues {
            return
        }
    }
}

func (p PropertiesDocument) Foreach(f func(value string, key string) bool) {
    for e := p.elems.Front(); e != nil; e = e.Next() {
        elem := e.Value.(*element)
        if ('=' == elem.typo) ||
            (':' == elem.typo) {
            continues := f(elem.value, elem.key)
            if !continues {
                return
            }
        }
    }
}

func (p PropertiesDocument) StringDefault(key string, def string) string {
    e, ok := p.props[key]
    if ok {
        return e.Value.(*element).value
    }
    
    return def
}

func (p PropertiesDocument) IntDefault(key string, def int64) int64 {
    e, ok := p.props[key]
    if ok {
        v, err := strconv.ParseInt(e.Value.(*element).value, 10, 64)
        if nil != err {
            return def
        }
        
        return v
    }
    
    return def
}

func (p PropertiesDocument) FloatDefault(key string, def float64) float64 {
    e, ok := p.props[key]
    if ok {
        v, err := strconv.ParseFloat(e.Value.(*element).value, 64)
        if nil != err {
            return def
        }
        
        return v
    }
    
    return def
}

func (p PropertiesDocument) BoolDefault(key string, def bool) bool {
    e, ok := p.props[key]
    if ok {
        v, err := strconv.ParseBool(e.Value.(*element).value)
        if nil != err {
            return def
        }
        
        return v
    }
    
    return def
}

func (p PropertiesDocument) ObjectDefault(key string, def interface{}, f func(k string, v string) (interface{}, error)) interface{} {
    e, ok := p.props[key]
    if ok {
        v, err := f(key, e.Value.(*element).value)
        if nil != err {
            return def
        }
        
        return v
    }
    
    return def
}

func (p PropertiesDocument) String(key string) string {
    return p.StringDefault(key, "")
}

func (p PropertiesDocument) Int(key string) int64 {
    return p.IntDefault(key, 0)
}

func (p PropertiesDocument) Float(key string) float64 {
    return p.FloatDefault(key, 0.0)
}

func (p PropertiesDocument) Bool(key string) bool {
    return p.BoolDefault(key, false)
}

func (p PropertiesDocument) Object(key string, f func(k string, v string) (interface{}, error)) interface{} {
    return p.ObjectDefault(key, nil, f)
}