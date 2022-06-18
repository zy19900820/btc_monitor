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
