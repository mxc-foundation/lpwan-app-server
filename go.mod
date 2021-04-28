module github.com/mxc-foundation/lpwan-app-server

go 1.15

require (
	cloud.google.com/go v0.56.0 // indirect
	github.com/DATA-DOG/go-sqlmock v1.4.1
	github.com/NickBall/go-aes-key-wrap v0.0.0-20170929221519-1c3aa3e4dfc5
	github.com/apex/log v1.1.0
	github.com/aws/aws-sdk-go v1.26.3
	github.com/brocaar/chirpstack-api/go/v3 v3.7.7
	github.com/brocaar/lorawan v0.0.0-20200726141338-ee070f85d494
	github.com/cespare/xxhash/v2 v2.1.1 // indirect
	github.com/coreos/go-oidc v2.2.1+incompatible
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/eclipse/paho.mqtt.golang v1.2.0
	github.com/elazarl/go-bindata-assetfs v1.0.1
	github.com/go-redis/redis/v7 v7.4.0
	github.com/gofrs/uuid v3.2.0+incompatible
	github.com/gogo/protobuf v1.3.1 // indirect
	github.com/golang/protobuf v1.5.2
	github.com/gopherjs/gopherjs v0.0.0-20190430165422-3e4dfb77656c // indirect
	github.com/gorilla/mux v1.7.3
	github.com/gorilla/websocket v1.4.1 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.1.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/jacobsa/crypto v0.0.0-20190317225127-9f44e2d11115 // indirect
	github.com/jmoiron/sqlx v1.2.0
	github.com/jung-kurt/gofpdf v1.16.2
	github.com/lestrrat-go/jwx v1.0.3
	github.com/lib/pq v1.2.0
	github.com/mmcloughlin/geohash v0.9.0
	github.com/pelletier/go-toml v1.6.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/pquerna/cachecontrol v0.0.0-20180517163645-1555304b9b35 // indirect
	github.com/pquerna/otp v1.2.0
	github.com/prometheus/client_golang v1.2.1
	github.com/prometheus/client_model v0.0.0-20191202183732-d1d2010b5bee // indirect
	github.com/prometheus/procfs v0.0.8 // indirect
	github.com/robertkrimen/otto v0.0.0-20191217063420-37f8e9a2460c
	github.com/robfig/cron v1.2.0
	github.com/rubenv/sql-migrate v0.0.0-20191213152630-06338513c237
	github.com/segmentio/kafka-go v0.3.6
	github.com/shopspring/decimal v1.2.0
	github.com/sirupsen/logrus v1.4.2
	github.com/smartystreets/assertions v1.0.0 // indirect
	github.com/spf13/afero v1.2.2 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/jwalterweatherman v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spf13/viper v1.6.2
	github.com/streadway/amqp v0.0.0-20190827072141-edfb9018d271
	github.com/stretchr/testify v1.5.1
	github.com/tmc/grpc-websocket-proxy v0.0.0-20200122045848-3419fae592fc
	golang.org/x/crypto v0.0.0-20200622213623-75b288015ac9
	golang.org/x/net v0.0.0-20210316092652-d523dce5a7f4
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	google.golang.org/genproto v0.0.0-20210426193834-eac7f76ac494
	google.golang.org/grpc v1.36.1
	google.golang.org/protobuf v1.26.0
	gopkg.in/sourcemap.v1 v1.0.5 // indirect
	gopkg.in/square/go-jose.v2 v2.5.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

// remove when https://github.com/tmc/grpc-websocket-proxy/pull/23 has been merged
// and grpc-websocket-proxy dependency has been updated to version including this fix.
replace github.com/tmc/grpc-websocket-proxy => github.com/brocaar/grpc-websocket-proxy v1.0.1

// replace github.com/brocaar/chirpstack-api/go/v3 => ../chirpstack-api/go
