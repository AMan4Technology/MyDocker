package templates_test

import (
    htmlTmp "html/template"
    "os"
    "testing"
    "text/tabwriter"
    "text/template"

    log "github.com/sirupsen/logrus"

    "MyDocker/templates"
)

func TestFPrintText(t *testing.T) {
    templates.RegisterTextTmp(InfosID, template.Must(template.ParseFiles(CurrentDir+InfosName)))
    templates.RegisterHtmlTmp(InfosID, htmlTmp.Must(htmlTmp.ParseFiles(CurrentDir+InfosHtmlName)))
    infos := []Info{
        {Name: "1", ID: "1"},
        {Name: "20000", ID: "2"},
        {Name: "3", ID: "30000"},
        {Name: "400000", ID: "400000"},
    }
    w := tabwriter.NewWriter(os.Stdout, 4, 1, 3, ' ', 0)
    templates.FPrintText(w, InfosID, InfosName, infos)
    if err := w.Flush(); err != nil {
        log.Errorf("Flush failed. %v", err)
    }
}

const (
    CurrentDir    = "./container/"
    InfosID       = "container.infos"
    InfosName     = "Infos"
    InfosHtmlName = InfosName + templates.HtmlName
)

type Info struct {
    Pid         string `json:"pid"`          // 容器的init进程在宿主机上的pid
    ID          string `json:"id"`           // 容器ID
    Name        string `json:"name"`         // 容器名
    Command     string `json:"command"`      // 容器内init进程的运行命令
    CreatedTime string `json:"created_time"` // 创建时间
    Status      string `json:"status"`       // 容器的状态
}
