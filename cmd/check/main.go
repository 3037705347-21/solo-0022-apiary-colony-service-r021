package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"time"

	"example.com/apiary-colony-service/internal/app"
)

var workflows = map[string]func(*client) error{
	"register-apiary":    checkRegisterApiary,
	"install-hive":       checkHiveInstall,
	"record-inspection":  checkRecordInspection,
	"replace-queen":      checkReplaceQueen,
	"advance-hive-state": checkAdvanceHiveState,
	"read-apiary-status": checkReadApiaryStatus,
}

func main() {
	workflow := flag.String("workflow", "", "要校验的 workflow 编号")
	flag.Parse()

	scenario, ok := workflows[*workflow]
	if !ok {
		fmt.Fprintf(os.Stderr, "未知 workflow：%q\n", *workflow)
		os.Exit(2)
	}

	base, stop, err := startServer()
	if err != nil {
		fmt.Fprintf(os.Stderr, "启动服务失败：%v\n", err)
		os.Exit(1)
	}
	defer stop()

	if err := scenario(newClient(base)); err != nil {
		fmt.Fprintf(os.Stderr, "workflow %s 校验失败：%v\n", *workflow, err)
		os.Exit(1)
	}
	fmt.Printf("workflow %s ok\n", *workflow)
}

func startServer() (string, func(), error) {
	application := app.New(app.Config{Addr: "127.0.0.1:0"})
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return "", nil, err
	}
	server := &http.Server{Handler: application.Handler(), ReadHeaderTimeout: 5 * time.Second}
	go func() {
		_ = server.Serve(listener)
	}()
	stop := func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = server.Shutdown(ctx)
	}
	return "http://" + listener.Addr().String(), stop, nil
}
