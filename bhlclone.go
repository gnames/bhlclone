package bhlclone

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/gnames/bhlindex/protob"
	"google.golang.org/grpc"
)

type funcGRPC func(protob.BHLIndexClient) error

func Client(f funcGRPC) error {
	conn, err := grpc.Dial("bhlrpc.globalnames.org:80", grpc.WithInsecure())
	if err != nil {
		return err
	}
	defer conn.Close()

	client := protob.NewBHLIndexClient(conn)

	err = f(client)
	if err != nil {
		return err
	}

	return nil
}

func VersionGRPC() funcGRPC {
	return func(c protob.BHLIndexClient) error {
		ver, err := c.Ver(context.Background(), &protob.Void{})
		if err != nil {
			return err
		}
		fmt.Println("bhlindex: " + ver.Value)
		fmt.Println("bhlclone: " + Version)
		fmt.Println()
		return nil
	}
}

func Titles() funcGRPC {
	return func(c protob.BHLIndexClient) error {
		fmt.Print(`
This method will download meta-data about BHL titles to titles.csv
in the current directory (~20MB). Is it OK? (Y/n): `)
		goAhead := askForConfirmation(true)
		if !goAhead {
			return nil
		}
		fmt.Println()
		f, err := os.Create("titles.csv")
		if err != nil {
			return err
		}

		w := csv.NewWriter(f)
		stream, err := c.Titles(context.Background(), &protob.TitlesOpt{})
		if err != nil {
			return err
		}
		importNum := 0

		for {
			importNum++
			if importNum%50000 == 0 {
				fmt.Printf("titles received: %d\n", importNum)
			}

			title, err := stream.Recv()
			if err == io.EOF {
				w.Flush()
				fmt.Printf("titles received: %d\n", importNum)
				fmt.Print(`
You can use "id" field to get pages information from titles with
"bhlclone pages".
`)
				return nil
			}
			if err != nil {
				return err
			}
			if importNum == 1 {
				headers := []string{"id", "archive_id", "path"}
				if err := w.Write(headers); err != nil {
					return err
				}
			}
			values := []string{strconv.Itoa(int(title.Id)), title.ArchiveId, title.Path}
			if err := w.Write(values); err != nil {
				return err
			}
		}
	}
}
