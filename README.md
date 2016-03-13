# service-database
Database service in Go for the GreenNav project    
This repository contains code which follows [Design plan for database service](https://github.com/Greennav/greennav.github.io/blob/master/wiki/Roadmap.md#design-plan)
## Setup Instructions:    
Install golang and make sure it is present in the environment $PATH:    
Check ``` go env     ``` to make sure everything is set.    
1. Get go-json-rest package 
    ```
    go get github.com/ant0ine/go-json-rest/rest
    ```
2. Fork the service-database repo and clone it to your $gopath/github.com/{username}/    
3. Go to the repository and run    
    ```
    go build main.go
    ./main
    ```     
  
  
If you want to test database connectivity uncomment the SetupPGDatabase() and ImportFromHttp() calls from the main.go

#### Database setup:     
```
CREATE TABLE info (
    id serial unique,
    version integer,
    "timestamp" text,
    changeset bigint,
    uid integer,
    "user" varchar(40),
    visible boolean
);

CREATE TABLE node (
    id bigint,
    lat double precision,
    lon double precision,
    tags json,
    infoid int,
    foreign key(infoid) references info(id)
);
```   

Current Status:    
The program can read the data from pbf file (file because the online file is around 15mb and takes some time to get fetch).     
The program identifies the nodes from the source and enter them into the database
Counter can be set to test with less no of nodes total nodes are around 2060543 which will take some time. Test for short dataset first    
Next Task :
Update code to store Ways and Relations