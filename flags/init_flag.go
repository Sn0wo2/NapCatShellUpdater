package flags

import (
	"flag"
	"github.com/Sn0wo2/NapCatShellUpdater/log"
	"path/filepath"
	"time"
)

var Config struct {
	Path           string
	Proxy          string
	Debug          bool
	NapCatPanelURL string
	NapCatToken    string
	Login          bool
	Version        string
	Sleep          time.Duration
}

func InitFlag() bool {
	path := "./"
	flag.StringVar(&path, "path", "./", "NapCat path")
	var err error
	Config.Path, err = filepath.Abs(path)
	if err != nil {
		Config.Path = path
		log.RPanic(err)
	}
	flag.StringVar(&Config.Proxy, "proxy", "", "HTTP Proxy")
	flag.BoolVar(&Config.Debug, "debug", true, "Enable debug logging")
	flag.StringVar(&Config.NapCatPanelURL, "ncpanel", "http://127.0.0.1:6099", "NapCat Panel URL")
	flag.StringVar(&Config.NapCatToken, "nctoken", "token", "NapCat Token")
	flag.BoolVar(&Config.Login, "login", true, "Login to NapCat Panel")
	flag.StringVar(&Config.Version, "version", "", "Update NapCat Version")
	flag.DurationVar(&Config.Sleep, "sleep", 30*time.Second, "Sleep time(Wait NapCat load)")
	flag.Parse()
	return true
}
