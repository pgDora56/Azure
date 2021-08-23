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


## Schedule data structure update
* `No`はいらなさそう
    * `Circle`と`ID`をUniqueキーとして扱う
    * RDBで日付と時間でソートさせれば良い
* `date`を独立させる
    * `end.date` が仕事してなかった
* `created`, `updated` の追加
    * calenderAPIでCreatedTimeとUpdatedTimeを取得できるので，そこから取得．
        * カレンダー情報を取得した上で，この時間が(5分＋マージン数秒)以内のものだけ更新処理をさせるようにする
    * これを活用してNew update viewを作る
* `deleted` の追加
    * 削除されたデータもすべて取得するように変更．calenderAPIのStatusで `cancelled` となっているかどうかで削除されたものであるかが確認できるため，これを活用して削除されたかを確認．開催日を超えたら通常のイベント同様にデータ自体を削除．

### 旧データ
```
"22ollf8qk3r8ldsrlol8f59iia": {
    "no": 9,
    "circle": "ikm",
    "id": "22ollf8qk3r8ldsrlol8f59iia",
    "title": "企画名",
    "description": "説明",
    "start": {
        "date": "",
        "time": "22:00"
    },
    "end": {
        "date": "",
        "time": "01:00"
    },
    "offline": false
},
```

### 新データ
```
"22ollf8qk3r8ldsrlol8f59iia": {
    "circle": "ikm",
    "id": "22ollf8qk3r8ldsrlol8f59iia",
    "title": "企画名",
    "description": "説明",
    "date": "2021/08/28(Sat)",
    "start": "22:00",
    "end": "01:00",
    "offline": false,
    "created": "2021/08/12 21:12:32",
    "updated": "2021/08/23 11:00:23(形式は適宜)",
    "deleted": false
},
```
