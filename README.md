# lunchtogether
Backend for a simple compiled Lunch menu application

## How to develop? 

Set up base go environment:
https://go.dev/doc/tutorial/getting-started#install

Add what you are planning to do as an Issue in Github.
Set yourself as assigned.

```
$ cd <go environment base>

$ git clone https://github.com/Jeberlen/lunchtogether.git , ssh clone or whatever you want! 

$ go mod tidy
````

And then just code! 

## How to run using docker? 

Build our docker images:

```
$ docker build --pull --rm -f "lunchtogether/Dockerfile.db" -t lunchtogetherdb:latest "lunchtogether"

$ docker build --pull --rm -f "lunchtogether/Dockerfile.server" -t lunchtogetherserver:latest "lunchtogether"
````

Run docker images in order:

```
$ docker run --rm -it -p 5432:5432/tcp lunchtogetherdb:latest

$ docker run --rm -it -p 8080:8080/tcp lunchtogetherserver:latest
````

## How to run locally?

Install postgres
Create a database called lunchtogetherdev
Run init.sql from root of repository
Change postgres@172.17.0.2 -> postgres@localhost

```
$ go server.go
```

#### Changes to schema.graphqls: 

Make you changes and run from project root: 

```
$ go run github.com/99designs/gqlgen generate
```
