package libquantum

type Options struct {
	Verbose bool
}

func NewOptions() *Options {
	return &Options{}
}
