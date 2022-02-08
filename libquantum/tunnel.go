package libquantum

import (
	"fmt"
)

type Endpoint struct {
	Host string
	Port int
}

// Method to return a Endpoint as a string
func (endpoint *Endpoint) String() string {
	return fmt.Sprintf("%s:%d", endpoint.Host, endpoint.Port)
}
