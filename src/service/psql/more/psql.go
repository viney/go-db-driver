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

const createSql = `
CREATE TABLE IF NOT EXISTS 
        tb_plan
(
    	plan_id SERIAL,
    	plan_price NUMERIC,
    	payment_days TEXT, 
    	periods TEXT,
    	plan_create_time TEXT,
    	plan_count_price NUMERIC,
    	periods_giving_first NUMERIC,
    	periods_giving_next NUMERIC,
        privilege_id INT,
        PRIMARY KEY(plan_id)
);
CREATE TABLE IF NOT EXISTS 
        tb_privilege
(
    	privilege_id SERIAL,
    	privilege_create_time TEXT,
        vip_level INTEGER,
        vip_rate NUMERIC,
        ospf_rate NUMERIC,
        i18n BOOLEAN,
        family_number INTEGER,
        family_rate NUMERIC,
        preferential_period TEXT,
        preferential_period_rate NUMERIC,
        plan_id int
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
	mysqldb, err := sql.Open("postgres", "host=127.0.0.1 port=4932 dbname=t2f user=t2f_admin password=[32t2f_admin15] sslmode=disable")
	checkErr(err)

	// new db
	db = mysqldb

	// create database table
	_, err = db.Exec(createSql)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		panic("psql err: " + err.Error())
	}
}

func main() {
	fmt.Println("hello")
}
