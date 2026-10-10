package sites

import "testing"

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "adds https scheme",
			input: "example.com",
			want:  "https://example.com",
		},
		{
			name:  "keeps https scheme",
			input: "https://example.com",
			want:  "https://example.com",
		},
		{
			name:  "keeps http scheme",
			input: "http://example.com",
			want:  "http://example.com",
		},
		{
			name:  "url with path and query",
			input: "example.com/path?foo=bar",
			want:  "https://example.com/path?foo=bar",
		},
		{
			name:  "url with port",
			input: "example.com:8080",
			want:  "https://example.com:8080",
		},
		{
			name:  "url with fragment",
			input: "https://example.com/path#section",
			want:  "https://example.com/path#section",
		},
		{
			name:    "unsupported ftp scheme",
			input:   "ftp://example.com",
			wantErr: true,
		},
		{
			name:    "unsupported file scheme",
			input:   "file:///tmp/test.txt",
			wantErr: true,
		},
		{
			name:    "empty input",
			input:   "",
			wantErr: true,
		},
		{
			name:    "empty host",
			input:   "https://",
			wantErr: true,
		},
		{
			name:    "missing host",
			input:   "https:///path",
			wantErr: true,
		},
		{
			name:    "invalid port",
			input:   "https://example.com:abc",
			wantErr: true,
		},
		{
			name:    "invalid IPv6 address",
			input:   "https://[::1",
			wantErr: true,
		},
		{
			name:    "invalid percent escape",
			input:   "https://example.com/%zz",
			wantErr: true,
		},
		{
			name:  "uppercase scheme",
			input: "HTTPS://example.com",
			want:  "https://example.com",
		},
		{
			name:  "scheme inside query",
			input: "example.com/login?next=https://other.com",
			want:  "https://example.com/login?next=https://other.com",
		},
		{
			name:    "no hostname",
			input:   "://example.com/",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeURL(tt.input)

			if (err != nil) != tt.wantErr {
				t.Fatalf("normalizeURL(%q) error = %v, wantErr %v",
					tt.input, err, tt.wantErr)
			}

			if err == nil && got != tt.want {
				t.Errorf("normalizeURL(%q) = %q, want %q",
					tt.input, got, tt.want)
			}
		})
	}
}
