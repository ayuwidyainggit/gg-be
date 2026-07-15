package adapter

import "testing"

func TestObsAdapterImpl_ResolvePublicFileURL(t *testing.T) {
	adapter := &ObsAdapterImpl{FileBaseUrl: "https://bucket.example.com"}

	tests := []struct {
		name    string
		input   string
		wantURL string
		wantOK  bool
	}{
		{name: "empty", input: "", wantOK: false},
		{name: "leading slash", input: "/folder/file.jpg", wantURL: "https://bucket.example.com/folder/file.jpg", wantOK: true},
		{name: "special chars", input: "folder/file a+b.jpg", wantURL: "https://bucket.example.com/folder/file%20a+b.jpg", wantOK: true},
		{name: "same host full url", input: "https://bucket.example.com/folder/file.jpg", wantURL: "https://bucket.example.com/folder/file.jpg", wantOK: true},
		{name: "host contains base host rejected", input: "https://evil-bucket.example.com/file.jpg", wantOK: false},
		{name: "external rejected", input: "https://evil.example.com/file.jpg", wantOK: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotURL, gotOK := adapter.ResolvePublicFileURL(tt.input)
			if gotOK != tt.wantOK {
				t.Fatalf("expected ok=%v, got %v", tt.wantOK, gotOK)
			}
			if gotURL != tt.wantURL {
				t.Fatalf("expected url=%q, got %q", tt.wantURL, gotURL)
			}
		})
	}
}
