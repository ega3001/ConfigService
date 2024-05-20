package cfgclients

import (
	"encoding/json"
	"path"
	"time"

	"github.com/go-zookeeper/zk"
)

type zkCfg struct {
	c *zk.Conn
}

func NewZKCfg(servers []string, duration time.Duration) (*zkCfg, error) {
	c, _, err := zk.Connect(servers, duration)
	if err != nil {
		return nil, err
	}
	return &zkCfg{c}, nil
}

func (zkc zkCfg) Create(key string, value map[string]any) error {
	valBytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = zkc.c.Create(
		key,
		valBytes,
		zk.FlagTTL,
		zk.WorldACL(zk.PermAll),
	)
	return err
}

func (zkc zkCfg) Children(key string) ([]string, error) {
	childs, _, err := zkc.c.Children(key)
	return childs, err
}

func (zkc zkCfg) Exists(key string) (bool, error) {
	exists, _, err := zkc.c.Exists(key)
	return exists, err
}

func (zkc zkCfg) Get(key string) (map[string]any, error) {
	res := map[string]any{}
	bytes, _, err := zkc.c.Get(key)
	if err != nil {
		return nil, err
	}
	return res, json.Unmarshal(bytes, &res)
}

func (zkc zkCfg) Set(key string, value map[string]any) error {
	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = zkc.c.Set(key, bytes, -1)
	return err
}

func (zkc zkCfg) Delete(key string) error {
	childs, err := zkc.Children(key)
	if err != nil {
		return err
	}
	for _, child := range childs {
		childFullKey := path.Join(key, child)
		zkc.Delete(childFullKey)
	}
	return zkc.c.Delete(key, -1)
}

func (zkc zkCfg) WatchNodeChanges(key string, callback func(map[string]any, error)) {
	for {
		var data map[string]any
		bytes, _, w, err := zkc.c.GetW(key)
		if err != nil {
			callback(nil, err)
		}
		err = json.Unmarshal(bytes, &data)

		callback(data, err)

		if event := <-w; event.Type == zk.EventNodeDeleted {
			return
		}
	}
}
