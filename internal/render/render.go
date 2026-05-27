package render

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"

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

func Markdown(src string) (*Result, error) {
	// Re-introduce hashing from the legacy gateway
	sum := sha256.Sum256([]byte(src))
	hash := hex.EncodeToString(sum[:])

	var buf bytes.Buffer
	if err := md.Convert([]byte(src), &buf); err != nil {
		return nil, err
	}

	return &Result{
		HTMLContent: buf.String(),
		Hash:        hash,
	}, nil
}
