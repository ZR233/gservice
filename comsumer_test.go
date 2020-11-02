package gservice

import (
	"github.com/hashicorp/consul/api"
	"testing"
)

func TestConsumer_Open(t *testing.T) {
	o := NewOptions()
	o.ConsulConfig = api.DefaultConfig()
	m, err := NewManager(o)
	if err != nil {
		t.Error(err)
		return
	}
	c := m.NewConsumer("digger/history_check/rpc")
	err = c.Open()
	if err != nil {
		t.Error(err)
		return
	}
	conn, err := c.GetConn()
	println(conn, err)
	conn, err = c.GetConn()
	println(conn, err)
	conn, err = c.GetConn()
	println(conn, err)
	conn, err = c.GetConn()
	println(conn, err)
	conn, err = c.GetConn()
	println(conn, err)
}
