package sql

import (
	"fmt"
	"log"
	"time"

	"btc_monitor/netApi"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"golang.org/x/xerrors"
)

var db *sqlx.DB

func createTable(name, password, ip string, port int, dbName string) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/?charset=utf8", name, password, ip, port)
	mysqlDb, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer mysqlDb.Close()

	err = mysqlDb.Ping()
	if err != nil {
		return err
	}

	createOrderTables(mysqlDb, dbName)

	return nil
}

//operation about mysql
func InitMysqlDB(name, password, ip string, port int, dbName string) error {
	err := createTable(name, password, ip, port, dbName)
	if err != nil {
		return err
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", name, password, ip, port, dbName)
	mysqlDb, err := sqlx.Open("mysql", dsn)
	//db.SetMaxIdleConns(conf.Mysql.Conn.MaxIdle)
	//db.SetMaxOpenConns(conf.Mysql.Conn.Maxopen)
	if err != nil {
		return err
	}

	db = mysqlDb
	db.SetConnMaxLifetime(30 * time.Second)

	return nil
}

const TABLE_ADDRESS_INFO = "t_address_info"

func CheckNetInit() (bool, error) {
	var count int
	querySql := "SELECT count(*) FROM " + TABLE_ADDRESS_INFO
	err := db.Get(&count, querySql)
	if err != nil {
		return false, err
	}
	if count == 0 {
		return true, nil
	} else if count > 200 {
		return false, nil
	} else {
		return false, xerrors.New("db count less 200")
	}

}

func InsertNetInfos(netInfos []netApi.LOCAL_BTC_ADDR_INFO) error {
	tx, err := db.Beginx() // 开启事务
	if err != nil {
		if tx != nil {
			tx.Rollback()
		}
		return err
	}

	timestamp := time.Now().Unix()
	timeStr := time.Now().Format("2006-01-02 15:04:05")
	for _, info := range netInfos {
		insertSql := "INSERT INTO " + TABLE_ADDRESS_INFO + " (address,alias,timestamp,time,value) " +
			" VALUES(?,?,?,?,?)"
		_, err = tx.Exec(insertSql, info.Addr, info.Alias, timestamp, timeStr, info.Count)
		if err != nil {
			log.Println(insertSql, " addr:", info.Addr, " alias:", info.Alias, " value:", info.Count)
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit() // 提交事务
	if err != nil {
		tx.Rollback() // 回滚
		return err
	}

	return nil
}
