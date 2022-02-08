package libquantum

import "fmt"

type Options struct {
	Verbose bool
}

func (options *Options) String() string {
	return fmt.Sprintf("\tVerbose: %t", options.Verbose)
}

func NewOptions() *Options {
	return &Options{}
}
