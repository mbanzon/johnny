package storage

// Page is a single page identified by a URL.
type Page struct {
	id     int64
	domain int64
	url    string
}

type PageStorage interface {
}
