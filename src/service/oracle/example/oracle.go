package main

import (
	"database/sql"
	_ "engine/go-oci8"
	"fmt"
	"os"
)

func main() {
	os.Setenv("NLS_LANG", "")

	db, err := sql.Open("oci8", "viney/admin@sqtcall")
	if err != nil {
		fmt.Println(err)
		return
	}
	rows, err := db.Query("select 3.14, 'foo', 100 from dual")
	if err != nil {
		fmt.Println(err)
		return
	}
	for rows.Next() {
		var f1 float64
		var f2 string
		var f3 int
		rows.Scan(&f1, &f2, &f3)
		println(f1, f2, f3) // 3.14 foo
	}
	rows.Close()
	_, err = db.Exec("create table foo(bar varchar2(256))")
	stmt, err := db.Prepare("insert into foo(bar) values(:1)")
	_, err = stmt.Exec("viney")

	//_, err = db.Exec("drop table foo")
	if err != nil {
		fmt.Println(err)
		return
	}

	db.Close()
}
