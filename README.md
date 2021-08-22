- generate grpc files
```shell
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative route-proto/route_guide.proto
```
- try to start the server but error:
```shell     
server/server.go:8:2: no required module provides package github.com/SarathLUN/grpc-go/examples/data; to add it:
        go get github.com/SarathLUN/grpc-go/examples/data
route-proto/route_guide_grpc.pb.go:7:2: no required module provides package google.golang.org/grpc; to add it:
        go get google.golang.org/grpc
route-proto/route_guide_grpc.pb.go:8:2: no required module provides package google.golang.org/grpc/codes; to add it:
        go get google.golang.org/grpc/codes
server/server.go:11:2: no required module provides package google.golang.org/grpc/credentials; to add it:
        go get google.golang.org/grpc/credentials
route-proto/route_guide_grpc.pb.go:9:2: no required module provides package google.golang.org/grpc/status; to add it:
        go get google.golang.org/grpc/status
server/server.go:12:2: no required module provides package google.golang.org/protobuf/proto; to add it:
        go get google.golang.org/protobuf/proto
route-proto/route_guide.pb.go:10:2: no required module provides package google.golang.org/protobuf/reflect/protoreflect; to add it:
        go get google.golang.org/protobuf/reflect/protoreflect
route-proto/route_guide.pb.go:11:2: no required module provides package google.golang.org/protobuf/runtime/protoimpl; to add it:
        go get google.golang.org/protobuf/runtime/protoimpl
```
- this error due to we are using `go mod`, so we need to run `go get` <u>without</u> flag `-u`:
```shell
 go get github.com/SarathLUN/grpc-go/examples/data
 go get google.golang.org/grpc/examples/data
 go get google.golang.org/grpc
 go get google.golang.org/grpc/codes
 go get google.golang.org/grpc/credentials
 go get google.golang.org/grpc/status
 go get google.golang.org/protobuf/proto
 go get google.golang.org/protobuf/reflect/protoreflect
 go get google.golang.org/protobuf/runtime/protoimpl
```
- try to run server again, but got anther error:
```shell                             
../../../../pkg/mod/google.golang.org/grpc@v1.40.0/credentials/credentials.go:31:2: missing go.sum entry for module providing package github.com/golang/protobuf/proto; to add:
        go mod download github.com/golang/protobuf
../../../../pkg/mod/google.golang.org/grpc@v1.40.0/internal/binarylog/method_logger.go:28:2: missing go.sum entry for module providing package github.com/golang/protobuf/ptypes; to add:
        go mod download github.com/golang/protobuf
../../../../pkg/mod/google.golang.org/genproto@v0.0.0-20200806141610-86f49bd18e98/googleapis/rpc/status/status.pb.go:28:2: missing go.sum entry for module providing package github.com/golang/protobuf/ptypes/any; to add:
        go mod download github.com/golang/protobuf
```
- so run `go mod download`
```shell
go mod download github.com/golang/protobuf
```
- now test server with client
```shell
go run server/server.go
2021/08/22 16:17:06 starting server
2021/08/22 16:17:06 server is listening on port: 10000

```
```shell
gor client/client.go     
2021/08/22 16:17:24 Getting feature for point (409146138 -746188906)
2021/08/22 16:17:24 name:"Berkshire Valley Management Area Trail, Jefferson, NJ, USA" location:{latitude:409146138 longitude:-746188906}

```
