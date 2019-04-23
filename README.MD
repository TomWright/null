[![Build Status](https://travis-ci.org/TomWright/null.svg?branch=master)](https://travis-ci.org/TomWright/null)
[![codecov](https://codecov.io/gh/TomWright/null/branch/master/graph/badge.svg)](https://codecov.io/gh/TomWright/null)
[![Documentation](https://godoc.org/github.com/TomWright/null?status.svg)](https://godoc.org/github.com/TomWright/null)

# null

```
go get -u github.com/tomwright/null
```

Easily scan/marshal between standard types and `null` values.

The null package provides a set of structs that can be used as both `scan` and `marshal` sources and destinations, with the sole exception being that an empty value will be marshal'd to a value of `null`, and a `null` value will be unmarshal'd to an empty value.

The following types are available:

- [`null.String`: `null.NewString("Some string")`](https://godoc.org/github.com/TomWright/null#String)
- [`null.Time`: `null.NewTime(time.Now())`](https://godoc.org/github.com/TomWright/null#Time)

# Determining validity

Each struct within the package has a property or func named `Valid`. `Valid` will be `true` if the provided value was present and not empty.
When retrieving a value from the struct, a `null` should be used whenever `Valid` is `false`.

# Additional types

Additional types will be added/maintained as needed or as requested. Feel free to submit any pull request you find useful or necessary.
