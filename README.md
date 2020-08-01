# バックエンド
## 概要

フレームワークはginを使っています。


## Installation

ポートは8080を使用しています。
今のところ`go run main.go`
多分バイナリファイル作る。

## API

- URI:``localhost:8080/pref`` 

    |API Name|Method|Request Parameters|
    |:---|:---:|:---|
    |Get Recommendation|GET|userId|
    |User data Update|POST|userId,words|

    - Response
        - 正常系
        ```json
        {
            "message": "success",
            "facilities": [
            {
                "id": "<hotelID>",
                "similarity": "<一致度を0~100で？>"
            }
          ]
        }
        ```
        - 異常系 404 userがいないかuserのprefが登録されていないとき
        ```json
        {
            "message": "user not found"
        }
        ```

- URI:``localhost:8080/similarity`` 

    |API Name|Method|Parameters|
    |:---|:---:|:---|
    |Get Recommendation|GET|userId,facilityId|
    - Response
        - 正常系
        ```json
         {
        "message": "success",
        "user": "<userId>",
        "facility": "<facilityId>",
        "similarity": "<一致度を0~100で？>"
        }
        ```
        - 異常系 404 userがいないかuserのprefが登録されていないかホテルが登録されていなかったとき
        ```json
        {
        "message": "user(facility) not found"
        }
        ```
