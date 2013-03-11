package main

import (
	"database/sql"
	_ "engine/go-oci8"
	"fmt"
	"os"
	"sync"
)

var (
	db  *sql.DB
	mux sync.Mutex
)

var userTableSql string = `
BEGIN
    BEGIN
         EXECUTE IMMEDIATE 'DROP TABLE user_profile';
    EXCEPTION
         WHEN OTHERS THEN
                IF SQLCODE != -942 THEN
                     RAISE;
                END IF;
    END;
    EXECUTE IMMEDIATE 'CREATE TABLE user_profile (userid NUMBER(10) NOT NULL, name VARCHAR(20) NOT NULL, created VARCHAR(20) NOT NULL)';
END;
`

func init() {
	mux.Lock()
	defer mux.Unlock()

	os.Setenv("NLS_LANG", "AMERICAN_AMERICA.ZHS16GBK")
	// check
	if db != nil {
		return
	}

	// open
	oracledb, err := sql.Open("oci8", "viney/admin@sqtcall")
	checkErr(err)

	// new db
	db = oracledb

	// create database table
	_, err = db.Exec(userTableSql)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		panic("oracle err:" + err.Error())
	}
	return
}

func main() {
	// insert
	insertSql := `insert into user_profile(userid,name,created) values(1,'viney','2013-03-06')`
	_, err := db.Exec(insertSql)
	checkErr(err)

	// update
	updateSql := `update user_profile set name='南海' where userid=1`
	_, err = db.Exec(updateSql)
	checkErr(err)

	// select
	querySql := `select userid,name,created from user_profile`
	rows, err := db.Query(querySql)

	type user struct {
		userid  int
		name    string
		created string
	}

	var u = &user{}
	for rows.Next() {
		err = rows.Scan(
			&u.userid,
			&u.name,
			&u.created)
		checkErr(err)
	}
	rows.Close()

	fmt.Println(*u)

	// delete
	deleteSql := `delete from user_profile where userid=1`
	_, err = db.Exec(deleteSql)
	checkErr(err)

	db.Close()
}
