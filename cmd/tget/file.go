// This file is explicitly for selecting files to be downloaded after adding a torrent
package main

import (
	"fmt"
	"log"
	"path"
	"sort"
	"strconv"
	"strings"

	"github.com/anacrolix/torrent"
	"github.com/pterm/pterm"
)

// takes a slice of indexes for files to be added
func SelectFileByIndex(t *torrent.Torrent, c *torrent.Client, idxs []int) error {
	files := t.Files()
	var chosen []*torrent.File

	for _, i := range idxs {
		if i > len(files) {
			return fmt.Errorf("file index out of range: index [%d] of file count [%d]\n", i, len(files))
		}

		x := files[i]
		x.Download()

		chosen = append(chosen, x)
	}

	FileProgress(chosen)
	return nil
}

func FileProgress(files []*torrent.File) {
	// Create a multi printer for managing multiple printers
	// if len(files) > 15 {
	// }

	// p, _ := pterm.DefaultProgressbar.WithTotal(100).WithTitle("Downloading stuff").Start()
	//
	// for i := 0; i < p.Total; i++ {
	// 	pc := FilePercent(files[i])
	// 	p.UpdateTitle("Downloading " + path.Base(files[i].Path()))
	// 	p.Increment().Current = int(pc)
	// }

	multi := pterm.DefaultMultiPrinter
	var pbs []*pterm.ProgressbarPrinter

	for _, f := range files {
		name := fmt.Sprintf("[DL] %s\n", path.Base(f.Path()))
		pb, err := pterm.DefaultProgressbar.WithTotal(100).WithWriter(multi.NewWriter()).Start(name)
		if err != nil {
			log.Println(err)
		}

		pbs = append(pbs, pb)
	}

	multi.Start()
	for !filesComplete(files) {

		for i, f := range files {
			prog := FilePercent(f)
			pbs[i].Increment().Current = int(prog)
		}
	}
}

func filesComplete(files []*torrent.File) bool {
	for _, f := range files {
		prog := FilePercent(f)
		if prog < 100 {
			return false
		}
	}
	return true
}

// Return progress of individual files in percentage
func FilePercent(file *torrent.File) float64 {
	ps := file.State()
	full := len(ps)
	var complete int

	for _, state := range ps {
		if state.Complete {
			complete++
		}
	}

	// pc := (float64(full) / float64(complete)) * 100
	pc := (float64(complete) / float64(full)) * 100
	log.Printf("file: %s %v", file.Path(), int(pc))
	return pc
	// return float64(full) / float64(complete) * 100
}

func removeZeroFiles(files []*torrent.File, idxs []int, t *torrent.Torrent) error {
	for _, i := range idxs {
		if i > len(files) {
			return fmt.Errorf("file index out of range: index [%d] of file count [%d]\n", i, len(files))
		}
	}
	return nil
}

// remove duplicate entries from slice
func removeDuplicate[T comparable](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

// takes aria2c like strings ex: --select-file=0,1,3,30-33
func ParseIndices(s string) ([]int, error) {
	var idxs []int
	cuts := strings.Split(s, ",")

	if len(cuts) <= 0 {
		return idxs, fmt.Errorf("unable to parse file index string, length is 0")
	}

	for _, i := range cuts {
		if i == "" {
			return idxs, fmt.Errorf("error: file index is an empty string")
		}

		if strings.Contains(i, "-") {
			splits := strings.Split(i, "-")

			if len(splits) != 2 {
				return idxs, fmt.Errorf("index file range can only contain two numbers")
			}

			if splits[0] > splits[1] {
				return idxs, fmt.Errorf("index file range must be in format [1-99] where N1 is less than N2, got: [%s]", i)
			}

			n1, err := strconv.Atoi(splits[0])
			if err != nil {
				return idxs, err
			}

			n2, err := strconv.Atoi(splits[1])
			if err != nil {
				return idxs, err
			}

			if n1 < 0 || n2 < 0 {
				return idxs, fmt.Errorf("file index cannot be negative, got range [%d] [%d]", n1, n2)
			}

			// TODO: check if this range is inclusive or not
			for x := n1; x <= n2; x++ {
				idxs = append(idxs, x)
			}
		} else {
			n, err := strconv.Atoi(i)
			if err != nil {
				return idxs, err
			}

			if n < 0 {
				return idxs, fmt.Errorf("file index cannot be negative, got [%d]", n)
			}

			idxs = append(idxs, n)
		}
	}

	if len(idxs) <= 0 {
		return idxs, fmt.Errorf("unable to parse file indices")
	}

	// sort by value
	sort.Ints(idxs)
	idxs = removeDuplicate(idxs)

	return idxs, nil
}
