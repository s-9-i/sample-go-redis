# sample-go-redis
[技術書典13](https://techbookfest.org/event/tbf13) にて執筆した[SGE Go Tech Book Vol.02](https://techbookfest.org/product/vKVim3NXwgyTr8mWxRZQQQ) の第1章「スマートフォンゲーム開発におけるRedisの活用」のサンプルコードになります。

## 各種サンプルコードの実行
```shell
$ cd sample-go-redis

# Redisの起動
$ docker-compose up -d

# サンプルコード実行
$ go run ./cache
$ go run ./ranking/simple
$ go run ./ranking/timestamp
$ go run ./sorter
$ go run ./latests
$ go run ./lock
```

## 免責事項
本書に記載された内容は、情報の提供のみを目的としています。
したがって、本書を用いた開発、製作、運用は、必ずご自身の責任と判断によって行ってください。
これらの情報による開発、製作、運用の結果について、著者はいかなる責任も負いません。
