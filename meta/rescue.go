package meta

import (
	"bytes"
	"errors"
	"fmt"
	"net"
	"net/http/httputil"
	"os"
	"runtime"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/simiancreative/simiango/logger"
)

var (
	dunno     = []byte("???")
	centerDot = []byte("·")
	dot       = []byte(".")
	slash     = []byte("/")
)

func RescuePanic(id RequestId, context interface{}) {
	r := recover()
	if r == nil {
		return
	}

	logger.Error(
		fmt.Sprintf("(%v) caused panic and was recovered", r),
		logger.Fields{
			"request_id": string(id),
			"stack":      string(stackAsBuf()),
			"context":    fmt.Sprintf("%+v", context),
		},
	)
}

func GinRecovery(buildResp func(*gin.Context, map[string]interface{})) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer recoverGinPanic(c, buildResp)
		c.Next()
	}
}

func recoverGinPanic(c *gin.Context, buildResp func(*gin.Context, map[string]interface{})) {
	if err := recover(); err != nil {
		// Check for a broken connection, as it is not really a
		// condition that warrants a panic stack trace.
		var brokenPipe bool
		if ne, ok := err.(*net.OpError); ok {
			var se *os.SyscallError
			if errors.As(ne, &se) {
				seStr := strings.ToLower(se.Error())
				if strings.Contains(seStr, "broken pipe") ||
					strings.Contains(seStr, "connection reset by peer") {
					brokenPipe = true
				}
			}
		}

		id, _ := c.Get("request_id")
		stack := stack()
		httpRequest, _ := httputil.DumpRequest(c.Request, false)
		headers := strings.Split(string(httpRequest), "\r\n")
		for idx, header := range headers {
			current := strings.Split(header, ":")
			if current[0] == "Authorization" {
				headers[idx] = current[0] + ": *"
			}
		}
		logger.Error(
			err.(error).Error(),
			logger.Fields{
				"request_id": id,
				"request":    headers,
				"stack":      stack,
			},
		)

		if brokenPipe {
			// If the connection is dead, we can't write a status to it.
			c.Error(err.(error)) //nolint: errcheck
			c.Abort()
			return
		}

		buildResp(c, map[string]interface{}{
			"request_id": id,
			"headers":    headers,
			"stack":      stack,
		})
	}
}

func stack() []string {
	skip := 3
	result := []string{}
	// As we loop, we open files and read them. These variables record the currently
	// loaded file.
	var lines [][]byte
	var lastFile string
	for i := skip; ; i++ { // Skip the expected number of frames
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		// Print this much at least.  If we can't find the source, it won't show.
		result = append(result, fmt.Sprintf("%s:%d (0x%x)", file, line, pc))
		if file != lastFile {
			data, err := os.ReadFile(file)
			if err != nil {
				continue
			}
			lines = bytes.Split(data, []byte{'\n'})
			lastFile = file
		}
		result = append(
			result,
			fmt.Sprintf("%s: %s", function(pc), source(lines, line)),
		)
	}

	return result
}

func stackAsBuf() []byte {
	buf := new(bytes.Buffer) // the returned data
	lines := stack()

	for _, line := range lines {
		fmt.Fprint(buf, line)
	}

	return buf.Bytes()
}

// source returns a space-trimmed slice of the n'th line.
func source(lines [][]byte, n int) []byte {
	n-- // in stack trace, lines are 1-indexed but our array is 0-indexed
	if n < 0 || n >= len(lines) {
		return dunno
	}
	return bytes.TrimSpace(lines[n])
}

// function returns, if possible, the name of the function containing the PC.
func function(pc uintptr) []byte {
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return dunno
	}
	name := []byte(fn.Name())
	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//	runtime/debug.*T·ptrmethod
	// and want
	//	*T.ptrmethod
	// Also the package path might contain dot (e.g. code.google.com/...),
	// so first eliminate the path prefix
	if lastSlash := bytes.LastIndex(name, slash); lastSlash >= 0 {
		name = name[lastSlash+1:]
	}
	if period := bytes.Index(name, dot); period >= 0 {
		name = name[period+1:]
	}
	name = bytes.Replace(name, centerDot, dot, -1)
	return name
}
