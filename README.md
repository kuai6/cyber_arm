# cyber_arm

To start working with project you should install golang dep utility:

`curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh`

-or-

`go get -u github.com/golang/dep/cmd/dep`

Then you can install all projects' dependencies using command:

`dep ensure -v`

To start application use command:

`go run main.go start`

Use following command to list all available arguments:

`go run main.go start --help`
