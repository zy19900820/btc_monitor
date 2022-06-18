package sql

import (
	"github.com/jmoiron/sqlx"
)

var db_create = `
create database if not exists btc;
`
var t_address_info = `
create table if not exists t_address_info
(
     id bigint(20) not NULL AUTO_INCREMENT primary key,
     address varchar(64) not null unique,
	 alias varchar(64) not null default "",
	 tag varchar(64) not null default "",
	 time varchar(64) not null default "",
	 timestamp bigint(20) not null,
	 value decimal(26, 8) not null default 0
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
`

func createOrderTables(db *sqlx.DB, dbName string) {
	db.MustExec(db_create)
	db.MustExec("USE " + dbName)
	db.MustExec(t_address_info)
}
