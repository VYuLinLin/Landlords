# Landlords
多人在线斗地主游戏持续开发中... (2024-04-22)

## front-end cocos creator 2.4.13

## Back-end golang 1.22

[斗地主游戏流程图](https://github.com/VYuLinLin/Landlords/blob/master/server/%E6%96%97%E5%9C%B0%E4%B8%BB%E6%B8%B8%E6%88%8F%E6%B5%81%E7%A8%8B%E5%9B%BE.png)

**Dependencies(See go.mod file for details)**

* Go 1.22
* github.com/astaxie/beego v1.12.3
* github.com/go-sql-driver/mysql v1.5.0
* github.com/googollee/go-socket.io v1.4.4
* github.com/gorilla/websocket v1.4.1

Quick Start

    git clone https://github.com/VYuLinLin/Landlords.git

    cd Landlords/server/src/landlords

    go run main.go
    or
    go build main.go
    main.exe

    Now visit http://localhost or http://127.0.0.1:80
