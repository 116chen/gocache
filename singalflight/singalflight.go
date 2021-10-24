package singalflight

import "sync"

// call 一次请求的结果
type call struct {
	wg    sync.WaitGroup
	value interface{}
	err   error
}

// Group 管理请求的并发
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
		// 不会有两个协程同时进入这里面，因为前面用mutext作了并发控制
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
