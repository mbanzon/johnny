package storage

// Domain represents a fully valid domain name. The top level domain is also stored for
// easy sorting.
type Domain struct {
	id   int64
	tld  string
	name string
}

type DomainStorage interface {
	New(d string) (Domain, error)
	GetFromName(name string) (Domain, error)
	GetFromID(id int64) (Domain, error)
}
