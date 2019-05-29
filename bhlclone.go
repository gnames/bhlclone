package bhlclone

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gnames/bhlindex/protob"
	"google.golang.org/grpc"
)

type funcGRPC func(protob.BHLIndexClient) error

type Option func(*BHLClone)

func TitleIds(titles []int) Option {
	return func(p *BHLClone) {
		p.Titles = titles
	}
}

func WithText(b bool) Option {
	return func(p *BHLClone) {
		p.Texts = b
	}
}

type BHLClone struct {
	Titles []int
	Texts  bool
}

func NewBHLClone(opts ...Option) *BHLClone {
	b := &BHLClone{}
	for _, opt := range opts {
		opt(b)
	}
	return b
}

func (b *BHLClone) toProtoOpts() *protob.PagesOpt {
	titleIds := make([]int32, len(b.Titles))
	for i, v := range b.Titles {
		titleIds[i] = int32(v)
	}
	return &protob.PagesOpt{
		TitleIds: titleIds,
		WithText: b.Texts,
	}
}

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

func Pages(opts ...Option) funcGRPC {
	b := NewBHLClone(opts...)
	return func(c protob.BHLIndexClient) error {
		var currentTitle string
		protoOpts := b.toProtoOpts()
		fp, err := os.Create("pages.csv")
		if err != nil {
			return err
		}
		fn, err := os.Create("names.csv")
		if err != nil {
			return err
		}
		if b.Texts {
			err = os.Mkdir("texts", os.ModePerm)
			if err != nil {
				return err
			}
		}
		wp := csv.NewWriter(fp)
		wn := csv.NewWriter(fn)

		stream, err := c.Pages(context.Background(), protoOpts)
		if err != nil {
			return err
		}
		importNum := 0

		for {
			importNum++
			if importNum%50000 == 0 {
				fmt.Printf("pages received: %d\n", importNum)
			}

			page, err := stream.Recv()
			if err == io.EOF {
				wp.Flush()
				fmt.Printf("Total pages received: %d\n", importNum)
				fmt.Print(`
You can find information about pages at pages.csv, and information about names
occurences at names.csv
`)
				return nil
			}
			if err != nil {
				return err
			}
			if importNum == 1 {
				headers := []string{
					"title_id",
					"page_id",
					"title_path",
					"char_offset",
				}
				if err := wp.Write(headers); err != nil {
					return err
				}
				headers = []string{
					"title_id",
					"page_id",
					"name",
					"odds",
					"path",
					"curated",
					"edit_distance",
					"edit_distance_stem",
					"source_id",
					"match",
					"offset_start",
					"offset_end",
				}
				if err := wn.Write(headers); err != nil {
					return err
				}
				currentTitle = page.TitleId
				if err = os.Mkdir("texts/"+currentTitle, os.ModePerm); err != nil {
					return err
				}

			}
			values := []string{
				page.TitleId,
				page.Id,
				page.TitlePath,
				strconv.Itoa(int(page.Offset)),
			}
			if err := wp.Write(values); err != nil {
				return err
			}
			for _, v := range page.Names {
				values := []string{
					page.TitleId,
					page.Id,
					v.Value,
					strconv.FormatFloat(float64(v.Odds), 'f', 4, 32),
					v.Path,
					strconv.FormatBool(v.Curated),
					strconv.Itoa(int(v.EditDistance)),
					strconv.Itoa(int(v.EditDistanceStem)),
					strconv.Itoa(int(v.SourceId)),
					v.Match.String(),
					string(v.OffsetStart),
					string(v.OffsetEnd),
				}
				if err := wn.Write(values); err != nil {
					return err
				}
				if b.Texts {
					if currentTitle != page.TitleId {
						currentTitle = page.TitleId
						err = os.Mkdir("texts/"+currentTitle, os.ModePerm)
						if err != nil {
							return err
						}
					}
					pageFile := filepath.Join("texts", currentTitle, page.Id+".txt")
					err = ioutil.WriteFile(pageFile, page.Text, 0664)
					if err != nil {
						return err
					}
				}
			}
		}
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
				fmt.Printf("Total titles received: %d\n", importNum)
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
