package cfg

import (
	"flag"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"aid.im/app/util/cache"
	"aid.im/app/util/logx"
)

var (
	Opt     *options
	SrvAddr string
	DBPath  string
)

type options struct {
	Mode       string         `yaml:"mode"`
	SrvAddr    string         `yaml:"srv_addr"`
	IndexPage  string         `yaml:"index_page"`
	LinkTTl    int64          `yaml:"link_ttl"`
	MinLinkLen int            `yaml:"min_link_len"`
	MaxLinkLen int            `yaml:"max_link_len"`
	Log        logx.LogConfig `yaml:"log"`
}

func InitCfg() {
	cfgPath := flag.String("c", "etc/cfg.yml", "config file")
	dbPath := flag.String("db", "url.db", "db file")
	flag.Parse()
	fmt.Println(*cfgPath, *dbPath)
	if *dbPath == "" || *cfgPath == "" {
		panic("cfg or db path cannot be empty!")
	}
	buf, err := ioutil.ReadFile(*cfgPath)
	if err != nil {
		panic(err)
	}
	Opt = &options{}
	err = yaml.Unmarshal(buf, Opt)
	fmt.Printf("%#v, %#v", Opt, err)
	if err != nil {
		panic(err)
	}
	DBPath = *dbPath
	SrvAddr = Opt.SrvAddr
}

func InitCache() {
	if err := cache.InitCache("memory", `{"interval":3600}`); err != nil {
		panic(err)
	}
}

func InitLog() {
	logx.InitDefaultLogger(&Opt.Log)
}
