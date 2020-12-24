package main

import (
	"context"
	"fmt"
	clienthttp "github.com/piusalfred/kitsvc/svc/client/http"
)

func main()  {
	instance := "localhost:9021"
	client,err := clienthttp.NewClient(instance)
	if err != nil {
		panic(err)
	}

	v,err := client.Version(context.Background())
	up, err := client.Uppercase(context.Background(),"hello world")

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(v)
	fmt.Println(up)
}


