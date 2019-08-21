// =================================================================
//
// Copyright (C) 2019 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	yaml "gopkg.in/yaml.v2"

	"github.com/spatialcurrent/go-counter/pkg/counter"
)

const (
	flagSplitBytes string = "bytes"
	flagSplitWords string = "words"
	flagSplitLines string = "lines"
	flagSkipErrors string = "skip-errors"
	flagSort       string = "sort"
	flagNumber     string = "number"
	flagMinimum    string = "minimum"
	flagMaximum    string = "maximum"
	flagCSV        string = "csv"
	flagJSON       string = "json"
	flagYAML       string = "yaml"
)

func initViper(cmd *cobra.Command) (*viper.Viper, error) {
	v := viper.New()
	err := v.BindPFlags(cmd.Flags())
	if err != nil {
		return v, errors.Wrap(err, "error binding flag set to viper")
	}
	v.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	v.AutomaticEnv() // set environment variables to overwrite config
	return v, nil
}

func checkDumpConfig(v *viper.Viper) error {
	s := 0
	if v.GetBool(flagSplitBytes) {
		s++
	}
	if v.GetBool(flagSplitWords) {
		s++
	}
	if v.GetBool(flagSplitLines) {
		s++
	}
	if s == 0 {
		return errors.New("must select one of bytes (-b), words (-w), or lines (-l) to count")
	}
	if s > 1 {
		return errors.New("must select only one of bytes(-b), words (-w), or lines (-l) to count")
	}
	return nil
}

func checkAllConfig(v *viper.Viper) error {
	s := 0
	if v.GetBool(flagSplitBytes) {
		s++
	}
	if v.GetBool(flagSplitWords) {
		s++
	}
	if v.GetBool(flagSplitLines) {
		s++
	}
	if s == 0 {
		return errors.New("must select one of bytes, words, or lines to count")
	}
	if s > 1 {
		return errors.New("must select only one of bytes, words, or lines to count")
	}
	return nil
}

func checkBottomConfig(v *viper.Viper) error {
	number := v.GetInt(flagNumber)
	if number <= 0 {
		return fmt.Errorf("number is %d, expecting v greater than zero", number)
	}

	max := v.GetInt(flagMaximum)
	if max == 0 {
		return fmt.Errorf("maximum is %d, expecting value not equal to zero", max)
	}

	return nil
}

func checkTopConfig(v *viper.Viper) error {
	number := v.GetInt(flagNumber)
	if number <= 0 {
		return fmt.Errorf("number is %d, expecting v greater than zero", number)
	}

	min := v.GetInt(flagMinimum)
	if min < 0 {
		return fmt.Errorf("minimum is %d, expecting value greater than or equal to zero", number)
	}

	return nil
}

func initSplitFlags(flag *pflag.FlagSet) {
	flag.BoolP(flagSplitBytes, "b", false, "count bytes")
	flag.BoolP(flagSplitWords, "w", false, "count words")
	flag.BoolP(flagSplitLines, "l", false, "count lines")
}

func initDumpsFlags(flag *pflag.FlagSet) {
	flag.BoolP(flagCSV, "c", false, "CSV output")
	flag.BoolP(flagJSON, "j", false, "JSON output")
	flag.BoolP(flagYAML, "y", false, "YAML output")
	flag.BoolP(flagSkipErrors, "e", false, "skip errors")
}

func initAllFlags(flag *pflag.FlagSet) {
	flag.BoolP(flagSort, "s", false, "sort values in descending order before the values are chosen")
	flag.BoolP(flagCSV, "c", false, "CSV output")
	flag.BoolP(flagJSON, "j", false, "JSON output")
	flag.BoolP(flagYAML, "y", false, "YAML output")
	flag.BoolP(flagSkipErrors, "e", false, "skip errors")
}

func initBottomFlags(flag *pflag.FlagSet) {
	flag.IntP(flagNumber, "n", 1, "number of values to return")
	flag.IntP(flagMaximum, "m", -1, "maximum value")
	flag.BoolP(flagSort, "s", false, "sort values in descending order before the values are chosen")
	flag.BoolP(flagCSV, "c", false, "CSV output")
	flag.BoolP(flagJSON, "j", false, "JSON output")
	flag.BoolP(flagYAML, "y", false, "YAML output")
	flag.BoolP(flagSkipErrors, "e", false, "skip errors")
}

func initTopFlags(flag *pflag.FlagSet) {
	flag.IntP(flagNumber, "n", 1, "number of values to return")
	flag.IntP(flagMinimum, "m", 0, "minimum value")
	flag.BoolP(flagSort, "s", false, "sort values in descending order before the values are chosen")
	flag.BoolP(flagCSV, "c", false, "CSV output")
	flag.BoolP(flagJSON, "j", false, "JSON output")
	flag.BoolP(flagYAML, "y", false, "YAML output")
	flag.BoolP(flagSkipErrors, "e", false, "skip errors")
}

func printObject(object interface{}, v *viper.Viper) error {
	if v.GetBool(flagJSON) {
		b, err := json.Marshal(object)
		if err != nil {
			return errors.Wrap(err, "error marshaling object to JSON")
		}
		fmt.Println(string(b))
		return nil
	}

	if v.GetBool(flagYAML) {
		b, err := yaml.Marshal(object)
		if err != nil {
			return errors.Wrap(err, "error marshaling object to YAML")
		}
		fmt.Println(string(b) + "\n")
		return nil
	}

	fmt.Println(fmt.Sprintf("%#v", object))
	return nil
}

func dropCarriageReturn(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

func scanLines(separator byte, dropCR bool) bufio.SplitFunc {
	return bufio.SplitFunc(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if i := bytes.IndexByte(data, separator); i >= 0 {
			// We have a full separator-terminated line.
			if dropCR {
				return i + 1, dropCarriageReturn(data[0:i]), nil
			}
			return i + 1, data[0:i], nil
		}
		// If we're at EOF, we have a final, non-terminated line. Return it.
		if atEOF {
			if dropCR {
				return len(data), dropCarriageReturn(data), nil
			}
			return len(data), data, nil
		}
		// Request more data.
		return 0, nil, nil
	})
}

func scanWords(delim func(r rune) bool) bufio.SplitFunc {
	return bufio.SplitFunc(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		// Skip leading spaces.
		start := 0
		for width := 0; start < len(data); start += width {
			var r rune
			r, width = utf8.DecodeRune(data[start:])
			if !delim(r) {
				break
			}
		}
		// Scan until space, marking end of word.
		for width, i := 0, start; i < len(data); i += width {
			var r rune
			r, width = utf8.DecodeRune(data[i:])
			if delim(r) {
				return i + width, data[start:i], nil
			}
		}
		// If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
		if atEOF && len(data) > start {
			return len(data), data[start:], nil
		}
		// Request more data.
		return start, nil, nil
	})
}

func countFiles(args []string, v *viper.Viper) (counter.Counter, error) {
	splitFunc := bufio.ScanBytes

	if v.GetBool(flagSplitWords) {
		splitFunc = scanWords(func(r rune) bool {
			if unicode.IsSpace(r) {
				return true
			}
			switch r {
			case '!', ',', ';', '(', ')', '{', '}', '[', ']', '<', '>', '`', '|', '=':
				return true
			}
			return false
		})
	}

	if v.GetBool(flagSplitLines) {
		splitFunc = scanLines('\n', true)
	}

	words, err := counter.CountFiles(args, splitFunc, v.GetBool(flagSkipErrors))
	if err != nil {
		return nil, errors.Wrap(err, "error counting words")
	}

	return words, nil
}

func main() {

	root := &cobra.Command{
		Use:                   "gocounter",
		Short:                 "gocounter",
		DisableFlagsInUseLine: true,
		Long:                  `gocounter is a super simple utility for counting words in files`,
	}

	dump := &cobra.Command{
		Use:                   "dump [-|stdin|FILE]...",
		Short:                 "dump",
		DisableFlagsInUseLine: true,
		Long:                  `dump frequency distribution of words`,
		RunE: func(cmd *cobra.Command, args []string) error {
			v, err := initViper(cmd)
			if err != nil {
				return errors.Wrap(err, "error initializing viper")
			}

			if len(args) == 0 {
				return cmd.Usage()
			}

			if errConfig := checkDumpConfig(v); errConfig != nil {
				return errConfig
			}

			words, err := countFiles(args, v)
			if err != nil {
				return errors.Wrap(err, "error counting words")
			}

			err = printObject(words, v)
			if err != nil {
				return errors.Wrap(err, "error printing values")
			}

			return nil
		},
	}
	initSplitFlags(dump.Flags())
	initDumpsFlags(dump.Flags())
	root.AddCommand(dump)

	all := &cobra.Command{
		Use:                   "all [-|stdin|FILE]...",
		Short:                 "all",
		DisableFlagsInUseLine: true,
		Long:                  `all the words`,
		RunE: func(cmd *cobra.Command, args []string) error {
			v, err := initViper(cmd)
			if err != nil {
				return errors.Wrap(err, "error initializing viper")
			}

			if len(args) == 0 {
				return cmd.Usage()
			}

			if errConfig := checkAllConfig(v); errConfig != nil {
				return errConfig
			}

			words, err := countFiles(args, v)
			if err != nil {
				return errors.Wrap(err, "error counting words")
			}

			values := words.All(v.GetBool(flagSort))

			err = printObject(values, v)
			if err != nil {
				return errors.Wrap(err, "error printing values")
			}

			return nil
		},
	}
	initSplitFlags(all.Flags())
	initAllFlags(all.Flags())
	root.AddCommand(all)

	bottom := &cobra.Command{
		Use:                   "bottom [-n N] [-m MAX] [-s] [-|stdin|FILE]...",
		Short:                 "bottom",
		DisableFlagsInUseLine: true,
		Long:                  `bottom returns the least frequent words`,
		RunE: func(cmd *cobra.Command, args []string) error {
			v, err := initViper(cmd)
			if err != nil {
				return errors.Wrap(err, "error initializing viper")
			}

			if len(args) == 0 {
				return cmd.Usage()
			}

			if errConfig := checkBottomConfig(v); errConfig != nil {
				return errConfig
			}

			words, err := countFiles(args, v)
			if err != nil {
				return errors.Wrap(err, "error counting words")
			}

			values := words.Bottom(v.GetInt(flagNumber), v.GetInt(flagMaximum), v.GetBool(flagSort))

			err = printObject(values, v)
			if err != nil {
				return errors.Wrap(err, "error printing values")
			}

			return nil
		},
	}
	initSplitFlags(bottom.Flags())
	initBottomFlags(bottom.Flags())
	root.AddCommand(bottom)

	top := &cobra.Command{
		Use:                   "top [-n N] [-m MIN] [-s] [-|stdin|FILE]...",
		Short:                 "top",
		DisableFlagsInUseLine: true,
		Long:                  `top returns the most frequent words`,
		RunE: func(cmd *cobra.Command, args []string) error {
			v, err := initViper(cmd)
			if err != nil {
				return errors.Wrap(err, "error initializing viper")
			}

			if len(args) == 0 {
				return cmd.Usage()
			}

			if errConfig := checkTopConfig(v); errConfig != nil {
				return errConfig
			}

			words, err := countFiles(args, v)
			if err != nil {
				return errors.Wrap(err, "error counting words")
			}

			values := words.Top(v.GetInt(flagNumber), v.GetInt(flagMaximum), v.GetBool(flagSort))

			err = printObject(values, v)
			if err != nil {
				return errors.Wrap(err, "error printing values")
			}

			return nil
		},
	}
	initSplitFlags(top.Flags())
	initTopFlags(top.Flags())
	root.AddCommand(top)

	if err := root.Execute(); err != nil {
		panic(err)
	}
}
