module github.com/NumexaHQ/captainCache

go 1.20

require (
	github.com/NumexaHQ/captainCache/numexa-common v0.0.0-00010101000000-000000000000
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/go-redis/redis/v8 v8.11.5
	github.com/gofiber/fiber/v2 v2.47.0
	github.com/lib/pq v1.10.9
	github.com/sashabaranov/go-openai v1.14.1
	golang.org/x/net v0.12.0
)

require (
	github.com/dustinkirkland/golang-petname v0.0.0-20230626224747-e794b9370d49
	golang.org/x/crypto v0.11.0
)

require github.com/stretchr/testify v1.7.1 // indirect

require (
	github.com/andybalholm/brotli v1.0.5 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/klauspost/compress v1.16.3 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/philhofer/fwd v1.1.2 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/savsgio/dictpool v0.0.0-20221023140959-7bf2e61cea94 // indirect
	github.com/savsgio/gotils v0.0.0-20230208104028-c358bd845dee // indirect
	github.com/sirupsen/logrus v1.9.3
	github.com/tinylib/msgp v1.1.8 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.47.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
)

replace github.com/NumexaHQ/captainCache/numexa-common => ../numexa-common
