package httputil

import (
	"context"
	"fmt"
	"io"
	"net/http"

	ioutil2 "github.com/tokamak-network/tokamak-thanos/op-service/ioutil"
)

// Downloader downloads files via HTTP with optional progress tracking.
type Downloader struct {
	Progressor ioutil2.Progressor
	MaxSize    int64
}

// Download fetches the given URL and writes the response body to the given writer.
func (d *Downloader) Download(ctx context.Context, url string, dest io.Writer) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("downloading %s: %w", url, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d for %s", resp.StatusCode, url)
	}

	var reader io.Reader = resp.Body
	if d.Progressor != nil && resp.ContentLength > 0 {
		reader = &progressReader{r: resp.Body, total: resp.ContentLength, prog: d.Progressor}
	}

	if _, err := io.Copy(dest, reader); err != nil {
		return fmt.Errorf("writing to file: %w", err)
	}
	return nil
}

type progressReader struct {
	r       io.Reader
	current int64
	total   int64
	prog    ioutil2.Progressor
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.r.Read(p)
	pr.current += int64(n)
	pr.prog(pr.current, pr.total)
	return n, err
}
