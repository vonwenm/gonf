package gonf

import (
	"io"
	"io/ioutil"
	"reflect"
	"strconv"
	"errors"
)

type Config struct {
	value string
	parent *Config
	map_ map[string]*Config
	array []*Config
}

func Read(r io.Reader) (*Config, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	input := string(b)
	tokens := make(chan token)
	l := newLexer(input, tokens)
	p := newParser(tokens)
	go l.lex()
	return p.parse()
}

func (c *Config) String(args ...interface{}) (string, error) {
	c, err := c.Get(args...)
	return c.value, err
}

func (c *Config) Int(args ...interface{}) (int, error) {
	s, err := c.String(args...)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(s)
}

func (c *Config) Get(args ...interface{}) (*Config, error) {
	var ok bool
	for _, a := range args {
		switch reflect.TypeOf(a).Kind() {
		case reflect.String:
			if c, ok = c.map_[a.(string)]; !ok {
				return nil, errors.New("key " + a.(string) + " not found")
			}
		case reflect.Int:
			if c = c.array[a.(int)]; c == nil {
				return nil, errors.New("index " + strconv.Itoa(a.(int)) + " not found")
			}
		}
	}
	return c, nil
}

func (c *Config) Map(s interface{}) error {
	t := reflect.TypeOf(s)

	if t.Kind() != reflect.Ptr {
		return errors.New("The argument to Map must be a pointer")
	}

	t = t.Elem()

	if t.Kind() != reflect.Struct {
		return errors.New("The argument to Map must be a struct")
	}

	v := reflect.ValueOf(s).Elem()

	c.rmap(t, v)

	return nil
}

func (c *Config) rmap(t reflect.Type, v reflect.Value) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if tag := field.Tag.Get("gonf"); tag != "" {
			f := v.FieldByName(field.Name)
			switch field.Type.Kind() {
			case reflect.String:
				if value, err := c.String(tag); err == nil {
					f.SetString(value)
				}
			case reflect.Int:
				if value, err := c.Int(tag); err == nil {
					f.SetInt(int64(value))
				}
			case reflect.Struct:
				if c, err := c.Get(tag); err == nil {
					c.rmap(f.Type(), f)
				}
			}
		}
	}
}