# service-database
Database service in Go for the GreenNav project    
This repository contains code which follows [Design plan for database service](https://github.com/Greennav/greennav.github.io/blob/master/wiki/Roadmap.md#design-plan)

Database setup:     
```
create table nodes(placename text unique,lon real, lat real);    
create table pcountry(placename text, countryname text, foreign key(placename) references nodes(placename));
```

APIs:    
1. Get all nodes data:     
        curl -X GET http://localhost:8888/node/ 
2. Get country name of the place     
        curl -X GET http://localhost:8888/country?place=Pune