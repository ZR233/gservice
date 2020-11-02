package gservice

import "io"

type Consumer struct {
	serviceName string
	manager     *Manager
	pool        connPool
	hostIter    int32
}

func (m *Manager) NewConsumer(serviceName string) (c *Consumer) {
	c = &Consumer{}
	c.serviceName = serviceName
	c.manager = m
	return
}

func (c *Consumer) Open() (err error) {
	err = c.newConnPool()
	return
}
func (c *Consumer) GetConn() (conn io.Closer, err error) {
	return c.pool.getOne()
}
