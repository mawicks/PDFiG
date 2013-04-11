package readers

import (
	"fmt"
	"bufio"
//	"os"
	"strings"
	"testing" )

func TestReverseReader (t *testing.T) {
	reader := bufio.NewReader(NewReverseReader(strings.NewReader("abcdefghijklmnopqrstuvwxyz")))
	var readString string
	fmt.Fscanf (reader, "%s", &readString)
	expected := "zyxwvutsrqponmlkjihgfedcba"
	if readString != expected {
		t.Errorf (`Expected: "%s"; got "%s"`, expected, readString)
	}
}

