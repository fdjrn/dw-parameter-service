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