module github.com/romanblanco/graffiti-ipfs

go 1.16

require (
	github.com/google/open-location-code/go v0.0.0-20201229230907-d47d9f9b95e9
	github.com/ipfs/go-ipfs-api v0.2.0 // indirect
	github.com/op/go-logging v0.0.0-20160315200505-970db520ece7
	github.com/romanblanco/go-ipfs-api v0.0.4-0.20210429201352-5d046310ffbd
	github.com/rwcarlsen/goexif v0.0.0-20190401172101-9e8deecbddbd
)

replace github.com/ipfs/go-ipfs-api => github.com/romanblanco/go-ipfs-api v0.0.4-0.20210416203454-8185da3731c7
