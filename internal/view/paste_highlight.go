package view

import (
	"html"
	"log"
	"strings"

	"github.com/Iyeyasu/bingo-paste/internal/model"
	"github.com/alecthomas/chroma"
	htmlf "github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
)

// HighlightPaste applies syntax highlighting to the given paste.
func HighlightPaste(paste *model.Paste) string {
	log.Printf("Highlighting syntax for paste %d", paste.ID)

	lexer := getLexer(paste.Syntax)
	style := getStyle("swapoff")
	formatter := getFormatter()

	it, err := lexer.Tokenise(nil, paste.RawContent)
	builder := new(strings.Builder)
	err = formatter.Format(builder, style, it)

	if err != nil {
		log.Println(err)
		return html.EscapeString(paste.RawContent)
	}

	return builder.String()
}

func getLexer(lang string) chroma.Lexer {
	log.Printf("Using %s lexer", lang)

	var lexer chroma.Lexer
	lexer = lexers.Get(lang)
	if lexer == nil {
		log.Printf("Failed to get lexer %s", lang)
		return lexers.Fallback
	}

	lexer = chroma.Coalesce(lexer)
	return lexer
}

func getFormatter() *htmlf.Formatter {
	log.Printf("Using html formatter")

	return htmlf.New(
		htmlf.Standalone(false),
		htmlf.WithLineNumbers(false),
		htmlf.WithClasses(true))
}

func getStyle(name string) *chroma.Style {
	log.Printf("Using style %s", name)

	style := styles.Get(name)
	if style == nil {
		log.Printf("Failed to find style %s", name)
		style = styles.Fallback
	}

	return style
}
