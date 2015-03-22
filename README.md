# Markdown Blog Generator

### How to use
```
git clone github.com/latermoon/mdblog
cd mdblog
export GOPATH=$GOPATH:($pwd)
go build src/mdblog.go

./mdblog -d website_dir -p localhost:8801
```

### Imported
```
go get github.com/PuerkitoBio/goquery
go get github.com/russross/blackfriday
go get github.com/nfnt/resize
go get github.com/howeyc/fsnotify
```

### URL
	http://blog.latermoon.me/first-blog.html
	http://blog.latermoon.me/2012/
	http://blog.latermoon.me/2012/06/
