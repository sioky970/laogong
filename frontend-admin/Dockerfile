# 构建阶段
FROM node:20-alpine AS builder

WORKDIR /app

# 安装 pnpm
RUN npm install -g pnpm

# 复制依赖文件
COPY package.json pnpm-lock.yaml pnpm-workspace.yaml ./
COPY packages ./packages

# 安装依赖（允许运行 build scripts）
RUN pnpm config set enable-pre-post-scripts true
RUN pnpm install --ignore-scripts
RUN pnpm approve-builds esbuild @parcel/watcher unrs-resolver || true
RUN pnpm install

# 安装缺失的依赖
RUN pnpm add -D @iconify/utils

# 复制源代码
COPY . .

# 构建生产版本
RUN pnpm build

# 运行阶段
FROM nginx:alpine

# 复制构建产物
COPY --from=builder /app/dist /usr/share/nginx/html

# 复制 nginx 配置
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
