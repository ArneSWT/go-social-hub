# Go Social Hub

Go Social Hub is a simple application for managing contacts using Go and SQLite. The application provides operations to create and read contacts.

Important: because go-sqlite3 is a CGO enabled package, you are required to set the environment variable CGO_ENABLED=1 and have a gcc compiler present within your path.

OR

use the Dockerfile:

```bash
docker build -t go-social-hub .

docker run -p 8080:8080 go-social-hub
```
