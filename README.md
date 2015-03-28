# Markdown Blog Generator

### How to use
```
git clone github.com/latermoon/mdblog
cd mdblog
export GOPATH=$GOPATH:($pwd)
go build src/mdblog.go

./mdblog -d website_dir -p localhost:8801

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
go get github.com/PuerkitoBio/goquery
go get github.com/russross/blackfriday
go get github.com/nfnt/resize
go get github.com/howeyc/fsnotify
go get github.com/go-martini/martini
go get github.com/martini-contrib/sessions
go get github.com/BurntSushi/toml
```

### URL
	http://blog.latermoon.me/first-blog.html
	http://blog.latermoon.me/private/wallpaper_1.html
	http://blog.latermoon.me/private/img/ok
