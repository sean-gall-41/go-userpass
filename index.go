package main

import(
  "net/http"
  "database/sql"
  "os"
  "log"
  "fmt"
  "golang.org/x/crypto/bcrypt"
  "github.com/go-sql-driver/mysql"
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
    // Return the new album's ID.
    return id, nil
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/index/" {
    http.Redirect(w, r, "/index/", http.StatusFound)
		return
	}
  switch r.Method {
    case "GET":
      http.ServeFile(w, r, "static/index.html")
    case "POST":
      if err := r.ParseForm(); err != nil {
        fmt.Fprintf(w, "indexHandler: %v", err)
        return
      }
      hash, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.MinCost)
      if err != nil {
        fmt.Fprintf(w, "indexHandler: %v", err)
        return
      }
      user := User {
        ID: 0,
        Email: r.FormValue("email"),
        Username: r.FormValue("username"),
        PassHash: hash,
      }
      fmt.Printf("Inserting user into db 'users'...\n")
      id, err := insertUser(user)
      if err != nil {
        fmt.Fprintf(w, "indexHandler: %v", err)
        return
      }
      fmt.Printf("Successfully inserted user into db!\n")
      user.ID = id
      fmt.Fprintf(w, "User: %v\n", user)
    default:
      fmt.Fprintf(w, "Only GET and POST methods are supported.")
  }
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

  //users, err := queryUsersByEmail("myemail@hotmail.com")
  //if err != nil {
  //  log.Fatal(err)
  //}
  //if len(users) > 1 {
  //  fmt.Errorf("Bruh you messed up the database: two users with same email? NAH lil bro")
  //  os.Exit(1)
  //}
  //if users == nil {
  //  fmt.Printf("no users found :(\n")
  //} else {
  //  fmt.Printf("user found: %v\n", users[0])
  //}

  http.HandleFunc("/index/", indexHandler)
  fmt.Printf("Listening on port 8080..\n")
  log.Fatal(http.ListenAndServe(":8080", nil))
}
