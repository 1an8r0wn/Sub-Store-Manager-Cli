package docker

import (
	"fmt"
)

type Dockerfile struct {
	Node            string
	NodeWithDataEnv string
	FE              string
}

func getNodeDockerfileStr(sh string, hash string) string {
	return fmt.Sprintf(`FROM node:16-alpine

WORKDIR /app

# 安装 nginx
RUN apk update && apk add --no-cache nginx=1.24.0-r7
# 写入 nginx 配置
RUN printf 'server {\n  listen 80;\n  server_name localhost;\n\n  location /%s/ {\n    proxy_set_header Host $http_host;\n    proxy_set_header X-Real-IP $remote_addr;\n    proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;\n    proxy_pass http://localhost:3000/;\n  }\n\n  location / {\n    return 404;\n  }\n}' > /etc/nginx/http.d/default.conf
# 创建启动脚本
RUN printf '#!/bin/sh\n\nnginx -g "daemon off;" &\n\n%s' > start.sh
# 赋予执行权限
RUN chmod +x start.sh

# 复制主程序
COPY . .

# 复制配置文件
RUN mkdir config

# 开放端口
EXPOSE 80

# 启动脚本
CMD ["sh", "start.sh"]
`, hash, sh)
}

var DockerfileStr = Dockerfile{
	FE: `FROM debian:bullseye-slim AS downloader

WORKDIR /app
RUN apt-get update && \
    apt-get install -y curl unzip && \
    rm -rf /var/lib/apt/lists/* && \
    curl -LJO https://sub-store-org.github.io/resource/ssm/nginx.conf && \
    curl -o dist.zip -LJ https://github.com/sub-store-org/Sub-Store-Front-End/releases/latest/download/dist.zip && \
    unzip dist.zip

FROM nginx:alpine AS runner

WORKDIR /app

COPY --from=downloader /app/dist ./www
COPY --from=downloader /app/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]`,
}

type dockerfileType string

var (
	dockerfileTypeFE             dockerfileType = "fe"
	dockerfileTypeNode           dockerfileType = "node"
	dockerfileTypeNodeOldVersion dockerfileType = "node-old-version"
)

// GetDockerfileStr returns the Dockerfile string based on the type provided.
func getDockerfileStr(t dockerfileType, hash string) string {
	switch t {
	case dockerfileTypeNodeOldVersion:
		return getNodeDockerfileStr("cd /app/config && node ../sub-store.bundle.js", hash)
	case dockerfileTypeNode:
		return getNodeDockerfileStr("SUB_STORE_DATA_BASE_PATH=/app/config node sub-store.bundle.js", hash)
	case dockerfileTypeFE:
		return DockerfileStr.FE
	default:
		return ""
	}
}
