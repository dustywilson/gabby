package domain

// Domain is a Gabby domain, which happens to map to DNS.
type Domain struct {
	name string
}

// String is the string representation of the domain itself (the domain name).
func (d *Domain) String() string {
	return d.name
}

// New returns a new Gabby domain
func New(name string) (*Domain, error) {
	return &Domain{
		name: name,
	}, nil
}
