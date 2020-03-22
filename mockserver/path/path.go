// Package path implements a path rewriter gin middleware.
// When url contains ?rewrite=1, it rewrites all http://...zhenai.com urls within body into /mock/....
// This is useful for human browsing without hitting real 3rd-party server.
package path

import (
	"bytes"
	"log"
	"net/http"
	"regexp"

	"github.com/gin-gonic/gin"
)

var re = regexp.MustCompile(`"http://(.*zhenai.com/[^"]*)"`)

// Rewrite implements the rewriter.
func Rewrite(c *gin.Context) {
	params := struct {
		Rewrite bool `form:"rewrite"`
	}{}
	err := c.Bind(&params)

	if err != nil || !params.Rewrite {
		if err != nil {
			log.Printf("Error binding params: %v", err)
		}
		c.Next()
		return
	}

	rw := &responseWriter{
		ResponseWriter: c.Writer,
	}
	c.Writer = rw
	c.Next()

	_, err = rw.writeAndFlush(re.ReplaceAll(rw.buffer.Bytes(), []byte("/mock/$1?rewrite=1")))
	if err != nil {
		log.Printf("error writing re-written data to real writer: %v", err)
		rw.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		rw.ResponseWriter.Flush()
	}
}

// responseWriter buffers all writes for future rewrite.
type responseWriter struct {
	// Embed a gin.ResponseWriter.
	// Once set responseWriter automatically implements interface gin.ResponseWriter.
	gin.ResponseWriter

	// buffer is the buffer we save all writes temporarily.
	buffer bytes.Buffer
}

func (w *responseWriter) Write(data []byte) (int, error) {
	return w.buffer.Write(data)
}

func (w *responseWriter) Flush() {
	// no-op. Must flush via flush() once all bytes are received and re-written.
}

func (w *responseWriter) writeAndFlush(data []byte) (int, error) {
	// Some handlers may have set Content-Length.
	// This must be unset because our length has changed in rewrite.
	w.Header().Del("Content-Length")

	n, err := w.ResponseWriter.Write(data)
	if err == nil {
		w.ResponseWriter.Flush()
	}
	return n, err
}
