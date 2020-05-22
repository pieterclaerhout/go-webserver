package respond

import (
	"fmt"
)

// asBytes returns a slice of bytes.
func asBytes(valuea ...interface{}) []byte {
	value := valuea[0]
	if value == nil {
		return []byte{}
	}

	switch val := value.(type) {
	case bool:
		if val == true {
			return []byte("true")
		}
		return []byte("false")
	case string:
		return []byte(val)
	case []byte:
		return val
	default:
		return []byte(fmt.Sprintf("%v", value))
	}
}
