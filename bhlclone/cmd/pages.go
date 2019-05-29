// Copyright © 2019 Dmitry Mozzherin <dmozzherin@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"os"

	"github.com/gnames/bhlclone"
	"github.com/spf13/cobra"
)

// pagesCmd represents the pages command
var pagesCmd = &cobra.Command{
	Use:   "pages",
	Short: "Download CSV files with information about pages and names in them",
	Long: `Text of BHL pages will be downloaded into a tree structure at the
current directory. Also csv files with metadata about pages and occurances of
scientific names in them.`,
	Run: func(cmd *cobra.Command, args []string) {
		opts := []bhlclone.Option{
			bhlclone.TitleIds(titleFlags(cmd)),
			bhlclone.WithText(textFlag(cmd)),
		}
		err := bhlclone.Client(bhlclone.Pages(opts...))
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(pagesCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// pagesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	pagesCmd.Flags().BoolP("text", "t", false, "Downloads texts of pages")
	pagesCmd.Flags().IntP("title_start", "s", 0, "Sets first title to download")
	pagesCmd.Flags().IntP("title_end", "e", 0, "Sets last title to download")
}

func titleFlags(cmd *cobra.Command) []int {
	start, err := cmd.Flags().GetInt("title_start")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	end, err := cmd.Flags().GetInt("title_end")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if start < 1 {
		return []int{}
	}
	if end <= start {
		return []int{start}
	}
	var titleIDs []int
	for i := start; i <= end; i++ {
		titleIDs = append(titleIDs, i)
	}
	return titleIDs
}

func textFlag(cmd *cobra.Command) bool {
	withText, err := cmd.Flags().GetBool("text")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return withText
}
