package httprpc

import "testing"

func TestNormalizePath(t *testing.T) {
	tests := []struct {
		path   string
		normal string
	}{
		{path: "", normal: "/"},
		{path: "/", normal: "/"},
		{path: "a", normal: "/a"},
		{path: "/a", normal: "/a"},
		{path: "/a/", normal: "/a"},
		{path: "/a/b", normal: "/a/b"},
		{path: "/a/b/", normal: "/a/b"},
		{path: "a/b", normal: "/a/b"},
		{path: "a/b/", normal: "/a/b"},
	}
	for _, tt := range tests {
		if got, want := normalizePath(tt.path), tt.normal; got != want {
			t.Errorf("normalizePath(%s): got %v, want %v", tt.path, got, want)
		}
	}
}

func TestSplitPath(t *testing.T) {
	tests := []struct {
		path   string
		class  string
		method string
	}{
		{path: "", class: "/", method: ""},
		{path: "/", class: "/", method: ""},
		{path: "/a", class: "/", method: "a"},
		{path: "/a/b", class: "/a", method: "b"},
		{path: "/a/b/", class: "/a", method: "b"},
		{path: "a/b", class: "/a", method: "b"},
		{path: "a/b/", class: "/a", method: "b"},
	}
	for _, tt := range tests {
		class, method := splitPath(tt.path)
		if got, want := class, tt.class; got != want {
			t.Errorf("splitPath(%s): class got %v, want %v", tt.path, got, want)
		}
		if got, want := method, tt.method; got != want {
			t.Errorf("splitPath(%s): method got %v, want %v", tt.path, got, want)
		}
		t.Logf("splitPath(%s): class got %v, method got %v", tt.path, class, method)
	}
}
