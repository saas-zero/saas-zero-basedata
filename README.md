基础数据服务

配置
菜单
API

go generate ./ent

goctl api go -api basedata_service.api -dir .

goctl rpc protoc basedata_service.proto --go_out=. --go-grpc_out=. --zrpc_out=. 

go run basedataservice.go -f etc/basedataservice.yaml

