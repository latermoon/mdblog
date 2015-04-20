package blog

type BlogConfig struct {
	Title       string
	Author      string
	Password    string
	Salt        string
	Server      string
	SessionName string // default: sess
	AuthKey     string // default: auth
}
