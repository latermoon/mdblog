package builder

type Articles []*Article

func (a Articles) Len() int {
	return len(a)
}

func (a Articles) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a Articles) Less(i, j int) bool {
	return a[i].info.ModTime().UnixNano() < a[j].info.ModTime().UnixNano()
}
