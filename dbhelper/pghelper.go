package dbhelper
import(
  _ "os/exec"
  "fmt"
	_ "bufio"
	_ "os"
	"database/sql"
	_ "github.com/lib/pq"
  "github.com/ant0ine/go-json-rest/rest"
)
var (
  DBNAME string
  TABLENAME string
  UNAME string
  PASSWD string
)
type NodeData struct {
  Placename string  `db:"placename"` 
  Lon float64 `db:"lon"`
  Lat float64 `db:"lat"`
}
func SetupPGDatabase() (db *sql.DB){
  DBNAME = "gntest"
  UNAME = "ram"
  PASSWD = "ram123"
	
	conn := fmt.Sprintf("postgres://%v:%v@localhost:5432/%v?sslmode=disable",UNAME,PASSWD,DBNAME) 
  // db, err := sql.Open("postgres", conn)
  db, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println(err)
	} 	
	return
}

func GetAllNodeData(w rest.ResponseWriter, r *rest.Request){
  DB := SetupPGDatabase()
  rows, err := DB.Query(`SELECT * FROM nodes`)
  if err != nil {
    fmt.Println(err)
  } else {
    d := NodeData{}
    data := []NodeData{}
    for rows.Next() {
      err = rows.Scan(&d.Placename,&d.Lon,&d.Lat)
      if err != nil {
        fmt.Println("Error in query",err)
      } else {
        fmt.Println("Place Name: ",d.Placename,"Longitude: ",d.Lon,"Latitude: ",d.Lat)
        data = append(data,d)
      }
    }
    w.WriteJson(&data)
  }
  
}
func GetCountryName(w rest.ResponseWriter, r *rest.Request){
  DB := SetupPGDatabase()
  pname := r.FormValue("place")
  type data struct {
    Name string
  }
  countrydata := data{}
  query := fmt.Sprintf("SELECT countryname FROM pcountry where placename='%v'",pname)
  fmt.Println(query)
  rows, err := DB.Query(query)
  if err != nil {
    fmt.Println(err)
  } else {
    for rows.Next() {
      err = rows.Scan(&countrydata.Name)
      if err != nil {
        fmt.Println("Error in query",err)
      } else {
        fmt.Println("Country Name: ",countrydata.Name)
        w.WriteJson(&countrydata)
      }
    }
  }
  
  
}