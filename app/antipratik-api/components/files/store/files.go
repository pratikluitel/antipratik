// Package store — file storage backend (local disk and Cloudflare R2).
package store

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/pratikluitel/antipratik/components/files"
	"github.com/pratikluitel/antipratik/config"
)

// NewFileStore returns a FileStore implementation based on the supplied config.
func NewFileStore(cfg config.StorageConfig) (files.FileStore, error) {
	switch cfg.Backend {
	case "r2":
		return newR2FileStore(cfg.R2)
	default:
		return newLocalFileStore(cfg.LocalDir), nil
	}
}

// ── Local implementation ──────────────────────────────────────────────────────

type localFileStore struct{ baseDir string }

func newLocalFileStore(dir string) *localFileStore { return &localFileStore{baseDir: dir} }

func (s *localFileStore) Put(_ context.Context, key string, r io.Reader, _ string) (err error) {
	dest := filepath.Join(s.baseDir, filepath.FromSlash(key))
	if err = os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return fmt.Errorf("localFileStore.Put mkdir: %w", err)
	}
	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("localFileStore.Put create: %w", err)
	}
	// Close errors on a writable file indicate un-flushed data; surface them
	// only when no prior error is already being returned.
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = fmt.Errorf("localFileStore.Put close: %w", cerr)
		}
	}()
	if _, err = io.Copy(f, r); err != nil {
		return fmt.Errorf("localFileStore.Put copy: %w", err)
	}
	return nil
}

func (s *localFileStore) Delete(_ context.Context, key string) error {
	path := filepath.Join(s.baseDir, filepath.FromSlash(key))
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("localFileStore.Delete: %w", err)
	}
	return nil
}

func (s *localFileStore) Get(_ context.Context, key string) (io.ReadSeekCloser, string, error) {
	path := filepath.Join(s.baseDir, filepath.FromSlash(key))
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, "", ErrFileNotFound
		}
		return nil, "", fmt.Errorf("localFileStore.Get: %w", err)
	}
	return f, contentTypeFromKey(key), nil
}

func (s *localFileStore) GetRange(_ context.Context, key, rangeHeader string) (io.ReadCloser, string, string, int64, error) {
	path := filepath.Join(s.baseDir, filepath.FromSlash(key))
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, "", "", 0, ErrFileNotFound
		}
		return nil, "", "", 0, fmt.Errorf("localFileStore.GetRange: %w", err)
	}
	fi, err := f.Stat()
	if err != nil {
		_ = f.Close()
		return nil, "", "", 0, fmt.Errorf("localFileStore.GetRange stat: %w", err)
	}
	ct := contentTypeFromKey(key)
	total := fi.Size()

	if rangeHeader == "" {
		return f, ct, "", total, nil
	}

	start, end, ok := parseByteRange(rangeHeader, total)
	if !ok {
		_ = f.Close()
		return nil, "", "", 0, fmt.Errorf("localFileStore.GetRange: unsatisfiable range %q", rangeHeader)
	}
	if _, err := f.Seek(start, io.SeekStart); err != nil {
		_ = f.Close()
		return nil, "", "", 0, fmt.Errorf("localFileStore.GetRange seek: %w", err)
	}
	length := end - start + 1
	cr := fmt.Sprintf("bytes %d-%d/%d", start, end, total)
	return &limitedReadCloser{Reader: io.LimitReader(f, length), Closer: f}, ct, cr, length, nil
}

// ── R2 implementation ─────────────────────────────────────────────────────────

type r2FileStore struct {
	client *s3.Client
	bucket string
}

func newR2FileStore(cfg config.R2Config) (*r2FileStore, error) {
	if cfg.Endpoint == "" {
		return nil, fmt.Errorf("r2 endpoint is required when backend=r2")
	}
	if cfg.Bucket == "" {
		return nil, fmt.Errorf("r2 bucket is required when backend=r2")
	}

	client := s3.New(s3.Options{
		BaseEndpoint: aws.String(cfg.Endpoint),
		Region:       "auto",
		Credentials:  aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, "")),
	})
	return &r2FileStore{client: client, bucket: cfg.Bucket}, nil
}

func (s *r2FileStore) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		return fmt.Errorf("r2FileStore.Delete: %w", err)
	}
	return nil
}

func (s *r2FileStore) Put(ctx context.Context, key string, r io.Reader, contentType string) error {
	// Buffer the reader so we can provide the content length.
	buf, err := io.ReadAll(r)
	if err != nil {
		return fmt.Errorf("r2FileStore.Put read: %w", err)
	}
	_, err = s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(s.bucket),
		Key:           aws.String(key),
		Body:          bytes.NewReader(buf),
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(int64(len(buf))),
	})
	if err != nil {
		return fmt.Errorf("r2FileStore.Put: %w", err)
	}
	return nil
}

func (s *r2FileStore) Get(ctx context.Context, key string) (io.ReadSeekCloser, string, error) {
	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		var nsk *types.NoSuchKey
		if errors.As(err, &nsk) {
			return nil, "", ErrFileNotFound
		}
		return nil, "", fmt.Errorf("r2FileStore.Get: %w", err)
	}
	defer func() { _ = out.Body.Close() }()
	ct := contentTypeFromKey(key)
	if out.ContentType != nil && *out.ContentType != "" {
		ct = *out.ContentType
	}
	buf, err := io.ReadAll(out.Body)
	if err != nil {
		return nil, "", fmt.Errorf("r2FileStore.Get read: %w", err)
	}
	return &bytesReadSeekCloser{bytes.NewReader(buf)}, ct, nil
}

func (s *r2FileStore) GetRange(ctx context.Context, key, rangeHeader string) (io.ReadCloser, string, string, int64, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}
	if rangeHeader != "" {
		input.Range = aws.String(rangeHeader)
	}
	out, err := s.client.GetObject(ctx, input)
	if err != nil {
		var nsk *types.NoSuchKey
		if errors.As(err, &nsk) {
			return nil, "", "", 0, ErrFileNotFound
		}
		return nil, "", "", 0, fmt.Errorf("r2FileStore.GetRange: %w", err)
	}
	ct := contentTypeFromKey(key)
	if out.ContentType != nil && *out.ContentType != "" {
		ct = *out.ContentType
	}
	var cr string
	if out.ContentRange != nil {
		cr = *out.ContentRange
	}
	var length int64 = -1
	if out.ContentLength != nil {
		length = *out.ContentLength
	}
	// out.Body is streamed directly — caller must close it.
	return out.Body, ct, cr, length, nil
}

// ── Helpers ───────────────────────────────────────────────────────────────────

// bytesReadSeekCloser wraps a *bytes.Reader to satisfy io.ReadSeekCloser.
type bytesReadSeekCloser struct{ *bytes.Reader }

func (bytesReadSeekCloser) Close() error { return nil }

// contentTypeFromKey returns a MIME type based on the file extension in key.
func contentTypeFromKey(key string) string {
	switch filepath.Ext(key) {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".webp":
		return "image/webp"
	case ".gif":
		return "image/gif"
	case ".mp3":
		return "audio/mpeg"
	case ".wav":
		return "audio/wav"
	case ".ogg":
		return "audio/ogg"
	case ".m4a":
		return "audio/mp4"
	case ".mp4":
		return "video/mp4"
	case ".webm":
		return "video/webm"
	case ".mov":
		return "video/quicktime"
	default:
		return "application/octet-stream"
	}
}

// limitedReadCloser pairs a limited reader with the underlying closer.
type limitedReadCloser struct {
	io.Reader
	io.Closer
}

// parseByteRange parses a "bytes=start-end" Range header against totalSize.
// Returns the resolved start and end byte positions (inclusive) and whether the range is satisfiable.
// Only single-range requests are supported; multi-range (comma-separated) returns false.
func parseByteRange(header string, total int64) (start, end int64, ok bool) {
	const prefix = "bytes="
	if !strings.HasPrefix(header, prefix) {
		return 0, 0, false
	}
	spec := header[len(prefix):]
	if strings.Contains(spec, ",") {
		return 0, 0, false // multi-range not supported
	}
	dash := strings.IndexByte(spec, '-')
	if dash < 0 {
		return 0, 0, false
	}
	startStr, endStr := spec[:dash], spec[dash+1:]
	if startStr == "" {
		// Suffix range: bytes=-N
		n, err := strconv.ParseInt(endStr, 10, 64)
		if err != nil || n <= 0 {
			return 0, 0, false
		}
		s := total - n
		if s < 0 {
			s = 0
		}
		return s, total - 1, true
	}
	s, err := strconv.ParseInt(startStr, 10, 64)
	if err != nil || s < 0 || s >= total {
		return 0, 0, false
	}
	if endStr == "" {
		return s, total - 1, true
	}
	e, err := strconv.ParseInt(endStr, 10, 64)
	if err != nil || e < s {
		return 0, 0, false
	}
	if e >= total {
		e = total - 1
	}
	return s, e, true
}
