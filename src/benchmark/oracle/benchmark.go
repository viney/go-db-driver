package main

import (
	"database/sql"
	_ "engine/goracle"
	"fmt"
	"testing"
)

var db *sql.DB

func bmSimpleQuery(b *testing.B) {
	var err error
	var num int64
	var str string
	var i int64
	var rows *sql.Rows
	for rep := 0; rep < 10000; rep++ {
		rows, err = db.Query("SELECT userid, name FROM user_profile")
		if err != nil {
			panic(err)
		}

		i = 0
		for rows.Next() {
			rows.Scan(&num, &str)
			if num != i {
				panic(fmt.Sprintf("Result didn't match: %d!=%d", num, i))
			}
			i++
		}

		if i != 100 {
			panic(fmt.Sprintf("Rows count doesn't match: %d!=100", i))
		}
	}
}

func bmPreparedQuery(b *testing.B) {
	stmt, err := db.Prepare("SELECT userid, name FROM user_profile")
	if err != nil {
		panic(err)
	}

	var num int64
	var str string
	var i int64
	var rows *sql.Rows
	for rep := 0; rep < 10000; rep++ {
		rows, err = stmt.Query()
		if err != nil {
			panic(err)
		}

		i = 0
		for rows.Next() {
			rows.Scan(&num, &str)
			if num != i {
				panic(fmt.Sprintf("Result didn't match: %d!=%d", num, i))
			}
			i++
		}

		if i != 100 {
			panic(fmt.Sprintf("Rows count doesn't match: %d!=100", i))
		}
	}

	stmt.Close()
}

func bmAutoQueryRow(b *testing.B) {
	var err error
	var num int64
	var str string
	var i int64
	for rep := 0; rep < 1000; rep++ {
		for i = 0; i < 100; i++ {
			err = db.QueryRow("SELECT * FROM user_profile WHERE userid=:1", i).Scan(&num, &str)
			if err != nil {
				panic(err)
			}
			if num != i {
				panic(fmt.Sprintf("Result didn't match: %d!=%d", num, i))
			}
		}
	}
}

func bmSimpleQueryRow(b *testing.B) {
	var err error
	var num int64
	var str string
	var i int64
	for rep := 0; rep < 1000; rep++ {
		for i = 0; i < 100; i++ {
			err = db.QueryRow(fmt.Sprintf("SELECT * FROM user_profile WHERE userid=%d", i)).Scan(&num, &str)
			if err != nil {
				panic(err)
			}
			if num != i {
				panic(fmt.Sprintf("Result didn't match: %d!=%d", num, i))
			}
		}
	}
}

func bmPreparedQueryRow(b *testing.B) {
	var err error
	var stmt *sql.Stmt
	var num int64
	var str string
	var i int64
	for rep := 0; rep < 1000; rep++ {
		stmt, err = db.Prepare("SELECT * FROM user_profile WHERE userid=:1")
		if err != nil {
			panic(err)
		}
		for i = 0; i < 100; i++ {
			stmt.QueryRow(i).Scan(&num, &str)
			if err != nil {
				panic(err)
			}
			if num != i {
				panic(fmt.Sprintf("Result didn't match: %d!=%d", num, i))
			}
		}
		stmt.Close()
	}
}

/*
func bmSimpleExec(b *testing.B) {
	var err error
	for i := 0; i < 100000; i++ {
		_, err = db.Exec("SET @test_var=1")
		if err != nil {
			panic(err)
		}
	}
}

func bmPreparedExec(b *testing.B) {
	stmt, err := db.Prepare("SET @test_var=1")
	if err != nil {
		panic(err)
	}

	for i := 0; i < 100000; i++ {
		stmt.Exec()
		if err != nil {
			panic(err)
		}
	}
	stmt.Close()
}
*/

func runBenchmark(name, driver, dsn string) {
	var br testing.BenchmarkResult
	var err error

	fmt.Println("*************************************************************")
	fmt.Println("   BENCHMARKING " + name)
	fmt.Println("*************************************************************")
	fmt.Println()

	db, err = sql.Open(driver, dsn)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec(createTableSql)
	if err != nil {
		panic(err)
	}

	// Insert data
	stmt, err := db.Prepare("INSERT INTO user_profile (userid, name) VALUES(:1, :2)")
	if err != nil {
		panic(err)
	}
	for i := 0; i < 100; i++ {
		_, err = stmt.Exec(i, "Test")
		if err != nil {
			panic(err)
		}
	}
	stmt.Close()

	fmt.Println("-------------------------------------------------------------")
	fmt.Println("   [10000 * Query 100 Rows]")
	fmt.Println("-------------------------------------------------------------")

	// SimpleQuery
	br = testing.Benchmark(bmSimpleQuery)
	fmt.Printf("SimpleQuery: %s [ %.0f queries/second ]\r\n", br.T.String(), (10000 / br.T.Seconds()))

	// PreparedQuery
	br = testing.Benchmark(bmPreparedQuery)
	fmt.Printf("PreparedQuery: %s [ %.0f queries/second ]\r\n", br.T.String(), (10000 / br.T.Seconds()))

	fmt.Println()
	fmt.Println("-------------------------------------------------------------")
	fmt.Println("   [100 * QueryRow] * 1000")
	fmt.Println("-------------------------------------------------------------")

	// AutoQueryRow
	br = testing.Benchmark(bmAutoQueryRow)
	fmt.Printf("AutoQueryRow: %s [ %.0f queries/second ]\r\n", br.T.String(), (100 * 1000 / br.T.Seconds()))

	// SimpleQueryRow
	br = testing.Benchmark(bmSimpleQueryRow)
	fmt.Printf("SimpleQueryRow: %s [ %.0f queries/second ]\r\n", br.T.String(), (100 * 1000 / br.T.Seconds()))

	// PreparedQueryRow
	br = testing.Benchmark(bmPreparedQueryRow)
	fmt.Printf("PreparedQueryRow: %s [ %.0f queries/second ]\r\n", br.T.String(), (100 * 1000 / br.T.Seconds()))

	/*fmt.Println()
	fmt.Println("-------------------------------------------------------------")
	fmt.Println("   [100000 * Exec]")
	fmt.Println("-------------------------------------------------------------")

	// SimpleExec
	br = testing.Benchmark(bmSimpleExec)
	fmt.Printf("SimpleExec: %s [ %.0f queries/second ]\r\n", br.T.String(), (100000 / br.T.Seconds()))

	// PreparedExec
	br = testing.Benchmark(bmPreparedExec)
	fmt.Printf("PreparedExec: %s [ %.0f queries/second ]\r\n", br.T.String(), (100000 / br.T.Seconds()))

	fmt.Println()
	fmt.Println()
	*/
}

var createTableSql string = `
BEGIN
    BEGIN
         EXECUTE IMMEDIATE 'DROP TABLE user_profile';
    EXCEPTION
         WHEN OTHERS THEN
                IF SQLCODE != -942 THEN
                     RAISE;
                END IF;
    END;
    EXECUTE IMMEDIATE 'CREATE TABLE user_profile (userid number(3) PRIMARY KEY, name varchar(4) NOT NULL)';
END;
`

func main() {
	runBenchmark("GORACLE [run 1]", "goracle", "viney/admin@sqtcall")
	// runBenchmark("GORACLE [run 2]", "goracle", "viney/admin@sqtcall")
}
