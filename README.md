# Timer1

- A simple CLI-app to track time spent on a task/project

### compile app

```bash
go build -o timer

# execute `./timer1` to check if app is working as expected
```

### install/uninstall app on debian/ubuntu

```bash
# install app (move compiled script to /usr/local/bin)
sudo mv ./timer /usr/local/bin

# uninstall app
sudo rm /usr/local/bin/timer
rm -rf ~/.timer1
```

### commands

```bash
timer1 start <project-name> # start timer
timer start xyz  # start timer for a new project named "xyz"

timer ls   # print time elapsed for currently running project

timer edit 10  # add 10 minutes to currently running session
timer edit -15  # remove 15 minutes from currently running sesion

timer stop  # start current project

timer restart  # restart currently stopped project
```
