package main

import(
  "log"
  "net/http"
)

func startServer() error {
  mux := http.NewServeMux()
  fileServer := http.FileServer(http.Dir("./ui/static/"))
  mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
  mux.HandleFunc("/login/", loginHandler)
  mux.HandleFunc("/login-success/", loginSuccessHandler)
  mux.HandleFunc("/register/", registerHandler)
  mux.HandleFunc("/password-reset-request/", requestResetPasswordHandler)
  mux.HandleFunc("/reset-password/", resetPasswordHandler)
  mux.HandleFunc("/username-forget/", usernameForgetHandler)
  log.Printf("Listening on port 8080..\n")
  return http.ListenAndServe(":8080", mux)
}
