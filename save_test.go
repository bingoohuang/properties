// nolint gomnd
package properties

import (
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/bingoohuang/gou/ran"
	"github.com/stretchr/testify/assert"
)

func TestPopulate(t *testing.T) {
	prop, _ := LoadMap(map[string]string{
		"key1":        "value1",
		"key2":        "yes",
		"key3":        "true",
		"XingMing":    "kongrong",
		"foo-bar":     "foobar",
		"NI-HAO":      "10s",
		"ta-hao":      "10s",
		"ORDER_PRICE": "100",
		"order_items": "10",
		"HelloWorld":  "10",
	})

	type MySub2 struct {
		XingMing string
	}

	type MySub struct {
		Key1 string `prop:"key1"`
		Key2 bool
		Key3 *bool
	}

	type my struct {
		MySub
		*MySub2
		FooBar     *string
		NiHao      time.Duration
		TaHao      *time.Duration
		xx         string
		YY         string
		OrderPrice int
		OrderItems int
		HelloWorld *int
	}

	var (
		m my
		x int
	)

	it := assert.New(t)

	err := prop.Populate(m, "prop")
	it.Error(err)

	err = prop.Populate(&x, "prop")
	it.Error(err)

	err = prop.Populate(&m, "prop")

	it.Nil(err)

	foobar := "foobar"
	HelloWorld := 10
	key3 := true
	taHao := 10 * time.Second

	it.Equal(my{
		MySub: MySub{
			Key1: "value1",
			Key2: true,
			Key3: &key3,
		},
		MySub2: &MySub2{
			XingMing: "kongrong",
		},
		FooBar:     &foobar,
		NiHao:      10 * time.Second,
		TaHao:      &taHao,
		xx:         "",
		YY:         "",
		OrderPrice: 100,
		OrderItems: 10,
		HelloWorld: &HelloWorld,
	}, m)

	prop, _ = LoadMap(map[string]string{
		"NI-HAO":      "10x",
		"ta-hao":      "10s",
		"ORDER_PRICE": "100",
		"order_items": "10",
		"HelloWorld":  "10",
	})

	type Myx struct {
		NiHao time.Duration
	}

	type Myy struct {
		*Myx
	}

	var (
		myx Myx
		myy Myy
	)

	it.Error(prop.Populate(&myx, "prop"))
	it.Error(prop.Populate(&myy, "prop"))
}

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
