package main

import(
  "os"
  "fmt"
  "encoding/json"
  "golang.org/x/crypto/bcrypt"
  "net/http"
  "github.com/joho/godotenv"
  "github.com/sean-gall-41/go-userpass/internal"
)

type serverResponse struct {
  Success bool `json:"success"`
  Message string `json:"message"`
}

func respond(w http.ResponseWriter, r *http.Request, response *serverResponse) {
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(http.StatusCreated)
  json.NewEncoder(w).Encode(response)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
    case "GET":
      internal.RenderTemplate(w, r, "register.tmpl")
    case "POST":
      registerRequest := struct {
        Email string `json:"email"`
        Username string `json:"username"`
        Password string `json:"password"`
      }{}
      err := json.NewDecoder(r.Body).Decode(&registerRequest)
      if err != nil {
        panic(err)
      }
      hash, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.MinCost)
      if err != nil {
        fmt.Fprintf(w, "registerHandler: %v", err)
        return
      }
      user := internal.User {
        ID: 0,
        Email: registerRequest.Email,
        Username: registerRequest.Username,
        PassHash: hash,
      }
      _, err = internal.InsertUser(user)
      var response serverResponse
      if err != nil {
        response = serverResponse {
          Success: false,
          Message: err.Error(),
        }
        respond(w, r, &response)
        return
      }
      response = serverResponse {
        Success: true,
        Message: "",
      }
      respond(w, r, &response)
    default:
      fmt.Fprintf(w, "Only GET and POST methods are supported.")
  }
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
    case "GET":
      internal.RenderTemplate(w, r, "login.tmpl")
    case "POST":
      loginRequest := struct {
        Username string `json:"username"`
        Password string `json:"password"`
      }{}
      err := json.NewDecoder(r.Body).Decode(&loginRequest)
      if err != nil {
        panic(err)
      }
      user, err := internal.QueryUsersByUsername(loginRequest.Username)
      var response serverResponse
      if err != nil {
        response = serverResponse {
          Success: false,
          Message: "Invalid Username",
        }
        respond(w, r, &response)
        return
      }
      err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(loginRequest.Password))
      if err != nil {
        response = serverResponse {
          Success: false,
          Message: "Incorrect password",
        }
        respond(w, r, &response)
        return
      }
      response = serverResponse {
        Success: true,
        Message: "",
      }
      respond(w, r, &response)
    default:
      fmt.Fprintf(w, "Only GET and POST methods are supported.")
  }
}

func loginSuccessHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method == "GET" {
      internal.RenderTemplate(w, r, "login-success.tmpl")
      return
  }
  fmt.Fprintf(w, "Only GET method is supported for handler loginSuccessHandler.")
}

func resetPasswordHandler(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
    case "GET":
      internal.RenderTemplate(w, r, "password-reset.tmpl")
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
  switch r.Method {
    case "GET":
      internal.RenderTemplate(w, r, "username-forget.tmpl")
    case "POST":
      usernameRequest := struct {
        Email string `json:"email"`
      }{}
      if err := json.NewDecoder(r.Body).Decode(&usernameRequest); err != nil {
        panic(err)
      }
      user, err := internal.QueryUsersByEmail(usernameRequest.Email)
      if err != nil {
        fmt.Fprintf(w, "usernameForgetHandler: %v", err)
        return
      }
      if err = godotenv.Load(); err != nil {
        fmt.Fprintf(w, "usernameForgetHandler: %v", err)
        return
      }
      email := Email {
        From: os.Getenv("EMAIL_FROM"),
        To: []string{user.Email},
        Subject: "Forgot username",
        Body: fmt.Sprintf(`<p>Hello Silly Little Human,</p>
                           <p>It would appear as though you have forgotten your username.
                              Allow me to help you, foolish child.</p>
                           <p>Your username is: <b>%s</b></p>`, user.Username),
      }
      var response serverResponse
      if err := sendEmail(&email); err != nil {
        response = serverResponse {
          Success: false,
          Message: "Something went wrong when sending an email",
        }
        respond(w, r, &response)
        return
      }
      response = serverResponse {
        Success: true,
        Message: "",
      }
      respond(w, r, &response)
    default:
      fmt.Fprintf(w, "Only GET and POST methods are supported.")
  }
}

