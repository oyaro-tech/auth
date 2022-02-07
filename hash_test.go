package auth

import "testing"

func TestComparePasswords(t *testing.T) {
	pwdHash := "$2a$10$GH8J3KJOwwclDHOgEs6qZ.KY18HADgN.hHHBao9oTZ13W7ian75Cm"
	input := "qwerty123"

	if !comparePasswords(input, pwdHash) {
		t.Errorf("ComparePasswords(\"%s\", \"%s\") = false; want true", input, pwdHash)
	}
}
