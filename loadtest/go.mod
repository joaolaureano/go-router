// The load suite is its own module on purpose: vegeta and its transitive
// dependencies are a test harness, not something a consumer of go-router
// should inherit by importing the router.
module github.com/joaolaureano/go-router/loadtest

go 1.26

require (
	github.com/joaolaureano/go-router v0.0.0
	github.com/tsenart/vegeta/v12 v12.13.0
)

require (
	github.com/influxdata/tdigest v0.0.1 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/rs/dnscache v0.0.0-20230804202142-fc85eb664529 // indirect
	golang.org/x/net v0.27.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/text v0.16.0 // indirect
)

replace github.com/joaolaureano/go-router => ..
