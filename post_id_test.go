package grapho

import "testing"

func Test_PostId_trimsSpaces(t *testing.T) {
	postId, err := NewPostId("  hello  ")
	if err != nil {
		t.Fatal(err)
	}

	if got, want := postId.String(), "hello"; got != want {
		t.Fatalf("postId = %q; want %q", got, want)
	}
}

func Test_PostId_onlyAllowsPosixPortableFilenameChars(t *testing.T) {
	tests := []struct {
		in  string
		err error
	}{
		{"hello-world", nil},
		{"hello", nil},
		{"hello-1", nil},
		{"hello_1", nil},
		{"11111", nil},
		{"hello world", ErrMalformedPostId},
		{"hello/world", ErrMalformedPostId},
		{"hello=world", ErrMalformedPostId},
		{"hello?world", ErrMalformedPostId},
		{"привет", ErrMalformedPostId},
	}

	for _, test := range tests {
		_, err := NewPostId(test.in)
		if err != test.err {
			t.Errorf("NewPostId(%q): %v", test.in, err)
		}
	}
}
