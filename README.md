# Database Service

This is going to be the database service/tool, that feeds the OSM data as well as other sources
 into different databases to use them for routing.

The code is currently written in Go and is able to extract routing data from PBF files and storing
it in a sqlite database.

## Setup Instructions:
  
Install golang and make sure it is present in the environment $PATH:    
Check ```go env``` to make sure everything is set.
  
1. Get this package
    ```bash
  go get -u github.com/Greennav/service-database
```

2. Fork the service-database repo and add new remote to your local copy of package
    ```bash
  cd $GOPATH/src/github.com/Greennav/service-database
  go get -u
  git remote add fork git@github.com:{GITHUB_USER}/service-database.git
```   

3. Go to the repository and run
    ```
  go build database-service.go
  ./database-service -i importer/monaco.osm.pbf -f test.db
  sqlite3 test.db
```
    ```SQL
  sqlite> select count(*) from nodes;
  978
```
  
4. Now you can push commits to your fork and then send PR
    ```bash
  git commit -a -m "Very important fix"
  git push fork
```  
  
## Current Status

The SQLite subpackage can import and write all nodes, ways and relations of a pbf file including tags. The nodes of a way are currently not imported somehow. Reading of the exported data from the database is completely missing; so is the REST-interface for queries and the PostgreSQL interface.

### Next Steps

- Make the import of the way's nodes work (sqlite/WriteWayNodes)
- Read the data from the database (sqlite/GetEverythingWithinCoordinates)
- Create a REST-interface and connect to sqlite/GetEverythingWithinCoordinates
