// Package store provides data store definitions
package store

// IPV4Entity Define the entity of IPv4 data.
type IPV4Entity struct {
	startIndex uint32
	endIndex   uint32
	metaIndex  uint32
}

// IPV6Entity Define the entity of IPv6 data.
type IPV6Entity struct {
	startIndex uint64
	endIndex   uint64
	metaIndex  uint32
}

// StartIndex start index
func (i *IPV4Entity) StartIndex() uint32 {
	return i.startIndex
}

// EndIndex end index
func (i *IPV4Entity) EndIndex() uint32 {
	return i.endIndex
}

// StartIndex start index
func (i *IPV6Entity) StartIndex() uint64 {
	return i.startIndex
}

// EndIndex end index
func (i *IPV6Entity) EndIndex() uint64 {
	return i.endIndex
}
