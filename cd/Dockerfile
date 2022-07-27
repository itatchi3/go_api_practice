FROM golang:1.18 as build
#FROMgolang:1.14-alpine

# コンテナログイン時のディレクトリ指定
WORKDIR /opt/cd_test
# ホストのファイルをコンテナの作業ディレクトリにコピー

COPY . .
# ADD . .

# ビルド
RUN go build -o app main.go

# 起動
CMD ["/opt/cd_test/app"]