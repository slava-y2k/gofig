module github.com/slava-y2k/gonfig/example

go 1.25.0

require (
	github.com/slava-y2k/gonfig v0.0.0-00010101000000-000000000000
	github.com/slava-y2k/gonfig/provider/keyvalue v0.0.0-00010101000000-000000000000
)

replace github.com/slava-y2k/gonfig => ../

replace github.com/slava-y2k/gonfig/provider/keyvalue => ../provider/keyvalue
