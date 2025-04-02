  
# 簡介（Introduction）
基於go與grpc的簡單user微服務

- 使用gorm連接mysql
- 使用jwt進行token認證
- 使用grpc進行微服務間的通訊
- 使用docker進行容器化部署
 
# install
 
## 先決條件
- 可連接的mysql資料庫
- config.yaml or add env

## 透過docker本地啟動
``` shell
docker run -it -p 50051:50051 -v your_config.yaml:/app/config.yaml itmrchow/todolist-user
```

## 環境變數
| 變數名稱              | 說明            | 預設值                                    |
| --------------------- | --------------- | ----------------------------------------- |
| APP_SERVER_NAME       | 服務名稱        | todolist-user                             |
| APP_SERVER_PORT       | 服務埠          | 50051                                     |
| APP_MYSQL_URL_SUFFIX  | mysql連接字串   | ?charset=utf8mb4&parseTime=True&loc=Local |
| APP_MYSQL_DB_ACCOUNT  | mysql帳號       |                                           |
| APP_MYSQL_DB_PASSWORD | mysql密碼       |                                           |
| APP_MYSQL_DB_HOST     | mysql主機       |                                           |
| APP_MYSQL_DB_PORT     | mysql埠         |                                           |
| APP_MYSQL_DB_NAME     | mysql資料庫     |                                           |
| APP_JWT_SECRET_KEY    | jwt密鑰         |                                           |
| APP_JWT_EXPIRE_AT     | jwt過期時間(hr) | 8                                         |


# 架構設計（Architecture Design）
## microservice
為什麼使用microservice架構？
- 針對user domain大量請求 , 可以對user service 來做水平擴容
- 如果未來有其他複雜的需求 , domain可以獨立成一個服務 , 如果是單體repo會越來越複雜
## grpc
- 高效能, 在同一connection同時發送多個請求
- 語言中立, 支援多種程式語言
- protocal, 定義接口提供強行別的資料結構, 減少運行錯誤

## 測試規劃


## features

## 部署與scaling

### On-premise deployment with scale plan.
### Cloud deployment with scale plan.

