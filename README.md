# Database Service

Database service in Go for the GreenNav project    
This repository contains code which follows the [Design plan for database service](https://github.com/Greennav/greennav.github.io/blob/master/wiki/Roadmap.md#design-plan)

## Setup Instructions:
  
Install golang and make sure it is present in the environment $PATH:    
Check ```go env``` to make sure everything is set.
  
1. Get go-json-rest package
    ```
    go get -u github.com/omniscale/imposm3
    go get -u github.com/mattn/go-sqlite3
    ```
2. Fork the service-database repo and clone it to your $gopath/github.com/{username}/    
3. Go to the repository and run
    ```
    go build service-database.go
    ./database-service -i importer/monaco.osm.pbf -f test.db
    sqlite3 test.db
    ```
    ```SQL
    sqlite> select count(*) from nodes;
    978
    ```    
    
## Current Status

The SQLite subpackage can import and write all nodes, ways and relations of a pbf file including tags. The nodes of a way are currently not imported somehow. Reading of the exported data from the database is completely missing; so is the REST-interface for queries and the PostgreSQL interface.

### Next Steps

- Make the import of the way's nodes work (sqlite/WriteWayNodes)
- Read the data from the database (sqlite/GetEverythingWithinCoordinates)
- Create a REST-interface and connect to sqlite/GetEverythingWithinCoordinates