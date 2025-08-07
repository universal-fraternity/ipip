// Package store provides data store definitions
package store

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"sort"
	"sync"

	"github.com/universal-fraternity/ipip/utils"
)

// Store Define the storage area for storing IP location data.
type Store struct {
	ipv4EntityList []*IPV4Entity
	ipv6EntityList []*IPV6Entity
	v6MateList     []*Meta
	v4MateList     []*Meta
	opt            Option
	v4Mu           sync.RWMutex
	v6Mu           sync.RWMutex
}

// NewStore returns a new store.
func NewStore() *Store {
	return &Store{}
}

// WithDataFiles Set data file
func (s *Store) WithDataFiles(fs []FileInfo) {
	if s != nil && len(fs) > 0 {
		s.opt.Files = fs
	}
}

// IPV4EntityCount Number of IPv4 instances
func (s *Store) IPV4EntityCount() int {
	if s != nil {
		return len(s.ipv4EntityList)
	}
	return 0
}

// IPV6EntityCount Number of IPv6 instances
func (s *Store) IPV6EntityCount() int {
	if s != nil {
		return len(s.ipv6EntityList)
	}
	return 0
}

// IPV4Entity Return the index IPV4Entity pointing to the entity list.
func (s *Store) IPV4Entity(i int) *IPV4Entity {
	if s != nil {
		if i < len(s.ipv4EntityList) && i >= 0 {
			return s.ipv4EntityList[i]
		}
	}
	return nil
}

// IPV6Entity Return the index IPV6 Entity pointing to the entity list.
func (s *Store) IPV6Entity(i int) *IPV6Entity {
	if s != nil && i < s.IPV6EntityCount() && i >= 0 {
		return s.ipv6EntityList[i]
	}
	return nil
}

// Search
func (s *Store) Search(addr net.IP) *Meta {
	if addr == nil {
		return nil
	}
	if utils.IsIPv4(addr.String()) {
		// IPv4
		s.v4Mu.RLock()
		defer s.v4Mu.RUnlock()

		ipIndex := binary.BigEndian.Uint32(addr.To4())
		if index := sort.Search(s.IPV4EntityCount(), func(i int) bool {
			// Find the index of the first IPIndex greater than or equal to the given IP
			return s.ipv4EntityList[i].StartIndex() >= ipIndex
		}); s.IPV4Entity(index) != nil {
			if s.IPV4Entity(index).StartIndex() != ipIndex {
				index -= 1
			}
			if index < 0 {
				return nil
			}

			v4Entity := s.IPV4Entity(index)
			if v4Entity.StartIndex() <= ipIndex && v4Entity.EndIndex() >= ipIndex {
				mi := v4Entity.metaIndex
				return s.v4MateList[mi]
			}
		}
	} else if utils.IsIPv6(addr.String()) {
		// IPV6
		s.v6Mu.RLock()
		defer s.v6Mu.RUnlock()

		ipIndex := binary.BigEndian.Uint64(addr.To16())
		if index := sort.Search(s.IPV6EntityCount(), func(i int) bool {
			// Find the index of the first IPIndex greater than or equal to the given IP
			return s.ipv6EntityList[i].StartIndex() >= ipIndex
		}); s.IPV6Entity(index) != nil {
			if s.IPV6Entity(index).StartIndex() != ipIndex {
				index -= 1
			}
			if index < 0 {
				return nil
			}

			v6Entity := s.IPV6Entity(index)
			if v6Entity.StartIndex() <= ipIndex && v6Entity.EndIndex() >= ipIndex {
				mi := v6Entity.metaIndex
				return s.v6MateList[mi]
			}
		}
	}
	return nil
}

// UnmarshalFrom Decompose and store from raeder.
func (s *Store) UnmarshalFrom(reader io.Reader, t int) error {
	if t == Unknown {
		return errors.New("unknown data type")
	}

	var err error
	ipv4List := make([]*IPV4Entity, 0)
	ipv6List := make([]*IPV6Entity, 0)
	metaTable := make(map[string]uint32)
	tmpMetaList := make([]*Meta, 0)
	iReader := bufio.NewReader(reader)
	for {
		var line []byte
		if line, err = iReader.ReadBytes('\n'); err != nil {
			break
		}
		rowMeta := &RowMeta{}
		if err = rowMeta.Unmarshal(line, t); err != nil {
			_, _ = fmt.Fprint(os.Stderr, "meta unmarshal error, ", err.Error(), string(line))
			continue
		}
		fp := rowMeta.Hash()
		if fp == "" {
			// Fingerprint calculation error
			_, _ = fmt.Fprint(os.Stderr, "Fingerprint calculation error")
			continue
		}
		var index uint32
		var ok bool
		if index, ok = metaTable[fp]; !ok {
			meta := &Meta{
				Country:        rowMeta.Country,
				Province:       rowMeta.Province,
				City:           rowMeta.City,
				Region:         rowMeta.Region,
				OwnerDomain:    rowMeta.OwnerDomain,
				IspDomain:      rowMeta.IspDomain,
				ChinaAdminCode: rowMeta.ChinaAdminCode,
				Latitude:       rowMeta.Latitude,
				Longitude:      rowMeta.Longitude,
				Timezone:       rowMeta.Timezone,
				CountryCode:    rowMeta.CountryCode,
				Asn:            rowMeta.Asn,
				UsageType:      rowMeta.UsageType,
				Line:           rowMeta.Line,
				Comment:        rowMeta.Comment,
				Type:           rowMeta.Type,
			}
			if s.opt.CB != nil {
				meta.Extends = s.opt.CB(meta)
			}
			index = uint32(len(tmpMetaList))
			tmpMetaList = append(tmpMetaList, meta)
			metaTable[fp] = index
		}

		ipObj := rowMeta.StartIPObj()
		if t == IPV6 {
			if rowMeta.Mode() != IPV6 {
				// Bad IP metadata
				_, _ = fmt.Fprint(os.Stderr, "Bad IP metadata, start_ip=", rowMeta.StartIP)
				continue
			}
			entity := &IPV6Entity{
				startIndex: binary.BigEndian.Uint64(ipObj.To16()),
				endIndex:   binary.BigEndian.Uint64(rowMeta.EndIpObj().To16()),
				metaIndex:  index,
			}
			ipv6List = append(ipv6List, entity)
		} else {
			if rowMeta.Mode() != IPV4 {
				// Bad IP metadata
				_, _ = fmt.Fprint(os.Stderr, "Bad IP metadata, start_ip=", rowMeta.StartIP)
				continue
			}
			entity := &IPV4Entity{
				startIndex: binary.BigEndian.Uint32(ipObj.To4()),
				endIndex:   binary.BigEndian.Uint32(rowMeta.EndIpObj().To4()),
				metaIndex:  index,
			}
			ipv4List = append(ipv4List, entity)
		}
	}
	if err != io.EOF {
		return fmt.Errorf("unmarshal entity list error, %s", err)
	}

	if len(ipv4List) > 0 {
		s.v4Mu.Lock()
		s.ipv4EntityList = ipv4List
		s.v4MateList = tmpMetaList
		s.v4Mu.Unlock()
	}
	if len(ipv6List) > 0 {
		s.v6Mu.Lock()
		s.ipv6EntityList = ipv6List
		s.v6MateList = tmpMetaList
		s.v6Mu.Unlock()
	}

	// Clear Memory
	ipv4List = nil
	ipv6List = nil
	metaTable = nil
	tmpMetaList = nil

	return nil
}

// LoadData load data
func (s *Store) LoadData(opt Option) error {
	if len(opt.Files) <= 0 {
		return errors.New("no incoming data file")
	}
	s.opt = opt
	return s.update()
}

// Update update data
func (s *Store) Update() error {
	return s.update()
}

func (s *Store) update() error {
	var err error

	for _, fn := range s.opt.Files {
		// open file by filename
		var fReader *os.File
		if fReader, err = os.Open(fn.Path); err != nil {
			return err
		}
		err = s.UnmarshalFrom(fReader, fn.Type)
		if err != nil {
			return err
		}

		_ = fReader.Close()
	}

	return nil
}
