package db
import "testing"
import "database/sql"
import "os"
import _ "github.com/lib/pq"
import "fmt"
const (
	host     = "localhost"
	port     = 5432
	user     = "root"
	password = "secret"
	dbname   = "simple_bank"
  )
var postgresqlDbInfo = fmt.Sprintf("host=%s port=%d user=%s "+
  "password=%s dbname=%s sslmode=disable",
  host, port, user, password, dbname)
var testQueries *Queries
func TestMain(m *testing.M){
	db ,_ := sql.Open("postgres",postgresqlDbInfo)
	testQueries = New(db)
	os.Exit(m.Run())
}