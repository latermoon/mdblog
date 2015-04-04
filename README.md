# Markdown Blog Generator

### How to use
```
git clone github.com/latermoon/mdblog
cd mdblog

export GOPATH=$GOPATH:($pwd)
go run src/mdblog.go website_dir

append to nginx.conf

server {
	listen    80;
	server_name    your.blog.name;

	location / {
		proxy_pass http://localhost:8801;
	}
}

```

### Imported
```
go get github.com/russross/blackfriday
go get github.com/nfnt/resize
go get github.com/PuerkitoBio/goquery
go get github.com/howeyc/fsnotify
go get github.com/BurntSushi/toml
go get github.com/golang/groupcache/lru
go get github.com/go-martini/martini
go get github.com/martini-contrib/sessions
```

### Website
```
blog.latermoon.me
	article
		private
		my-first-markdown-blog.md
		about-me.md
	static
	template
```
