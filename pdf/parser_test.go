package pdf_test

import (
	"bytes"
	"fmt"
	"github.com/mawicks/pdfdig/pdf"
	"strings"
	"testing" )

var mockFile pdf.File = pdf.NewMockFile(100, 1)

func TestParser (t *testing.T) {
	testParse := func(source string, expected string) {
		parser := pdf.NewParser (strings.NewReader(source))
		o,err := parser.Scan (mockFile)
		if (err != nil) {
			t.Errorf(`Scan() of "%s" returned error: %v`, source, err)
		}
		testOneObject (t, fmt.Sprintf(`Scan of "%s"`, pdf.AsciiFromBytes([]byte(source))), o, mockFile, expected)
	}

	testParseFail := func(source string, prefix string) {
		parser := pdf.NewParser (strings.NewReader(source))
		_,err := parser.Scan (mockFile)
		if (err == nil) {
			t.Errorf(`Scan() of "%s" did NOT return an error`, source)
		}
		context := parser.GetContext()
		if !bytes.Equal(context,[]byte(prefix)) {
			t.Errorf(`Scan() of "%s" returned "%s" as error context instead of "%s".`, source,
				pdf.AsciiFromBytes(context),
				pdf.AsciiFromBytes([]byte(prefix)))
		}
	}

	testParse ("null", "null")
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

	// White space tests.
	testParse ("[ 1 % Ignore me \n 2 ]", "[1 2]")
	testParse ("[ 1 % Ignore me \r 2 ]", "[1 2]")
	testParse ("[ 1 % Ignore me \n\r 2 ]", "[1 2]")
	testParse ("[ 1 \n\r 2 ]", "[1 2]")

	// Parsing sequences of integers is tricky because of the
	// possibility that one of them ends in "R", e.g., "1 0 R".
	// Make sure these works properly.
	testParse ("[100]", "[100]")
	testParse ("[201 202]", "[201 202]")
	testParse ("[301 302 303]", "[301 302 303]")
	testParse ("[401 402 403 404]", "[401 402 403 404]")

	testParse ("[100]", "[100]")
	testParse ("[201 202 R]", "[201 202 R]")
	testParse ("[301 302 303 R]", "[301 302 303 R]")
	testParse ("[401 402 403 404 R]", "[401 402 403 404 R]")

	testParse ("[100]", "[100]")
	testParse ("[201 /name]", "[201 /name]")
	testParse ("[301 /name 303 304]", "[301 /name 303 304]")

	testParse ("[100]", "[100]")
	testParse ("[201 3.14]", "[201 3.14]")
	testParse ("[301 3.14 303 304]", "[301 3.14 303 304]")

	testParse ("[100]", "[100]")
	testParse ("[201 202 R 203]", "[201 202 R 203]")
	testParse ("[301 302 303 R 304]", "[301 302 303 R 304]")
	testParse ("[401 402 403 404 R 405]", "[401 402 403 404 R 405]")

	testParseFail("  /a#", "  /a#")
	testParseFail("  /a#(123)", "  /a#(")
	testParseFail("falxe  ", "falxe")

}

