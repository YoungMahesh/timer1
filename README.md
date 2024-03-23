# Timer1

- A simple CLI-app to track time spent on a task/project

### compile app

```bash
go build -o timer1

# execute `./timer1` to check if app is working as expected
```

### install/uninstall app on debian/ubuntu

```bash
# install app (move compiled script to /usr/local/bin)
sudo mv ./timer1 /usr/local/bin

# uninstall app
sudo rm /usr/local/bin/timer1
```

### commands

```bash
timer1 start project1  # start timer for a new project

timer1 ls   # print time elapsed for currently running project

# print time elapsed for currently project and stop/delete current project
timer1 stop
```
