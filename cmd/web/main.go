package main

import(
  "net/http"
  "database/sql"
  "os"
  "log"
  "fmt"
  "github.com/go-sql-driver/mysql"
)

func main() {
  dbName := "go-userpass-users"
  fmt.Printf("Attempting to connect to database %s\n", dbName)
  cfg := mysql.Config {
    User: os.Getenv("DBUSER"),
    Passwd: os.Getenv("DBPASS"),
    Net: "unix",
    Addr: "/var/run/mysqld/mysqld.sock",
    DBName: dbName,
    AllowNativePasswords: true,
  }
  var err error
  db, err = sql.Open("mysql", cfg.FormatDSN())
  if err != nil {
    log.Fatal(err)
  }
  pingErr := db.Ping()
  if pingErr != nil {
    log.Fatal(pingErr)
  }
  fmt.Printf("Connection Succesful.\n")

  mux := http.NewServeMux()
  fileServer := http.FileServer(http.Dir("./ui/static/"))
  mux.Handle("/static/", http.StripPrefix("/static", fileServer))
  mux.HandleFunc("/login/", loginHandler)
  mux.HandleFunc("/login-success/", loginSuccessHandler)
  mux.HandleFunc("/register/", registerHandler)
  mux.HandleFunc("/password-reset/", resetPasswordHandler)
  mux.HandleFunc("/username-forget/", usernameForgetHandler)
  log.Printf("Listening on port 8080..\n")
  log.Fatal(http.ListenAndServe(":8080", mux))
}

