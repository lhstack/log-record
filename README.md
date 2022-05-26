## 基于golang实现的日志收集
```markdown
通过监听socket连接，接收客户端提交的日志数据，将其存储到数据库，并创建http服务器，用于暴露http接口供前端使用
socket日志编码协议，参考main.go和command.go文件
```
![图片](./1.png)