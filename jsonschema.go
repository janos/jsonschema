package jsonschema // import "resenje.org/jsonschema"

import (
	"reflect"
	"strings"
	"unicode"
)

var jsonType = map[reflect.Kind]string{
	reflect.Bool:    "bool",
	reflect.Int:     "integer",
	reflect.Int8:    "integer",
	reflect.Int16:   "integer",
	reflect.Int32:   "integer",
	reflect.Int64:   "integer",
	reflect.Uint:    "integer",
	reflect.Uint8:   "integer",
	reflect.Uint16:  "integer",
	reflect.Uint32:  "integer",
	reflect.Uint64:  "integer",
	reflect.Float32: "number",
	reflect.Float64: "number",
	reflect.String:  "string",
	reflect.Slice:   "array",
	reflect.Struct:  "object",
	reflect.Map:     "object",
}

type property struct {
	Description          string              `json:"description,omitempty"`
	Type                 string              `json:"type,omitempty"`
	Items                *property           `json:"items,omitempty"`
	Properties           map[string]property `json:"properties,omitempty"`
	AdditionalProperties bool                `json:"additionalProperties,omitempty"`
	Required             []string            `json:"required,omitempty"`
}

type Schema struct {
	Schema string `json:"$schema,omitempty"`
	property
}

func New(variable interface{}) Schema {
	d := Schema{
		Schema: "http://json-schema.org/schema#",
	}
	d.read(reflect.ValueOf(variable).Type())
	return d
}

func (p *property) read(t reflect.Type) {
	kind := t.Kind()
	p.Type = jsonType[kind]

	switch kind {
	case reflect.Slice:
		pn := &property{}
		pn.read(t.Elem())
		p.Items = pn
	case reflect.Map:
		if jsType := jsonType[t.Elem().Kind()]; jsType != "" {
			p.Properties = make(map[string]property, 0)
			pn := &property{}
			pn.read(t.Elem())
			p.Properties[".*"] = *pn
		} else {
			p.AdditionalProperties = true
		}
	case reflect.Struct:
		p.Type = "object"
		p.Properties = make(map[string]property, 0)
		p.AdditionalProperties = false

		count := t.NumField()
		for i := 0; i < count; i++ {
			field := t.Field(i)

			tags := field.Tag
			name := strings.Split(tags.Get("json"), ",")[0]
			if name == "-" {
				continue
			}
			if name == "" {
				if unicode.IsLower(rune(field.Name[0])) {
					continue
				}
				name = field.Name
			}

			pn := &property{
				Description: tags.Get("description"),
			}
			pn.read(field.Type)
			p.Properties[name] = *pn

			if strings.Contains(tags.Get("minion"), "required") {
				p.Required = append(p.Required, name)
			}
		}
	case reflect.Ptr:
		p.read(t.Elem())
	}
}
