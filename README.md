# Need For Typing

## V1.1.0 Feature

### 登录
1. **注册**
2. **登录**
3. **邮箱验证**
4. **忘记密码**

### 素材和音乐
1. **获取七牛云上传token & 获取下载资源url**（已支持七牛云存储）
2. **获取歌曲与素材的详细信息**

### 文章生成（gRPC + SSE）
1. **基于 gRPC 生成文章**
2. **使用 SSE（Server-Sent Events）推送文章流式内容**
3. **支持 gRPC 客户端与 Python 服务器交互**

### 分数
1. **上传总分**
2. **获取个人所有最佳成绩**
3. **获取某首歌的所有玩家成绩**

### 联机
1. **识别 1P / 2P**
2. **加入房间，识别房间**
3. **分数实时同步**

### 部署指南
```sh
git clone git@github.com:muxi-mini-project/2025-Need-For-Typing-Backend.git
go mod tidy
go run main.go
