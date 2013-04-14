package pdf_test

import (
	"github.com/mawicks/goPDF/pdf"
	"strings"
	"testing" )

func TestParser (t *testing.T) {
	reader := strings.NewReader("null true false /foo /a#20b#20c (abc) (a(bc)) (\\061) /a#")

	o,err := pdf.Scan (reader)
	if (err != nil) {
		t.Error(`Scan() of "null ..." returned error:`, err)
	}
	testOneObject (t, "null", o, nil, "null")

	o,err = pdf.Scan (reader)
	if (err != nil) {
		t.Error(`Scan() of "true ..." returned error:`, err)
	}
	testOneObject (t, "true", o, nil, "true")

	o,err = pdf.Scan (reader)
	if (err != nil) {
		t.Error(`Scan() of "false ..." returned error:`, err)
	}
	testOneObject (t, "false", o, nil, "false")

	o,err = pdf.Scan (reader)
	if (err != nil) {
		t.Error(`Scan() of "/foo ..." returned error:`, err)
	}
	testOneObject (t, "/foo", o, nil, "/foo")

	o,err = pdf.Scan (reader)
	if (err != nil) {
		t.Error(`Scan() of "/a#20#b#20c ..." returned error:`, err)
	}
	testOneObject (t, "/a#20b#20c", o, nil, "/a#20b#20c")

	o,err = pdf.Scan (reader)
	if (err != nil) {
		t.Error(`Scan() of "(abc) ..." returned error:`, err)
	}
	testOneObject (t, "(abc)", o, nil, "(abc)")

	o,err = pdf.Scan (reader)
	if (err != nil) {
		t.Error(`Scan() of "(a(bc)) ..." returned error:`, err)
	}
	testOneObject (t, "(abc)", o, nil, "(a\\(bc\\))")

	o,err = pdf.Scan (reader)
	if (err != nil) {
		t.Error(`Scan() of "(\\061) ..." returned error:`, err)
	}
	testOneObject (t, "(\\061)", o, nil, "(1)")

	o,err = pdf.Scan (reader)
	if (err == nil) {
		t.Error(`Scan() of "/a#" did NOT return error:`, err)
	}

	reader = strings.NewReader("falxe ")
	o,err = pdf.Scan (reader)
	if (err == nil) {
		t.Error(`Scan() of "falxe" did NOT return error:`, err)
	}

	// Make sure end of file doesn't generate an erro
	reader = strings.NewReader("false")
	o,err = pdf.Scan (reader)
	if (err != nil) {
		t.Error(`Scan() of "false" returned error:`, err)
	}

	// Make sure end of file doesn't generate an error
	reader = strings.NewReader("/foo")
	o,err = pdf.Scan (reader)
	if (err != nil) {
		t.Error(`Scan() of "/foo" returned error:`, err)
	}

}

