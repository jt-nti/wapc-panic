package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"runtime/debug"

	"github.com/wapc/wapc-go"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: hello <name>")
		return
	}
	name := os.Args[1]
	ctx := context.Background()
	code, err := ioutil.ReadFile("testdata/hello/target/wasm32-unknown-unknown/debug/hello.wasm")
	if err != nil {
		panic(err)
	}

	module, err := wapc.New(consoleLog, code, hostCall)
	if err != nil {
		panic(err)
	}
	defer module.Close()

	instance, err := module.Instantiate()
	if err != nil {
		panic(err)
	}
	defer instance.Close()

	result, err := instance.Invoke(ctx, "hello", []byte(name))
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s\n", string(result))
}

func consoleLog(msg string) {
	fmt.Println(msg)
}

func hostCall(ctx context.Context, binding, namespace, operation string, payload []byte) (result []byte, err error) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("[Recovering from panic: %v \nStack: %s \n", r, string(debug.Stack()))
			err = fmt.Errorf("Operation panicked: %s %s %s", binding, namespace, operation)
		}
	}()

	// Route the payload to any custom functionality accordingly.
	// You can even route to other waPC modules!!!
	switch namespace {
	case "hello":
		switch operation {
		case "echo":
			fmt.Printf("hostCall echo: %s\n", string(payload))
			if string(payload) == "panic" {
				fmt.Println("Panic!")
				panic("Doomed!")
			}
			return payload, nil // echo
		}
	}
	return []byte("default"), nil
}
