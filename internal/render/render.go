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

// We keep WithUnsafe() so Goldmark renders harmless HTML like <kbd> or <details>.
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

// Created a strict sanitizer policy specifically designed for user-generated content.
var sanitizer = bluemonday.UGCPolicy()

func Markdown(src string) (*Result, error) {
	sum := sha256.Sum256([]byte(src))
	hash := hex.EncodeToString(sum[:])

	var buf bytes.Buffer
	if err := md.Convert([]byte(src), &buf); err != nil {
		return nil, err
	}

	// XSS PROTECTION: Sanitize the raw HTML output before returning it to the router
	safeHTML := sanitizer.Sanitize(buf.String())

	return &Result{
		HTMLContent: safeHTML,
		Hash:        hash,
	}, nil
}
