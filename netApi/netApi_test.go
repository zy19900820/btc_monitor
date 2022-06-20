package netApi

import (
	"log"
	"testing"
)

func TestGetLatestAddrInfo(t *testing.T) {
	infos, err := GetLatestAddrInfo()
	log.Println(infos)
	log.Println(err)
}

func TestGetMulAddressInfo(t *testing.T) {
	var strs []string
	strs = append(strs, "a")
	strs = append(strs, "b")
	GetMulAddressInfo(strs)
}

func TestParse(t *testing.T) {
	Parse()
}
