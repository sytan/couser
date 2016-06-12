package common

import (
	"flag"
	"testing"
)

var (
	str    string
	length int    //for TestSubStrByLen
	prefix string //for TestSubStrByLen

	split    string //for TestLastSplit
	backward bool   //for TestLastSplit
)

func init() {
	flag.StringVar(&str, "s", "", "string")
	flag.StringVar(&prefix, "pre", "", "string")
	flag.IntVar(&length, "l", 0, "length of sub string")
	flag.StringVar(&split, "split", "", "spliter for TestLastSplit")
	flag.BoolVar(&backward, "b", false, "backward for TestLastSplit")
	flag.Parse()
}

//TestSubStrByLen implements get certain lenth of sub string
func TestSubStrByLen(t *testing.T) {
	subStr := SubStrByLen(str, length, prefix)
	t.Log(str, length, prefix)
	t.Log(subStr)
}

func TestLastSplit(t *testing.T) {
	lastSplit := LastSplit(str, split, backward)
	t.Log(str, split, lastSplit)
}
