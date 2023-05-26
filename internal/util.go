package internal

import (
  "fmt"
  "log"
  "strings"
  "net/http"
  "html/template"
  "crypto/rand"
  "math/big"
)

const tokenCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type TemplateData struct {
  Title string
}
// /lost-password/ -> lost password
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmplName string) {
  title := strings.Replace(tmplName, ".tmpl", "", -1)
  title = strings.Replace(title, "-", " ", -1)
  tmplData := TemplateData { Title: title }
  files := [] string {
    "./ui/templates/html-base.tmpl",
    "./ui/templates/styles-base.tmpl",
    "./ui/templates/" + tmplName,
  }
  t, err := template.ParseFiles(files...)
  if err != nil {
    log.Print(err.Error())
    http.Error(w, err.Error(), http.StatusInternalServerError)
    return
  }
  if err := t.ExecuteTemplate(w, "html-base", &tmplData); err != nil {
    http.Error(w, err.Error(), http.StatusInternalServerError)
  }
}

func GeneratePasswordResetToken() (string, error) {
  token := make([]byte, 64)
  charSetLen := big.NewInt(int64(len(tokenCharset)))
  for i := 0; i < 64; i++ {
    randID, err := rand.Int(rand.Reader, charSetLen)
    if err != nil {
      return "", err
    }
    token[i] = tokenCharset[randID.Int64()]
  }
  return string(token), nil
}

func FormatTimeZone(offset int) string {
  hours := offset / 3600
  minutes := (offset % 3600) / 60

  timeZone := fmt.Sprintf("%+03d:%02d", hours, minutes)
  return timeZone
}
