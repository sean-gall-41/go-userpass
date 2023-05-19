package main

import(
  "fmt"
  "golang.org/x/crypto/bcrypt"
  "net/http"
)

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/register/" {
    http.Redirect(w, r, "/register/", http.StatusFound)
		return
	}
  switch r.Method {
    case "GET":
      http.ServeFile(w, r, "./ui/html/register.html")
    case "POST":
      if err := r.ParseForm(); err != nil {
        fmt.Fprintf(w, "registerHandler: %v", err)
        return
      }
      hash, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.MinCost)
      if err != nil {
        fmt.Fprintf(w, "registerHandler: %v", err)
        return
      }
      user := User {
        ID: 0,
        Email: r.FormValue("email"),
        Username: r.FormValue("username"),
        PassHash: hash,
      }
      //fmt.Printf("Inserting user into db 'users'...\n")
      //id, err := insertUser(user)
      //if err != nil {
      //  fmt.Fprintf(w, "indexHandler: %v", err)
      //  return
      //}
      //fmt.Printf("Successfully inserted user into db!\n")
      //user.ID = id
      fmt.Fprintf(w, "User: %v\n", user)
    default:
      fmt.Fprintf(w, "Only GET and POST methods are supported.")
  }
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login/" {
    http.Redirect(w, r, "/login/", http.StatusFound)
		return
	}
  switch r.Method {
    case "GET":
      http.ServeFile(w, r, "./ui/html/login.html")
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
      //fmt.Printf("Inserting user into db 'users'...\n")
      //id, err := insertUser(user)
      //if err != nil {
      //  fmt.Fprintf(w, "indexHandler: %v", err)
      //  return
      //}
      //fmt.Printf("Successfully inserted user into db!\n")
      //user.ID = id
      fmt.Fprintf(w, "User: %v\n", user)
    default:
      fmt.Fprintf(w, "Only GET and POST methods are supported.")
  }
}

