package main

import "github.com/urfave/cli"

var (
    ti = cli.BoolFlag{
        Name:  "ti",
        Usage: "enable tty"}
    d = cli.BoolFlag{
        Name:  "d",
        Usage: "detach container"}
    m = cli.StringFlag{
        Name:  "m",
        Usage: "memory limit"}
    cpushare = cli.StringFlag{
        Name:  "cpushare",
        Usage: "cpushare limit"}
    cpuset = cli.StringFlag{
        Name:  "cpuset",
        Usage: "cpuset limit"}
    v = cli.StringFlag{
        Name:  "v",
        Usage: "volume"}
    name = cli.StringFlag{
        Name:  "name",
        Usage: "container name"}
    e = cli.StringSliceFlag{
        Name:  "e",
        Usage: "set environment"}
    net = cli.StringFlag{
        Name:  "net",
        Usage: "container network"}
    p = cli.StringSliceFlag{
        Name:  "p",
        Usage: "port mapping"}
    driver = cli.StringFlag{
        Name:  "driver",
        Usage: "network driver"}
    subnet = cli.StringFlag{
        Name:  "subnet",
        Usage: "subnet cidr"}
)
