package proof

import (
	"archive/tar"
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Extraction limits guard against hostile or runaway archives.
const (
	maxSnapshotBytes   = 256 << 20 // 256 MiB total extracted
	maxSnapshotEntries = 50000
)

// ResolveCommit resolves a revision (such as HEAD or a branch name) to a full
// commit SHA within repoDir.
func ResolveCommit(ctx context.Context, repoDir, revision string) (string, error) {
	if revision == "" {
		revision = "HEAD"
	}
	out, err := runGit(ctx, repoDir, "rev-parse", "--verify", revision+"^{commit}")
	if err != nil {
		return "", fmt.Errorf("resolving revision %q: %w", revision, err)
	}
	sha := strings.TrimSpace(out)
	if len(sha) != 40 {
		return "", fmt.Errorf("resolving revision %q: unexpected object id %q", revision, sha)
	}
	return sha, nil
}

// ExtractSnapshot exports the committed tree at commit into a fresh temporary
// directory using `git archive`, applying strict safety checks. It never
// touches the working checkout, so dirty and untracked files are excluded.
// The caller is responsible for removing the returned directory.
func ExtractSnapshot(ctx context.Context, repoDir, commit string) (string, error) {
	archive, err := runGitBytes(ctx, repoDir, "archive", "--format=tar", commit)
	if err != nil {
		return "", fmt.Errorf("archiving commit %s: %w", commit, err)
	}

	dest, err := os.MkdirTemp("", "agentproof-snapshot-")
	if err != nil {
		return "", fmt.Errorf("creating snapshot directory: %w", err)
	}

	if err := extractTar(bytes.NewReader(archive), dest); err != nil {
		_ = os.RemoveAll(dest)
		return "", err
	}
	return dest, nil
}

// extractTar safely extracts a tar stream into dest.
func extractTar(r io.Reader, dest string) error {
	tr := tar.NewReader(r)
	var totalBytes int64
	var entries int

	destAbs, err := filepath.Abs(dest)
	if err != nil {
		return fmt.Errorf("resolving snapshot destination: %w", err)
	}

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("reading archive: %w", err)
		}

		// Skip safe tar metadata records such as git's pax_global_header,
		// which carries the commit id and is not a file to extract.
		switch hdr.Typeflag {
		case tar.TypeXGlobalHeader, tar.TypeXHeader:
			continue
		}

		entries++
		if entries > maxSnapshotEntries {
			return fmt.Errorf("archive exceeds %d entries", maxSnapshotEntries)
		}

		target, err := safeJoin(destAbs, hdr.Name)
		if err != nil {
			return err
		}

		switch hdr.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(target, 0o755); err != nil {
				return fmt.Errorf("creating directory %q: %w", hdr.Name, err)
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
				return fmt.Errorf("creating parent of %q: %w", hdr.Name, err)
			}
			remaining := maxSnapshotBytes - totalBytes
			if remaining < 0 {
				remaining = 0
			}
			written, err := writeRegular(target, tr, hdr.FileInfo().Mode().Perm(), remaining+1)
			if err != nil {
				return err
			}
			totalBytes += written
			if totalBytes > maxSnapshotBytes {
				return fmt.Errorf("archive exceeds %d bytes", maxSnapshotBytes)
			}
		case tar.TypeSymlink, tar.TypeLink:
			return fmt.Errorf("archive entry %q is a link, which is rejected", hdr.Name)
		case tar.TypeChar, tar.TypeBlock, tar.TypeFifo:
			return fmt.Errorf("archive entry %q is a special file, which is rejected", hdr.Name)
		default:
			return fmt.Errorf("archive entry %q has unsupported type %q", hdr.Name, string(rune(hdr.Typeflag)))
		}
	}
	return nil
}

func writeRegular(target string, r io.Reader, mode os.FileMode, limit int64) (int64, error) {
	if mode == 0 {
		mode = 0o644
	}
	f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return 0, fmt.Errorf("creating file %q: %w", target, err)
	}
	defer f.Close()
	n, err := io.Copy(f, io.LimitReader(r, limit))
	if err != nil {
		return n, fmt.Errorf("writing file %q: %w", target, err)
	}
	return n, nil
}

// safeJoin joins name onto base, rejecting absolute paths and any entry that
// traverses parent directories rather than silently sanitising it.
func safeJoin(base, name string) (string, error) {
	if name == "" {
		return "", fmt.Errorf("archive entry has empty name")
	}
	slash := filepath.ToSlash(name)
	if strings.HasPrefix(slash, "/") {
		return "", fmt.Errorf("archive entry %q is an absolute path", name)
	}
	for _, seg := range strings.Split(slash, "/") {
		if seg == ".." {
			return "", fmt.Errorf("archive entry %q traverses parent directories", name)
		}
	}
	target := filepath.Join(base, filepath.FromSlash(slash))
	rel, err := filepath.Rel(base, target)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("archive entry %q escapes the snapshot root", name)
	}
	return target, nil
}

func runGit(ctx context.Context, dir string, args ...string) (string, error) {
	out, err := runGitBytes(ctx, dir, args...)
	return string(out), err
}

func runGitBytes(ctx context.Context, dir string, args ...string) ([]byte, error) {
	cmd := exec.CommandContext(ctx, "git", args...)
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg != "" {
			return nil, fmt.Errorf("git %s: %s", strings.Join(args, " "), msg)
		}
		return nil, fmt.Errorf("git %s: %w", strings.Join(args, " "), err)
	}
	return stdout.Bytes(), nil
}
