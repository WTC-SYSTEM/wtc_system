module github.com/WTC-SYSTEM/wtc_system/api_gateway

go 1.18

require (
	github.com/WTC-SYSTEM/wtc_system/libs/apperror v0.0.0-20220715211119-2850083e0199
	github.com/WTC-SYSTEM/wtc_system/libs/logging v0.0.0-20220702174402-8bce3a3a771f
	github.com/WTC-SYSTEM/wtc_system/libs/utils v0.0.0-20220715213612-340a1653668b
	github.com/cristalhq/jwt/v3 v3.1.0
	github.com/fatih/structs v1.1.0
	github.com/google/uuid v1.3.0
	github.com/gorilla/mux v1.8.0
	github.com/ilyakaznacheev/cleanenv v1.2.6
)

require (
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-redis/redis/v9 v9.0.0-beta.1
)

require (
	github.com/BurntSushi/toml v1.1.0 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/joho/godotenv v1.4.0 // indirect
	github.com/kr/pretty v0.3.0 // indirect
	github.com/rogpeppe/go-internal v1.8.0 // indirect
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
	golang.org/x/sys v0.0.0-20220412211240-33da011f77ad // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b // indirect
	olympos.io/encoding/edn v0.0.0-20201019073823-d3554ca0b0a3 // indirect
)

//replace github.com/WTC-SYSTEM/wtc_system/libs/apperror => ../../libs/apperror
