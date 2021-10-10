# michelangelo

#### 介绍
`michelangelo` 将持续更新.
该项目是一个工具箱.

#### 组件
##### 优雅启动
`graceful` 可以开启一个可控的应用
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

##### 日志
`log` 记录日志
``` go
	New(nil)
	log.Print(GetLevel())
	Debug("debug")
	Info("info")
	Error("error")
```

