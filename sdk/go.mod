module github.com/jxo-me/plus-core/sdk/v2

go 1.20

replace (
	github.com/jxo-me/plus-core/core/v2 => ../core
	github.com/jxo-me/plus-core/pkg/v2 => ../pkg
	github.com/jxo-me/rabbitmq-go => ../../rabbitmq-go
)

require (
	github.com/apache/rocketmq-client-go/v2 v2.1.2
	github.com/casbin/casbin/v2 v2.79.0
	github.com/go-redsync/redsync/v4 v4.11.0
	github.com/gogf/gf-jwt/v2 v2.1.0
	github.com/gogf/gf/v2 v2.5.7
	github.com/google/uuid v1.4.0
	github.com/json-iterator/go v1.1.12
	github.com/jxo-me/gf-metrics v0.1.4
	github.com/jxo-me/gfbot v0.1.16
	github.com/jxo-me/plus-core/core/v2 v2.0.29
	github.com/jxo-me/plus-core/pkg/v2 v2.0.35
	github.com/jxo-me/rabbitmq-go v1.0.13
	github.com/nsqio/go-nsq v1.1.0
	github.com/redis/go-redis/v9 v9.3.0
	github.com/zegl/goriak/v3 v3.2.4
	golang.org/x/text v0.14.0
)

require (
	github.com/BurntSushi/toml v1.2.0 // indirect
	github.com/basho/backoff v0.0.0-20150307023525-2ff7c4694083 // indirect
	github.com/basho/riak-go-client v0.0.0-20170327205844-5587c16e0b8b // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/bits-and-blooms/bitset v1.8.0 // indirect
	github.com/casbin/govaluate v1.1.0 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/clbanning/mxj/v2 v2.7.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/emirpasic/gods v1.12.0 // indirect
	github.com/ergo-services/ergo v1.999.224 // indirect
	github.com/fatih/color v1.15.0 // indirect
	github.com/fsnotify/fsnotify v1.7.0 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang-jwt/jwt/v4 v4.3.0 // indirect
	github.com/golang/mock v1.4.4 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/golang/snappy v0.0.1 // indirect
	github.com/gorilla/websocket v1.5.1 // indirect
	github.com/grokify/html-strip-tags-go v0.0.1 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/konsorten/go-windows-terminal-sequences v1.0.1 // indirect
	github.com/magiconair/properties v1.8.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.4 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/prometheus/client_golang v1.16.0 // indirect
	github.com/prometheus/client_model v0.3.0 // indirect
	github.com/prometheus/common v0.42.0 // indirect
	github.com/prometheus/procfs v0.10.1 // indirect
	github.com/rabbitmq/amqp091-go v1.9.0 // indirect
	github.com/richardlehane/mscfb v1.0.4 // indirect
	github.com/richardlehane/msoleps v1.0.3 // indirect
	github.com/rivo/uniseg v0.4.4 // indirect
	github.com/sirupsen/logrus v1.4.0 // indirect
	github.com/tidwall/gjson v1.13.0 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	github.com/xuri/efp v0.0.0-20230802181842-ad255f2331ca // indirect
	github.com/xuri/excelize/v2 v2.8.0 // indirect
	github.com/xuri/nfp v0.0.0-20230819163627-dc951e3ffe1a // indirect
	go.opentelemetry.io/otel v1.14.0 // indirect
	go.opentelemetry.io/otel/sdk v1.14.0 // indirect
	go.opentelemetry.io/otel/trace v1.14.0 // indirect
	go.uber.org/atomic v1.5.1 // indirect
	golang.org/x/crypto v0.15.0 // indirect
	golang.org/x/lint v0.0.0-20190930215403-16217165b5de // indirect
	golang.org/x/net v0.18.0 // indirect
	golang.org/x/sys v0.14.0 // indirect
	golang.org/x/term v0.14.0 // indirect
	golang.org/x/time v0.4.0 // indirect
	golang.org/x/tools v0.6.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	stathat.com/c/consistent v1.0.0 // indirect
)
