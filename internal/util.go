package internal

import (
  "log"
  "strings"
  "net/http"
  "html/template"
)

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
