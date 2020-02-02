package util

import (
	"html"
	"strings"

	"github.com/Iyeyasu/bingo-paste/internal/util/log"
	"github.com/alecthomas/chroma"
	htmlf "github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

// FormatCode applies syntax highlighting to the given string.
func FormatCode(lang string, content string) string {
	log.Debugf("Highlighting syntax using language %s", lang)

	lexer := getLexer(lang)
	style := getStyle("swapoff")
	formatter := getFormatter()

	it, err := lexer.Tokenise(nil, content)
	builder := new(strings.Builder)
	err = formatter.Format(builder, style, it)

	if err != nil {
		log.Errorf("Failed to highlight syntax: %s", err.Error())
		return html.EscapeString(content)
	}

	return builder.String()
}

func getLexer(lang string) chroma.Lexer {
	log.Debugf("Using lexer '%s'", lang)

	lexer := lexers.Get(lang)
	if lexer == nil {
		log.Warnf("Failed to get lexer %s. Falling back to default", lang)
		return lexers.Fallback
	}

	lexer = chroma.Coalesce(lexer)
	return lexer
}

func getFormatter() *htmlf.Formatter {
	log.Debugf("Using html formatter")

	return htmlf.New(
		htmlf.Standalone(false),
		htmlf.WithLineNumbers(false),
		htmlf.WithClasses(true))
}

func getStyle(name string) *chroma.Style {
	log.Debugf("Using style '%s'", name)

	style := styles.Get(name)
	if style == nil {
		log.Warnf("Failed to find style %s. Falling back to default", name)
		return styles.Fallback
	}

	return style
}
