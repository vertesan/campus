package downloader

import "sync"

type SafeCounter struct {
  mutex sync.Mutex
  num   int
}

func (c *SafeCounter) Increase() {
  c.mutex.Lock()
  c.num++
  c.mutex.Unlock()
}

func (c *SafeCounter) Value() int {
  c.mutex.Lock()
  defer c.mutex.Unlock()
  return c.num
}

func (c *SafeCounter) Clear() {
  c.mutex.Lock()
  c.num = 0
  c.mutex.Unlock()
}
