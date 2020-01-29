# gRPCExample
## Basic gRPC example with client and server

1. Clone project.
1. You might need to install dependencies like protobuf and gRPC (will add later).
1. To genereate:
> cd echo && protoc --go_out=plugins=grpc:. --grpc-gateway_out=logtostderr=true,grpc_api_configuration=./router.yaml:.  ./echo.proto
1. First run the server, by default it will listen on port 8082.
1. You run the client to get basic response back.



