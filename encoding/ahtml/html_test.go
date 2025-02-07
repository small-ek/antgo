package ahtml

import (
	"log"
	"testing"
)

func TestHtml(t *testing.T) {
	var src = `<p>Test paragraph.</p><!-- Comment -->  <a href="#fragment">Other text</a>`
	var src2 = `A 'quote' "is" <b>bold</b>`
	var src3 = `A 'quote' "is" <b>bold</b>`
	log.Println(StripTags(src))
	var data = Entities(src2)

	log.Println(data)
	log.Println(EntitiesDecode(data))

	var data2 = SpecialChars(src3)
	log.Println(data2)
	log.Println(SpecialCharsDecode(data2))
}
