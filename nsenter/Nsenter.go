package nsenter

/*
#define _GNU_SOURCE
#include <unistd.h>
#include <errno.h>
#include <sched.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <fcntl.h>

// 这里的__attribute__((constructor))指的是，一旦这个包被引用，那么这个函数就会被自动执行
// 类似于构造函数，会在程序一启动的时候运行
__attribute__((constructor)) void enter_namespace(void) {
    char *mydocker_pid;
    mydocker_pid = getenv("mydocker_pid");
    if (!mydocker_pid) {
        return;
    }
    char *mydocker_cmd;
    mydocker_cmd = getenv("mydocker_cmd");
    if (!mydocker_cmd) {
        return;
    }
    char nspath[1024];
    char *namespaces[]={"ipc","uts","net","pid","mnt"};
    int i;
    for (i=0; i<5; i++) {
        sprintf(nspath,"/proc/%s/ns/%s",mydocker_pid,namespaces[i]);
        int fd = open(nspath,O_RDONLY);
        setns(fd,0);
        close(fd);
    }
    int res = system(mydocker_cmd);
    exit(0);
    return;
}
*/
import "C"
