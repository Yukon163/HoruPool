package confbox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
)

// copy from https://github.com/tucnak/store/blob/master/store.go

type MarshalFunc func(v interface{}) ([]byte, error)

type UnmarshalFunc func(data []byte, v interface{}) error

type format struct {
	m  MarshalFunc
	um UnmarshalFunc
}

var (
	formats = map[string]format{}
)

func init() {
	formats["json"] = format{m: json.Marshal, um: json.Unmarshal}
	formats["yaml"] = format{m: yaml.Marshal, um: yaml.Unmarshal}
	formats["yml"] = format{m: yaml.Marshal, um: yaml.Unmarshal}
}

// Register 注册序列化反序列化方法
func Register(extension string, m MarshalFunc, um UnmarshalFunc) {
	formats[extension] = format{m, um}
}

// Load 加载
func Load(path string, v interface{}) error {

	if format, ok := formats[extension(path)]; ok {
		return LoadWith(path, v, format.um)
	}

	panic("store: unknown configuration format")
}

// Save 保存
func Save(path string, v interface{}) error {

	if format, ok := formats[extension(path)]; ok {
		return SaveWith(path, v, format.m)
	}

	panic("store: unknown configuration format")
}

func LoadWith(path string, v interface{}, um UnmarshalFunc) error {
	data, err := ioutil.ReadFile(path)

	if err != nil {
		// There is a chance that file we are looking for
		// just doesn't exist. In this case we are supposed
		// to create an empty configuration file, based on v.
		empty := reflect.New(reflect.TypeOf(v))
		if innerErr := Save(path, &empty); innerErr != nil {
			// Smth going on with the file system... returning error.
			return err
		}

		v = empty

		return nil
	}

	if err := um(data, v); err != nil {
		return fmt.Errorf("store: failed to unmarshal %s: %v", path, err)
	}

	return nil
}

func SaveWith(path string, v interface{}, m MarshalFunc) error {
	var b bytes.Buffer

	if data, err := m(v); err == nil {
		b.Write(data)
	} else {
		return fmt.Errorf("store: failed to marshal %s: %v", path, err)
	}

	b.WriteRune('\n')

	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}

	if err := ioutil.WriteFile(path, b.Bytes(), os.ModePerm); err != nil {
		return err
	}
	return nil
}

func extension(path string) string {
	for i := len(path) - 1; i >= 0; i-- {
		if path[i] == '.' {
			return path[i+1:]
		}
	}

	return ""
}
