package resp

import "testing"

func TestParseErr(t *testing.T) {
	errCode1 := BadRequest
	errCode2 := UserNotFound

	module1, code1 := ParseErrorCode(errCode1)
	t.Log("1:", module1, code1)
	module2, code2 := ParseErrorCode(errCode2)
	t.Log("2:", module2, code2)
}
