### gmapproxy

This is a simple GRPC service to Find Nearby Search from Google Place API.

The search result will automatically stored in MYSQL for 720 Hours / 3 Days for reduce api call limit and optimize multiple types search.

Any data more than 720 Hours will be Purged due to Cache Policy in Google Place Term Of Services.

### Install

First you need to install ProtocolBuffers 3.0.0-beta-3 or later.
See install guide here: https://github.com/google/protobuf

Then, go get -u as usual the following packages:

```
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway
go get -u github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
go get -u github.com/golang/protobuf/protoc-gen-go
```

After grpc library and platform installed, you need to install app dependency.

```
go get -u -v .
```

To Start GRPC Service run:

```sh
go run grpc-server/main.go
```

To Start JSON GRPC Gateway run:

```sh
go run main.go
```

### Development

Dont edit place/place_service.pb.\* directly
its generated by running `sh makeproto` command