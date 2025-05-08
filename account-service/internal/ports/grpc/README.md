### Regenerating GRPC code 

```bash
protoc --go_out=. --go_opt=paths=source_relative     --go-grpc_out=. --go-grpc_opt=paths=source_relative  account.proto
```
