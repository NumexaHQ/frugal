module github.com/NumexaHQ/monger

go 1.20

replace github.com/NumexaHQ/captainCache => ../auth

replace github.com/NumexaHQ/captainCache/numexa-common => ../numexa-common

require (
	github.com/ClickHouse/clickhouse-go/v2 v2.11.0
	github.com/NumexaHQ/captainCache v0.0.0-00010101000000-000000000000
	github.com/NumexaHQ/captainCache/numexa-common v0.0.0-00010101000000-000000000000
	github.com/pkoukk/tiktoken-go v0.1.5
	github.com/sirupsen/logrus v1.9.3
	gorm.io/gorm v1.25.2
)

require (
	github.com/ClickHouse/ch-go v0.53.0 // indirect
	github.com/dlclark/regexp2 v1.10.0 // indirect
	github.com/go-faster/city v1.0.1 // indirect
	github.com/go-faster/errors v0.6.1 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.2 // indirect
	github.com/hashicorp/go-retryablehttp v0.7.4 // indirect
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
	github.com/andybalholm/brotli v1.0.5
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.7.3
	github.com/klauspost/compress v1.16.3 // indirect
	golang.org/x/sys v0.10.0 // indirect
	gorm.io/driver/clickhouse v0.5.1
)
