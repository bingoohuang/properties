package properties

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/bingoohuang/gou/str"
	"github.com/bingoohuang/strcase"
)

// Populate populates the properties to the structure's field.
func (p Doc) Populate(b interface{}, tag string) error {
	v := reflect.ValueOf(b)
	vt := v.Type()

	if vt.Kind() != reflect.Ptr {
		return errors.New("only argument of pointer of structure  supported")
	}

	v = v.Elem()
	vt = v.Type()

	if vt.Kind() != reflect.Struct {
		return errors.New("only argument of pointer of structure  supported")
	}

	for i := 0; i < vt.NumField(); i++ {
		structField := vt.Field(i)
		if structField.PkgPath != "" { // bypass non-exported fields
			continue
		}

		fieldType := structField.Type
		fieldPtr := fieldType.Kind() == reflect.Ptr

		if fieldPtr {
			fieldType = fieldType.Elem()
		}

		field := v.Field(i)

		if fieldType.Kind() == reflect.Struct {
			err := p.parseStruct(fieldType, tag, fieldPtr, field)
			if err != nil {
				return err
			}

			continue
		}

		switch fieldType.Kind() {
		case reflect.String:
			p.parseString(structField, tag, fieldPtr, field)
		case reflect.Int:
			p.parseInt(structField, tag, fieldPtr, field)
		case reflect.Bool:
			p.parseBool(structField, tag, fieldPtr, field)
		default:
			if fieldType == reflect.TypeOf(time.Duration(0)) {
				if err := p.parseDuration(structField, tag, fieldPtr, field); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (p Doc) parseStruct(fieldType reflect.Type, tag string, fieldPtr bool, field reflect.Value) error {
	fv := reflect.New(fieldType)
	if err := p.Populate(fv.Interface(), tag); err != nil {
		return err
	}

	if fieldPtr {
		field.Set(fv)
	} else {
		field.Set(fv.Elem())
	}

	return nil
}

func (p Doc) parseString(structField reflect.StructField, tag string, fieldPtr bool, field reflect.Value) {
	if v, ok := p.tryGet(structField, tag); ok {
		if !fieldPtr {
			field.Set(reflect.ValueOf(v))
		} else {
			field.Set(reflect.ValueOf(&v))
		}
	}
}

func (p Doc) parseInt(structField reflect.StructField, tag string, fieldPtr bool, field reflect.Value) {
	if v, ok := p.tryGet(structField, tag); ok {
		vi := str.ParseInt(v)

		if !fieldPtr {
			field.Set(reflect.ValueOf(vi))
		} else {
			field.Set(reflect.ValueOf(&vi))
		}
	}
}

func (p Doc) parseBool(structField reflect.StructField, tag string, fieldPtr bool, field reflect.Value) {
	if v, ok := p.tryGet(structField, tag); ok {
		v = strings.ToLower(v)
		vi := v == "true" || v == "yes" || v == "on" || v == "1"

		if !fieldPtr {
			field.Set(reflect.ValueOf(vi))
		} else {
			field.Set(reflect.ValueOf(&vi))
		}
	}
}

func (p Doc) parseDuration(structField reflect.StructField, tag string, fieldPtr bool, field reflect.Value) error {
	if v, ok := p.tryGet(structField, tag); ok {
		d, err := time.ParseDuration(v)
		if err != nil {
			return err
		}

		if !fieldPtr {
			field.Set(reflect.ValueOf(d))
		} else {
			field.Set(reflect.ValueOf(&d))
		}
	}

	return nil
}

func (p Doc) tryGet(structField reflect.StructField, tag string) (string, bool) {
	if value, ok := p.Get(structField.Name); ok {
		return value, true
	}

	if tagValue := structField.Tag.Get(tag); tagValue != "" {
		if value, ok := p.Get(tagValue); ok {
			return value, true
		}
	}

	if value, ok := p.Get(strcase.ToCamelLower(structField.Name)); ok {
		return value, true
	}

	if value, ok := p.Get(strcase.ToSnake(structField.Name)); ok {
		return value, true
	}

	if value, ok := p.Get(strcase.ToSnakeUpper(structField.Name)); ok {
		return value, true
	}

	if value, ok := p.Get(strcase.ToKebab(structField.Name)); ok {
		return value, true
	}

	if value, ok := p.Get(strcase.ToKebabUpper(structField.Name)); ok {
		return value, true
	}

	return "", false
}

// String gives the whole properties as a string
func (p Doc) String() string {
	s, _ := p.Export()

	return s
}

// ExportFile saves the doc to file.
func (p Doc) ExportFile(file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	err = p.Save(w)

	if err != nil {
		return err
	}

	w.Flush()

	return nil
}

// Export saves the doc to file or stream.
func (p Doc) Export() (string, error) {
	buf := bytes.NewBufferString("")
	err := p.Save(buf)

	return buf.String(), err
}

// Save saves the doc to file or stream.
func (p Doc) Save(writer io.Writer) error {
	var err error

	p.Accept(func(typo byte, value string, key string) bool {
		if isProperty(typo) {
			_, err = fmt.Fprintf(writer, "%s%c%s\n", key, typo, value)
		} else {
			_, err = fmt.Fprintln(writer, value)
		}

		return nil == err
	})

	return err
}
