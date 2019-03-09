package templates

import (
    "html/template"
    "io"

    log "github.com/sirupsen/logrus"
)

var htmlTmpWithID = make(map[string]*template.Template)

func RegisterHtmlTmp(id string, tmp *template.Template) {
    if htmlTmpWithID[id] != nil {
        log.Fatalf("Html tmp %s is exist.", id)
    }
    htmlTmpWithID[id] = tmp
}

func FPrintHtml(writer io.Writer, id, name string, data interface{}) {
    if err := textTmpWithId[id].ExecuteTemplate(writer, name, data); err != nil {
        log.Errorf("ExecuteTmp %s.%s failed. %v", id, name, err)
        return
    }
}
