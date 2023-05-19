package main

import(
  "net/http"
  "database/sql"
  "log"
  "fmt"
  //"github.com/go-sql-driver/mysql"
)

var db *sql.DB

type User struct {
  ID int64
  Email string
  Username string
  PassHash []byte
}

func queryUsersByEmail(email string) ([]User, error) {
  var users []User
  rows, err := db.Query(
    "SELECT * FROM users WHERE email = ?",
    email,
  )
  if err != nil {
    return nil, fmt.Errorf("queryUserByEmail %q: %v", email, err)
  }
  defer rows.Close()
  for rows.Next() {
    var usr User
    if err := rows.Scan(&usr.ID, &usr.Email, &usr.Username, &usr.PassHash); err != nil {
      return nil, fmt.Errorf("queryUserByEmail %q: %v", email, err)
    }
    users = append(users, usr)
  }
  if err := rows.Err(); err != nil {
    return nil, fmt.Errorf("queryUserByEmail %q: %v", email, err)
  }
  return users, nil
}

func insertUser(usr User) (int64, error) {
    result, err := db.Exec("INSERT INTO users (email, username, password_hash) VALUES (?, ?, ?)",
      usr.Email,
      usr.Username,
      usr.PassHash,
    )
    if err != nil {
        return 0, fmt.Errorf("insertUser: %v", err)
    }
    // Get the new user's generated ID for the client.
    id, err := result.LastInsertId()
    if err != nil {
        return 0, fmt.Errorf("insertUser: %v", err)
    }
    // Return the new users's ID.
    return id, nil
}

func main() {
  //dbName := "go-userpass-users"
  //fmt.Printf("Attempting to connect to database %s\n", dbName)
  //cfg := mysql.Config {
  //  User: os.Getenv("DBUSER"),
  //  Passwd: os.Getenv("DBPASS"),
  //  Net: "unix",
  //  Addr: "/var/run/mysqld/mysqld.sock",
  //  DBName: dbName,
  //  AllowNativePasswords: true,
  //}
  //var err error
  //db, err = sql.Open("mysql", cfg.FormatDSN())
  //if err != nil {
  //  log.Fatal(err)
  //}
  //pingErr := db.Ping()
  //if pingErr != nil {
  //  log.Fatal(pingErr)
  //}
  //fmt.Printf("Connection Succesful.\n")
  mux := http.NewServeMux()
  fileServer := http.FileServer(http.Dir("./ui/static/"))
  mux.Handle("/static/", http.StripPrefix("/static", fileServer))
  mux.HandleFunc("/index/", indexHandler)
  log.Printf("Listening on port 8080..\n")
  log.Fatal(http.ListenAndServe(":8080", mux))
}

