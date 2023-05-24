package main

import(
  "os"
  "log"
  "net/http"
  "github.com/sean-gall-41/go-userpass/internal"
)

func main() {
  if _, err := internal.StartMySQL(); err != nil {
    log.Fatal(err)
    os.Exit(-1)
  }

  mux := http.NewServeMux()
  fileServer := http.FileServer(http.Dir("./ui/static/"))
  mux.Handle("/static/", http.StripPrefix("/static", fileServer))
  mux.HandleFunc("/login/", loginHandler)
  mux.HandleFunc("/login-success/", loginSuccessHandler)
  mux.HandleFunc("/register/", registerHandler)
  mux.HandleFunc("/password-reset/", resetPasswordHandler)
  mux.HandleFunc("/username-forget/", usernameForgetHandler)
  log.Printf("Listening on port 8080..\n")
  log.Fatal(http.ListenAndServe(":8080", mux))
}

