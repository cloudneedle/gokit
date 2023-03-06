package config

type Type int

const (
	ETCD Type = iota // ETCD
	YAML             // YAML
)

// Client 配置客户端
type Client struct {
	t    Type   // 配置类型
	path string // 配置路径:如果是ETCD,则为ETCD的地址,多个地址用逗号分隔;如果是YAML,则为YAML文件的路径
	Etcd *Etcd  // ETCD客户端,如果是YAML,则为nil
}

// ClientOption 配置客户端选项
type ClientOption func(*Client)

// WithType 设置配置类型
func WithType(t Type) ClientOption {
	return func(c *Client) {
		c.t = t
	}
}

// WithPath 设置配置路径
func WithPath(path string) ClientOption {
	return func(c *Client) {
		c.path = path
	}
}

// New 创建配置客户端
func New(opts ...ClientOption) (*Client, error) {
	cli := new(Client)
	for _, opt := range opts {
		opt(cli)
	}

	if cli.t == ETCD {
		// 创建ETCD客户端
		etcd, err := newEtcdClient(cli.path)
		if err != nil {
			return nil, err
		}
		cli.Etcd = etcd
	}
	return cli, nil
}

// Close 关闭配置客户端
func (c *Client) Close() error {
	if c.Etcd != nil {
		return c.Etcd.Close()
	}
	return nil
}
