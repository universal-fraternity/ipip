// Package store provides data store definitions test file
package store

import (
	"net"
	"testing"
)

func TestLoad(t *testing.T) {
	st := NewStore()
	t.Log(st.LoadData(Option{
		Files: []FileInfo{{
			Path: "testdata/v6.txt",
			Type: IPV6,
		}},
	}))
}

func TestSearch(t *testing.T) {
	st := NewStore()
	t.Log(st.LoadData(Option{
		Files: []FileInfo{{Path: "testdata/v6.txt", Type: IPV6}, {Path: "testdata/v4.txt", Type: IPV4}},
	}))
	t.Log(st.Search(net.ParseIP("1.55.77.18")))
	t.Log(st.Search(net.ParseIP("1.55.29.242")))
	t.Log(st.Search(net.ParseIP("2001:506:100:4a::4000:0")))
	t.Log(st.Search(net.ParseIP("2001:506:100:40::2:1")))
}

func BenchmarkStore_Search(b *testing.B) {
	st := NewStore()
	if err := st.LoadData(Option{
		Files: []FileInfo{{Path: "testdata/v6.txt", Type: IPV6},
			{Path: "testdata/v4.txt", Type: IPV4}},
	}); err != nil {
		b.Error(err)
	}
	addr := net.ParseIP("1.54.192.168")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		st.Search(addr)
	}
}
