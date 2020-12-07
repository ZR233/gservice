package gservice

import (
	"fmt"
	"io"
	"sync"
	"time"
)

// 判断连接是否可用
type ConnTestFunc func(closer io.Closer) error

// ConnFactory is a function to create new connections.
type ConnFactory func(host string) (io.Closer, error)

type Conn struct {
	io.Closer
	service
}

func newConn(conn io.Closer, s service) *Conn {
	c := &Conn{
		Closer: conn,
	}
	c.service = s
	return c
}
func (c *Conn) connTest() (err error) {
	err = c.testFun(c.Closer)
	if err != nil {
		err = fmt.Errorf("conn [%s] fail\n%w", c.host, err)
	}
	return
}

type connPool struct {
	err        error
	firstFlush chan error
	consumer   *Consumer
	connPool   map[string]*Conn
	connIter   int
	sync.Mutex
}

func (c *Consumer) newConnPool() (err error) {
	c.pool = connPool{
		firstFlush: make(chan error),
		consumer:   c,
		connPool:   map[string]*Conn{},
	}

	err = c.pool.flushConn()
	return
}

func (c *connPool) setConnList(once *sync.Once) {
	var (
		err  error
		conn io.Closer
	)
	defer func() {
		once.Do(func() {
			c.firstFlush <- err
		})
	}()

	entry, _, err := c.consumer.manager.consulClient.Health().Service(c.consumer.serviceName, "", true, nil)
	if err != nil {
		err = fmt.Errorf("get service info fail\n%w", err)
		c.err = err
		return
	}
	if len(entry) == 0 {
		err = fmt.Errorf("service[%s] no alive node", c.consumer.serviceName)
		c.err = err
		return
	}

	c.Lock()
	defer c.Unlock()

	serviceList := map[string]*service{}
	for _, e := range entry {
		host := fmt.Sprintf("%s:%d", e.Node.Address, e.Service.Port)
		s := &service{}
		s.fromTags(e.Service.Tags)
		serviceList[host] = s
	}

	newPool := map[string]*Conn{}

	//填入已存在的连接
	for host, conn := range c.connPool {
		if _, ok := serviceList[host]; ok {
			newPool[host] = conn
		} else {
			_ = conn.Close()
		}
	}
	//填入新连接
	for host, service := range serviceList {
		if _, ok := c.connPool[host]; !ok {
			service.host = host
			conn, err = service.factory(host)
			if err != nil {
				err = fmt.Errorf("can't conn to [%s]\n%w", host, ErrConn)
				c.err = err
				continue
			}
			newPool[host] = newConn(conn, *service)
		}
	}
	c.connPool = newPool
}
func (c *connPool) getOne() (conn io.Closer, err error) {
	c.Lock()
	defer c.Unlock()

	var connList []*Conn
	var deleteConnList []string
	for key, conn := range c.connPool {
		if err := conn.connTest(); err != nil {
			c.err = err
			_ = conn.Close()
			deleteConnList = append(deleteConnList, key)
		} else {
			connList = append(connList, conn)
		}
	}
	for _, key := range deleteConnList {
		delete(c.connPool, key)
	}

	if len(connList) == 0 {
		err = fmt.Errorf("service[%s] no alive node\n%w", c.consumer.serviceName, c.err)
		return
	}

	if c.connIter >= len(connList) {
		c.connIter = 0
	}

	conn = connList[c.connIter].Closer
	c.connIter++
	return
}

func (c *connPool) flushConn() (err error) {
	once := &sync.Once{}
	c.firstFlush = make(chan error)
	go func() {
		for {
			if c.consumer.manager.ctx.Err() != nil {
				return
			}
			c.setConnList(once)
			time.Sleep(time.Second)
		}
	}()
	err = <-c.firstFlush
	return
}
