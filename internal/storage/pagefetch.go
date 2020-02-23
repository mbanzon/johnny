package storage

// PageFetch is used to store a fetch result of a Page.
type PageFetch struct {
	id            int64
	pageID        int64
	FetchTime     int64
	FetchDuration int64
	ResponseCode  int
	Size          int64
	ContentType   string
}
