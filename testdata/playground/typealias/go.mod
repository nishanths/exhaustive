module goplay

go 1.16

replace (
	bar => gorm.io/gorm v1.21.0
	foo => gorm.io/gorm v1.20.12
)

require (
	github.com/gorilla/websocket v1.4.2
	golang.org/x/time v0.0.0-20210611083556-38a9dc6acbc6
)
