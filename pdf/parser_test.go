package pdf_test

import (
	"bytes"
	"fmt"
	"github.com/mawicks/goPDF/pdf"
	"strings"
	"testing" )

func TestParser (t *testing.T) {
	testParse := func(source string, expected string) {
		o,err,_ := pdf.Scan (strings.NewReader(source))
		if (err != nil) {
			t.Errorf(`Scan() of "%s" returned error: %v`, source, err)
		}
		testOneObject (t, fmt.Sprintf(`Scan of "%s"`, pdf.AsciiFromBytes([]byte(source))), o, nil, expected)
	}

	testParseFail := func(source string, prefix string) {
		_,err,context := pdf.Scan (strings.NewReader(source))
		if (err == nil) {
			t.Errorf(`Scan() of "%s" did NOT return an error`, source)
		}
		if !bytes.Equal(context,[]byte(prefix)) {
			t.Errorf(`Scan() of "%s" returned "%s" as error context instead of "%s".`, source,
				pdf.AsciiFromBytes(context),
				pdf.AsciiFromBytes([]byte(prefix)))
		}
	}

	testParse (" \n  null", "null")
	testParse ("true", "true")
	testParse ("false", "false")
	testParse ("/foo", "/foo")
	testParse ("/a#20b#20c", "/a#20b#20c")
	testParse ("(abc)", "(abc)")
	testParse ("(a(bc))", "(a\\(bc\\))")
	testParse ("(\\061)", "(1)")
	testParse ("[]", "[]")
	testParse ("[true false]", "[true false]")
	testParse ("<<>>", "<<>>")
	testParse ("<</foo true>>", "<</foo true>>")
	testParse ("<302031>", "(0 1)")
	testParse ("123.456", "123.456")
	testParse ("-54321", "-54321")
	testParse ("<</Length 5>>\nstream\nabcde\nendstream", "<</Length 5>>\nstream\nabcde\nendstream")

	testParseFail("  /a#", "  /a#")
	testParseFail("  /a#(123)", "  /a#(")
	testParseFail("falxe  ", "falxe")

}

