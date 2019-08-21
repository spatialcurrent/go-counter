// =================================================================
//
// Copyright (C) 2018 Spatial Current, Inc. - All Rights Reserved
// Released as open source under the MIT License.  See LICENSE file.
//
// =================================================================

package counter

import (
	"bufio"
	"os"

	"github.com/pkg/errors"
)

// CountFiles generates a frequency distribution for the tokens found in the files
// files is a list of paths to files.
// splitFunc is a bufio.SplitFunc, which can be bufio.ScanBytes, bufio.ScanWords, bufio.ScanLines, or a custom function.
// If skipErrors is set to true, then errors opening or reading files are skipped and not returned.
func CountFiles(files []string, splitFunc bufio.SplitFunc, skipErrors bool) (Counter, error) {

	tokens := New()

	for _, f := range files {
		if f == "-" || f == "stdin" {
			scanner := bufio.NewScanner(bufio.NewReader(os.Stdin))
			scanner.Split(splitFunc)
			for scanner.Scan() {
				tokens.Increment(scanner.Text())
			}
			if err := scanner.Err(); err != nil {
				if !skipErrors {
					return nil, errors.Wrap(err, "error scanning for tokens")
				}
			}
		} else {
			r, err := os.Open(f)
			if err != nil {
				if skipErrors {
					continue
				}
				return nil, errors.Wrapf(err, "error opening file %q for reading", f)
			}
			defer r.Close()
			fi, err := r.Stat()
			if err != nil {
				if skipErrors {
					continue
				}
				return nil, errors.Wrapf(err, "error stating file %q for reading", f)
			}
			if fi.Mode().IsRegular() {
				scanner := bufio.NewScanner(bufio.NewReader(r))
				scanner.Split(splitFunc)
				for scanner.Scan() {
					tokens.Increment(scanner.Text())
				}
				if err := scanner.Err(); err != nil {
					if !skipErrors {
						return nil, errors.Wrap(err, "error scanning for tokens")
					}
				}
				r.Close()
			}
		}
	}

	return tokens, nil
}
