package server

import (
	"github.com/BurntSushi/toml"
)

/*
# Blog Title
Title = "latermoon's blog"

# Default Author
Author = "latermoon"

# Password (in MD5) is use for accessing the /private/ path
Password = "6d4db5ff0c117864a028"

# Run the http server on a given host and port
Server = "127.0.0.1:8801"
*/

type BlogConfig struct {
	Title    string
	Author   string
	Password string
	Server   string
}

func NewBlogConfig(filename string) (*BlogConfig, error) {
	cfg := &BlogConfig{}
	_, err := toml.DecodeFile(filename, cfg)
	return cfg, err
}
