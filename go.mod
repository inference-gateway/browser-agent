module github.com/inference-gateway/browser-agent

go 1.26.4

require (
	github.com/inference-gateway/adk v0.23.0
	github.com/sethvargo/go-envconfig v1.3.1
	github.com/spf13/cobra v1.10.2
	go.uber.org/zap v1.28.0
	gopkg.in/yaml.v3 v3.0.1
	github.com/maxbrunsfeld/counterfeiter/v6 v6.12.2 // indirect
)

tool (
	github.com/maxbrunsfeld/counterfeiter/v6
)
