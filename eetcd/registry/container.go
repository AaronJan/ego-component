package registry

import (
	"github.com/gotomicro/ego/core/conf"
	"github.com/gotomicro/ego/core/elog"

	"github.com/gotomicro/ego-component/eetcd"
)

type Option func(c *Container)

type Container struct {
	config *Config
	name   string
	logger *elog.Component
	client *eetcd.Component
}

func DefaultContainer() *Container {
	return &Container{
		logger: elog.EgoLogger.With(elog.FieldMod("client.egrpc")),
	}
}

func Load(key string) *Container {
	c := DefaultContainer()
	var config = DefaultConfig()
	if err := conf.UnmarshalKey(key, &config); err != nil {
		c.logger.Panic("parse Config error", elog.FieldErr(err), elog.FieldKey(key))
		return c
	}

	c.config = config
	c.name = key

	return c
}

func WithClientEtcd(etcdClient *eetcd.Component) Option {
	return func(c *Container) {
		c.client = etcdClient
	}
}

// Build ...
func (c *Container) Build(options ...Option) *Component {
	for _, option := range options {
		option(c)
	}
	if c.client == nil {
		c.logger.Panic("client etcd nil", elog.FieldKey("use WithClientEtcd method"))
	}
	return newComponent(c.name, c.config, c.logger, c.client)
}
