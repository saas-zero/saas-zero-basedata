基于zero构建的多租户微服务版本

基础数据服务  

地址：https://github.com/saas-zero/saas-zero-basedata

项目根目录
go generate ./ent

进入api
goctl api go -api basedata_service.api -dir . -style go_zero

进入grpc
goctl rpc protoc basedata_service.proto --go_out=. --go-grpc_out=. --zrpc_out=. -m --style go_zero

进对应目录
go run basedataservice.go -f etc/basedataservice.yaml

