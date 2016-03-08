package dbhelper
import(
  _ "os/exec"
  "fmt"
	"bufio"
	"os"
)
func SetupPGDatabase() (DBNAME, TABLENAME, UNAME, PASSWD string){
  // go get github.com/rnubel/pgmgr
	fmt.Println("Hello")
	//later this can be read from file so no need to enter same data everytime
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Enter the database name: ")
	DBNAME,_ = reader.ReadString('\n')
	fmt.Println("Enter the Table name: ")
	TABLENAME,_ = reader.ReadString('\n')
	fmt.Println("Enter the user name: ")
	UNAME,_ = reader.ReadString('\n')
	fmt.Println("Enter the password: ")
	PASSWD,_ = reader.ReadString('\n')
	return
}