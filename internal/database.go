package internal

import(
  "os"
  "log"
  "fmt"
  "time"
  "database/sql"
  "golang.org/x/crypto/bcrypt"
  "github.com/joho/godotenv"
  "github.com/go-sql-driver/mysql"
  timefmt "github.com/itchyny/timefmt-go"
)

var db *sql.DB

type User struct {
  ID int64
  Email string
  Username string
  PassHash []byte
}

type Token struct {
  ID int64
  TokenHash []byte
  TimeCreated string
}

func StartMySQL() error {
  if err := godotenv.Load(); err != nil {
    return err
  }
  log.Print("Connecting to database...\n")
  cfg := mysql.Config {
    User: os.Getenv("DB_USER"),
    Passwd: os.Getenv("DB_PASS"),
    Net: os.Getenv("DB_NET"),
    Addr: os.Getenv("DB_ADDR"),
    DBName: os.Getenv("DB_NAME"),
    AllowNativePasswords: true,
  }
  var err error
  db, err = sql.Open("mysql", cfg.FormatDSN())
  if err != nil {
    return err
  }
  if err := db.Ping(); err != nil {
    return err
  }
  log.Print("Connection Succesful.\n")
  return nil
}

func QueryUsersByUsername(username string) (*User, error) {
  var user User
  row := db.QueryRow(
    "SELECT * FROM users WHERE username = ? LIMIT 1",
    username,
  )
  if err := row.Scan(&user.ID, &user.Email, &user.Username, &user.PassHash); err != nil {
    return nil, fmt.Errorf("queryUsersByUserAndPassHash: %v", err)
  }
  return &user, nil
}

func QueryUsersByEmail(email string) (*User, error) {
  var usr User
  row := db.QueryRow(
    "SELECT * FROM users WHERE email = ? LIMIT 1",
    email,
  )
  if err := row.Scan(&usr.ID, &usr.Email, &usr.Username, &usr.PassHash); err != nil {
    return nil, fmt.Errorf("queryUserByEmail %q: %v", email, err)
  }
  return &usr, nil
}

func UserExists(usr User) bool {
  row := db.QueryRow(
    "SELECT id FROM users WHERE email = ? AND username = ? LIMIT 1",
    usr.Email,
    usr.Username,
  )
  var id int64
  if err := row.Scan(&id); err != nil {
    return false
  }
  return true
}

func InsertUser(usr User) (int64, error) {
    if UserExists(usr) {
      return 0, fmt.Errorf("User already exists!")
    }
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

func userTokenExists(userID int64) bool {
  row := db.QueryRow("SELECT id FROM tokens WHERE id = ? LIMIT 1", userID)
  var id int64
  if err := row.Scan(&id); err != nil {
    return false
  }
  if id == userID {
    return true
  }
  return false
}

func deleteUserToken(userID int64) error {
  if _, err := db.Exec("DELETE FROM tokens WHERE id = ?", userID); err != nil {
    return err
  }
  return nil
}

func ValidateToken(token string) bool {
  rows, err := db.Query("SELECT * FROM tokens")
  if err != nil {
    return false
  }
  defer rows.Close()
  var dbToken Token
  for rows.Next() {
    if err := rows.Scan(&dbToken.ID, &dbToken.TokenHash, &dbToken.TimeCreated); err != nil {
      return false
    }
    if err := bcrypt.CompareHashAndPassword(dbToken.TokenHash, []byte(token)); err == nil {
      parsedCreateDate, err := timefmt.Parse(dbToken.TimeCreated, "%Y-%m-%d %H:%M:%S")
      if err != nil {
        return false
      }
      // lil hacky way of getting the current time zone when LoadLocation("") fails
      _, timeZoneOffset := time.Now().Zone()
      location := time.FixedZone("current-time-zone", timeZoneOffset)
      parsedCreateDate = parsedCreateDate.In(location)
      expirationTime := parsedCreateDate.Add(time.Hour)
      if time.Now().After(expirationTime) {
        return false
      }
      return true
    }
  }
  return false
}

func InsertTokenIntoTokens(token string, userID int64) error {
  if userTokenExists(userID) {
    deleteUserToken(userID)
  }
  hashedToken, err := bcrypt.GenerateFromPassword([]byte(token), bcrypt.MinCost)
  if err != nil {
    return fmt.Errorf("Could not hash token")
  }
  if _, err := db.Exec("INSERT INTO tokens (id, token_hash, time_created) VALUES (?, ?, UTC_TIMESTAMP())",
    userID,
    hashedToken,
  ); err != nil {
    return fmt.Errorf("Could not insert token into tokens")
  }
  return nil
}

