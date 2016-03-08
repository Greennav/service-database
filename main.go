package main
import(
  "github.com/sitaramshelke/service-database/dbhelper"
	"database/sql"
	_ "github.com/lib/pq"
  "fmt"
)
var (
  DBNAME string
  TABLENAME string
  UNAME string
  PASSWD string
)
func main(){
  DBNAME,TABLENAME,UNAME,PASSWD = dbhelper.SetupPGDatabase()
  conn := fmt.Sprintf("postgres://%s:%s/%s?sslmode=disable",UNAME,PASSWD,DBNAME) 
  // db, err := sql.Open("postgres", conn)
  _, err := sql.Open("postgres", conn)
	if err != nil {
		fmt.Println(err)
	} else {
    fmt.Println("Connected to "+DBNAME)
  }
  fmt.Println("Tablename",TABLENAME)
  
}
