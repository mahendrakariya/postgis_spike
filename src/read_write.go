package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

const READ_MAX = 900
const WRITE_MAX = 1000

func main() {
	db, err := sql.Open("postgres", "user=postgres password=123456 dbname=cabspike sslmode=disable")
	defer db.Close()

	if err != nil {
		log.Fatal(err)
	}

	finished := make(chan int)

	for i := 0; i < READ_MAX; i++ {
		go readDatabase(db, i, finished)
	}

	// for i := 0; i < WRITE_MAX; i++ {
	// go writeDatabase(db)
	// }
	var total_read_time int64
	for i := 0; i < READ_MAX; i++ {
		start := time.Now()
		<-finished
		total_read_time += time.Since(start).Nanoseconds()
		start = time.Now()
	}
	avg_time_ns := total_read_time / READ_MAX
	fmt.Printf("Avg read time: %v ms\n", avg_time_ns / int64(time.Millisecond))
	fmt.Printf("Total read time for %v connections: %v ms\n", READ_MAX, total_read_time/int64(time.Millisecond)) 
}

func readDatabase(db *sql.DB, i int, finished chan int) {
	rows, err := db.Query("SELECT count(*) from drivers")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var count int
	for rows.Next() {
		err = rows.Scan(&count)
	}
	finished <- count
	//fmt.Printf("Rows %v", count)
}

func writeDatabase(db *sql.DB) {
	result, err := db.Exec("UPDATE drivers set geog='POINT(77.6130574652 12.9018253405)' where driver_id=1")
	if err != nil {
		log.Fatal(err)
	}
	res, err := result.RowsAffected()
	fmt.Printf("RESULT %v", res)
}
