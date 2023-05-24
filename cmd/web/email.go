package main

import (
  "os"
  "fmt"
  "strconv"
  "crypto/tls"
  "gopkg.in/gomail.v2"
)

type Email struct {
  From string
  To []string
  Cc []string
  Subject string
  Body string
  Attachments []string
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

