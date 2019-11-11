package headerconverter

import (
	"net/http"
)

type Converter struct {
	Header string
}

func (c *Converter) GetToken(r *http.Request) []byte {
	if s := r.Header.Get(c.Header); s != "" {
		return []byte(s)
	}

	return nil
}

func (c *Converter) SetToken(r *http.Request, rw http.ResponseWriter, token []byte) {
	if token == nil {
		rw.Header().Del(c.Header)
	} else {
		rw.Header().Set(c.Header, string(token))
	}
}
