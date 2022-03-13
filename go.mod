module github.com/jmisasa/swlbuddy

go 1.17

require (
	github.com/jmisasa/swlbuddy/rigctl v0.0.0
	github.com/mattn/go-gtk v0.0.0-20191030024613-af2e013261f5
	github.com/reiver/go-telnet v0.0.0-20180421082511-9ff0b2ab096e
)

require (
	github.com/mattn/go-pointer v0.0.1 // indirect
	github.com/reiver/go-oi v1.0.0 // indirect
	github.com/stretchr/testify v1.7.0 // indirect
)

replace github.com/jmisasa/swlbuddy/rigctl => ./rigctl
