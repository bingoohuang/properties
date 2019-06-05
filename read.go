package properties

import "strconv"

// StrOr retrieves the string value by key.
// If the line is not exist, the def will be returned.
func (p Doc) StrOr(key, def string) string {
	if val, ok := p.Get(key); ok {
		return val
	}

	return def
}

// IntOr retrieves the int value by key.
// If the line is not exist, the def will be returned.
func (p Doc) IntOr(key string, def int) int {
	if val, ok := p.Get(key); ok {
		if v, err := strconv.Atoi(val); err == nil {
			return v
		}
	}

	return def
}

// Int64Or retrieves the int64 value by key.
// If the line is not exist, the def will be returned.
func (p Doc) Int64Or(key string, def int64) int64 {
	if val, ok := p.Get(key); ok {
		if v, err := strconv.ParseInt(val, 10, 64); err == nil {
			return v
		}
	}

	return def
}

// Uint64Or Same as Int64Or, but the return type is uint64.
func (p Doc) Uint64Or(key string, def uint64) uint64 {
	if val, ok := p.Get(key); ok {
		if v, err := strconv.ParseUint(val, 10, 64); err == nil {
			return v
		}
	}

	return def
}

// Float64Or   retrieve the float64 value by key.
// If the line is not exist, the def will be returned.
func (p Doc) Float64Or(key string, def float64) float64 {
	if val, ok := p.Get(key); ok {
		if v, err := strconv.ParseFloat(val, 64); err == nil {
			return v
		}
	}

	return def
}

// BoolOr   retrieve the bool value by key.
// If the line is not exist, the def will be returned.
// This function mapping "1", "t", "T", "true", "TRUE", "True" as true.
// This function mapping "0", "f", "F", "false", "FALSE", "False" as false.
// If the line is not exist of can not map to value of bool,the def will be returned.
func (p Doc) BoolOr(key string, def bool) bool {
	if val, ok := p.Get(key); ok {
		if v, err := strconv.ParseBool(val); err == nil {
			return v
		}
	}

	return def
}

// ObjectOr maps the value of the key to any object.
// The f is the customized mapping function.
// Return def if the line is not exist of f have a error returned.
func (p Doc) ObjectOr(key string, def interface{}, f func(k, v string) (interface{}, error)) interface{} {
	if val, ok := p.Get(key); ok {
		if v, err := f(key, val); err == nil {
			return v
		}
	}

	return def
}

// Str same as StrOr but the def is "".
func (p Doc) Str(key string) string {
	return p.StrOr(key, "")
}

// Int is same as IntOr but the def is 0 .
func (p Doc) Int(key string) int {
	return p.IntOr(key, 0)
}

// Int64 is same as Int64Or but the def is 0 .
func (p Doc) Int64(key string) int64 {
	return p.Int64Or(key, 0)
}

// Uint64 same as Uint64Or but the def is 0 .
func (p Doc) Uint64(key string) uint64 {
	return p.Uint64Or(key, 0)
}

// Float64 same as Float64Or but the def is 0.0 .
func (p Doc) Float64(key string) float64 {
	return p.Float64Or(key, 0.0)
}

// Bool same as BoolOr but the def is false.
func (p Doc) Bool(key string) bool {
	return p.BoolOr(key, false)
}

// Object is same as ObjectOr but the def is nil.
//
// Notice: If the return value can not be assign to nil, this function will panic/
func (p Doc) Object(key string, f func(k, v string) (interface{}, error)) interface{} {
	return p.ObjectOr(key, nil, f)
}
