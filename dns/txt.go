package dns

import (
	"strings"
	"time"

	"github.com/dustywilson/httpdig"
)

// TXT is a DNS TXT RR
type TXT struct {
	TTL   time.Duration
	Value string
}

// LookupTXT looks up a TXT RR
func LookupTXT(name string) ([]TXT, error) {
	r, err := httpdig.Query(name, "TXT")
	if err != nil {
		return nil, err
	}
	txts := make([]TXT, len(r.Answer))
	for i, a := range r.Answer {
		txts[i] = TXT{
			TTL:   a.TTL,
			Value: strings.TrimPrefix(strings.TrimSuffix(a.Data, `"`), `"`),
		}
	}
	return txts, nil
}
