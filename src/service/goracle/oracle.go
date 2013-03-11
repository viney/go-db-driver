package main

import (
	"database/sql"
	_ "engine/goracle"
	"fmt"
	"os"
	"sync"
	"time"
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
    EXECUTE IMMEDIATE 'CREATE TABLE user_profile (userid NUMBER(10) PRIMARY KEY, name VARCHAR(20) NOT NULL, created VARCHAR(20) NOT NULL)';
END;
`

func init() {
	mux.Lock()
	defer mux.Unlock()

	// os.Setenv("NLS_LANG", "AMERICAN_AMERICA.ZHS16GBK")
	os.Setenv("NLS_LANG", "AMERICAN_AMERICA.AL32UTF8")

	// check
	if db != nil {
		return
	}

	// open
	oracledb, err := sql.Open("goracle", "viney/admin@sqtcall")
	checkErr(err)

	// new db
	db = oracledb

	// create database table
	result, err := db.Exec(userTableSql)
	resultErr(result, err)
}

func checkErr(err error) {
	if err != nil {
		panic("oracle err:" + err.Error())
	}
	return
}

func resultErr(res sql.Result, err error) (lastInsertId, rowsAffected int64) {
	if err != nil {
		panic("oracle err: " + err.Error())
	}

	lastInsertId, err = res.LastInsertId()
	if err != nil {
		//panic("res.LastInsertId: " + err.Error())
	}

	rowsAffected, err = res.RowsAffected()
	if err != nil {
		panic("res.RowsAffected: " + err.Error())
	}
	fmt.Println(lastInsertId, rowsAffected)

	return lastInsertId, rowsAffected
}

func main() {
	// insert
	tx, err := db.Begin()
	checkErr(err)

	insertSql := `insert into user_profile(userid,name,created) values(:1,:2,:3)`
	stmt, err := tx.Prepare(insertSql)
	result, err := stmt.Exec(1, "viney", time.Now().Format("2006-01-02 15:04:05"))
	resultErr(result, err)

	// update
	updateSql := `update user_profile set name=:1 where userid=1`
	stmt, err = tx.Prepare(updateSql)
	result, err = stmt.Exec("南海")
	resultErr(result, err)

	// select
	querySql := `select userid,name,created from user_profile`
	rows, err := tx.Query(querySql)

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
	/*	deleteSql := `delete from user_profile where userid=1`
		_, err = db.Exec(deleteSql)
		checkErr(err) */

	err = tx.Commit()
	checkErr(err)

	db.Close()
}
