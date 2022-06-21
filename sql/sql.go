package sql

import (
	"fmt"
	"log"
	"strings"
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
const TABLE_TOP_DAILY_INFO = "t_top_daily_info"

func CheckNetInit() (bool, error) {
	var count int
	querySql := "SELECT count(*) FROM " + TABLE_ADDRESS_INFO
	err := db.Get(&count, querySql)
	if err != nil {
		return false, err
	}
	if count == 0 {
		return true, nil
	} else if count >= 200 {
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

func UpdateNetInfos(netInfos []netApi.LOCAL_BTC_ADDR_INFO) error {
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
		updateSql := "update " + TABLE_ADDRESS_INFO + " set value = ?, timestamp = ?, time = ? where address = ? "
		result, err := tx.Exec(updateSql, info.Count, timestamp, timeStr, info.Addr)
		if err != nil {
			log.Println(updateSql, " value:", info.Count, " timestamp:", timestamp, " time:", timeStr, " address:", info.Addr)
			tx.Rollback()
			return err
		}
		affectRow, _ := result.RowsAffected()
		if affectRow == 0 {
			insertSql := "INSERT INTO " + TABLE_ADDRESS_INFO + " (address,alias,timestamp,time,value) " +
				" VALUES(?,?,?,?,?)"
			_, err = tx.Exec(insertSql, info.Addr, info.Alias, timestamp, timeStr, info.Count)
			if err != nil {
				log.Println(insertSql, " addr:", info.Addr, " alias:", info.Alias, " value:", info.Count)
				//这里不知道为什么会报错
				if strings.Contains(err.Error(), "Duplicate entry") {
					continue
				}
				tx.Rollback()
				return err
			}
		}
	}

	err = tx.Commit() // 提交事务
	if err != nil {
		tx.Rollback() // 回滚
		return err
	}

	return nil
}

func GetSqlTop200Addr() ([]string, error) {
	var adds []string
	querySql := "SELECT address FROM " + TABLE_ADDRESS_INFO + " order by value desc limit 200"
	err := db.Select(&adds, querySql)
	if err != nil {
		log.Println(querySql)
		return adds, err
	}
	return adds, nil
}

func UpdateDaily(num int) error {
	dateTime := time.Now().Format("2006-01-02")
	var count int
	querySql := "SELECT count(*) FROM " + TABLE_TOP_DAILY_INFO + " WHERE time = ?"
	err := db.Get(&count, querySql, dateTime)
	if err != nil {
		log.Println(querySql)
		return err
	}
	if count > 3 {
		return xerrors.New("今天已经更新")
	}

	var values []float64
	querySql = "SELECT value FROM " + TABLE_ADDRESS_INFO + " WHERE alias = \"\" order by value desc limit ?"
	err = db.Select(&values, querySql, num)
	if err != nil {
		log.Println(querySql)
		return err
	}
	total := 0.0
	for _, value := range values {
		total += value
	}

	insertSql := "INSERT INTO " + TABLE_TOP_DAILY_INFO + " (tag,time,value) " +
		" VALUES(?,?,?)"
	_, err = db.Exec(insertSql, num, dateTime, total)
	if err != nil {
		log.Println(insertSql, " tag:", num, " time:", dateTime, " value:", total)
		return err
	}
	return nil
}
