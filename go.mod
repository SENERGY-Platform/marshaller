module github.com/SENERGY-Platform/marshaller

go 1.16

require (
	github.com/SENERGY-Platform/converter v0.0.0-20200114152115-63d2a698b3c4
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/clbanning/mxj v1.8.4
	github.com/coocood/freecache v1.1.0
	github.com/google/uuid v1.1.1
	github.com/julienschmidt/httprouter v1.3.0
)

//replace github.com/SENERGY-Platform/converter => ../converter-service
