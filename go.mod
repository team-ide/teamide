module teamide

go 1.18

//replace github.com/team-ide/go-dialect => ../go-dialect
//
//replace github.com/team-ide/go-tool => ../go-tool

replace github.com/dop251/goja => github.com/team-ide/goja v1.0.4

//
//replace github.com/team-ide/goja v1.0.3 => ../goja

//replace github.com/dop251/goja => ../goja

replace go.uber.org/zap v1.27.0 => github.com/team-ide/zap v1.27.0

require (
	github.com/PuerkitoBio/goquery v1.8.1
	github.com/apache/thrift v0.17.0
	github.com/creack/pty v1.1.21
	github.com/dop251/goja v0.0.0-20240516125602-ccbae20bcec2
	github.com/gin-gonic/gin v1.10.0
	github.com/google/uuid v1.6.0
	github.com/gorilla/websocket v1.5.1
	github.com/jacobsa/go-serial v0.0.0-20180131005756-15cf729a72d4
	github.com/mssola/user_agent v0.6.0
	github.com/pkg/sftp v1.13.6
	github.com/shirou/gopsutil/v3 v3.23.12
	github.com/tealeg/xlsx/v3 v3.3.12
	github.com/team-ide/cron v1.0.1
	github.com/team-ide/go-dialect v1.9.30
	github.com/team-ide/go-tool v1.2.48
	go.mongodb.org/mongo-driver v1.17.3
	go.uber.org/zap v1.27.0
	golang.org/x/crypto v0.28.0
	golang.org/x/net v0.25.0
	gopkg.in/natefinch/lumberjack.v2 v2.2.1
	gopkg.in/yaml.v3 v3.0.1
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	gitee.com/opengauss/openGauss-connector-go-pq v1.0.7 // indirect
	github.com/Shopify/sarama v1.38.1 // indirect
	github.com/andybalholm/cascadia v1.3.1 // indirect
	github.com/bytedance/sonic v1.11.6 // indirect
	github.com/bytedance/sonic/loader v0.1.1 // indirect
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/cloudwego/base64x v0.1.4 // indirect
	github.com/cloudwego/iasm v0.2.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/dlclark/regexp2 v1.11.4 // indirect
	github.com/eapache/go-resiliency v1.3.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20230111030713-bf00bc1b83b6 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/frankban/quicktest v1.14.6 // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/gin-contrib/sse v0.1.0 // indirect
	github.com/go-logfmt/logfmt v0.5.1 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.20.0 // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
	github.com/go-sourcemap/sourcemap v2.1.3+incompatible // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/go-zookeeper/zk v1.0.4 // indirect
	github.com/goccy/go-json v0.10.2 // indirect
	github.com/godror/godror v0.37.0 // indirect
	github.com/godror/knownpb v0.1.0 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/btree v1.0.0 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/pprof v0.0.0-20230207041349-798e818bf904 // indirect
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.7.6 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.3 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.16.7 // indirect
	github.com/klauspost/cpuid/v2 v2.2.7 // indirect
	github.com/kr/fs v0.1.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/kr/text v0.2.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/lib/pq v1.10.9 // indirect
	github.com/lufia/plan9stats v0.0.0-20211012122336-39d0f177ccd0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-sqlite3 v1.14.16 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/olivere/elastic/v7 v7.0.32 // indirect
	github.com/pelletier/go-toml/v2 v2.2.2 // indirect
	github.com/peterbourgon/diskv/v3 v3.0.1 // indirect
	github.com/pierrec/lz4/v4 v4.1.17 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/power-devops/perfstat v0.0.0-20210106213030-5aafc221ea8c // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/rogpeppe/fastuuid v1.2.0 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/shabbyrobe/xmlwriter v0.0.0-20200208144257-9fca06d00ffa // indirect
	github.com/shoenig/go-m1cpu v0.1.6 // indirect
	github.com/shopspring/decimal v1.4.0 // indirect
	github.com/team-ide/go-driver v1.3.6 // indirect
	github.com/team-ide/go-interpreter v0.1.2 // indirect
	github.com/tjfoc/gmsm v1.4.1 // indirect
	github.com/tklauser/go-sysconf v0.3.12 // indirect
	github.com/tklauser/numcpus v0.6.1 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.2.12 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	github.com/yusufpapurcu/wmi v1.2.3 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	golang.org/x/arch v0.8.0 // indirect
	golang.org/x/sync v0.8.0 // indirect
	golang.org/x/sys v0.26.0 // indirect
	golang.org/x/text v0.19.0 // indirect
	golang.org/x/xerrors v0.0.0-20220907171357-04be3eba64a2 // indirect
	google.golang.org/protobuf v1.34.1 // indirect
)
