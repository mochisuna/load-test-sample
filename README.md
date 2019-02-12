# load-test-sample
負荷試験用のパッケージです。  
コンテナを起動するとDBとサーバーと静的ページが立ち上がります。  
サンプルデータの登録・取得・リダイレクトを順に実行できることを調べるテストなどで使える・・・かも。


# Usage
```
$ make dstart
$ make dmigrate
```

# endpoint
1. Create: `POST /v1/users`
1. Refer: `GET /v1/users/{user_id}`
1. Redirect: `POST /v1/display`
  ```
  > redirect sample body
  {
    "id": "hogehoge",
    "name":"fugafuga",
    "secret_key": "piyopiyo"
  }
  ```
