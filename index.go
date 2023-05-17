package main

import(
  "database/sql"
  "os"
  "log"
  "fmt"
  "github.com/go-sql-driver/mysql"
)

var db *sql.DB

type User struct {
  ID int64
  Email string
  Username string
  PassHash string
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

  users, err := queryUsersByEmail("myemail@hotmail.com")
  if err != nil {
    log.Fatal(err)
  }
  if len(users) > 1 {
    fmt.Errorf("Bruh you messed up the database: two users with same email? NAH lil bro")
    os.Exit(1)
  }
  if users == nil {
    fmt.Printf("no users found :(\n")
  } else {
    fmt.Printf("user found: %v\n", users[0])
  }
}
