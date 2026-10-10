package proof

import (
	"archive/tar"
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func buildTar(t *testing.T, entries []*tar.Header, bodies map[string]string) []byte {
	t.Helper()
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	for _, h := range entries {
		if err := tw.WriteHeader(h); err != nil {
			t.Fatal(err)
		}
		if body, ok := bodies[h.Name]; ok {
			if _, err := tw.Write([]byte(body)); err != nil {
				t.Fatal(err)
			}
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func TestExtractTarSkipsPAXGlobalHeader(t *testing.T) {
	dest := t.TempDir()
	data := buildTar(t, []*tar.Header{
		{Name: "pax_global_header", Typeflag: tar.TypeXGlobalHeader, Size: 0},
		{Name: "file.txt", Typeflag: tar.TypeReg, Size: 5, Mode: 0o644},
	}, map[string]string{"file.txt": "hello"})

	if err := extractTar(bytes.NewReader(data), dest); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	got, err := os.ReadFile(filepath.Join(dest, "file.txt"))
	if err != nil || string(got) != "hello" {
		t.Fatalf("file not extracted: %q err=%v", got, err)
	}
}

func TestExtractTarRejectsTraversal(t *testing.T) {
	dest := t.TempDir()
	data := buildTar(t, []*tar.Header{
		{Name: "../escape.txt", Typeflag: tar.TypeReg, Size: 1, Mode: 0o644},
	}, map[string]string{"../escape.txt": "x"})
	if err := extractTar(bytes.NewReader(data), dest); err == nil {
		t.Fatal("expected traversal rejection")
	}
}

func TestExtractTarRejectsSymlink(t *testing.T) {
	dest := t.TempDir()
	data := buildTar(t, []*tar.Header{
		{Name: "link", Typeflag: tar.TypeSymlink, Linkname: "/etc/passwd"},
	}, nil)
	if err := extractTar(bytes.NewReader(data), dest); err == nil {
		t.Fatal("expected symlink rejection")
	}
}

func TestExtractTarEntryLimit(t *testing.T) {
	dest := t.TempDir()
	var headers []*tar.Header
	for i := 0; i < maxSnapshotEntries+1; i++ {
		headers = append(headers, &tar.Header{Name: "d/" + itoa(i) + ".txt", Typeflag: tar.TypeReg, Size: 0, Mode: 0o644})
	}
	data := buildTar(t, headers, nil)
	if err := extractTar(bytes.NewReader(data), dest); err == nil {
		t.Fatal("expected entry limit rejection")
	}
}

func itoa(i int) string {
	if i == 0 {
		return "0"
	}
	var b []byte
	for i > 0 {
		b = append([]byte{byte('0' + i%10)}, b...)
		i /= 10
	}
	return string(b)
}
