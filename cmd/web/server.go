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
  mux.HandleFunc("/register-success/", registerSuccessHandler)
  mux.HandleFunc("/password-reset-request/", requestResetPasswordHandler)
  mux.HandleFunc("/password-reset-request-success/", requestResetPasswordSuccessHandler)
  mux.HandleFunc("/reset-password/", resetPasswordHandler)
  mux.HandleFunc("/reset-password-success/", resetPasswordSuccessHandler)
  mux.HandleFunc("/username-forget/", usernameForgetHandler)
  mux.HandleFunc("/username-forget-success/", usernameForgetSuccessHandler)
  log.Printf("Listening on port 8080..\n")
  return http.ListenAndServe(":8080", mux)
}
