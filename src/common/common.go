//Package common implements common unilities
//Date : 2015-12-26
package common

import (
	"math"
	"strings"
)

//format str number with certain length , if not enough , fill with "0"
//-length means backward
func SubStrByLen(str string, length int, prefix string) (subStr string) {
	strLen := len(str)
	absLength := int(math.Abs(float64(length))) //get absolute length ignore the negtive or positive
	if length >= 0 {
		if strLen >= length {
			subStr = str[0:length]
		} else {
			subStr = str
		}
	} else {
		if strLen >= absLength {
			subStr = str[strLen-absLength : strLen]
		} else {
			subStr = str
		}

	}
	var preStr string
	for i := 0; i < absLength-strLen; i++ {
		preStr += prefix
	}
	return (preStr + subStr)
}

//Get last split sub string
func LastSplit(s, split string, backward ...bool) string {
	var lastIndex int
	splitStr := strings.Split(s, split)
	length := len(splitStr)
	if length == 0 {
		lastIndex = 0
	} else {
		lastIndex = length - 1
	}

	if len(backward) != 0 {
		if backward[0] == true {
			lastIndex = 0
		}
	}

	return splitStr[lastIndex]
}
