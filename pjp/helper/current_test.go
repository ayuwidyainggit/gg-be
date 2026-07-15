package helper

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetParentCustomerId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	ctx, _ := gin.CreateTestContext(httptest.NewRecorder())

	if _, ok := GetParentCustomerId(ctx); ok {
		t.Fatal("expected missing parent customer id to return false")
	}

	ctx.Set("parentCustomerId", "C220010001")
	got, ok := GetParentCustomerId(ctx)
	if !ok {
		t.Fatal("expected parent customer id to exist")
	}
	if got != "C220010001" {
		t.Fatalf("got %q, want %q", got, "C220010001")
	}
}
