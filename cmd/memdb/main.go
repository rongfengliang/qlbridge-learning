package main

import (
	"database/sql"
	"database/sql/driver"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/araddon/qlbridge/datasource/memdb"
	_ "github.com/araddon/qlbridge/qlbdriver"
	"github.com/araddon/qlbridge/schema"

	u "github.com/araddon/gou"
	"github.com/araddon/qlbridge/expr/builtins"
)

var (
	logging = "info"
)

func init() {

	flag.StringVar(&logging, "logging", "info", "logging [ debug,info ]")
	flag.Parse()
	u.SetupLogging(logging)
	u.SetColorOutput()
}

func main() {
	builtins.LoadAllBuiltins()
	inrow := []driver.Value{"dalong", 2222, "v2"}
	inrow2 := []driver.Value{"dalong1", 2222, "v2"}
	memdb, err := memdb.NewMemDbData("demoapp", [][]driver.Value{inrow, inrow2}, []string{"name", "age", "version"})
	if err != nil {
		log.Fatalln("memdb error", err.Error())
	}
	schema.RegisterSourceAsSchema("demoapp", memdb)
	db, err := sql.Open("qlbridge", "demoapp")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()
	// some insert ops
	_, err = db.Exec(`insert into demoapp(name,age,version) values('dalongdemo',333,'v3'))`)
	if err != nil {
		log.Fatalln("insert errpr")
	}
	// query
	rows, err := db.Query("select name,age,version,now() as now from demoapp")
	if err != nil {
		u.Errorf("could not execute query: %v", err)
		return
	}
	defer rows.Close()
	cols, _ := rows.Columns()
	readCols := make([]interface{}, len(cols))
	writeCols := make([]string, len(cols))
	for i := range writeCols {
		readCols[i] = &writeCols[i]
	}
	fmt.Printf("\n\nScanning through memdb: (%v)\n\n", strings.Join(cols, ","))
	for rows.Next() {
		rows.Scan(readCols...)
		fmt.Println(strings.Join(writeCols, ", "))
	}
	fmt.Println("")
}
