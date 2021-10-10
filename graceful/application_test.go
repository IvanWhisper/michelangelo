package graceful

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestApplication_Run_Std(t *testing.T) {
	app := NewIApplication(WithCmd("go"))
	err := app.Run("env")
	if err != nil {
		t.Error(err)
	}
}

func TestApplication_Run_PrintCh(t *testing.T) {
	app := NewIApplication(WithCmd("go"), WithPrintCh())
	go func() {
		for v := range app.GetPrintCh() {
			print("打印输出：", v)
		}
	}()
	err := app.Run("env")
	if err != nil {
		t.Error(err)
	}
}

func ExampleApplication_Run() {
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
