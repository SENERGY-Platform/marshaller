module github.com/SENERGY-Platform/marshaller

go 1.18

require (
	github.com/SENERGY-Platform/converter v0.0.0-20220727105823-2d11887c48ce
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/clbanning/mxj v1.8.4
	github.com/coocood/freecache v1.1.0
	github.com/google/uuid v1.1.1
	github.com/julienschmidt/httprouter v1.3.0
)

require (
	github.com/Knetic/govaluate v3.0.0+incompatible // indirect
	github.com/RyanCarrier/dijkstra v1.1.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.0.3 // indirect
	gopkg.in/go-playground/colors.v1 v1.2.0 // indirect
)

//replace github.com/SENERGY-Platform/converter => ../converter-service
