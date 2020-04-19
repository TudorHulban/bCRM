# bCRM
backend CRM

### golangci-lint install:
a. Download
```
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ./
```
b. Run

### Docker container for PostgreSQL
```
docker run --name P1 -d -p 5432:5432 -e POSTGRES_PASSWORD=pp postgres:alpine
```