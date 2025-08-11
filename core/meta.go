// Package core provides data core logics for handling IPV4/6 address.
package core

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/universal-fraternity/ipip/utils"
)

const (
	Unknown = iota
	IPV4
	IPV6
)

// RowMeta define row metadata
type RowMeta struct {
	StartIP        string      // The starting point of the IP segment is divided into IP addresses
	EndIP          string      // The termination point of the IP segment is divided into IP addresses
	Country        string      // Location Information - Country
	Province       string      // Location Information - Province
	City           string      // Location Information - City
	Region         string      // Location Information - Region
	OwnerDomain    string      // Owner Domain Name
	IspDomain      string      // Operator Domain Name
	ChinaAdminCode int32       // Administrative division code of China
	Latitude       float64     // Latitude
	Longitude      float64     // Longitude
	Timezone       string      // Time Zone ID
	CountryCode    string      // Country code
	Asn            []int64     // AS ID number
	UsageType      string      // Application scenarios
	Line           string      // national line
	Comment        *string     // Other remarks information
	Type           *string     // Network type
	Extends        interface{} // Extended Information
	startIPObj     net.IP      // IP starting position
	endIPObj       net.IP      // IP end position
}

// Meta Definition of metadata results (location information)
type Meta struct {
	Country        string      // Location Information - Country
	Province       string      // Location Information - Province
	City           string      // Location Information - City
	Region         string      // Location Information - Region
	OwnerDomain    string      // Owner Domain Name
	IspDomain      string      // Operator Domain Name
	ChinaAdminCode int32       // Administrative division code of China
	Latitude       float64     // Latitude
	Longitude      float64     // Longitude
	Timezone       string      // Time Zone ID
	CountryCode    string      // Country code
	Asn            []int64     // AS ID number
	UsageType      string      // Application scenarios
	Line           string      // national line
	Comment        *string     // Other remarks information
	Type           *string     // Network type
	Extends        interface{} // Extended Information
}

// NewMeta Return a new meta
func NewMeta() *Meta {
	return &Meta{}
}

// WithExtends Set extension information
func (m *Meta) WithExtends(en interface{}) {
	m.Extends = en
}

// IsEmpty If all fields are empty, return true.
func (m Meta) IsEmpty() bool {
	return m.Country == "" &&
		m.Province == "" &&
		m.City == "" &&
		m.Region == "" &&
		m.OwnerDomain == "" &&
		m.IspDomain == "" &&
		m.Latitude == 0 &&
		m.Longitude == 0 &&
		m.Timezone == "" &&
		m.CountryCode == "" &&
		len(m.Asn) == 0
}

// IsEmpty If all fields are empty, return true.
func (r *RowMeta) IsEmpty() bool {
	return r.StartIP == "" &&
		r.EndIP == "" &&
		r.Country == "" &&
		r.Province == "" &&
		r.City == "" &&
		r.Region == "" &&
		r.OwnerDomain == "" &&
		r.IspDomain == "" &&
		r.Latitude == 0 &&
		r.Longitude == 0 &&
		r.Timezone == "" &&
		r.CountryCode == "" &&
		len(r.Asn) == 0
}

// Unmarshal Parse detailed information into meta format.
func (r *RowMeta) Unmarshal(buffer []byte, t int) error {
	if t == IPV6 {
		return r.UnmarshalV6(buffer)
	} else if t == IPV4 {
		return r.UnmarshalV4(string(buffer))
	}
	return errors.New("unknown data type")
}

// UnmarshalV4 Parse detailed information into meta format.
func (r *RowMeta) UnmarshalV4(row string) error {
	var e error
	if r != nil {
		for i, item := range strings.Split(strings.TrimSuffix(row, "\n"), "\t") {
			switch i {
			case 0:
				r.StartIP = item
				r.startIPObj = net.ParseIP(item)
			case 1:
				r.EndIP = item
				r.endIPObj = net.ParseIP(item)
			case 2:
				r.Country = item
			case 3:
				r.Province = utils.RefineOutput(item)
			case 4:
				r.City = utils.RefineOutput(item)
			case 5:
				r.Region = utils.RefineOutput(item)
			case 6:
				r.OwnerDomain = utils.RefineOutput(item)
			case 7:
				r.IspDomain = utils.RefineOutput(item)
			case 8:
				r.ChinaAdminCode, _ = utils.String2Int32(item)
			case 9:
				r.Latitude, _ = utils.String2Float64(item)
			case 10:
				r.Longitude, _ = utils.String2Float64(item)
			case 11:
				r.Timezone = utils.RefineOutput(item)
			case 12:
				r.CountryCode = utils.RefineOutput(item)
			case 13:
				if r.Asn, e = parseAsn(item); e != nil {
					return e
				}
			case 14:
				r.UsageType = utils.RefineOutput(item)
			case 15:
				r.Line = utils.RefineOutput(item)
			}
		}
	} else {
		return errors.New("meta is null")
	}
	return e
}

// UnmarshalV6 Parse detailed information into meta format.
func (r *RowMeta) UnmarshalV6(buffer []byte) error {
	var e error
	if r != nil {
		for i, item := range strings.Split(strings.TrimSuffix(string(buffer), "\n"), "\t") {
			switch i {
			case 0:
				r.StartIP = item
				_, ipv6Net, err := net.ParseCIDR(item)
				if err != nil {
					fmt.Println("ParseCIDR failed. err=", err.Error())
					break
				}
				r.startIPObj = ipv6Net.IP
				r.endIPObj = utils.LastIP(ipv6Net)
			case 1:
				r.Country = item
			case 2:
				r.Province = utils.RefineOutput(item)
			case 3:
				r.City = utils.RefineOutput(item)
			case 4:
				r.Region = utils.RefineOutput(item)
			case 5:
				r.OwnerDomain = utils.RefineOutput(item)
			case 6:
				r.IspDomain = utils.RefineOutput(item)
			case 7:
				r.ChinaAdminCode, _ = utils.String2Int32(item)
			case 8:
				r.Latitude, _ = utils.String2Float64(item)
			case 9:
				r.Longitude, _ = utils.String2Float64(item)
			case 10:
				r.Timezone = utils.RefineOutput(item)
			case 11:
				r.CountryCode = utils.RefineOutput(item)
			case 12:
				if r.Asn, e = parseAsn(item); e != nil {
					return e
				}
			case 13:
				r.UsageType = utils.RefineOutput(item)
			case 14:
				r.Line = utils.RefineOutput(item)
			}
		}
	} else {
		return errors.New("meta is null")
	}
	return e
}

// Hash hash calculation
func (r *RowMeta) Hash() string {
	if r.IsEmpty() {
		return ""
	}
	h1 := sha1.New()
	source := fmt.Sprintf("%s%s%s%s%s%s%d%g%g%s%s%d%s%s", r.Country, r.Province, r.City, r.Region,
		r.OwnerDomain, r.IspDomain, r.ChinaAdminCode, r.Latitude, r.Longitude,
		r.Timezone, r.CountryCode, r.Asn, r.UsageType, r.Line)
	_, _ = io.WriteString(h1, source)
	return string(h1.Sum(nil))
}

// Mode Type judgment
func (r *RowMeta) Mode() int {
	if r.startIPObj.To4() != nil {
		return IPV4
	} else if r.startIPObj.To16() != nil {
		return IPV6
	}
	return Unknown
}

// StartIPObj Starting IP
func (r *RowMeta) StartIPObj() net.IP {
	return r.startIPObj
}

// EndIpObj End IP
func (r *RowMeta) EndIpObj() net.IP {
	return r.endIPObj
}

// String  Format output
func (m *Meta) String() string {
	if m != nil {
		return fmt.Sprintf("country:%s province:%s city:%s region:%s"+
			"owner_domain:%s isp_domain:%s china_admin_code:%d latitude:%g Longitude:%g"+
			"timezone:%s country_code:%s asn:%d usage_type:%s line:%s",
			m.Country, m.Province, m.City, m.Region,
			m.OwnerDomain, m.IspDomain, m.ChinaAdminCode, m.Latitude, m.Longitude,
			m.Timezone, m.CountryCode, m.Asn, m.UsageType, m.Line)
	}
	return ""
}

// MarshalString Serialize data entities into strings.
func (m *Meta) MarshalString() (string, error) {
	var line string
	if m != nil {
		line = fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%d\t%g\t%g\t%s\t%s\t%d\t%s\t%s",
			m.Country, m.Province, m.City, m.Region,
			m.OwnerDomain, m.IspDomain, m.ChinaAdminCode, m.Latitude, m.Longitude,
			m.Timezone, m.CountryCode, m.Asn, m.UsageType, m.Line)
	}
	return line, nil
}

// Marshal Serialize data entities into strings.
func (m *Meta) Marshal() ([]byte, error) {
	line, err := m.MarshalString()
	return []byte(line), err
}

func parseAsn(s string) ([]int64, error) {
	asn := make([]int64, 0)
	array := strings.Split(s, ",")
	for _, item := range array {
		val, err := utils.String2Int64(item)
		if err != nil {
			return nil, err
		}
		asn = append(asn, val)
	}
	return asn, nil
}
