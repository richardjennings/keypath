package keypath

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Unpack converts string keys into nested map structures
func Unpack(v map[string]string) (interface{}, error) {
	var r interface{}
	// sorted keys for ordered iteration
	keys := sortKeys(v)
	for j := range keys {
		if err := unpackKey(v, keys[j], &r); err != nil {
			return nil, err
		}
	}
	return r, nil
}

func unpackKey(v map[string]string, k string, rr *interface{}) error {
	parts := strings.Split(k, ".")
	if len(parts) > 1 {
		if err := unpackKeyParts(parts, rr, v[k]); err != nil {
			return err
		}
	} else {
		if err := unpackKeySinglePart(k, rr, v[k]); err != nil {
			return err
		}
	}
	return nil
}

func unpackKeyParts(parts []string, rr *interface{}, vv string) error {
	var curList, nextList *int
	var isMap bool
	prr := rr
	// for each part of the path
	for i := 0; i < len(parts); i++ {
		curList, nextList, isMap = keyMeta(rr, parts, i)

		// if key is numeric and current value is not a map,
		if curList != nil && !isMap {
			// create
			if *rr == nil {
				*rr = make([]interface{}, *curList+1)
			}
			// update pointer
			if i < len(parts)-1 {
				rr = &((*rr).([]interface{})[*curList])
			}
			continue
		}

		if *rr == nil {
			// map did not exist yet
			*rr = make(map[string]interface{})
		}

		if _, ok := (*rr).([]interface{}); ok {
			// if trying to add a string key as an index to a list throw an error
			return fmt.Errorf("string key %s in list", parts[i])
		}
		if _, ok := (*rr).(map[string]interface{})[parts[i]]; !ok {
			// key of map has not been set. create either map or list
			if nextList != nil {
				(*rr).(map[string]interface{})[parts[i]] = make([]interface{}, *nextList+1)
			} else {
				(*rr).(map[string]interface{})[parts[i]] = make(map[string]interface{})
			}
		}
		if i < len(parts)-1 {
			// move rr deeper
			a := (*rr).(map[string]interface{})[parts[i]]
			rr = &a
		}
	}
	setKeyValue(rr, vv, curList, prr, parts)

	return nil
}

func keyMeta(rr *interface{}, parts []string, i int) (*int, *int, bool) {
	var curList, nextList *int
	// is current value a map ?
	_, isMap := (*rr).(map[string]interface{})
	// is current key numeric ?
	if iv, err := strconv.Atoi(parts[i]); err == nil {
		curList = &iv
	}
	// is next key numeric ?
	if i < len(parts)-1 {
		if iv, err := strconv.Atoi(parts[i+1]); err == nil {
			nextList = &iv
		}
	}
	return curList, nextList, isMap
}

func setKeyValue(rr *interface{}, vv string, index *int, prr *interface{}, parts []string) {
	switch (*rr).(type) {
	case map[string]interface{}:
		(*rr).(map[string]interface{})[parts[len(parts)-1]] = vv
	case []interface{}:
		if len((*rr).([]interface{})) < *index+1 {
			// currently required when reverse sorting alpha-num does not
			// result in the bigger number first, e.g. 101, 21.
			// sorting by numeric part could remove the need for this.
			a := (*rr).([]interface{})
			a = append(a, make([]interface{}, (*index+1)-len(a))...)
			a[*index] = vv
			replace(prr, parts[:len(parts)-1], a)
		} else {
			(*rr).([]interface{})[*index] = vv
		}
	}
}

func unpackKeySinglePart(k string, rr *interface{}, vv string) error {
	if *rr == nil {
		// 1 part only, either map or list
		if iv, err := strconv.Atoi(k); err == nil {
			*rr = make([]interface{}, iv+1)
			(*rr).([]interface{})[iv] = vv
		} else {
			*rr = map[string]interface{}{
				k: vv,
			}
		}
	} else {

		// 1 part only, either map or list
		_, isMap := (*rr).(map[string]interface{})
		if iv, err := strconv.Atoi(k); err == nil && !isMap {

			// if exists already must be string type
			if vvv, ok := (*rr).([]interface{}); ok {
				if _, ok := vvv[iv].(string); !ok {
					return errors.New("type mismatch")
				}
			}

			// numeric type can be a map key in map context
			(*rr).([]interface{})[iv] = vv
		} else {
			// if exists already must be string type
			if vvv, ok := (*rr).(map[string]interface{})[k]; ok {
				if _, ok := vvv.(string); !ok {
					return errors.New("type mismatch")
				}
			}

			//
			(*rr).(map[string]interface{})[k] = vv
		}
	}
	return nil
}

func sortKeys(m map[string]string) []string {
	var keys []string
	for k := range m {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))
	return keys
}

func replace(r *interface{}, parts []string, required []interface{}) {
	rr := r
	for i, p := range parts {
		v, err := strconv.Atoi(p)
		if i == len(p)-1 {
			// at point to modify
			if err == nil {
				(*rr).([]interface{})[v] = required
			} else {
				(*rr).(map[string]interface{})[p] = required
			}
			return
		}
		if err == nil {
			*rr = (*rr).([]interface{})[v]
		} else {
			*rr = (*rr).(map[string]interface{})[p]
		}
	}
}
