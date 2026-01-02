module blog

go 1.23

require (
	github.com/golang-jwt/jwt/v5 v5.3.0
	github.com/joho/godotenv v1.5.1
	github.com/mattn/go-sqlite3 v1.14.20
	golang.org/x/crypto v0.31.0
)

require github.com/aws/aws-lambda-go v1.51.1 // indirect

// Development tools (run: go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate)
// tools: github.com/sqlc-dev/sqlc/cmd/sqlc@latest
