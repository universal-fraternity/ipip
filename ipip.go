// Package ipip IPIP Data Source Core Driver Test Code
package ipip

import (
	"net"
	"sync"

	"github.com/universal-fraternity/ipip/core"
)

var (
	defaultStore *core.Store
	once         sync.Once
)

// Store output core.Store.
type Store = core.Store

// Meta output core.Meta。
type Meta = core.Meta

// Option output core.Option。
type Option = core.Option

// FileInfo output core.FileInfo
type FileInfo = core.FileInfo

// Init init core.Store and load data
func Init(opt Option) error {
	once.Do(func() {
		defaultStore = core.NewStore()
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
