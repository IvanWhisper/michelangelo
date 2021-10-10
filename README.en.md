# michelangelo

#### Description
`michelangelo` will keep iterating.
It is a toolkit.

#### Components
##### Graceful
Package `graceful` run a controlled application
``` go
	app := NewIApplication(
		WithCmd("go"),
		WithPrintCh(),
		WithName("go app"),
		WithContext(context.Background()),
		WithWorkPath("./"),
		WithTimeOut(5*time.Second))
	err := app.Run("env")
	if err != nil {
		fmt.Print(err)
	}
```

##### Log
Package `log` record log
``` go
	New(nil)
	log.Print(GetLevel())
	Debug("debug")
	Info("info")
	Error("error")
```
