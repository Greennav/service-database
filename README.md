# service-database
Database service in Go for the GreenNav project    
This repository contains code which follows [Design plan for database service](https://github.com/Greennav/greennav.github.io/blob/master/wiki/Roadmap.md#design-plan)

Database setup:     
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