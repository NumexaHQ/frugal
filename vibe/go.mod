module github.com/NumexaHQ/vibe

go 1.20

replace github.com/NumexaHQ/captainCache => ../auth

replace github.com/NumexaHQ/captainCache/numexa-common => ../numexa-common

require (
	github.com/ClickHouse/clickhouse-go/v2 v2.11.0
	github.com/NumexaHQ/captainCache v0.0.0-00010101000000-000000000000
	github.com/gofiber/fiber/v2 v2.48.0
	github.com/sirupsen/logrus v1.9.3
	gorm.io/gorm v1.25.2
)

require (
	github.com/ClickHouse/ch-go v0.53.0 // indirect
	github.com/NumexaHQ/captainCache/numexa-common v0.0.0-00010101000000-000000000000 // indirect
	github.com/go-faster/city v1.0.1 // indirect
	github.com/go-faster/errors v0.6.1 // indirect
	github.com/hashicorp/go-version v1.6.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/paulmach/orb v0.9.0 // indirect
	github.com/pierrec/lz4/v4 v4.1.17 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/segmentio/asm v1.2.0 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	go.opentelemetry.io/otel v1.16.0 // indirect
	go.opentelemetry.io/otel/trace v1.16.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

require (
	github.com/andybalholm/brotli v1.0.5 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/klauspost/compress v1.16.3 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/mattn/go-runewidth v0.0.14 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.48.0
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	gorm.io/driver/clickhouse v0.5.1
)
