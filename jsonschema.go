// Copyright (c) 2015 Janoš Guljaš <janos@resenje.org>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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

// Property defines a standard json-schema.org Property.
type Property struct {
	Description          string              `json:"description,omitempty"`
	Type                 string              `json:"type,omitempty"`
	Items                *Property           `json:"items,omitempty"`
	Properties           map[string]Property `json:"properties,omitempty"`
	AdditionalProperties bool                `json:"additionalProperties,omitempty"`
	Required             []string            `json:"required,omitempty"`
}

// Schema encapsulates Property and adds a new description field $schema.
type Schema struct {
	Schema string `json:"$schema,omitempty"`
	Property
}

// New returns a Schema from a provided interface{}.
func New(variable interface{}) Schema {
	d := Schema{
		Schema: "http://json-schema.org/schema#",
	}
	d.read(reflect.ValueOf(variable).Type())
	return d
}

func (p *Property) read(t reflect.Type) {
	kind := t.Kind()
	p.Type = jsonType[kind]

	switch kind {
	case reflect.Slice:
		pn := &Property{}
		pn.read(t.Elem())
		p.Items = pn
	case reflect.Map:
		if jsType := jsonType[t.Elem().Kind()]; jsType != "" {
			p.Properties = make(map[string]Property)
			pn := &Property{}
			pn.read(t.Elem())
			p.Properties[".*"] = *pn
		} else {
			p.AdditionalProperties = true
		}
	case reflect.Struct:
		p.Type = "object"
		p.Properties = make(map[string]Property)
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

			pn := &Property{
				Description: tags.Get("description"),
			}
			pn.read(field.Type)
			p.Properties[name] = *pn

			if strings.Contains(tags.Get("jsonschema"), "required") {
				p.Required = append(p.Required, name)
			}
		}
	case reflect.Ptr:
		p.read(t.Elem())
	}
}
