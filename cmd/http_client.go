package cmd

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/cockroachdb/errors"
)

func FetchJSON(ctx context.Context, url string) (map[string]any, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	conn := http.DefaultClient
	req.Header.Set("User-Agent", "tcpmon/0.1.0")

	rsp, err := conn.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if rsp.Body != nil {
		defer rsp.Body.Close()
	}

	buf, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	val := make(map[string]any)
	err = json.Unmarshal(buf, &val)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return val, nil
}
