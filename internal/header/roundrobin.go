package header

import "sync/atomic"

type counter struct {
	value uint64
}

func (c *counter) Next(max int) int {
	if max == 0 {
		return 0
	}
	v := atomic.AddUint64(&c.value, 1)
	return int(v % uint64(max))
}
