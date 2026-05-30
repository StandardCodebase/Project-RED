package render

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

type Result struct {
	HTMLContent string
	Hash        string
}

var md = goldmark.New(
	goldmark.WithExtensions(
		extension.GFM,
		extension.Table,
		extension.Typographer,
	),
	goldmark.WithParserOptions(
		parser.WithAutoHeadingID(),
	),
	goldmark.WithRendererOptions(
		html.WithUnsafe(),
	),
)

var sanitizer = func() *bluemonday.Policy {
	p := bluemonday.UGCPolicy()
	p.AllowAttrs("class").OnElements("code", "pre", "span", "div", "p", "li", "input")
	p.AllowAttrs("id").OnElements("h1", "h2", "h3", "h4", "h5", "h6")
	p.AllowAttrs("align").OnElements("td", "th")
	p.AllowAttrs("type", "checked", "disabled", "readonly").OnElements("input")
	p.AllowRelativeURLs(true)
	p.AllowDataURIImages()
	return p
}()

func Markdown(src string) (*Result, error) {
	sum := sha256.Sum256([]byte(src))
	hash := hex.EncodeToString(sum[:])

	var buf bytes.Buffer
	if err := md.Convert([]byte(src), &buf); err != nil {
		return nil, err
	}

	safeHTML := sanitizer.Sanitize(buf.String())

	return &Result{
		HTMLContent: safeHTML,
		Hash:        hash,
	}, nil
}
