package main

import (
	"container/list"
	"database/sql"
	_ "engine/pq"
	"fmt"
	"sync"
	"time"
)

var (
	db  *sql.DB
	mux sync.Mutex
)

func init() {
	mux.Lock()
	defer mux.Unlock()

	// check
	if db != nil {
		return
	}

	// open
	mysqldb, err := sql.Open("postgres", "host=localhost port=5432 user=postgres password=admin dbname=t2f sslmode=disable")
	checkErr(err)

	// new db
	db = mysqldb
}

func checkErr(err error) {
	if err != nil {
		panic("psql err: " + err.Error())
	}
}

var (
	users = list.New()
)

func main() {
	// select
	querySql := `select uid,goods_id from tb_user_goods where uid=$1`
	rows, err := db.Query(querySql, "28613428968424")
	checkErr(err)

	type usergoods struct {
		uid     string
		goodsid int
	}

	for rows.Next() {
		ug := &usergoods{}
		err = rows.Scan(
			&ug.uid,
			&ug.goodsid,
		)
		checkErr(err)
	    t := time.Now()
		users.PushFront(ug)
	    fmt.Println(time.Now().Sub(t))
	}

	for e := users.Front(); e != nil; e = e.Next() {
		fmt.Println(e.Value.(*usergoods).uid, e.Value.(*usergoods).goodsid)
	}
}
