package sql

import (
	"log"
	"testing"
)

func getDBConf() (string, string, string, int, string) {
	return "root", "Xjz!1234", "127.0.0.1", 3306, "btc"
}
func TestChangeUserImage(t *testing.T) {
	err := InitMysqlDB(getDBConf())
	if err != nil {
		t.Error(err)
		return
	}

	bNeedInit, err := CheckNetInit()
	log.Println(bNeedInit)
	log.Println(err)
}
