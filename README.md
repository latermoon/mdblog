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
go get gopkg.in/fsnotify.v1
github.com/howeyc/fsnotify
```

### Builder
```
builder := NewBlogBuilder()
builder.AddArticles(...)
builder.RebuildAll()

```