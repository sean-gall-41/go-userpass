package main

import(
  "fmt"
  "database/sql"
)

var db *sql.DB

type User struct {
  ID int64
  Email string
  Username string
  PassHash []byte
}

func queryUsersByUsername(username string) (*User, error) {
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

