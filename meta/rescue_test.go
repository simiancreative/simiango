package meta

import (
	"bytes"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/simiancreative/simiango/logger"
)

func TestRescuePanic(t *testing.T) {
	logger.Enable()
	defer RescuePanic("test_id", "test_context")
	panic("Test panic")
}

func TestGinRecovery(t *testing.T) {
	r := gin.Default()
	r.Use(GinRecovery(func(c *gin.Context, m map[string]interface{}) {
		c.JSON(http.StatusOK, m)
	}))
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, world!")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected HTTP response code 200, got: %d", w.Code)
	}
}

func TestRecoverGinPanic(t *testing.T) {
	logger.Enable()

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		defer recoverGinPanic(c, func(c *gin.Context, m map[string]interface{}) {
			c.JSON(http.StatusOK, m)
		})
		panic(errors.New("Test error"))
	})
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, world!")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("Expected HTTP response code 200, got: %d", w.Code)
	}
}

func TestRecoverGinPanic_BrokenPipe(t *testing.T) {
	logger.Enable()

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		defer recoverGinPanic(c, func(c *gin.Context, m map[string]interface{}) {
			c.JSON(http.StatusOK, m)
		})
		panic(&net.OpError{Err: &os.SyscallError{Err: errors.New("broken pipe")}})
	})
	r.GET("/test", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, world!")
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("Expected HTTP response code 500, got: %d", w.Code)
	}
}

func TestStack(t *testing.T) {
	stack := stack()
	if len(stack) == 0 {
		t.Fatalf("Expected stack trace, got empty")
	}
}

func TestStackAsBuf(t *testing.T) {
	buf := stackAsBuf()
	if len(buf) == 0 {
		t.Fatalf("Expected stack trace, got empty")
	}
}

func TestSource(t *testing.T) {
	lines := [][]byte{[]byte("line1"), []byte("line2"), []byte("line3")}
	source := source(lines, 2)
	if !bytes.Equal(source, []byte("line2")) {
		t.Fatalf("Expected 'line2', got '%s'", source)
	}
}

func TestFunction(t *testing.T) {
	fn := function(1)
	if len(fn) == 0 {
		t.Fatalf("Expected function name, got empty")
	}
}
