package test

import (
	"github.com/small-ek/ginp/encoding/html"
	"log"
	"testing"
)

func TestHtml(t *testing.T) {
	var src = `<p>Test paragraph.</p><!-- Comment -->  <a href="#fragment">Other text</a>`
	var src2 = `A 'quote' "is" <b>bold</b>`
	var src3 = `A 'quote' "is" <b>bold</b>`
	log.Println(html.StripTags(src))
	var data = html.Entities(src2)

	log.Println(data)
	log.Println(html.EntitiesDecode(data))

	var data2 = html.SpecialChars(src3)
	log.Println(data2)
	log.Println(html.SpecialCharsDecode(data2))
}
