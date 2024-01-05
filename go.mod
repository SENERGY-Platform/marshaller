module github.com/SENERGY-Platform/marshaller

go 1.21

toolchain go1.21.0

require (
	github.com/SENERGY-Platform/converter v0.0.4
	github.com/clbanning/mxj v1.8.4
	github.com/julienschmidt/httprouter v1.3.0
)

require (
	github.com/SENERGY-Platform/models/go v0.0.0-20230824080159-16585960df38
	github.com/SENERGY-Platform/service-commons v0.0.0-20231115074650-7021aeae60e4
)

require (
	github.com/IBM/sarama v1.42.1 // indirect
	github.com/Knetic/govaluate v3.0.0+incompatible // indirect
	github.com/RyanCarrier/dijkstra v1.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/eapache/go-resiliency v1.5.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20230731223053-c322873962e3 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.5.0 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/hashicorp/go-multierror v1.1.1 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.7.6 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.4 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/klauspost/compress v1.17.4 // indirect
	github.com/lucasb-eyer/go-colorful v1.2.0 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible // indirect
	github.com/pierrec/lz4/v4 v4.1.19 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/segmentio/kafka-go v0.4.47 // indirect
	golang.org/x/crypto v0.17.0 // indirect
	golang.org/x/net v0.19.0 // indirect
	gopkg.in/go-playground/colors.v1 v1.2.0 // indirect
)

//replace github.com/SENERGY-Platform/converter => ../converter-service
//replace github.com/SENERGY-Platform/models/go => ../models/go
