module github.com/slava-y2k/gofig/example

go 1.25.0

require (
	github.com/slava-y2k/gofig v0.0.0-00010101000000-000000000000
	github.com/slava-y2k/gofig/provider/keyvalue v0.0.0-00010101000000-000000000000
)

replace github.com/slava-y2k/gofig => ../

replace github.com/slava-y2k/gofig/provider/keyvalue => ../provider/keyvalue
