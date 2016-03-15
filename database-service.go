package main

import (
	"flag"
	"fmt"
	"log"

	//"github.com/GreenNav/service-database/database"
	"github.com/GreenNav/service-database/database/sqlite"
	"github.com/GreenNav/service-database/importer"
)

var dbTypeFlagDescription = "Determines which database to use. Can be sqlite or postgres"
var dbFileFlagDescription = "Sets the name of the database file to import to"
var dbSchemaFlagDescription = "File with the database setup queries"
var importFlagDescription = "Pbf file to import"

var dbTypeFlag = flag.String("dbtype", "sqlite", dbTypeFlagDescription)
var dbFileFlag = flag.String("dbfile", "gn.db", dbFileFlagDescription)
var dbSchemaFlag = flag.String("dbschema", "gn.db", dbSchemaFlagDescription)
var importFlag = flag.String("import", "", importFlagDescription)

func init() {
	flag.StringVar(dbTypeFlag, "d", "sqlite", dbTypeFlagDescription+" (shorthand)")
	flag.StringVar(dbFileFlag, "f", "gn.db", dbFileFlagDescription+" (shorthand)")
	flag.StringVar(dbSchemaFlag, "s", "./schemata/sqlite.sql", dbSchemaFlagDescription+" (shorthand)")
	flag.StringVar(importFlag, "i", "", importFlagDescription+" (shorthand)")
}

func main() {
	flag.Parse()

	if *importFlag == "" {
		log.Fatal("Please use the -import flag to import a .pbf file")
	}

	fmt.Println("Creating the database")
	db, err := sqlite.CreateEmpty(*dbFileFlag, *dbSchemaFlag)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Importing the file")
	err = importer.WriteToDatabase(*importFlag, db)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Closing database")
	db.Close()
	fmt.Println("Done!")
}
