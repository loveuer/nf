package opt

type _new struct {
	Template          string
	DisableInitScript bool
}

type config struct {
	Debug         bool
	DisableUpdate bool
	New           _new
}

var Cfg = &config{}

var TemplateMap = map[string]string{
	"ultone": "https://gitea.loveuer.com/loveuer/ultone.git",
}
