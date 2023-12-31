# Parameter Service
---
### Deskripsi
Service ini berfungsi untuk meng-handle manajemen dan/atau konfigurasi parameter. Seperti:
    - Manajemen Voucher
    
### Service Type
    - RESTFul API

### RestAPI Endpoint
    Fetch All Voucher Data:
    - POST  | /api/v1/params/voucher/all (paginated)

    Fetch Voucher Detail by ID:
    - GET   | /api/v1/params/voucher/:id

    Fetch Voucher Detail by Code:
    - POST   | /api/v1/params/voucher/code

    Create new Voucher:
    - POST  | /api/v1/params/voucher

    Delete Voucher by ID (soft deletion):
    - DEL   | /api/v1/params/voucher/:id

    Update Voucher by ID:
    - PUT   | /api/v1/params/voucher/:id

### Build Docker Image
    docker build -t dw-voucher:1.0.0 -f Dockerfile .
    
### Available Environment Value:
    - DATABASE_MONGODB_URI : conncetion uri to mongodb cluster
        
        example: mongodb+srv://<user>:<password>@<cluster-host>/?retryWrites=true&w=majority

    - DATABASE_MONGODB_DB_NAME : Database Name used for parameter service

        example: dw-mdw-parameter

### Docker Run Command
    docker run -d -p 9000:9000 --name dw-voucher-service --env "DATABASE_MONGODB_DB_NAME=dev-mdw-parameter" --restart unless-stopped dw-voucher:1.0.0