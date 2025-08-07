// Package ipip IPIP数据源核心驱动
package ipip

import (
	"net"
	"sync"

	"github.com/universal-fraternity/ipip/store"
)

var (
	defaultStore *store.Store
	once         sync.Once
)

// Store output store.Store.
type Store = store.Store

// Meta output store.Meta。
type Meta = store.Meta

// Option output store.Option。
type Option = store.Option

// FileInfo output store.FileInfo
type FileInfo = store.FileInfo

// Load load data
func Load(opt Option) error {
	once.Do(func() {
		defaultStore = store.NewStore()
	})
	return defaultStore.LoadData(opt)
}

// Update update data
func Update(fs ...FileInfo) error {
	return update(fs...)
}

func update(fs ...FileInfo) error {
	if len(fs) > 0 {
		defaultStore.WithDataFiles(fs)
	}
	return defaultStore.Update()
}

// Search meta by address .
func Search(addr string) *Meta {
	return defaultStore.Search(net.ParseIP(addr))
}
