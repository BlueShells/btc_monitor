# 使用官方的 Go 作为基础镜像
FROM golang:1.20-alpine

# 设置工作目录
WORKDIR /app

# 复制当前目录中的所有文件到工作目录
COPY . .

# 获取依赖
RUN go mod tidy

# 编译 Go 程序
RUN go build -o btc_block_height_exporter .

# 暴露端口
EXPOSE 9091

# 运行编译好的二进制文件
CMD ["./btc_block_height_exporter"]