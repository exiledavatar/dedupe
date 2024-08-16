package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/urfave/cli/v2"
	"golang.org/x/exp/slices"
)

func main() {

	var kvSeparator, iSeparator, vSeparator, itype, out string
	var trimWhitespace bool
	items := []string{}
	dedupedItems := []string{}

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

			items = strings.Split(input, iSeparator)
			// fmt.Println(len(items))
			if trimWhitespace {
				for i, v := range items {
					items[i] = strings.TrimSpace(v)
				}
			}
			// fmt.Println(len(items))

			switch {
			case slices.Contains([]string{"item", "i"}, itype):
				dedupedItems = dedupeItems(items)
				out = strings.Join(dedupedItems, iSeparator)
			case slices.Contains([]string{"key", "k"}, itype):
				check := map[string]bool{}
				for _, v := range items {
					key, _, found := strings.Cut(v, kvSeparator)
					if trimWhitespace {
						key = strings.TrimSpace(key)
					}

					switch _, ok := check[key]; {
					case !found:
						dedupedItems = append(dedupedItems, v)
					case !ok:
						check[key] = true
						dedupedItems = append(dedupedItems, v)
					}
				}
				out = strings.Join(dedupedItems, iSeparator)

			case slices.Contains([]string{"key-value", "kv"}, itype):
				// this dedupe's every value set and every key=value(s) set

				check := map[string]string{}
				for _, item := range items {
					key, value, found := strings.Cut(item, kvSeparator)
					if found {
						values := strings.Split(value, vSeparator)
						if trimWhitespace {
							item = strings.TrimSpace(item)
							key = strings.TrimSpace(key)
							for i, v := range values {
								values[i] = strings.TrimSpace(v)
							}
						}
						values = dedupeItems(values)
						value = strings.Join(values, vSeparator)
					}

					switch v, ok := check[key]; {
					case !found:
						// fmt.Println("not found")
						dedupedItems = append(dedupedItems, item)
					case ok && v == value:
						// fmt.Println("duplicate")
						// this is a dupliace - do nothing
					default:
						// everything else is either a new key or new key=value pair
						// fmt.Println("new")
						check[key] = value
						item = fmt.Sprint(key, kvSeparator, value)
						dedupedItems = append(dedupedItems, item)
					}
				}
				// fmt.Println(len(dedupedItems))
				out = strings.Join(dedupedItems, iSeparator)

			default:
				out = "Something's not right - check your arguments"
			}

			out = string(regexp.MustCompile(`(?s)(.*)\n*$`).ReplaceAll([]byte(out), []byte("$1\n")))

			fmt.Print(out)
			return nil
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "key-value-separator",
				Aliases:     []string{"kv"},
				Value:       "=",
				Usage:       "separator between key=value",
				Destination: &kvSeparator,
			},
			&cli.StringFlag{
				Name:        "item-separator",
				Aliases:     []string{"i"},
				Value:       "\n",
				Usage:       "separator between each line or key=value set",
				Destination: &iSeparator,
			},
			&cli.StringFlag{
				Name:        "value-separator",
				Aliases:     []string{"v"},
				Value:       ",",
				Usage:       "separator between each value within a key=value0,value2...",
				Destination: &vSeparator,
			},
			&cli.StringFlag{
				Name:        "type",
				Aliases:     []string{"t"},
				Value:       "key-value",
				Usage:       "specify key-value, key, or item (alternatively kv, k, i)",
				Destination: &itype,
			},
			&cli.BoolFlag{
				Name:        "trim-whitespace",
				Aliases:     []string{"trim", "tw"},
				Value:       true,
				Usage:       "if true, dedupe will trim and ignore leading/trailing whitespace after parsing, but before comparing.",
				Destination: &trimWhitespace,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}

}

func dedupeItems(values []string) []string {
	check := map[string]bool{}
	deduped := []string{}
	for _, v := range values {
		if _, ok := check[v]; !ok && v != "" {
			check[v] = true
			deduped = append(deduped, v)
		}
	}
	return deduped
}

// func dedupeString(value, separator string, trimWhiteSpace bool) string {
// 	parsed := strings.Split(value, separator)
// 	check := map[string]bool{}
// 	deduped := []string{}
// 	for _, v := range parsed {
// 		if trimWhiteSpace {
// 			v = strings.TrimSpace(v)
// 		}
// 		if _, ok := check[v]; !ok {
// 			check[v] = true
// 			deduped = append(deduped, v)
// 		}
// 	}
// 	return strings.Join(deduped, separator)

// }
