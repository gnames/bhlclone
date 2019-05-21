package bhlclone

import (
	"context"
	"fmt"

	"github.com/gnames/bhlindex/protob"
	"google.golang.org/grpc"
)

func Client() error {
	conn, err := grpc.Dial("bhlrpc.globalnames.org:80", grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := protob.NewBHLIndexClient(conn)

	ver, err := client.Ver(context.Background(), &protob.Void{})
	if err != nil {
		return err
	}

	fmt.Println(ver.Value)
	return nil
}
