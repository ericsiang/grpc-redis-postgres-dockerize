# Golang  gRPC project

### 使用說明
* 啟動服務
    直接使用 docker compose 指令
    ```
    docker-compose up -d
    ```

    啟動後會產生三個 docker container 服務 
    * go gRPC server
    * redis
    * postgreSQL
* 連線到 postgreSQL，手動建立 users 資料表
    ```
    CREATE TABLE users (
        id 		SERIAL  PRIMARY KEY ,
        name     varchar(80)    NOT NULL,
        email    varchar(80)      NOT NULL
    );
    ```
    
* 使用 grpcurl cli 測試 gRPC server 
  
    * 查看 gRPC server 服務清單
        ```
        grpcurl  -plaintext 127.0.0.1:50051  list
        ```
        command response:
        ```
            grpc.reflection.v1.ServerReflection
            grpc.reflection.v1alpha.ServerReflection
            proto.UserService
        ```
    * 查看 proto.UserService 內的方法清單
        ```
        grpcurl  -plaintext 127.0.0.1:50051  list proto.UserService
        ```
        command response:
        ```
        proto.UserService.CreateUser
        proto.UserService.GetUser
        ```
    * 檢視方法定義
        ```
        grpcurl -plaintext 127.0.0.1:50051 describe proto.UserService.GetUser
        ```
        command response:
        ```
        proto.UserService.GetUser is a method:
        rpc GetUser ( .proto.GetUserRequest ) returns ( .proto.User );
        ```
    * 檢視請求參數
        ```
        grpcurl -plaintext 127.0.0.1:50051 describe proto.GetUserRequest
        ```
        command response:
        ```
        proto.GetUserRequest is a message:
        message GetUserRequest {
        int64 id = 1;
        }
        ```
    * 請求服務 GetUser 方法
        ```
        grpcurl -d '{"id":1}' -plaintext 127.0.0.1:50051 proto.UserService.GetUser
        ```
    * 請求服務 CreateUser 方法，一樣去查詢該方法所需參數


### 使用到的 package
<table>
    <th>package</th>
    <th>說明</th>
    <th>操作說明</th>
    <tr>
        <td><a href="https://github.com/spf13/viper" target="_blank">zap</a></td>
        <td>Zap 是一個快速、結構化、級別化的日誌庫，由 Uber 開發</td>
        <td>-</td>
    </tr>
    <tr>
        <td><a href="https://github.com/joho/godotenv" target="_blank">godotenv</a></td>
        <td>用Go 語言讀取專案內.env 環境變數</td>
        <td>-</td>
    </tr>
    <tr>
        <td><a href="https://github.com/lib/pq" target="_blank">pq
</a></td>
        <td>pq 是連接 PostgreSQL 的庫</td>
        <td> - </td>
    </tr>
    <tr>
        <td><a href="https://github.com/redis/go-redis/v9" target="_blank">go-redis</a></td>
        <td>go-redis 是 Redis 客户端库</td>
        <td> - </td>
    </tr>
    <tr>
        <td><a href="https://github.com/grpc/grpc-go" target="_blank">gRPC-Go</a></td>
        <td>建立 gRPC server</td>
        <td> - </td>
    </tr>
    <tr>
        <td><a href="https://github.com/grpc/grpc-go" target="_blank">google.golang.org/grpc/reflection</a></td>
        <td>使用 grpcurl cli 測試 gRPC server </td>
        <td> - </td>
    </tr>
    <tr>
        <td><a href="https://github.com/grpc-ecosystem/go-grpc-middleware" target="_blank">go-grpc-middleware</a></td>
        <td>Go gRPC 用的 middleware </td>
        <td> - </td>
    </tr>
    
</table>