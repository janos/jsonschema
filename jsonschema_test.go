// Copyright (c) 2015, 2016 Janoš Guljaš <janos@resenje.org>
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jsonschema

import (
	"encoding/json"
	"testing"
)

func Test(t *testing.T) {
	for i, test := range []struct {
		v interface{}
		s string
	}{
		{
			v: true,
			s: `{"$schema":"http://json-schema.org/schema#","type":"bool"}`,
		},
		{
			v: int(0),
			s: `{"$schema":"http://json-schema.org/schema#","type":"integer"}`,
		},
		{
			v: int8(0),
			s: `{"$schema":"http://json-schema.org/schema#","type":"integer"}`,
		},
		{
			v: int16(0),
			s: `{"$schema":"http://json-schema.org/schema#","type":"integer"}`,
		},
		{
			v: int32(0),
			s: `{"$schema":"http://json-schema.org/schema#","type":"integer"}`,
		},
		{
			v: int64(0),
			s: `{"$schema":"http://json-schema.org/schema#","type":"integer"}`,
		},
		{
			v: "",
			s: `{"$schema":"http://json-schema.org/schema#","type":"string"}`,
		},
		{
			v: uint(0),
			s: `{"$schema":"http://json-schema.org/schema#","type":"integer"}`,
		},
		{
			v: uint8(0),
			s: `{"$schema":"http://json-schema.org/schema#","type":"integer"}`,
		},
		{
			v: uint16(0),
			s: `{"$schema":"http://json-schema.org/schema#","type":"integer"}`,
		},
		{
			v: uint32(0),
			s: `{"$schema":"http://json-schema.org/schema#","type":"integer"}`,
		},
		{
			v: uint64(0),
			s: `{"$schema":"http://json-schema.org/schema#","type":"integer"}`,
		},
		{
			v: float32(0),
			s: `{"$schema":"http://json-schema.org/schema#","type":"number"}`,
		},
		{
			v: float64(0),
			s: `{"$schema":"http://json-schema.org/schema#","type":"number"}`,
		},
		{
			v: "",
			s: `{"$schema":"http://json-schema.org/schema#","type":"string"}`,
		},
		{
			v: new(string),
			s: `{"$schema":"http://json-schema.org/schema#","type":"string"}`,
		},
		{
			v: []string{},
			s: `{"$schema":"http://json-schema.org/schema#","type":"array","items":{"type":"string"}}`,
		},
		{
			v: struct {
				Test string
			}{},
			s: `{"$schema":"http://json-schema.org/schema#","type":"object","properties":{"Test":{"type":"string"}}}`,
		},
		{
			v: struct {
				Test string `json:"test"`
			}{},
			s: `{"$schema":"http://json-schema.org/schema#","type":"object","properties":{"test":{"type":"string"}}}`,
		},
		{
			v: struct {
				Test string `jsonschema:"required"`
			}{},
			s: `{"$schema":"http://json-schema.org/schema#","type":"object","properties":{"Test":{"type":"string"}},"required":["Test"]}`,
		},
		{
			v: struct {
				Test string `json:"-"`
			}{},
			s: `{"$schema":"http://json-schema.org/schema#","type":"object"}`,
		},
		{
			v: struct {
				test string
			}{},
			s: `{"$schema":"http://json-schema.org/schema#","type":"object"}`,
		},
		{
			v: struct {
				Test1 struct {
					Test2 string
					Test3 []int
				}
			}{},
			s: `{"$schema":"http://json-schema.org/schema#","type":"object","properties":{"Test1":{"type":"object","properties":{"Test2":{"type":"string"},"Test3":{"type":"array","items":{"type":"integer"}}}}}}`,
		},
		{
			v: map[string]string{},
			s: `{"$schema":"http://json-schema.org/schema#","type":"object","properties":{".*":{"type":"string"}}}`,
		},
		{
			v: map[string]error{},
			s: `{"$schema":"http://json-schema.org/schema#","type":"object","additionalProperties":true}`,
		},
	} {
		s, err := json.Marshal(New(test.v))
		if err != nil {
			t.Errorf("test case %d: %s", i, err)
		}
		if string(s) != test.s {
			t.Errorf("test case %d: expected %s, got %s", i, test.s, s)
		}
	}
}
