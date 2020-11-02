package gservice

import (
	"context"
	"fmt"
	"github.com/hashicorp/consul/api"
)

type Options struct {
	ConsulConfig *api.Config
	ConsulClient *api.Client
}

func NewOptions() *Options {
	o := &Options{}
	return o
}

type Manager struct {
	options      *Options
	consulClient *api.Client
	ctx          context.Context
	cancel       context.CancelFunc
}

func NewManager(options *Options) (m *Manager, err error) {
	m = &Manager{
		options: options,
	}
	m.ctx, m.cancel = context.WithCancel(context.Background())
	if m.options.ConsulClient != nil {
		m.consulClient = m.options.ConsulClient
	} else if m.options.ConsulConfig != nil {
		m.consulClient, err = api.NewClient(m.options.ConsulConfig)
	} else {
		err = fmt.Errorf("there is no consul config")
		return
	}
	return
}

func (m *Manager) Close() error {
	m.cancel()
	return nil
}
