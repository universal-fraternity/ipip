// Package ipip IPIP Data Source Core Driver Test Code
package ipip

import (
	"testing"

	"github.com/universal-fraternity/ipip/store"
)

func TestSearch(t *testing.T) {
	if err := Load(Option{
		Files: []FileInfo{{Path: "store/testdata/v6.txt", Type: store.IPV6},
			{Path: "store/testdata/v4.txt", Type: store.IPV4}},
		CB: nil,
	}); err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(Search("1.55.29.242"))
}
