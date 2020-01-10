module github.com/SENERGY-Platform/marshaller

go 1.13

require (
	github.com/SENERGY-Platform/converter v0.0.0-20200110094538-195b91237d88
	github.com/SENERGY-Platform/device-manager v0.0.0-20191216115423-e51ad8574d21 // indirect
	github.com/SmartEnergyPlatform/jwt-http-router v0.0.0-20190722084820-0e1fe0dc7a07
	github.com/cespare/xxhash v1.1.0 // indirect
	github.com/clbanning/mxj v1.8.4
	github.com/coocood/freecache v1.1.0
	github.com/google/uuid v1.1.1
	github.com/lucasb-eyer/go-colorful v1.0.3
	gopkg.in/go-playground/colors.v1 v1.2.0
)

//replace github.com/SENERGY-Platform/converter => ../converter-service
