package cookiesessionconverter

import (
	"net/http"

	"fknsrs.biz/p/cookiesession"
	"fknsrs.biz/p/requesttoken"
)

type Converter struct {
	requesttoken.Converter
	Store *cookiesession.Store
}

func (c Converter) GetSessionAndToken(r *http.Request) ([]byte, []byte) {
	if s := c.Store.Get(r); s.Valid {
		return s.SID.Bytes(), c.Converter.GetToken(r)
	}

	return nil, c.Converter.GetToken(r)
}
