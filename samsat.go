package samsat

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/go-querystring/query"
)

const url = "http://samsat-pkb.jakarta.go.id/cek-ranmor+pajak-dki/"

// Request is samsat form datajjj
type Request struct {
	Nopa   string `url:"nopa"`
	Noph   string `url:"noph"`
	N1k    string `url:"n1k"`
	Tombol string `url:"tombol"`
	Flag   int32  `url:"flag"`
}

type Response struct {
	Brand     string `json:"brand"`
	Model     string `json:"model"`
	Color     string `json:"color"`
	Owner     string `json:"owner"`
	SellPrice string `json:"sell_price"`
	TaxStatus string `json:"tax_status"`
}

func Search(nopa, noph string) (*Response, error) {
	// send form data
	reqBody := &Request{
		Nopa:   nopa,
		Noph:   noph,
		Tombol: "Proses",
		N1k:    "",
		Flag:   2,
	}
	payload, err := query.Values(reqBody)
	if err != nil {
		return nil, err
	}

	res, err := http.PostForm(url, payload)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	// start parsing vehicle info
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return nil, err
	}

	sel := "#print-content > div > table > tbody > tr:nth-child(%d) > td:nth-child(%d)"
	return &Response{
		Brand:     strings.TrimSpace(doc.Find(fmt.Sprintf(sel, 2, 2)).Text()),
		Model:     strings.TrimSpace(doc.Find(fmt.Sprintf(sel, 2, 4)).Text()),
		Color:     strings.TrimSpace(doc.Find(fmt.Sprintf(sel, 3, 2)).Text()),
		Owner:     strings.TrimSpace(doc.Find(fmt.Sprintf(sel, 1, 2)).Text()),
		SellPrice: strings.TrimSpace(doc.Find(fmt.Sprintf(sel, 5, 3)).Text()),
		TaxStatus: strings.TrimSpace(doc.Find(fmt.Sprintf(sel, 11, 2)).Text()),
	}, nil
}
