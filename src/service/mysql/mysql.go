package main

import (
	"database/sql"
	// _ "engine/go-mysql-driver/mysql"
	_ "engine/mymysql/godrv"
	"fmt"
	"sync"
	"time"
)

var (
	db  *sql.DB
	mux sync.Mutex
)

var userTableSql string = `
    create table if not exists user_profile
    (
        uid int primary key auto_increment,
        name varchar(20) not null,
        created varchar(20) not null
    )type=myisam
`

func init() {
	mux.Lock()
	defer mux.Unlock()

	// check
	if db != nil {
		return
	}

	// open
	//mysqldb, err := sql.Open("mymysql", "root:newding@tcp(192.168.1.22:3306)/test?charset=utf8")
	mysqldb, err := sql.Open("mymysql", "tcp:192.168.1.22:3306*test/root/newding")
	checkErr(err)

	// new db
	db = mysqldb

	// create database table
	_, err = db.Exec(userTableSql)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		panic("mysql err:" + err.Error())
	}
}

func main() {
	// insert
	insertSql := `insert into user_profile values(null,?,?)`
	stmt, err := db.Prepare(insertSql)
	checkErr(err)

	res, err := stmt.Exec("viney", time.Now().Format("2006-01-02 15:04:05"))
	checkErr(err)

	i, err := res.LastInsertId()
	checkErr(err)
	fmt.Println("exec insert,last insert id: " + fmt.Sprint(i))

	// update
	updateSql := `update user_profile set name=? where uid=?`
	stmt, err = db.Prepare(updateSql)
	checkErr(err)

	res, err = stmt.Exec("viney.chow", i)
	checkErr(err)

	affect, err := res.RowsAffected()
	checkErr(err)

	fmt.Println("exec update,rows affected: " + fmt.Sprint(affect))

	// select
	querySql := `select * from user_profile where uid=?`
	rows, err := db.Query(querySql, i)

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
	rows.Close()

	fmt.Println(*u)

	// delete
	deleteSql := `delete from user_profile where uid=?`
	stmt, err = db.Prepare(deleteSql)
	checkErr(err)

	res, err = stmt.Exec(i)
	checkErr(err)

	affect, err = res.RowsAffected()
	checkErr(err)

	fmt.Println("exec delete,rows affected: " + fmt.Sprint(affect))

	db.Close()
}
