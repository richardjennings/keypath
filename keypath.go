package keypath

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// unpack converts string keys into nested map structures
func unpack(v map[string]string) (interface{}, error) {

	var keys []string
	for k := range v {
		keys = append(keys, k)
	}
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))

	var r interface{}
	var rr *interface{}
	prr := &r

	var curList *int
	var nextList *int

	// iterate keys
	for j := range keys {
		k := keys[j]
		vv := v[k]

		// reset rr to top on iteration
		rr = prr

		// split path by seperator
		if p := strings.Split(k, "."); len(p) > 1 {

			// for each part of the path
			for i := 0; i < len(p); i++ {

				// reset pointers to lists
				curList = nil
				nextList = nil

				// determine current and next type if exists
				if iv, err := strconv.Atoi(p[i]); err == nil {
					curList = &iv
				}

				_, isMap := (*rr).(map[string]interface{})

				if i < len(p)-1 {
					if iv, err := strconv.Atoi(p[i+1]); err == nil {
						nextList = &iv
					}
				}

				if curList != nil && !isMap {
					if rr == nil {
						*rr = []interface{}{}
					}
					// list
					if *rr == nil {
						*rr = make([]interface{}, *curList+1)
					}

					if i < len(p)-1 {
						rr = &((*rr).([]interface{})[*curList])
					}
				} else {
					// map
					if *rr == nil {
						// map did not exist yet
						*rr = make(map[string]interface{})
					}

					if _, ok := (*rr).([]interface{}); ok {
						// if trying to add a string key as an index to a list,
						// throw an error
						return nil, fmt.Errorf("string key %s in list", p[i])
					}
					if _, ok := (*rr).(map[string]interface{})[p[i]]; !ok {
						// key of map not exists
						// create either map or list
						if nextList != nil {
							(*rr).(map[string]interface{})[p[i]] = make([]interface{}, *nextList+1)
						} else {
							(*rr).(map[string]interface{})[p[i]] = make(map[string]interface{})
						}
					}
					if i < len(p)-1 {
						// move rr deeper
						a := (*rr).(map[string]interface{})[p[i]]
						rr = &a
					}
				}
			}

			// set value
			switch (*rr).(type) {
			case map[string]interface{}:
				(*rr).(map[string]interface{})[p[len(p)-1]] = vv
			case []interface{}:
				if len((*rr).([]interface{})) < *curList+1 {
					// currently required when reverse sorting alpha-num does not
					// result in the bigger number first, e.g. 101, 21.
					// sorting by numeric part could remove the need for this.
					a := (*rr).([]interface{})
					a = append(a, make([]interface{}, (*curList+1)-len(a))...)
					a[*curList] = vv
					replace(prr, p[:len(p)-1], a)
				} else {
					(*rr).([]interface{})[*curList] = vv
				}
			}

		} else {
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
							return nil, errors.New("type mismatch")
						}
					}

					// numeric type can be a map key in map context
					(*rr).([]interface{})[iv] = vv
				} else {
					// if exists already must be string type
					if vvv, ok := (*rr).(map[string]interface{})[k]; ok {
						if _, ok := vvv.(string); !ok {
							return nil, errors.New("type mismatch")
						}
					}

					//
					(*rr).(map[string]interface{})[k] = vv
				}
			}
		}
	}
	return r, nil
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
