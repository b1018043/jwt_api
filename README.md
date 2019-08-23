# jwt_api
## 説明
1. jwtを使った認証付きのapi
2. todoを記録したりしている
## 使い方
```
go build
./jwt_api.exe
```
入力値
```
// path="/login", method="post"
{
    "email":"email",
    "password":"password"
}
```
```
// path="/signup", method="post"
{
    "username":"username",
    "email":"email",
    "password":"password"
}
```
```
// path="/private", method="get"
```
