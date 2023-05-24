package main

import(
  "os"
  "fmt"
  "strconv"
  "golang.org/x/crypto/bcrypt"
  "crypto/tls"
  "net/http"
  "gopkg.in/gomail.v2"
  "github.com/joho/godotenv"
  "github.com/sean-gall-41/go-userpass/internal"
)

func registerHandler(w http.ResponseWriter, r *http.Request) {
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
      user := internal.User {
        ID: 0,
        Email: r.FormValue("email"),
        Username: r.FormValue("username"),
        PassHash: hash,
      }
      _, err = internal.InsertUser(user)
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
  switch r.Method {
    case "GET":
      http.ServeFile(w, r, "./ui/html/login.html")
    case "POST":
      if err := r.ParseForm(); err != nil {
        fmt.Fprintf(w, "indexHandler: %v", err)
        return
      }
      user, err := internal.QueryUsersByUsername(r.FormValue("username"))
      if err != nil {
        fmt.Fprintf(w, "indexHandler: Could not find username!")
        return
      }
      err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(r.FormValue("password")))
      if err != nil {
        //TODO: handle by reporting error to user: incorrect password for user
        fmt.Fprintf(w, "indexHandler: Given password for user '%s' is incorrect!", user.Username)
        return
      }
      http.Redirect(w, r, "/login-success/", http.StatusFound)
    default:
      fmt.Fprintf(w, "Only GET and POST methods are supported.")
  }
}

func loginSuccessHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method == "GET" {
      http.ServeFile(w, r, "./ui/html/login-success.html")
      return
  }
  fmt.Fprintf(w, "Only GET method is supported for handler loginSuccessHandler.")
}

func resetPasswordHandler(w http.ResponseWriter, r *http.Request) {
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

type Email struct {
  From string
  To []string
  Cc []string
  Subject string
  Body string
  Attachments []string
}

func usernameForgetHandler(w http.ResponseWriter, r *http.Request) {
  switch r.Method {
    case "GET":
      http.ServeFile(w, r, "./ui/html/username-forget.html")
    case "POST":
      if err := r.ParseForm(); err != nil {
        fmt.Fprintf(w, "usernameForgetHandler: %v", err)
        return
      }
      user, err := internal.QueryUsersByEmail(r.FormValue("email"))
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
      if err := sendEmail(&email); err != nil {
        fmt.Fprintf(w, "usernameForgetHandler: %v", err)
        return
      }
    default:
      fmt.Fprintf(w, "Only GET and POST methods are supported.")
  }
}

func sendEmail(email *Email) error {
  m := gomail.NewMessage()
  m.SetHeader("From", email.From)
  m.SetHeader("To", email.To...)
  if len(email.Cc) != 0 {
    m.SetHeader("Cc", email.Cc...)
  }
  m.SetHeader("Subject", email.Subject)
  m.SetBody("text/html", email.Body)
  if len(email.Attachments) != 0 {
    for _, file := range email.Attachments {
      m.Attach(file) // what if this fails?
    }
  }
  port, err := strconv.Atoi(os.Getenv("SMTP_PORT"));
  if err != nil {
    return fmt.Errorf("sendEmail: %v", err)
  }
  d := gomail.NewDialer(
    os.Getenv("SMTP_HOST"),
    port,
    os.Getenv("SMTP_USER"),
    os.Getenv("SMTP_PASS"),
  )
  d.TLSConfig = &tls.Config{InsecureSkipVerify: true}
  if err := d.DialAndSend(m); err != nil {
    return fmt.Errorf("sendEmail: %v", err)
  }
  return nil
}

