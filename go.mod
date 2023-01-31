module github.com/SENERGY-Platform/marshaller

go 1.19

require (
	github.com/SENERGY-Platform/converter v0.0.0-20230131075057-f47029811afa
	github.com/clbanning/mxj v1.8.4
	github.com/julienschmidt/httprouter v1.3.0
)

require (
	github.com/SENERGY-Platform/models/go v0.0.0-20230105115534-8edcf0271764
	github.com/patrickmn/go-cache v2.1.0+incompatible
)

require (
	github.com/Knetic/govaluate v3.0.0+incompatible // indirect
	github.com/RyanCarrier/dijkstra v1.1.0 // indirect
	github.com/lucasb-eyer/go-colorful v1.0.3 // indirect
	gopkg.in/go-playground/colors.v1 v1.2.0 // indirect
)

//replace github.com/SENERGY-Platform/converter => ../converter-service
