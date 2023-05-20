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
      _, err = insertUser(user)
      if err != nil {
        fmt.Fprintf(w, "indexHandler: %v", err)
        return
      }
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
      user, err := queryUsersByUsername(r.FormValue("username"))
      if err != nil {
        //TODO: handle by reporting error to user: could not find username
        fmt.Fprintf(w, "indexHandler: %v", err)
        return
      }
      err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(r.FormValue("password")))
      if err != nil {
        //TODO: handle by reporting error to user: incorrect password for user
        fmt.Fprintf(w, "indexHandler: %v", err)
        return
      }
      http.Redirect(w, r, "/login-success/", http.StatusFound)
    default:
      fmt.Fprintf(w, "Only GET and POST methods are supported.")
  }
}

func loginSuccessHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login-success/" {
    http.Redirect(w, r, "/login-success/", http.StatusFound)
		return
	}
  if r.Method == "GET" {
      http.ServeFile(w, r, "./ui/html/login-success.html")
      return
  }
  fmt.Fprintf(w, "Only GET method is supported for handler loginSuccessHandler.")
}

func resetPasswordHandler(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/password-reset/" {
    http.Redirect(w, r, "/password-reset/", http.StatusFound)
		return
  }
  switch r.Method {
    case "GET":
      http.ServeFile(w, r, "./ui/html/password-reset.html")
    case "POST":
      if err := r.ParseForm(); err != nil {
        fmt.Fprintf(w, "resetPasswordHandler: %v", err)
        return
      }
      // TODO: finish writing this method
    default:
      fmt.Fprintf(w, "Only GET and POST methods are supported.")
  }
}

func usernameForgetHandler(w http.ResponseWriter, r *http.Request) {
  if r.URL.Path != "/username-forget/" {
    http.Redirect(w, r, "/username-forget/", http.StatusFound)
		return
  }
  switch r.Method {
    case "GET":
      http.ServeFile(w, r, "./ui/html/username-forget.html")
    case "POST":
      if err := r.ParseForm(); err != nil {
        fmt.Fprintf(w, "usernameForgetHandler: %v", err)
        return
      }
      // TODO: finish writing this method
    default:
      fmt.Fprintf(w, "Only GET and POST methods are supported.")
  }
}

