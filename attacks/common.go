package attacks

import "io"

type countingReadCloser struct {
	rc    io.ReadCloser
	bytes int64
}

func (c *countingReadCloser) Read(p []byte) (int, error) {
	n, err := c.rc.Read(p)
	c.bytes += int64(n)
	return n, err
}

func (c *countingReadCloser) Close() error {
	return c.rc.Close()
}
