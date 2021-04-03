# Intro Quiz Portal Square "Azure"

[Intro Quiz Portal Square "Azure"](https://azure.mocho.ml)

## 概要
* Googleカレンダーと連携、複数のイントロイベントを一括で見ることができる

## 機能
* 団体ごとにカレンダーIDをJSONで一括管理する
    * 別途、Azure単発カレンダーを設ける。単発イベはここに追加。
    * JSONでは、以下の内容を扱う。
        * 略称(半角4/全角2程度が上限?) - カレンダー表示用
        * 正式名称 - サークル詳細表示用
        * サークル等概要 - サークル詳細表示用
        * 任意の項目
            * Twitterだけは需要が多そうなので特化させる
            * 注意するべき事柄のための `Warning` 句を用意
* 5分ごとにカレンダーを走査して、直近3ヶ月を取得。JSONとして保存する。
    * 要求に応じて、保存されたものからデータを取得。これによりGoogle Calender APIの使用回数を減らす。っていうか時間結構掛かる印象が。そこへの対応ってのもある。
            * 多分取得で2~3秒くらい。
    * 適当にDB使っても良いんだけど、Goの速さと、Jsonバッケージの使いやすさに甘えた形。速度的にまずくなったら修正を検討。

![Implementation Image](implementation_image.PNG "Implementation Image")

## JSON

### circles.json

```json
"aiq": {
    "simple": "AIQ", 
    "name": "AIQオンライン",
    "overview": "サークル概要",
    "detail":[
        {
            "item": "Twitter",
            "value": "aiq_list"
        },
        {
            "item": "公式サイト",
            "link": "https://aaaaaaaaa.jp/bbbbbbb",
            "value": "Website"
        }
    ],
    "contact": "xxx@example.com",
    "url": "xxxxxxxxxxxxxxxx@group.calendar.google.com"
}
```

