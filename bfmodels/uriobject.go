package bfmodels

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

// TODO Make better struct and xml parser method
type URIObject struct {
	Type         string
	Uri          string
	UrlThumbnail string
	Body         struct {
		Text string
		Link struct {
			Href string
			Text string
		}
		Title string
		Swift struct {
			B64 string
		}
		File struct {
			OriginalName string
			Size         int
		}
	}
}

type B64Payload struct {
	Type        string          `json:"type"`
	Summary     string          `json:"summary"`
	Attachments []B64Attachment `json:"attachments"`
}

type B64Attachment struct {
	ContentType string           `json:"contentType"`
	Content     B64AttachContent `json:"content"`
}

type B64AttachContent struct {
	Title    string     `json:"title"`
	Text     string     `json:"text"`
	SubTitle string     `json:"subtitle"`
	Images   []B64Image `json:"images"`
	Image    B64Image   `json:"image"`
	Buttons  B64Buttons `json:"buttons"`
	Media    B64Media   `json:"media"`
	Tap      B64Tap     `json:"tap"`
}

type B64Media struct {
	URL string `json:"url"`
}

type B64Image struct {
	URL string `json:"url"`
	Alt string `json:"alt"`
	Tap B64Tap `json:"tap"`
}

type B64Buttons struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	Value string `json:"value"`
}

type B64Tap struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

func (u *URIObject) DecodeB64() (*B64Payload, error) {
	var res = new(B64Payload)
	sDec, _ := base64.StdEncoding.DecodeString(u.Body.Swift.B64)

	err := json.Unmarshal(sDec, &res)

	return res, err
}

func NewURIObjectFromText(in string) (*URIObject, error) {
	res := &URIObject{}
	buf, err := html.Parse(bytes.NewReader([]byte(in)))
	if err != nil {
		return nil, err
	}

	toNormalstring := func(in string) string {
		out := strings.Replace(in, "\\\"", "", 2)
		return out
	}

	var f func(n *html.Node)

	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if n.Data == "uriobject" {
				for _, attr := range n.Attr {
					if attr.Key == "type" {
						res.Type = toNormalstring(attr.Val)
					}
					if attr.Key == "url_thumbnail" {
						res.UrlThumbnail = toNormalstring(attr.Val)
					}
					if attr.Key == "uri" {
						res.Uri = toNormalstring(attr.Val)
					}
				}

				if n.FirstChild != nil {
					res.Body.Text = n.FirstChild.Data
				}
			}

			if n.Data == "a" {
				for _, attr := range n.Attr {
					if attr.Key == "href" {
						res.Body.Link.Href = toNormalstring(attr.Val)
					}
				}

				if n.FirstChild != nil {
					res.Body.Link.Text = n.FirstChild.Data
				}
			}

			if n.Data == "originalname" {
				for _, attr := range n.Attr {
					if attr.Key == "v" {
						res.Body.File.OriginalName = toNormalstring(attr.Val)
					}
				}
			}

			if n.Data == "filesize" {
				for _, attr := range n.Attr {
					if attr.Key == "v" {
						size, _ := strconv.Atoi(toNormalstring(attr.Val))
						res.Body.File.Size = size
					}
				}
			}

			if n.Data == "title" {
				if n.FirstChild != nil {
					res.Body.Title = n.FirstChild.Data
				}
			}

			if n.Data == "swift" {
				for _, attr := range n.Attr {
					if attr.Key == "b64" {
						res.Body.Swift.B64 = toNormalstring(attr.Val)
					}
				}
			}

		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(buf)
	return res, nil
}
