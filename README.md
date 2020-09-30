# planexdownload
planexの温度計WS-USB01-THPを使ってクラウドにアップされた環境データを
定期的にダウンロードするプログラムです。
planexのクラウドはデータの保存期間が一か月のためその期間内にデータをローカルに
保存することによって保存期間後もデータを閲覧することができます。

# How to Use
リポジトリのクローン後ビルドしてください

```
go build planexdownload
```

config.yamlに必要な設定情報を入力してください。
config.yamlはmainパッケージと同一ディレクトリに配置しておいてください。それ以外だと認識できません。
設定情報は配列にして複数入力しても構いません

```
-  
    #製品名です 例:WH-USB01-THP
    DeviceName: "REPLACE"
    #各デバイスに付けた名前です
    NicName   : "REPLACE"
    #MACアドレスです　製品裏面に記載があります
    MacAddress: "REPLACE"
    #API叩くようのトークンです　PLANEXのサイトにログインすると確認できます
    Token     : "REPLACE"
    #ダウンロードデータの保存先パスです
    SavePath : "REPLACE"
    #データ取得日を現在日付から何日の形式で指定します
    #下記の例だと現在の日付から三日前のデータを一日分取得します
    From : -3
    To : -3
```

あとは実行するだけでOK
処理失敗時のLOGはmainパッケージのあるディレクトリに吐き出されます。
