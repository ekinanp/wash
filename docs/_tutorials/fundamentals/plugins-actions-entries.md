---
title: Understanding plugins, actions, and entries
---
This tutorial introduces you to Wash, which is a shell environment layered on top of the system shell. You’ll learn about plugins, actions and entries while navigating through the Docker plugin.

First, let’s start up the Wash shell.

```
bash-5.0$ ./wash
Welcome to Wash!
  Wash includes several built-in commands: wexec, find, list, meta, tail.
  See commands run with wash via 'whistory', and logs with 'whistory <id>'.
Try 'help'
wash . ❯ ls
aws        docker     gcp        kubernetes
```

Each thing listed here is a ‘directory’. Each 'directory' is a Wash plugin, and each Wash plugin serves as an adapter between a given vendor and Wash. Plugins are the sole mode of interaction between a Wash user and a given vendor's API. For example, the `docker` directory lets you interact with Docker resources like containers and volumes. Similarly, the `aws` directory lets you interact with AWS resources like EC2 instances and S3 objects.

Everything in Wash is an entry, including resources. Each plugin provides a hierarchical view of all its entries. With Wash, you navigate through a plugin's API in the same way that you would navigate the Linux command line -- by changing directories, listing the contents of a directory, and performing actions. 

Each entry supports a specific set of Wash actions.  For example, you could use `list` to list  an entry's children, `exec` to execute a command on the entry, or `read` to read an entry's contents. Note that entries which support `list` are represented as directories in the shell. All other entries are represented as files.

Let's walk through an example to see what all of this looks like in practice. We'll navigate through the Docker plugin and interact with some of its entries.

```
wash . ❯ cd docker
wash docker ❯
```

Notice how the prompt changed from `wash . >` to `wash docker >`. This is a useful way to track your current location when navigating through a given plugin.

```
wash docker ❯ ls
containers volumes
```

From the output, we see that the Docker plugin lets us interact with containers and volumes via the `containers` and `volumes` entries. These entries support the `list` action, so they're represented as directories. Let’s try interacting with some containers.

```
wash docker ❯ cd containers
wash docker/containers ❯ ls
sleepy_heisenberg
k8s_POD_compose-6c67d745f6-q54n8_docker_57a0f7e9-c41c-11e9-9d31-025000000001_1
k8s_POD_compose-api-57ff65b8c7-gpk27_docker_579e832e-c41c-11e9-9d31-025000000001_0
k8s_POD_coredns-fb8b8dccf-2mdnw_kube-system_bfbca97d-c3e6-11e9-9d31-025000000001_0
k8s_POD_coredns-fb8b8dccf-nsrj4_kube-system_bfbdb38c-c3e6-11e9-9d31-025000000001_0
k8s_POD_etcd-docker-desktop_kube-system_3773efb8e009876ddfa2c10173dba95e_0
k8s_POD_kube-apiserver-docker-desktop_kube-system_7c4f3d43558e9fadf2d2b323b2e78235_0
k8s_POD_kube-controller-manager-docker-desktop_kube-system_9c58c6d32bd3a2d42b8b10905b8e8f54_0
k8s_POD_kube-proxy-v4fc5_kube-system_bfc80631-c3e6-11e9-9d31-025000000001_0
k8s_POD_kube-scheduler-docker-desktop_kube-system_124f5bab49bf26c80b1c1be19641c3e8_0
k8s_POD_redis_default_4d21ee44-c5c7-11e9-9d31-025000000001_0
k8s_compose_compose-6c67d745f6-q54n8_docker_57a0f7e9-c41c-11e9-9d31-025000000001_0
k8s_compose_compose-api-57ff65b8c7-gpk27_docker_579e832e-c41c-11e9-9d31-025000000001_0
k8s_coredns_coredns-fb8b8dccf-2mdnw_kube-system_bfbca97d-c3e6-11e9-9d31-025000000001_0
k8s_coredns_coredns-fb8b8dccf-nsrj4_kube-system_bfbdb38c-c3e6-11e9-9d31-025000000001_0
k8s_etcd_etcd-docker-desktop_kube-system_3773efb8e009876ddfa2c10173dba95e_0
k8s_kube-apiserver_kube-apiserver-docker-desktop_kube-system_7c4f3d43558e9fadf2d2b323b2e78235_0
k8s_kube-controller-manager_kube-controller-manager-docker-desktop_kube-system_9c58c6d32bd3a2d42b8b10905b8e8f54_0
k8s_kube-proxy_kube-proxy-v4fc5_kube-system_bfc80631-c3e6-11e9-9d31-025000000001_0
k8s_kube-scheduler_kube-scheduler-docker-desktop_kube-system_124f5bab49bf26c80b1c1be19641c3e8_0
k8s_redis_redis_default_4d21ee44-c5c7-11e9-9d31-025000000001_0
```

Note that all the entries listed here are Docker containers. We can use the `winfo` command to see a Docker container’s supported actions.

```
wash docker/containers ❯ winfo sleepy_heisenberg
Path: /Users/enis.inan/Library/Caches/wash/mnt452975753/docker/containers/sleepy_heisenberg
Name: sleepy_heisenberg
CName: sleepy_heisenberg
Actions:
- list
- exec
Attributes:
  atime: 2019-08-21T21:30:07-07:00
  crtime: 2019-08-21T21:30:07-07:00
  ctime: 2019-08-21T21:30:07-07:00
  mtime: 2019-08-21T21:30:07-07:00
```

It looks like Docker containers support the `exec` action. That means we can execute commands on them. Let's find out some information about `sleepy_heisenburg`.

```
wash docker/containers ❯ wexec sleepy_heisenberg uname
Linux
```

Nice! Also notice that the Docker container supports the `list` action. That means it’s modeled as a directory, so we can cd into it.

```
wash docker/containers ❯ cd sleepy_heisenberg
wash docker/containers/sleepy_heisenberg ❯ ls
fs            log           metadata.json
```

Everything listed here is specific to the `sleepy_heisenberg` container. For example, `log` represents the `sleepy_heisenberg` container’s log. Let's look at its supported actions with `wsinfo`:

```
wash docker/containers/sleepy_heisenberg ❯ winfo log
Path: /Users/enis.inan/Library/Caches/wash/mnt452975753/docker/containers/sleepy_heisenberg/log
Name: log
CName: log
Actions:
- read
- stream
Attributes: {}
```

Notice that `log` doesn’t support the `list` action. That means it’s represented as a file, so if we try to cd into it, the action fails:

```
wash docker/containers/sleepy_heisenberg ❯ cd log
cd: not a directory: log
```

However, `log` does support `read` and `stream`, so we can cat and tail its contents.

```
wash docker/containers/sleepy_heisenberg ❯ cat log
root@a72c5ec85c0c:/# ls
bin  boot  dev  etc  home  lib  lib64  media  mnt  opt  proc  root  run  sbin  srv  sys  tmp  usr  var
root@a72c5ec85c0c:/# uname
Linux
wash docker/containers/sleepy_heisenberg ❯ tail -f log
===> log <===
root@a72c5ec85c0c:/# ls
bin  boot  dev  etc  home  lib  lib64  media  mnt  opt  proc  root  run  sbin  srv  sys  tmp  usr  var
root@a72c5ec85c0c:/# uname
Linux
```

(Hit `Ctrl+C` to cancel `tail -f`).

That’s enough navigation for now. Let’s go back to the Wash root.

```
wash docker/containers/sleepy_heisenberg ❯ cd $W
wash . ❯
```

Notice that our prompt changed back to `wash . >`. That means we are indeed back at the Wash root. The `W` environment variable stores the Wash root’s absolute path, so you can invoke cd `$W` anytime you want to go back to the Wash root.

# Exercises
1. You can tab-complete entries! Try using tab-completion to type `ls docker/containers/k8s_POD_compose-6c67d745f6-q54n8_docker_57a0f7e9-c41c-11e9-9d31-025000000001_1`.

2. You can also glob entries! The following parts are meant to show you some interesting things that you can do with globbing.

    1. What’s the output of each of the following commands? Try to figure out the answer without invoking the command. Also, it’s OK to give a high-level overview of the output (i.e. you don’t have to worry about whitespace and about getting the paths right). For example, something like “prints out all the plugins” is an acceptable answer for the command `echo *` (as well as a more specific answer like “prints out aws, docker, kubernetes, and gcp”).

        1. `echo docker/*`
        2. `echo docker/containers/*`
        3. `echo docker/containers/k8s*`
        4. `echo docker/containers/k8s*compose*`
        5. `echo docker/volumes/*`

        <details><summary>Expand to reveal answers</summary><ol>
          <li>The containers and volumes directories.</li>
          <li>All Docker containers.</li>
          <li>All Docker containers that start with <code>k8s</code></li>
          <li>All Docker containers that contain the <code>k8s*compose</code> string</li>
          <li>All Docker volumes</li>
        </ol></details>

    1. How would you tail every container’s log file? Hint: The invocation is of the form `tail -f <glob>`.

        <details><summary>Expand to reveal answer</summary>
        <code>tail -f docker/containers/*/log</code></details>

1. This exercise is broken up into several parts.

    1. We saw three entries when we ls’ed a Docker container: `log`, `fs`, and `metadata.json`. We already know that the `log` entry represents the container’s log. What do you think `fs` represents? Hint: Try cd’ing into it and ls’ing stuff.

        <details><summary>Expand to reveal answer</summary>
        <code>fs</code> represents the root directory of the container’s filesystem. It lets you navigate through the container as if you were logged onto it via something like SSH. As you’ll soon see, this lets you do some pretty cool stuff.</details>

    1. Inside the same container’s ‘directory’, what command lets you read its `/var/log/messages` file? What command lets you tail it? Hint: `cat` lets you read a file. `tail -f` lets you tail it.


        <details><summary>Expand to reveal answer</summary>
        <code>cat fs/var/log/messages</code> lets you read the <code>/var/log/messages</code> file. <code>tail -f fs/var/log/messages</code> lets you tail the <code>/var/log/messages</code> file. Thus, you can read/tail a Docker container’s log files as if you were logged onto it.</details>

    1. What command lets you tail every container’s `/var/log/messages` file? Hint: See Exercise 2b’s answer.

        <details><summary>Expand to reveal answer</summary>
        <code>tail -f docker/containers/*/fs/var/log/messages</code>. Thus, you can tail log files on multiple containers.</details>

    1. Again inside the same container’s ‘directory’, what command lets you tail every file with the `.log` extension in its `/var/log` directory? Hint: The glob `**/*.log` matches every file with the `.log` extension, including subdirectories.

        <details><summary>Expand to reveal answer</summary>
        <code>tail -f fs/var/log/**/*.log</code>. This exercise is meant to remind you that everything in Wash is an entry, including a container’s files and directories. That means you can still glob them just like you would if you were logged onto the container.</details>
