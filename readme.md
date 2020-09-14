# How to run the application locally

- run `go get -u -f ./...` to install go libraries
- setup the postgres server, and update the connStr in line 20 from db/postgres.go
- use the below script to create the schema
```
create table clicks (
    color int,
    created_at TIMESTAMP
)
```
- run the application using `go run cmd/clpsec-app-server/main.go`
- play the game via localhost:8081/ws and localhost:8081/client; refresh the first page to reset the counters.

### References 
- [How to Build a Simple Web App with React, Graphql, and Go](https://medium.com/@chrischuck35/how-to-build-a-simple-web-app-in-react-graphql-go-e71c79beb1d)
- [Getting Started with React and Hello World Example](https://www.golangprograms.com/getting-started-with-react-and-hello-world-example.html)
- [Go with react](https://medium.com/@rocketlaunchr.cloud/go-with-react-de5ee4f01df9)
- [golang websocket](https://gist.github.com/rogerwelin/4d9891b2440e62e88ce34f9f7bdc31b7)

### Contact
Please contact with [mac](macma.hk@gmail.com) if you have further questions
