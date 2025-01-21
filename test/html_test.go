package test

import (
	"github.com/small-ek/antgo/encoding/ahtml"
	"log"
	"testing"
)

func TestHtml(t *testing.T) {
	var src = `<p>Test paragraph.</p><!-- Comment -->  <a href="#fragment">Other text</a>`
	var src2 = `A 'quote' "is" <b>bold</b>`
	var src3 = `A 'quote' "is" <b>bold</b>`
	log.Println(ahtml.StripTags(src))
	var data = ahtml.Entities(src2)

	log.Println(data)
	log.Println(ahtml.EntitiesDecode(data))

	var data2 = ahtml.SpecialChars(src3)
	log.Println(data2)
	log.Println(ahtml.SpecialCharsDecode(data2))
}
