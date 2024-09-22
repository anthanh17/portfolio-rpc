# rpc-portfolio

## REQUIREMENT 🙏

1. `protoc`

```
apt install -y protobuf-compiler # linux
brew install protobuf # macos
protoc --version
```

2. `go plugins`

```
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

3. `sqlc`

> sqlc generates golang code from SQL.

#### Here's how it works:

- You write queries in SQL.
- You run sqlc to generate code with type-safe interfaces to those queries.
- You write application code that calls the generated code.

```
brew install sqlc # macos
sudo snap install sqlc # ubuntu
```

#### Generates golang code from SQL

```
make sqlc
```

4. `golang-migrate`

> Database migrations using CLI

```
brew install golang-migrate
```

5. `make`
   > Run self-defined script code

```
brew install make # macos
```

## How to use

```
make server
```
