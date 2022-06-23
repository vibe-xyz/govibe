package models

import (
	"encoding/hex"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"unicode"

	"github.com/ethereum/go-ethereum/crypto"
)

type TagOptions struct {
	Skip      bool
	Name      string
	Omitempty bool
	Omitzero  bool
}

// chk element in array
func InArray(needle interface{}, hystack interface{}) bool {
	if harr, ok := ToSlice(hystack); ok {
		for _, item := range harr {
			if item == needle {
				return true
			}
		}
	}
	return false
}

// convert to array
func ToSlice(arr interface{}) ([]interface{}, bool) {
	v := reflect.ValueOf(arr)
	if v.Kind() != reflect.Slice {
		return nil, false
	}
	l := v.Len()
	ret := make([]interface{}, l)
	for i := 0; i < l; i++ {
		ret[i] = v.Index(i).Interface()
	}
	return ret, true
}

func LowerCaseWithUnderscores(name string) string {
	newName := []rune{}
	for i, c := range name {
		if i == 0 {
			newName = append(newName, unicode.ToLower(c))
		} else {
			if unicode.IsUpper(c) {
				newName = append(newName, '_')
				newName = append(newName, unicode.ToLower(c))
			} else {
				newName = append(newName, c)
			}
		}
	}
	return string(newName)
}

func GetTagOptions(tag reflect.StructTag, tagname string) TagOptions {
	t := tag.Get(tagname)
	if t == "-" {
		return TagOptions{Skip: true}
	}
	var opts TagOptions
	parts := strings.Split(t, ",")
	opts.Name = strings.Trim(parts[0], " ")
	for _, s := range parts[1:] {
		switch strings.Trim(s, " ") {
		case "omitempty":
			opts.Omitempty = true
		case "omitzero":
			opts.Omitzero = true
		}
	}
	return opts
}

func GetSvrmark(svrname string, serverid ...string) string {
	hostname, _ := os.Hostname()
	if pidx := strings.Index(string(hostname), "."); pidx > 0 {
		hostname = string([]byte(hostname)[:pidx-1])
	}
	if len(serverid) > 0 && len(serverid[0]) > 0 {
		return fmt.Sprintf("%s-%s", svrname, serverid[0])
	}
	pid := os.Getpid()
	return fmt.Sprintf("%s-%s-%d", hostname, svrname, pid)
}

func ChecksumAddress(address string) string {
	address = strings.Replace(strings.ToLower(address), "0x", "", 1)
	crypto.Keccak256([]byte(address))
	_hash := hex.EncodeToString(crypto.Keccak256([]byte(address)))
	_address := "0x"

	for k, v := range address {
		l, _ := strconv.ParseInt(string(_hash[k]), 16, 16)
		if l > 7 {
			_address += strings.ToUpper(string(v))
		} else {
			_address += string(v)
		}
	}
	return _address
}

func If(condition bool, trueVal, falseVal interface{}) interface{} {
	if condition {
		return trueVal
	}
	return falseVal
}
