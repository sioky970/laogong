# 构建阶段
FROM golang:1.23-alpine AS builder

WORKDIR /app

# 安装依赖
RUN apk add --no-cache gcc musl-dev

# 复制源代码
COPY . .

# 下载依赖
RUN go mod tidy && go mod download

# 构建
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main .

# 运行阶段
FROM alpine:latest

WORKDIR /app

# 安装 ca-certificates 和字体（用于图片渲染）
RUN apk --no-cache add ca-certificates fontconfig ttf-dejavu

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .

# 创建上传目录
RUN mkdir -p /app/uploads

EXPOSE 8083

CMD ["./main"]
