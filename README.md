# KeyPath

## What

Transforms `map[string]string{}` into `interface{}` using the map keys as definitions of structure.

## Semantics

Numeric keys indicate a list:
```
map[string]string{
    "1": "1",
}
=> 
[]interface{}{
    1: "1",
}
```

Non-Numeric keys specify a map:
```
map[string]string{
    "a": "1",
}
=> 
map[string]interface{}{
    "a": "1",
}
```

If Non-Numeric and Numeric keys are mixed, the Numeric Keys
are treated as Non-Numeric:
```
map[string]string{
    "0": "1",
    "a": "2",
}
=> 
map[string]interface{}{
    "0": "1",
    "a": "2",
}
```

Types can be nested arbitrarily:
```
map[string]string{
    "0.a.0": "1",
    "0.a.1": "2",
    "0.a.2": "3",
    "0.b.c": "4",
    "a.b.c": "5",
}
=> 
map[string]interface{}{
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
}
note: the "a.b.c" key changes the structure from a list to a map.
```

## Errors

An error is returned if a type previously specified is changed, for example:

```
map[string]string{
    "a":   "1",
    "a.b": "2",
}
```
results in a `errors.New("type mismatch")`
