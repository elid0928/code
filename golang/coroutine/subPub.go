package coroutine

import (
	"sync"
	"time"
)

// 构建一个发布者对象
type publisher struct {
	// 需要一个🔓
	lock        sync.RWMutex
	buffer      int                      // 缓存大小
	timeout     time.Duration            //超时时间
	subscribers map[subscriber]topicFunc // 订阅者信息
}

// 订阅者
// 主题过滤器
type (
	subscriber chan interface{}
	topicFunc  func(v interface{}) bool
)

// 构建一个发布者对象
func NewPublisher(timeout time.Duration, buf int) *publisher {

	return &publisher{
		buffer:      buf,
		timeout:     timeout,
		subscribers: make(map[subscriber]topicFunc),
	}
}

// 订阅全部主题?
func (p *publisher) Subscribe() chan interface{} {
	return p.SubscribeTopic(nil)
}

// 订阅符合过滤器的主题
func (p *publisher) SubscribeTopic(topic topicFunc) chan interface{} {
	// 加锁
	p.lock.Lock()
	ch := make(chan interface{}, p.buffer)
	p.subscribers[ch] = topic
	p.lock.Unlock()
	return ch
}

// 取消订阅
func (p *publisher) Evict(sub chan interface{}) bool {
	p.lock.Lock()
	delete(p.subscribers, sub)
	close(sub)
	p.lock.Unlock()
	return true
}

// 发送消息
// 发布一个主题
func (p *publisher) Publish(v interface{}) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	var wg sync.WaitGroup
	for sub, topic := range p.subscribers {
		wg.Add(1)
		go p.sendTopic(sub, topic, v, &wg)
	}
	wg.Wait()
}

// 关闭发布者对象，同时关闭所有的订阅者管道。
func (p *publisher) Close() {
	p.lock.Lock()
	defer p.lock.Unlock()

	for sub := range p.subscribers {
		delete(p.subscribers, sub)
		close(sub)
	}
}

// 发送主题，可以容忍一定的超时
func (p *publisher) sendTopic(
	sub subscriber, topic topicFunc, v interface{}, wg *sync.WaitGroup,
) {
	defer wg.Done()
	if topic != nil && !topic(v) {
		return
	}

	select {
	case sub <- v:
	case <-time.After(p.timeout):
	}
}
