module github.com/SENERGY-Platform/marshaller

go 1.17

require (
	github.com/SENERGY-Platform/converter v0.0.0-20220302112210-0325026a0e9f
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/clbanning/mxj v1.8.4
	github.com/coocood/freecache v1.1.0
	github.com/google/uuid v1.1.1
	github.com/julienschmidt/httprouter v1.3.0
)

require (
	github.com/RyanCarrier/dijkstra v1.1.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.0.3 // indirect
	gopkg.in/go-playground/colors.v1 v1.2.0 // indirect
)

//replace github.com/SENERGY-Platform/converter => ../converter-service
