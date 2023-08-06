# go_app_template
This is a basic template to start a new golang application based on clean architecture 

## Setup of the component:

Must have golang installed version >= 1.18

make file consists of 4 steps: generate, test, build, run you can run all of them

```make all```

Run test cases run

```make test```
or

```go test -v ./tests```

Start the http server on port 9090:

```make run```

## Run By Docker

```
 docker build -t go_app .
 
 docker run -p 9090:9090 go_app
 
```

