package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var opts *options

func NewCfg(cfgFile string) {
	buf, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		panic(err)
	}
	opts = &options{}
	err = yaml.Unmarshal(buf, opts)
	if err != nil {
		panic(err)
	}
}

type options struct {
	Host        string `yaml:"host"`
	DB          string `yaml:"db"`
	LinkTimeout int64  `yaml:"link_timeout"`
	MinLinkLen  int    `yaml:"min_link_len"`
	MaxLinkLen  int    `yaml:"max_link_len"`
}
