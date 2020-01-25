package util

import (
	"regexp"

	fmt_util "github.com/Iyeyasu/bingo-paste/internal/util/fmt"
	log "github.com/sirupsen/logrus"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
	"github.com/tdewolff/minify/v2/js"
	"github.com/tdewolff/minify/v2/svg"
)

// Minify minifies HTML and any embedded CSS, JavaScript, or SVGs.
func Minify(body string) string {
	log.Debug("Minifying HTML")

	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("image/svg+xml", svg.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)

	minified, err := m.String("text/html", body)
	if err != nil {
		log.Fatal("Failed to minify HTML template")
		return ""
	}

	oldSize := fmt_util.FormatByteSize(int64(len(body)))
	newSize := fmt_util.FormatByteSize(int64(len(minified)))
	log.Trace(minified)
	log.Debugf("Minified HTML: %s -> %s", oldSize, newSize)

	return minified
}
