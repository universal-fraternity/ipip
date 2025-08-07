// Package store provides data store definitions
package store

// CallBackFunc Callback function format definition
type CallBackFunc func(meta *Meta) interface{}

// FileInfo File configuration information
type FileInfo struct {
	Path string
	Type int
}

// Option config option
type Option struct {
	Files []FileInfo
	CB    CallBackFunc
}
