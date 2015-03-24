package builder

type ArticleInfos []*ArticleInfo

func (a ArticleInfos) Len() int {
	return len(a)
}

func (a ArticleInfos) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ArticleInfos) Less(i, j int) bool {
	return a[i].Date.UnixNano() < a[j].Date.UnixNano()
}
