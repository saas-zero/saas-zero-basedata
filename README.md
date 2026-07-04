# SaaS-Zero Basedata

基于 zero 构建的多租户微服务版本 — 基础数据服务

## 生成 Ent CRUD 代码

```bash
go generate ./ent
```

## 生成 HTTP API 代码（驼峰命名）

```bash
cd api
goctl api go -api basedata_service.api -dir . -style goZero
```

## 生成 gRPC 代码（驼峰命名）

```bash
cd rpc
goctl rpc protoc basedata_service.proto --go_out=. --go-grpc_out=. --zrpc_out=. -m --style goZero
```

## 启动服务

```bash
cd rpc
go run basedataservice.go -f etc/basedataservice.yaml
```

## 注意事项

- ent 生成的包名为结构体全小写（`SysUser` → `sysuser`），不可配置
- goctl 使用 `--style goZero` 确保文件名和方法名驼峰
- 返回前端的 ID 字段使用 string 类型，避免 int64 精度丢失
- 数据库列名 snake_case，JSON 字段 camelCase
