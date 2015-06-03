package raymond

import "reflect"

// DataFrame represents a private data frame
//
// Cf. private variables documentation at: http://handlebarsjs.com/block_helpers.html
type DataFrame struct {
	parent *DataFrame
	data   map[string]interface{}
}

// NewDataFrame instanciates a new private data frame
func NewDataFrame() *DataFrame {
	return &DataFrame{
		data: make(map[string]interface{}),
	}
}

// Copy instanciates a new private data frame, with parent set to self
func (p *DataFrame) Copy() *DataFrame {
	result := NewDataFrame()

	for k, v := range p.data {
		result.data[k] = v
	}

	result.parent = p

	return result
}

// NewIterDataFrame instanciates a new private data frame, with parent set to self an with iterable data set (@index, @key, @first, @last)
func (p *DataFrame) NewIterDataFrame(length int, i int, key interface{}) *DataFrame {
	result := p.Copy()

	result.Set("index", i)
	result.Set("key", key)
	result.Set("first", i == 0)
	result.Set("last", i == length-1)

	return result
}

// Set sets a data value
func (p *DataFrame) Set(key string, val interface{}) {
	p.data[key] = val
}

// Get gets a data value
func (p *DataFrame) Get(key string) interface{} {
	return p.Find([]string{key})
}

// Find gets a deep data value
//
// @todo This is NOT consistent with the way we resolve data in template (cf. `evalDataPathExpression()`) ! FIX THAT !
func (p *DataFrame) Find(parts []string) interface{} {
	data := p.data

	for i, part := range parts {
		val := data[part]
		if val == nil {
			return nil
		}

		if i == len(parts)-1 {
			// found
			return val
		}

		valValue := reflect.ValueOf(val)
		if valValue.Kind() != reflect.Map {
			// not found
			return nil
		}

		// continue
		data = mapStringInterface(valValue)
	}

	// not found
	return nil
}

// mapStringInterface converts any `map` to `map[string]interface{}`
func mapStringInterface(value reflect.Value) map[string]interface{} {
	result := make(map[string]interface{})

	for _, key := range value.MapKeys() {
		result[strValue(key)] = value.MapIndex(key).Interface()
	}

	return result
}
