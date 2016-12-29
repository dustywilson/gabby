package gabby

import (
	"encoding/base64"
	"errors"
	"net"
	"net/url"
	"strings"
	"sync"

	"github.com/dedis/crypto/abstract"
	"github.com/dedis/crypto/ed25519"
	"github.com/dedis/crypto/random"
	uuid "github.com/satori/go.uuid"
)

// Domain is a Gabby domain, which maps to DNS
type Domain struct {
	domain          string
	key             *Key
	myPrivateKey    abstract.Scalar
	serverPublicKey abstract.Point
	wsURL           *url.URL
	sync.RWMutex
}

// ErrNoRecords means that no records were found during a record lookup
var ErrNoRecords = errors.New("no records were found")

// NewDomain returns a new *Domain from a string version of the domain
func NewDomain(domain string) *Domain {
	return &Domain{
		domain: domain,
	}
}

// updateDomainDetails fetches info about the domain from DNS
func (d *Domain) updateDomainDetails() error {
	d.Lock()
	defer d.Unlock()
	_, servers, err := net.LookupSRV("servers", "gabby", d.domain) // "_servers._gabby.$domain"
	if err != nil {
		return err
	}
	for _, srvV := range servers {
		wss, err := net.LookupTXT(strings.Join([]string{"_wss._gabby", srvV.Target}, ".")) // "_wss._gabby.$srvV" which might not actually be in $domain
		if err != nil {
			continue
		}
		for _, wssV := range wss {
			wsURL, err := url.Parse(wssV)
			if err != nil {
				continue
			}
			if wsURL.Scheme != "wss" { // yes, wss. not ws.
				continue
			}
			keys, err := net.LookupTXT(strings.Join([]string{"_key._gabby", srvV.Target}, ".")) // "_key._gabby.$srvV" which might not actually be in $domain
			if err != nil {
				continue
			}
			for _, keyB64 := range keys { // note: it's best if there's only one key listed; remembering old keys for key rolling is a problem on the server, not the client
				keyBytes, err := base64.StdEncoding.DecodeString(keyB64)
				if err != nil {
					continue
				}

				suite := ed25519.NewAES128SHA256Ed25519(true)
				serverKey := suite.Point()
				err = serverKey.UnmarshalBinary(keyBytes)
				if err != nil {
					continue
				}

				d.myPrivateKey = suite.Scalar().Pick(random.Stream)
				d.key = NewKey(
					uuid.UUID{}, // ???: what do we use here?
					suite.Point().Mul(serverKey, d.myPrivateKey),
				)
				d.serverPublicKey = serverKey
				d.wsURL = wsURL
			}
		}
	}
	return ErrNoRecords
}
