package templates

import (
    "io"
    "text/template"

    log "github.com/sirupsen/logrus"
)

var textTmpWithId = make(map[string]*template.Template)

func RegisterTextTmp(id string, tmp *template.Template) {
    if textTmpWithId[id] != nil {
        log.Fatalf("Text tmp %s is exist.", id)
    }
    textTmpWithId[id] = tmp
}

func FPrintText(writer io.Writer, id, name string, data interface{}) {
    if err := textTmpWithId[id].ExecuteTemplate(writer, name, data); err != nil {
        log.Errorf("ExecuteTmp %s.%s failed. %v", id, name, err)
        return
    }
}
