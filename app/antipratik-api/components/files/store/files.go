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

func (s *localFileStore) GetRange(_ context.Context, key string, r *files.ParsedRange) (io.ReadCloser, string, string, int64, error) {
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

	if r == nil {
		return f, ct, "", total, nil
	}

	start, end, ok := resolveRange(r, total)
	if !ok {
		_ = f.Close()
		return nil, "", "", 0, fmt.Errorf("localFileStore.GetRange: unsatisfiable range")
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

func (s *r2FileStore) GetRange(ctx context.Context, key string, r *files.ParsedRange) (io.ReadCloser, string, string, int64, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(key),
	}
	if r != nil {
		input.Range = aws.String(parsedRangeToS3Header(r))
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

// limitedReadCloser pairs a limited reader with the underlying closer.
type limitedReadCloser struct {
	io.Reader
	io.Closer
}
