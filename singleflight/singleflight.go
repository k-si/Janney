package singleflight

import "sync"

type Call struct {
	wg  sync.WaitGroup
	val interface{}
	err error
}

type Group struct {
	mu    sync.Mutex
	calls map[string]*Call
}

func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()

	if g.calls == nil {
		g.calls = make(map[string]*Call)
	}

	// 其他线程访问，同时等待相同请求的结果
	if c, ok := g.calls[key]; ok {
		g.mu.Unlock()
		c.wg.Wait()
		return c.val, c.err
	}

	c := new(Call)
	c.wg.Add(1)
	g.calls[key] = c
	g.mu.Unlock()

	// 处理请求
	c.val, c.err = fn()
	c.wg.Done()

	g.mu.Lock()
	delete(g.calls, key)
	g.mu.Unlock()

	return c.val, c.err
}
