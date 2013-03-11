package main

import (
	"database/sql"
	_ "engine/pq"
	"fmt"
	"sync"
)

var (
	db  *sql.DB
	mux sync.Mutex
)

var userTableSql string = `
    create table if not exists user_profile
    (
        uid serial,
        name varchar(20) not null,
        created varchar(20) not null,
        primary key(uid)
    )
`

func init() {
	mux.Lock()
	defer mux.Unlock()

	// check
	if db != nil {
		return
	}

	// open
	mysqldb, err := sql.Open("postgres", "host=localhost port=4932 user=t2f_admin password=[32t2f_admin15] dbname=t2f sslmode=disable")
	checkErr(err)

	// new db
	db = mysqldb

	// create database table
	_, err = db.Exec(userTableSql)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		fmt.Println("psql err: " + err.Error())
	}
}

func main() {
	// insert
	insertSql := `insert into user_profile(name,created) values('admin1','2012-11-07') returning uid`
	//	result, err := db.Exec(insertSql)
	rows, err := db.Query(insertSql)
	checkErr(err)

	var a int
	for rows.Next() {
		err = rows.Scan(&a)
		checkErr(err)
		fmt.Println(a)
	}

	/*fmt.Println("------------------------")
	i, err := result.LastInsertId()
	checkErr(err)
	fmt.Println(i)
	fmt.Println("------------------------")*/

	// update
	/*	updateSql := `update user_profile set name=$1 where name=$2`
		stmt, err := db.Prepare(updateSql)
		checkErr(err)

		res, err := stmt.Exec("viney.chow", "viney")
		checkErr(err)

		affect, err := res.RowsAffected()
		checkErr(err)

		fmt.Println("exec update,rows affected: " + fmt.Sprint(affect))

		// select
		querySql := `select * from user_profile where name=$1`
		rows, err := db.Query(querySql, "viney.chow")

		type user struct {
			uid     int
			name    string
			created string
		}

		var u = &user{}
		for rows.Next() {
			err = rows.Scan(
				&u.uid,
				&u.name,
				&u.created)
			checkErr(err)
		}

		fmt.Println(*u)

		// delete
		deleteSql := `delete from user_profile where name=$1`
		stmt, err = db.Prepare(deleteSql)
		checkErr(err)

		res, err = stmt.Exec("viney.chow")
		checkErr(err)

		affect, err = res.RowsAffected()
		checkErr(err)

		fmt.Println("exec delete,rows affected: " + fmt.Sprint(affect))*/
}
