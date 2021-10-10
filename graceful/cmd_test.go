package graceful

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestNewIApplication(t *testing.T) {
	app := NewIApplication(WithCmd("go"), WithPrintCh())
	err := app.Run("env")
	if err != nil {
		t.Error(err)
	}
}

func ExampleNewIApplication() {
	app := NewIApplication(
		WithCmd("go"),
		WithPrintCh(),
		WithName("go应用"),
		WithContext(context.Background()),
		WithWorkPath("./"),
		WithTimeOut(5*time.Second))
	err := app.Run("env")
	if err != nil {
		fmt.Print(err)
	}
}
