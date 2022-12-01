package file

import "testing"

func TestMatchFile(t *testing.T) {
	if !MatchFile(OpModeDecrypt, "/path/foo.crypted") {
		t.Fatal("does not match file to decrypt")
	}

	if MatchFile(OpModeDecrypt, "/path/foo") {
		t.Fatal("matches file erronously to decrypt")
	}

	if !MatchFile(OpModeEncrypt, "/path/foo.docx") {
		t.Fatal("does not match .docx file to encrypt")
	}

	if MatchFile(OpModeEncrypt, "/path/foo.crypted") {
		t.Fatal("matches crypted file to encrypt")
	}
}
