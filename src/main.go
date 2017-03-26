package main 
import{
	
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
}

const{
	TaskerAuthHeader = "X-Tasker-Authentication"
}

func main(){
	var (
		port          = os.Getenv("PORT")
		mysqlDatabase = os.Getenv("MYSQL_DATABASE")
		mysqlPassword = os.Getenv("MYSQL_PASSWORD")
		mysqlUser     = os.Getenv("MYSQL_USER")
		mysqlHost     = os.Getenv("MYSQL_HOST")
		mysqlPort     = os.Getenv("MYSQL_PORT")
		versionFile   = os.Getenv("VERSION_FILE")
	)
	dbConnStr := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", mysqlUser, mysqlPassword, mysqlHost, mysqlPort, mysqlDatabase)
	db, err := sql.Open("mysql", dbConnStr)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()
	ctx := &tasker.TaskerContext{
		Tasks: &tasker.TaskStore{db},
	}
	// Register help HTTP endpoints with handlers
	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		err = db.Ping()
		if err != nil {
			log.Println(err)
			io.WriteString(w, "error: connection to database\n")
			return
		}

		io.WriteString(w, "pong\n")
	})

	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, versionFile)
	})

	// Register tasker HTTP endpoints with handlers.
	http.HandleFunc("/tasks", tasker.NewHandler(ctx, tasker.TasksHandler))

	log.Printf("listening on port %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}