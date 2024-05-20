package node

type CfgClient interface {
	Create(key string, value map[string]any) error
	Children(key string) ([]string, error)
	Exists(key string) (bool, error)
	Get(key string) (map[string]any, error)
	Set(key string, value map[string]any) error
	Delete(key string) error
	WatchNodeChanges(key string, callback func(map[string]any, error))
}
