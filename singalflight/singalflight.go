package singalflight

import "sync"

type call struct {
	wg    sync.WaitGroup
	value interface{}
	err   error
}

type Group struct {
	mu sync.Mutex
	m  map[string]*call
}

// Do 防止并发请求过快而导致缓存击穿，缓存雪崩现象
func (this *Group) Do(key string, f func() (interface{}, error)) (interface{}, error) {
	this.mu.Lock()
	if this.m == nil {
		this.m = make(map[string]*call)
	}
	if c, ok := this.m[key]; ok {
		this.mu.Unlock()
		c.wg.Wait()
		return c.value, c.err
	}
	c := new(call)
	c.wg.Add(1)
	this.m[key] = c
	this.mu.Unlock()

	c.value, c.err = f()
	c.wg.Done()

	this.mu.Lock()
	delete(this.m, key)
	this.mu.Unlock()

	return c.value, c.err
}
