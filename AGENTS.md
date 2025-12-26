# MCP 服务工具

## 开发规范及强制性要求

以下为项目开发过程中必须遵守的强制性规范：

| 序号 | 要求类别     | 具体要求                                                                                 |
|----|----------|--------------------------------------------------------------------------------------|
| 1  | 开发语言     | Golang 的版本为 1.25.5+                                                                  |
| 2  | 日志管理     | 采用 Zap 做为项目统一的日志管理工具                                                                 |
| 3  | 错误判断     | 非nil类的error判断不允许使用双等号（`==`），统一采用标准库`errors.Is`方法判断error类型                            |
| 4  | 变量命名     | 变量声明不允许与当前软件包或引用软件包同名，也不能跟Go语言保留字或标准库方法名称同名                                          |
| 5  | 业务注释     | 所有业务实现方法必须采用添加符合 Golang 风格和标准的中/英文注释                                                 |
| 6  | API 注释   | 所有的 API 或前端 Web 控制器必须添加符合 Golang 风格和标准的中/英文注释                                        |
| 7  | 结构体注释    | 所有结构体及成员变量必须添加符合 Golang 风格和标准的中/英文注释                                                 |
| 8  | 配置一致性    | 系统配置文件与配置工具声明的数据项或成员变量必须一致，非必填项配置要做默认值配置                                             |
| 9  | 单元测试     | 所有的业务实现或工具方法必须创建对应的单元测试，覆盖率必须在 100% 以上                                               |
| 10 | gRPC 同步  | gRPC 通信类功能必须同步修改 Proto 文件内容                                                          |
| 11 | API 参数   | 明确定义的或已知固定的 API 参数，必须使用外部的结构体进行传输交互，不允许采用内嵌结构和 Map 传值                                |
| 12 | API 路由   | API 必须正确定义路由的提交方式：`GET` / `PUT` / `POST` / `DELETE` / `HEAD` / `OPTIONS`             |
| 13 | API 响应   | API 响应必须明确正确的 HTTP 状态码                                                               |
| 14 | 结构体声明    | 所有的结构体都必须放在pkg/types，不允许在其他包中声明结构体，不能使用内嵌结构体                                         |
| 15 | 常量声明     | 对于大于一次的多次字内容引用，要采用常量声明，不允许直接使用字面量                                                    |
| 16 | 禁止批量变更   | 禁止使用命令行或脚本对代码文件进行批量修改操作                                                              |
| 17 | 代码变更     | 完成代码的变更后，要通过 `gofmt` 格式化代码，然后编译后再进行相关测试或功能验证                                         |
| 18 | 结构体预热    | pkg/types下的结构体必须添加到pkg/utils/json/pretouch.go进行JSON预热                                |
| 19 | 打包构建     | 打包编译时必须添加编译标签：`go build -tags="sonic"`。Linux 下需配置 `GOEXPERIMENT="jsonv2,greenteagc"` |
| 20 | 重复功能代码处理 | 在实现新功能前，要先检查当前项目是否已经存在类似的功能代码，要优先复用已有的功能代码，在全局范围内禁止出现重复的功能代码                         |
| 21 | JSON编译码  | 对于JSON数据的编译码，统一采用`pkg/utils/json/json`处理                                             |

## 代码组织结构与命名规范

### 目录结构约束

- API层：在`internal/api`目录下，必须为每个功能模块创建独立的子目录。所有 API 相关的代码文件必须按模块功能拆分存放至对应模块目录中。
- Service层：在`internal/services`目录下，并为其下每个功能模块创建对应的子目录，用于存放该模块的业务逻辑（Service层）代码文件。

### 文件命名约束

- 所有API控制器（Controller）文件 的文件名必须以 handlers_ 为前缀。
- 对应API结构体（控制器）的名称必须以 Handler 结尾。

### 代码组织约束

- 禁止所有 API 模块共用或继承同一个基类。
- 必须为每个独立的功能模块创建专属的、以 handlers_ 开头的文件，并在其中定义以 Handler 结尾的专属结构体。

### 本项目的强制性依赖库

  | 序号 | 依赖库                                    | 版本                   | 用途                 | 备注                         |
  |----|----------------------------------------|----------------------|--------------------|----------------------------|
  | 1  | Go                                     | 1.25.5+              | 开发语言               | -                          |
  | 2  | github.com/hashicorp/raft              | v1.7.3               | Raft 一致性协议         | -                          |
  | 3  | github.com/hashicorp/raft-boltdb/v2    | v2.3.1               | Raft 日志存储          | 需要编译 tag: hashicorpmetrics |
  | 4  | github.com/docker/docker               | v28.5.2+incompatible | Docker API 客户端     | -                          |
  | 5  | github.com/containerd/containerd/v2    | v2.2.1               | containerd API 客户端 | -                          |
  | 6  | go.uber.org/zap                        | v1.27.1              | 结构化日志              | -                          |
  | 7  | gopkg.in/natefinch/lumberjack.v2       | v2.2.1               | 日志轮转               | -                          |
  | 8  | github.com/shirou/gopsutil/v4          | v4.25.11             | 系统资源采集             | -                          |
  | 9  | gorm.io/gorm                           | v1.31.1              | ORM 框架             | -                          |
  | 10 | gorm.io/driver/mysql                   | v1.6.0               | MySQL 驱动           | -                          |
  | 11 | gorm.io/driver/postgres                | v1.6.0               | PostgreSQL 驱动      | -                          |
  | 12 | github.com/redis/go-redis/v9           | v9.17.2              | Redis 客户端          | -                          |
  | 13 | go.etcd.io/bbolt                       | v1.4.3               | 嵌入式 KV 存储          | -                          |
  | 14 | github.com/gin-gonic/gin               | v1.11.0              | HTTP Web 框架        | -                          |
  | 15 | github.com/golang-jwt/jwt/v5           | v5.3.0               | JWT 认证             | -                          |
  | 16 | github.com/go-ldap/ldap/v3             | v3.4.12              | LDAP 客户端           | -                          |
  | 17 | github.com/NVIDIA/go-nvml              | v0.13.0-1            | NVIDIA GPU 监控      | -                          |
  | 18 | github.com/google/uuid                 | v1.6.0               | UUID 生成            | -                          |
  | 19 | golang.org/x/crypto                    | v0.46.0              | 加密库                | -                          |
  | 20 | github.com/opencontainers/runtime-spec | v1.3.0               | OCI 运行时规范          | -                          |
  | 21 | github.com/bytedance/sonic             | v1.14.2              | 高性能 JSON 序列化/反序列化  | -                          |
  | 22 | github.com/goccy/go-json               | v0.10.5              | 高性能 JSON 库         | -                          |
  | 23 | github.com/json-iterator/go            | v1.1.12              | 高性能 JSON 迭代器       | -                          |
  | 24 | google.golang.org/grpc                 | v1.77.0              | gRPC 框架            | -                          |
  | 25 | google.golang.org/protobuf             | v1.36.11             | Protocol Buffers   | -                          |
  | 26 | github.com/quic-go/quic-go             | v0.57.1              | QUIC/HTTP3 协议支持    | -                          |
  | 27 | github.com/gin-contrib/cors            | v1.7.6               | CORS 跨域支持          | -                          |
  | 28 | github.com/gin-contrib/gzip            | v1.2.5               | Gzip 压缩中间件         | -                          |
  | 29 | github.com/gin-contrib/pprof           | v1.5.3               | pprof 性能分析         | -                          |
  | 30 | github.com/gin-contrib/requestid       | v1.0.5               | 请求 ID 追踪           | -                          |
  | 31 | github.com/gin-contrib/zap             | v1.1.6               | Gin Zap 日志集成       | -                          |
  | 32 | github.com/lf4096/gin-compress         | v0.1.0               | Gin 压缩中间件          | -                          |
  | 33 | github.com/swaggo/swag                 | v1.16.6              | Swagger 文档生成       | -                          |
  | 34 | github.com/swaggo/gin-swagger          | v1.6.1               | Gin Swagger 集成     | -                          |
  | 35 | github.com/swaggo/files                | v1.0.1               | Swagger 静态文件       | -                          |
  | 36 | github.com/stretchr/testify            | v1.11.1              | 单元测试框架             | -                          |
  | 37 | github.com/docker/go-connections       | v0.6.0               | Docker 连接工具        | -                          |
  | 38 | github.com/hashicorp/go-hclog          | v1.6.3               | HashiCorp 日志库      | -                          |
  | 39 | github.com/goccy/go-yaml               | v1.19.1              | YAML 解析            | -                          |
  | 40 | github.com/modelcontextprotocol/go-sdk | v1.2.0               | MCP工具库             | -                          |


必须严格遵循上面的开发规范及相关的要求。核心要点总结：按功能模块严格分离 API 与 Services 代码目录，遵循统一的 handlers_ 文件命名与 **Handler 结构体命名规范，并确保各模块间的独立性。