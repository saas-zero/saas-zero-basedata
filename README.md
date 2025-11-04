基于zero构建的多租户微服务版本

基础数据服务  

地址：https://github.com/saas-zero/saas-zero-basedata

go generate ./ent

goctl api go -api basedata_service.api -dir .

goctl rpc protoc basedata_service.proto --go_out=. --go-grpc_out=. --zrpc_out=. 

go run basedataservice.go -f etc/basedataservice.yaml

