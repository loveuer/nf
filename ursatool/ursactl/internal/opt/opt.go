package opt

type _new struct {
	Template          string
	DisableInitScript bool
}

type config struct {
	Debug         bool
	DisableUpdate bool
	Version       bool
	New           _new
}

var Cfg = &config{}

var TemplateMap = map[string]string{
	"ultone": "https://gitea.loveuer.com/loveuer/ultone.git",
}

const README = "# {{ .project_name }}\n\n### 启动\n- `go run . --help`\n- `go run .`\n\n### 构建\n- `go build -ldflags '-s -w' -o dist/{{ .project_name}}_app .`\n- `docker build -t <image> -f Dockerfile .`"
