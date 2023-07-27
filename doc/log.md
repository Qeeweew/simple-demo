# log包使用提醒

1. log基本使用

   ``````go
   // 有四种模式：调试（Debug）,错误信息（Error）,提醒（Info）,警告（Warn）
   
   // 具体到使用就是log.Logger.你想要输出的模式(Debug Error Info Warn).(你要输出的内容)
   
   // 举个例子
   log.Logger.Error("数据插入失败")
   ``````

2. 日志输出

   如果你开启的gin的模式是debug开发模式的话，只会在终端输出，其他模式会在项目文件下创建`cmd/runtime/log/`，输出的内容全部在cmd/runtime/log里面的文件里。

3. 注意事项

   在记录日志时，不要输出一些敏感信息，例如密码、密钥等。
