package keypath

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnpack(t *testing.T) {

	for _, tc := range []struct {
		d   string
		v   map[string]string
		e   interface{}
		err error
	}{

		{
			d: "map",
			v: map[string]string{
				"a": "1",
			},
			e: map[string]interface{}{
				"a": "1",
			},
		},
		{
			d: "map, map map",
			v: map[string]string{
				"a":   "1",
				"b.c": "2",
			},
			e: map[string]interface{}{
				"a": "1",
				"b": map[string]interface{}{
					"c": "2",
				},
			},
		},
		{
			d: "map, map map map, map map map",
			v: map[string]string{
				"a":     "1",
				"b.c.d": "2",
				"b.c.e": "3",
			},
			e: map[string]interface{}{
				"a": "1",
				"b": map[string]interface{}{
					"c": map[string]interface{}{
						"d": "2",
						"e": "3",
					},
				},
			},
		},
		{
			d: "list",
			v: map[string]string{
				"1": "1",
			},
			e: []interface{}{
				1: "1",
			},
		},
		{
			d: "list . list",
			v: map[string]string{
				"1.2": "1",
			},
			e: []interface{}{
				1: []interface{}{
					2: "1",
				},
			},
		},
		{
			d: "list . list . map",
			v: map[string]string{
				"1.2.test": "1",
			},
			e: []interface{}{
				1: []interface{}{
					2: map[string]interface{}{
						"test": "1",
					},
				},
			},
		},
		{
			d: "map . map . list",
			v: map[string]string{
				"a.1": "1",
			},
			e: map[string]interface{}{
				"a": []interface{}{
					1: "1",
				},
			},
		},
		{
			d: "map . map . list",
			v: map[string]string{
				"a.b.1": "1",
			},
			e: map[string]interface{}{
				"a": map[string]interface{}{
					"b": []interface{}{
						1: "1",
					},
				},
			},
		},
		{
			d: "grow slice",
			v: map[string]string{
				"f.2": "5",
				"f.5": "6",
			},
			e: map[string]interface{}{
				"f": []interface{}{
					2: "5",
					5: "6",
				},
			},
		},
		{
			d: "reverse sort prefers map",
			v: map[string]string{
				"f.5": "7",
				"f.g": "8",
			},
			e: map[string]interface{}{
				"f": map[string]interface{}{
					"5": "7",
					"g": "8",
				},
			},
		},
		{
			d: "modifying existing list structure errors",
			v: map[string]string{
				"0":   "1",
				"0.b": "2",
			},
			e:   nil,
			err: errors.New("type mismatch"),
		},
		{
			d: "modifying existing map structure errors",
			v: map[string]string{
				"a":   "1",
				"a.b": "2",
			},
			e:   nil,
			err: errors.New("type mismatch"),
		},
		{
			d: "struct can need to be grown",
			v: map[string]string{
				"a.101": "1",
				"a.20":  "2",
			},
			e: map[string]interface{}{
				"a": []interface{}{
					20:  "2",
					101: "1",
				},
			},
		},
		{
			d: "arbitrary nesting",
			v: map[string]string{
				"0.a.0": "1",
				"0.a.1": "2",
				"0.a.2": "3",
				"0.b.c": "4",
				"a.b.c": "5",
			},
			e: map[string]interface{}{
				"0": map[string]interface{}{
					"a": []interface{}{
						0: "1",
						1: "2",
						2: "3",
					},
					"b": map[string]interface{}{
						"c": "4",
					},
				},
				"a": map[string]interface{}{
					"b": map[string]interface{}{
						"c": "5",
					},
				},
			},
		},
		{
			d: "",
			v: map[string]string{
				"a":     "1",
				"0":     "2",
				"b.c.d": "4",
				"b.c.e": "5",
				"f.2":   "6",
				"f.5":   "7",
				"f.g":   "8",
			},
			e: map[string]interface{}{
				"a": "1",
				"0": "2",
				"b": map[string]interface{}{
					"c": map[string]interface{}{
						"d": "4",
						"e": "5",
					},
				},
				"f": map[string]interface{}{
					"2": "6",
					"5": "7",
					"g": "8",
				},
			},
		},
	} {
		a, err := Unpack(tc.v)
		assert.Equal(t, tc.e, a, tc.d)
		assert.Equal(t, tc.err, err)
	}
}
