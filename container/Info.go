package container

import (
    htmlTmp "html/template"
    "path/filepath"
    "text/template"

    "MyDocker/templates"
)

func init() {
    templates.RegisterTextTmp(InfosID, template.Must(template.ParseFiles(
        filepath.Join(TemplateDir, InfosName))))
    templates.RegisterHtmlTmp(InfosID, htmlTmp.Must(htmlTmp.ParseFiles(
        filepath.Join(TemplateDir, InfosHTMLName))))
}

const (
    DefaultInfoDir      = "/var/run/mydocker/"
    DefaultInfoLocation = DefaultInfoDir + "%s/"
    ConfigName          = "config.json"
    LogFile             = "log/container.log"

    TemplateDir   = templates.BaseURL + "container/"
    InfosID       = "container.Infos"
    InfosName     = "Infos"
    InfosHTMLName = InfosName + templates.HTMLName
)

type Info struct {
    Pid          string   `json:"pid"`           // 容器的init进程在宿主机上的pid
    ID           string   `json:"id"`            // 容器ID
    Name         string   `json:"name"`          // 容器名
    Command      string   `json:"command"`       // 容器内init进程的运行命令
    CreatedTime  string   `json:"created_time"`  // 创建时间
    Status       string   `json:"status"`        // 容器的状态
    Volume       string   `json:"volume"`        // 容器的数据卷
    PortMappings []string `json:"port_mappings"` // host和容器的端口映射组合集
}
