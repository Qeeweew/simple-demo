# simple-demo

## 抖音项目服务端简单示例

具体功能内容参考飞书说明文档

项目使用的go版本是go1.20

项目需要同时运行cmd中的main.go和router.go

```shell
cd cmd
go run .
```

### 项目目录说明
- cmd 存放启动 server 相关
- common 存放公共的基础模块，比如全局错误码，封装的工具（为了简化开发，将一系列逻辑封装成一个函数）
    - config 配置模块，包括配置文件读取的实现，项目以 `app.yaml` 作为配置文件，这种存在个体差异的配置通过 `.gitignore` 设置了 commit 时忽略，其中的 `app.example.yaml` 代替忽略的 `app.yaml` 揭示了配置文件的格式
    - db 对mysql和redis数据库进行初始化的代码实现
    - log 封装的日志模块，应该仅用于与业务相关的部分
    - model 对应数据库中的表结构
    - result 封装的返回结果模块
- controller 负责对发送的数据进行处理，调用 service 中实现的业务逻辑，并根据所调用的 service 层的函数的返回值判断并处理返回一次请求的响应信息（比如 200 OK 跟上相应的数据）
- doc 一些文档
- middleware 业务中间件模块，比如用户鉴权
- service 业务逻辑的实现，比如一次点赞，在 service 层要调用 db 包，对数据库进行操作，并检查是否产生错误，若无返回 nil ，有则抛给调用该 service 函数的 controller，将错误交给 controller，并在 controller 中判断接下来应该怎么做，比如直接将错误信息和错误代码返回给用户
- util 封装了一些辅助实现代码逻辑而构成的小工具
- test 业务测试模块