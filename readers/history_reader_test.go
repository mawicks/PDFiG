package readers

import (
	"strings"
	"testing" )

func TestHistoryReader (t *testing.T) {
	reader := NewHistoryReader(strings.NewReader("abcdefghijklmnopqrstuvwxyz"),4)

	check := func (expected string) {
		if h:=string(reader.GetHistory()); h != expected {
			t.Errorf (`Expected: "%s"; got "%s"`, expected, h)
		}
	}

	readAndCheck := func (expected string) {
		reader.ReadByte(); check(expected);
	}
	unreadAndCheck := func (expected string) {
		reader.UnreadByte(); check(expected);
	}

	check ("")

	readAndCheck ("a")
	readAndCheck ("ab")
	readAndCheck ("abc")
	readAndCheck ("abcd")
	readAndCheck ("bcde")
	readAndCheck ("cdef")

	unreadAndCheck("cde")
	unreadAndCheck("cd")

	b := make([]byte,4); reader.Read(b); check ("efgh")
}

