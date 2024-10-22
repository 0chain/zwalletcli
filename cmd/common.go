package cmd

import (
	"fmt"
	"github.com/spf13/pflag"
	"log"
	"os"
	"sort"
	"strings"
)

func PrintError(v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
}

func ExitWithError(v ...interface{}) {
	fmt.Fprintln(os.Stderr, v...)
	os.Exit(1)
}

func setupInputMap(flags *pflag.FlagSet, sKeys, sValues string) map[string]string {
	var err error
	var keys []string
	if flags.Changed(sKeys) {
		keys, err = flags.GetStringSlice(sKeys)
		if err != nil {
			log.Fatal(err)
		}
	}

	var values []string
	if flags.Changed(sValues) {
		values, err = flags.GetStringSlice(sValues)
		if err != nil {
			log.Fatal(err)
		}
	}

	input := make(map[string]string)
	if len(keys) != len(values) {
		log.Fatal("number " + sKeys + " must equal the number " + sValues)
	}
	for i := 0; i < len(keys); i++ {
		v := strings.TrimSpace(values[i])
		k := strings.TrimSpace(keys[i])
		input[k] = v
	}
	return input
}

func printMap(outMap map[string]string) {
	keys := make([]string, 0, len(outMap))
	for k := range outMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Println(k, "\t", outMap[k])
	}
}
