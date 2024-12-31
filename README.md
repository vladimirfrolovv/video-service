# Video-service

## Upload video to s3(minio)

### Start minio
docker-compose up -d
### Start Service
go run cmd/service/main.go


### Example request
url -X POST -F "file=@/path.mp4" http://localhost:8080/upload

