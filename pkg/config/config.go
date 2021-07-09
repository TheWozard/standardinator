package config

import (
	"TheWozard/standardinator/pkg/path"
	"TheWozard/standardinator/pkg/path/jsonpath"
)

type OutputConfig struct {
	Name string `json:"Name"`
	Path string `json:"For"`
}

func (c OutputConfig) GetCongfigMatchPath() (path.Parsed, error) {
	return jsonpath.Parse(c.Path)
}
