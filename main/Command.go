package main

import (
    "errors"
    "fmt"
    "os"

    log "github.com/sirupsen/logrus"
    "github.com/urfave/cli"

    "MyDocker/network"

    "MyDocker/cgroups/subsystems"
    "MyDocker/container"
)

var (
    initCommand = cli.Command{
        Name: "init",
        Usage: `Init container process run user's process in container.
Do not call it outside.`,
        /* 1. 获取传递过来的command参数
           2. 执行容器初始化操作 */
        Action: func(ctx *cli.Context) error {
            log.Info("init come on")
            return container.RunContainerInitProcess()
        }}
    runCommand = cli.Command{
        Name: "run",
        Usage: `Create a container with namespace and cgroups limit.
example: mydocker run -ti [image] [command]`,
        Flags: []cli.Flag{ti, d, m, cpushare, cpuset, v, name, e, net, p},
        /* 这里是run命令执行的真正函数
           1. 判断参数是否包含command
           2. 获取用户指定的command
           3. 调用 Run function 去准备启动容器 */
        Action: func(ctx *cli.Context) error {
            if len(ctx.Args()) == 0 {
                return errors.New("missing container command")
            }
            var (
                tty    = ctx.Bool(ti.Name)
                detach = ctx.Bool(d.Name)
            )
            if tty && detach {
                return errors.New("ti and d arguments can not both provided")
            }
            log.Infof("CreateTty %v", tty)
            Run(tty, ctx.Args().Tail(), ctx.StringSlice(e.Name), ctx.StringSlice(p.Name),
                &subsystems.ResourceConfig{
                    MemoryLimit: ctx.String(m.Name),
                    CpuShare:    ctx.String(cpushare.Name),
                    CpuSet:      ctx.String(cpuset.Name),
                }, ctx.String(v.Name), ctx.String(name.Name), ctx.Args().Get(0),
                ctx.String(net.Name))
            return nil
        }}
    commitCommand = cli.Command{
        Name: "commit",
        Usage: `commit a container into image.
example: mydocker commit [container] [image]`,
        Action: func(ctx *cli.Context) error {
            if len(ctx.Args()) == 0 {
                return errors.New("missing container name")
            }
            commitContainer(ctx.Args().Get(0), ctx.Args().Get(1))
            return nil
        }}
    listCommand = cli.Command{
        Name:  "ps",
        Usage: "list all the containers",
        Action: func(ctx *cli.Context) error {
            ListContainers()
            return nil
        }}
    logCommand = cli.Command{
        Name:  "logs",
        Usage: "print logs of a container",
        Action: func(ctx *cli.Context) error {
            if len(ctx.Args()) == 0 {
                return errors.New("please input the container name")
            }
            logOfContainer(ctx.Args().Get(0))
            return nil
        }}
    execCommand = cli.Command{
        Name:  "exec",
        Usage: "exec a command into container",
        Action: func(ctx *cli.Context) error {
            if os.Getenv(EnvExecPid) != "" {
                log.Infof("Pid %s callback.", os.Getgid())
                return nil
            }
            if len(ctx.Args()) < 2 {
                return errors.New("missing container name or command")
            }
            ExecContainer(ctx.Args().Get(0), ctx.Args().Tail())
            return nil
        }}
    stopCommand = cli.Command{
        Name:  "stop",
        Usage: "stop a container",
        Action: func(ctx *cli.Context) error {
            if len(ctx.Args()) == 0 {
                return errors.New("missing container name")
            }
            stopContainer(ctx.Args().Get(0))
            return nil
        }}
    removeCommand = cli.Command{
        Name:  "rm",
        Usage: "remove unused container",
        Action: func(ctx *cli.Context) error {
            if len(ctx.Args()) == 0 {
                return errors.New("missing container name")
            }
            removeContainer(ctx.Args().Get(0))
            return nil
        }}
    networkCommand = cli.Command{
        Name:  "network",
        Usage: "container network commands",
        Subcommands: []cli.Command{
            createCommand,
            listCommandOfNetwork,
            removeCommandOfNetwork}}
)

// subCommands of network
var (
    createCommand = cli.Command{
        Name:  "create",
        Usage: "create a container network",
        Flags: []cli.Flag{driver, subnet},
        Action: func(ctx *cli.Context) error {
            if len(ctx.Args()) == 0 {
                return errors.New("missing network name")
            }
            network.Init()
            if err := network.Cerate(ctx.Args().Get(0), ctx.String(driver.Name), ctx.String(subnet.Name));
              err != nil {
                return fmt.Errorf("create network error: %v", err)
            }
            return nil
        }}
    listCommandOfNetwork = cli.Command{
        Name:  "list",
        Usage: "list container networks",
        Action: func(ctx *cli.Context) error {
            network.Init()
            network.List()
            return nil
        }}
    removeCommandOfNetwork = cli.Command{
        Name:  "remove",
        Usage: "remove container network",
        Action: func(ctx *cli.Context) error {
            if len(ctx.Args()) == 0 {
                return errors.New("missing network name")
            }
            network.Init()
            if err := network.Delete(ctx.Args().Get(0)); err != nil {
                return fmt.Errorf("remove network error: %v", err)
            }
            return nil
        }}
)
