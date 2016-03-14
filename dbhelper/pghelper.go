package dbhelper
import(
  _ "os/exec"
  "fmt"
	_ "bufio"
	_ "os"
	"database/sql"
	_ "github.com/lib/pq"
  "github.com/ant0ine/go-json-rest/rest"
  "github.com/qedus/osmpbf"
  "encoding/json"
)
var (
  db *sql.DB
  err error
  DBNAME string
  TABLENAME string
  UNAME string
  PASSWD string
)
var Vehicles = []string{
  "Fiat Fiorino",
  "Smart Roadster",
  "Sam",
  "Citysax",
  "MUTE",
  "Spyder-S",
  "Think",
  "Luis",
  "STROMOS",
  "Karabag Fiat 500E",
  "Lupower Fiat 500E",
}
type NodeData struct {
  Placename string  `db:"placename"` 
  Lon float64 `db:"lon"`
  Lat float64 `db:"lat"`
  Countryname string `db:"countryname"`
}
func SetupPGDatabase() {
  DBNAME = "gntest"   //Your database name here
  UNAME = "ram"     //Username here
  PASSWD = "ram123" //password here
	
	conn := fmt.Sprintf("postgres://%v:%v@localhost:5432/%v?sslmode=disable",UNAME,PASSWD,DBNAME) 
  // db, err := sql.Open("postgres", conn)
  db, err = sql.Open("postgres", conn)
	if err != nil {
		fmt.Println(err)
	} 	
	return
}

func GetAllNodeData(w rest.ResponseWriter, r *rest.Request){
  
  rows, err := db.Query(`SELECT * FROM nodes`)
  if err != nil {
    fmt.Println(err)
  } else {
    single_node:= NodeData{}
    data := []NodeData{}
    for rows.Next() {
      err = rows.Scan(&single_node.Placename,&single_node.Lon,&single_node.Lat)
      if err != nil {
        fmt.Println("Error in query",err)
      } else {
        fmt.Println("Place Name: ",single_node.Placename,"Longitude: ",single_node.Lon,"Latitude: ",single_node.Lat)
        data = append(data,single_node)
      }
    }
    w.WriteJson(&data)
  }
  
}
func GetNodeData(w rest.ResponseWriter, r *rest.Request){
  pname := r.FormValue("place")
  cname := r.FormValue("country")
  query := fmt.Sprintf("select n.placename,n.lon,n.lat,c.countryname from nodes as n,pcountry as c where c.placename='%v' and c.countryname = '%v'",pname,cname)
  rows, err := db.Query(query)
  if err != nil {
    fmt.Println(err)
  } else {
    single_node := NodeData{}
    for rows.Next() {
      err = rows.Scan(&single_node.Placename,&single_node.Lon,&single_node.Lat,&single_node.Countryname)
      if err != nil {
        fmt.Println("Error in query",err)
      } else {
        fmt.Println("Place Name: ",single_node.Placename,"Longitude: ",single_node.Lon,"Latitude: ",single_node.Lat,"Countryname: ",single_node.Countryname)
        w.WriteJson(&single_node)
      }
    }
    
  }
  
}
func GetCountryName(w rest.ResponseWriter, r *rest.Request){
  placename := r.FormValue("place")
  type countryname struct {
    Name string
  }
  country := countryname{}
  countrydata := []countryname{}
  
  query := fmt.Sprintf("SELECT countryname FROM pcountry where placename='%v'",placename)
  fmt.Println(query)
  rows, err := db.Query(query)
  if err != nil {
    fmt.Println(err)
  } else {
    for rows.Next() {
      err = rows.Scan(&country.Name)
      if err != nil {
        fmt.Println("Error in query",err)
      } else {
        fmt.Println("Country Name: ",country.Name)
        countrydata = append(countrydata,country)
      }
    }
  }
  w.WriteJson(&countrydata)
  
}

func InsertNode(node *osmpbf.Node) {
  var lastNodeId int
  var lastInfoId int
  jsonstring,e := json.Marshal(node.Tags)
  // fmt.Println("Json: ",string(jsonstring))
  info := node.Info
  tstamp := fmt.Sprintf("%v",info.Timestamp)
  err = db.QueryRow("INSERT INTO info(version,timestamp,changeset,uid,\"user\",visible) VALUES($1,$2,$3,$4,$5,$6) returning id;", info.Version,tstamp,info.Changeset,info.Uid,info.User,info.Visible).Scan(&lastInfoId)
  if err != nil {
    fmt.Println("Error in insert node-info: ",err)
  }
  if e != nil {
    fmt.Println("Error in marshal tags: ",e)
  }
  err = db.QueryRow("INSERT INTO node(id,lat,lon,tags,infoid) VALUES($1,$2,$3,$4,$5) returning id;", node.ID, node.Lat,node.Lon,jsonstring,lastInfoId).Scan(&lastNodeId)
  if err != nil {
    fmt.Println("Error in insert node-node: ",err)
  }
  
  
}

func GetAllVehicles(w rest.ResponseWriter, r *rest.Request){
  
  vehicleslist := map[string]interface{}{}
  vehicleslist["Vehicles"] = Vehicles
  w.WriteJson(&vehicleslist)
  
}