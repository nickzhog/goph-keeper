protoc:
	@protoc internal/proto/gophkeeper.proto  --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative
certs:
	@openssl genrsa -out private.key 2048
	@openssl req -new -nodes -x509 -sha256 -subj "/CN=localhost" -addext "subjectAltName = DNS:localhost" -key private.key -out cert.crt -days 365
	@cp cert.crt cmd/client/.
	@cp cert.crt cmd/test_client/.
	@cp cert.crt cmd/server/.
clean:
	@find . -type f -name "*.log" -delete