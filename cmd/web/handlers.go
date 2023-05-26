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

func requestResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
    case "GET":
      internal.RenderTemplate(w, r, "request-password-reset.tmpl")
    case "POST":
      passwordResetRequest := struct {
        Email string `json:"email"`
      }{}
      if err := json.NewDecoder(r.Body).Decode(&passwordResetRequest); err != nil {
        panic(err)
      }
      user, err := internal.QueryUsersByEmail(passwordResetRequest.Email)
      if err != nil {
        response := serverResponse {
          Success: false,
          Message: "Internal database error.",
        }
        respond(w, r, &response)
        return
      }
      if err = godotenv.Load(); err != nil {
        response := serverResponse {
          Success: false,
          Message: "Failure loading environment.",
        }
        respond(w, r, &response)
        return
      }
      // generate the token, insert into db, return the hashed token
      // to send as part of URL to user
      token, err := internal.GeneratePasswordResetToken()
      if err != nil {
        response := serverResponse {
          Success: false,
          Message: "Error generating token.",
        }
        respond(w, r, &response)
        return
      }
      if err := internal.InsertTokenIntoTokens(
        token,
        user.ID,
      ); err != nil {
        response := serverResponse {
          Success: false,
          Message: "Error inserting token into db.",
        }
        respond(w, r, &response)
        return
      }
      // TODO: create a route to handle this page
      passwordResetURL := fmt.Sprintf("/reset-password?token=%v", token)

      email := Email {
        From: os.Getenv("EMAIL_FROM"),
        To: []string{user.Email},
        Subject: "Reset Password Request",
        Body: fmt.Sprintf(`<p>Hello Silly Little Human,</p>
                           <p>It would appear as though you have forgotten your password.
                              Allow me to help you, foolish child.</p>
                           <p>Follow this link to reset it: <a href="%s">%s</a></p>`,
                           passwordResetURL,
                           "Reset Password",
        ),
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

func resetPasswordHandler(w http.ResponseWriter, r *http.Request) {
  var response serverResponse
  token := r.URL.Query().Get("token")
  validToken := internal.ValidateToken(token)
  if !validToken {
    response = serverResponse {
      Success: false,
      Message: "invalid token. Maybe it's expired?",
    }
    respond(w, r, &response)
    return
  }
  switch r.Method {
    case "GET":
      internal.RenderTemplate(w, r, "reset-password.tmpl")
      //TODO: down the road do a redirect
    case "POST":
      passwordResetRequest := struct {
        Password string `json:"password"`
      }{}
      if err := json.NewDecoder(r.Body).Decode(&passwordResetRequest); err != nil {
        panic(err)
      }
      userID, err := internal.GetUserIDfromToken(token)
      if err != nil {
        response = serverResponse {
          Success: false,
          Message: err.Error(),
        }
        respond(w, r, &response)
        return
      }
      if err := internal.UpdateUserPassword(userID, passwordResetRequest.Password); err != nil {
        response = serverResponse {
          Success: false,
          Message: err.Error(),
        }
        respond(w, r, &response)
        return
      }
      response = serverResponse {
        Success: true,
        Message: "Successfully updated password",
      }
      respond(w, r, &response)
      return
      //TODO: down the road, re-route to a success page, then redirect to login
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

