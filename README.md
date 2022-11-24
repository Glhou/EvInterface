# EvInterface

最後に日本語版

## Installation

Run the server with ```run go *.go```, you can also build the server and run the executable file afterward.

You can choose the port of the server on the `server.go` file.

After that plase initialize the database (only for the first time or if there is no database) by accessing to the `/reset` route (`http://localhost:PORT/reset`).

Then you can access to the interface on `http://localhost:PORT/`.

## Actions
All the actions are conducted using POST requests.
1) ### Add a Token

    You can add a token by using the route `/token`, here is the body :

    ```json
    {
        "TokenId": "your token id" (string),
        "TokenLat" : latitude (float64),
        "TokenLon" :  longitude (float64),
        "TokenPrice" : price (int)
    }
    ```

2) ### Add a Bid / Car

    You can add a bid by using the route `/bid`, here is the body :

    ```json
    {
    "CarId": "your car id" (string),
    "CarEnergy" : energy_level (int),
    "CarRadius" : radius_in_km (float),
    "CarLat" : latitude (float64),
    "CarLon": longitude (float64),
    "Price" : price (int),
    "TokenId" : "The tokenId of the token the car bid on" (string)
    }
    ```

3) ### Auction

   1) #### If there is a winner

        You can send the winner and the token and we will remove all the others bidder (the link of the winner will be removed after 30 seconds). For that use the route `/auction` here is the boddy :

        ```json
            {
                "WinnerCarId" : "the id of the car who won" (string),
                "TokenId" : "the id of the token that the car won" (string)
            }
        ```

    2) #### If there is no bidder after 30 minutes

        You can remove the token with the same route `/auction` but using this body :

        ```json
            {
                "WinnerCarId" : "-1" (string),
                "TokenId" : "the id of the token that you want to remove" (string)
            }
        ```

## 日本語版

### インストール

サーバーは ``run go *.go`` で実行します。また、サーバーをビルドしてから実行ファイルを実行することも可能です。

サーバーのポート番号は `server.go` ファイルで選択することができます。

その後、`/reset` ルート (`http://localhost:PORT/reset`) にアクセスして、データベースを初期化する (初回のみ、またはデータベースが存在しない場合)。

その後、`http://localhost:PORT/` のインターフェースにアクセスしてください。


### アクション
すべてのアクションはPOSTリクエストで行われます。
1) #### トークンの追加

    トークンを追加するには、ルート `/token` を使用します。

    ```json
    {
        "TokenId": "あなたのトークンID" (string),
        "TokenLat" : 緯度 (float64),
        "TokenLon" : 経度 (float64),
        "TokenPrice" : 価格 (int)
    }
    ```

2) #### 入札の追加 / 車の追加

    入札を追加するには、ルート `/bid` を使用します。

    ```json
    {
    "CarId": "あなたの車のID" (string),
    "CarEnergy" : エネルギーレベル (int),
    "CarRadius" : radius_in_km (float),
    "CarLat" : 緯度 (float64),
    "CarLon" : 経度 (float64,。
    "Price" : 価格 (int),
    "TokenId" : "その車が入札したトークンのtokenId" (string)
    }
    ```

3) #### オークション

   1) ##### 落札者がいる場合

        落札者とトークンを送信すると、他の入札者をすべて削除します (落札者のリンクは30秒後に削除されます)。そのためには、ルート `/auction` を使用します。

        ```json
            {
                "WinnerCarId" : "落札した車のID" (string),
                "TokenId" : "その車が勝ったトークンのID" (string)
            }
        ```

    2) ##### 30分経過しても落札者がいない場合

        同じルート `/auction` でトークンを削除することができますが、次のようなボディを使用します。

        ```json
            {
                "WinnerCarId" : "-1" (string),
                "TokenId" : "削除したいトークンのID" (string)
            }
        ```