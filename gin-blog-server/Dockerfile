FROM golang:1.25-alpine as BUILDER
WORKDIR /gvb


#这部分是为了支持vips处理图片转换，测试---恢复需要删除
RUN apk add --no-cache \
    gcc \
    musl-dev \
    vips-dev \
    pkgconfig \
    libraw-dev \
    libheif-dev



COPY go.mod go.mod
COPY go.sum go.sum
RUN go env -w GO111MODULE=on \
    && go env -w GOPROXY=https://goproxy.cn,direct \
    && go mod download
COPY . .

#RUN cd cmd && go build -o server . 
#上面是原版的，下面是加了cgo,支持vips,恢复需要删除下面，恢复上面
RUN cd cmd && CGO_ENABLED=1 go build -o server . 


FROM alpine:3.19

#同样是为了支持vips,不行删了这行
# 添加运行时依赖（包含 DNG 和 HEIC 支持）
RUN apk add --no-cache \
    vips \
    ca-certificates \
    libraw \
    libheif

ENV WORK_PATH /gvb
WORKDIR ${WORK_PATH}
COPY --from=0 ${WORK_PATH}/cmd/server .
COPY --from=0 ${WORK_PATH}/config.docker.yml .
COPY --from=0 ${WORK_PATH}/assets/ip2region.xdb ./assets/ip2region.xdb
COPY --from=0 ${WORK_PATH}/assets/templates ./assets/templates
EXPOSE 8765

ENTRYPOINT ./server -c config.docker.yml