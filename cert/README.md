# 使用 openssl 创建证书

```bash
openssl req -newkey rsa:2048 -new -nodes -x509 -days 3650 -keyout key.pem -out cert.pem -subj "/CN=$1"
```
