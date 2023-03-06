package config

import (
	"context"
	clientv3 "go.etcd.io/etcd/client/v3"
	"strings"
	"time"
)

type Etcd struct {
	host   []string
	client *clientv3.Client
}

// newEtcdClient 创建ETCD客户端
func newEtcdClient(path string) (*Etcd, error) {
	host := strings.Split(path, ",")
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   host,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	return &Etcd{host: host, client: cli}, nil
}

// Get 获取配置
func (e *Etcd) Get(key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	resp, err := e.client.Get(ctx, key)
	if err != nil {
		return "", err
	}
	// 由调用层判断是否存在
	if len(resp.Kvs) == 0 {
		return "", nil
	}
	return string(resp.Kvs[0].Value), nil
}

// Watch 监听配置
func (e *Etcd) Watch(key string, callback func(key string, value string)) error {
	rch := e.client.Watch(context.Background(), key)
	for wresp := range rch {
		for _, ev := range wresp.Events {
			callback(string(ev.Kv.Key), string(ev.Kv.Value))
		}
	}
	return nil
}

// GetPrefix 获取前缀配置
func (e *Etcd) GetPrefix(prefix string) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	resp, err := e.client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}
	// 由调用层判断是否存在
	if len(resp.Kvs) == 0 {
		return nil, nil
	}
	m := make(map[string]string)
	for _, v := range resp.Kvs {
		m[string(v.Key)] = string(v.Value)
	}
	return m, nil
}

// WatchPrefix 监听前缀配置
func (e *Etcd) WatchPrefix(prefix string, callback func(map[string]string)) error {
	rch := e.client.Watch(context.Background(), prefix, clientv3.WithPrefix())
	for wresp := range rch {
		m := make(map[string]string)
		for _, ev := range wresp.Events {
			m[string(ev.Kv.Key)] = string(ev.Kv.Value)
		}
		callback(m)
	}
	return nil
}

// Put 设置配置
func (e *Etcd) Put(key, value string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := e.client.Put(ctx, key, value)
	if err != nil {
		return err
	}
	return nil
}

// Delete 删除配置
func (e *Etcd) Delete(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := e.client.Delete(ctx, key)
	if err != nil {
		return err
	}
	return nil
}

// Close 关闭ETCD客户端
func (e *Etcd) Close() error {
	return e.client.Close()
}
