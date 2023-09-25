package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slices"
)

func main() {

	var separator string

	app := &cli.App{
		Action: func(cCtx *cli.Context) error {
			input := cCtx.Args().First()

			if slices.Contains([]string{"", "-"}, input) {
				b, err := ioutil.ReadAll(os.Stdin)
				if err != nil {
					return err
				}
				input = string(b)
			}

			fmt.Print(dedupe(input, separator))
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "separator",
				Aliases:     []string{"s"},
				Value:       ":",
				Destination: &separator,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}

func dedupe(value, separator string) string {
	parsed := strings.Split(value, separator)
	check := map[string]bool{}
	deduped := []string{}
	for _, v := range parsed {
		if _, ok := check[v]; !ok {
			check[v] = true
			deduped = append(deduped, v)
		}
	}
	return strings.Join(deduped, separator)
}
