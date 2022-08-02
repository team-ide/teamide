package node

import (
	"errors"
	"hash/crc32"
	"sync"
)

var (
	MessageListenerPoolStop = errors.New("消息监听器池已停止")
	MessageListenerNull     = errors.New("消息监听器池暂无监听器")
)

type MessageListenerPool struct {
	listenerMu sync.Mutex
	listeners  []*MessageListener
	timeout    int64
	isStop     bool
	getIndex   int
	fromNodeId string
	toNodeId   string
}

func (this_ *MessageListenerPool) Remove(listener *MessageListener) {
	listener.stop()
	if this_.isStop {
		return
	}
	this_.listenerMu.Lock()
	defer this_.listenerMu.Unlock()
	var list []*MessageListener

	var listeners = this_.listeners
	for _, one := range listeners {
		if one != listener {
			list = append(list, one)
		}
	}
	this_.listeners = list
}

func (this_ *MessageListenerPool) Put(listener *MessageListener) (size int) {
	if this_.isStop {
		return
	}
	this_.listenerMu.Lock()
	defer this_.listenerMu.Unlock()
	this_.listeners = append(this_.listeners, listener)
	size = len(this_.listeners)
	return
}

func (this_ *MessageListenerPool) Stop() {
	this_.isStop = true
	this_.listenerMu.Lock()
	defer this_.listenerMu.Unlock()

	var listeners = this_.listeners
	this_.listeners = []*MessageListener{}
	for _, one := range listeners {
		one.stop()
	}
}

func (this_ *MessageListenerPool) Clean() {
	this_.listenerMu.Lock()
	defer this_.listenerMu.Unlock()

	var listeners = this_.listeners
	this_.listeners = []*MessageListener{}
	for _, one := range listeners {
		one.stop()
	}
}

func (this_ *MessageListenerPool) getTimeout() (timeout int64) {
	timeout = this_.timeout
	return timeout * 1000

}
func (this_ *MessageListenerPool) get(key string) (listener *MessageListener, err error) {

	if this_.isStop {
		err = MessageListenerPoolStop
		return
	}

	this_.listenerMu.Lock()
	defer this_.listenerMu.Unlock()

	var list = this_.listeners
	var size = len(list)
	if size == 0 {
		err = MessageListenerNull
		return
	}
	if key != "" {
		hashCode := int(crc32.ChecksumIEEE([]byte(key)))
		listener = list[hashCode%size]
		return
	}

	listener = list[this_.getIndex%size]

	this_.getIndex++
	if this_.getIndex >= size {
		this_.getIndex = 0
	}

	return
}

func (this_ *MessageListenerPool) Do(key string, do func(listener *MessageListener) (err error)) (err error) {

	if this_.isStop {
		err = MessageListenerPoolStop
		return
	}
	listener, err := this_.get(key)
	if err != nil {
		return
	}
	err = do(listener)
	return
}
