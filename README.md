# clipboard share

Windows  Linux(x11)设备间共享剪切板

## 用法示例

```bash
# 运行服务端
  clipboard_share -host 192.168.1.1 -port 8080 -isServer
# 运行客户端
  clipboard_share -host 192.168.1.1 -port 8080
# 运行服务端并使用TLS加密通信
  clipboard_share -host 192.168.1.1 -port 8080 -isServer -useTls
# 运行客户端并使用TLS加密通信
  clipboard_share -host 192.168.1.1 -port 8080 -useTls
```

## TODO List

| 功能 | 状态|
|---|---|
| 连接密码 | TODO |
| 传输加密 | DONE |
| 文本复制 | DONE |
| PNG复制 | DONE |
| 文件复制 | TODO |

## 注意

证书(cert.pem)和密钥(key.pem)需要放置在cert目录下
