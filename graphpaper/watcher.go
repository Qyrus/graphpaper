package graphpaper

import (
	"os"
	"fmt"
	"log"
	"path/filepath"
	"time"
	"strings"
)

type visitor int64

func (v visitor) VisitDir(path string, f *os.FileInfo) bool {
	return true
}

func (v visitor) VisitFile(path string, f *os.FileInfo) {
	if f.Mtime_ns > int64(v) {
		parts := strings.Split(path, "/")
		file := parts[len(parts)-1]
		dir := parts[len(parts)-2]
		ext := filepath.Ext(file)

		if ext == ".gpr" {

			metric := file[:(len(file) - 4)]
			node := dir

			rawfile, err := OpenFile(path)
			if err != nil {
				log.Fatalln("fatal: Failed to open file", err)
			}
			defer rawfile.Close()

			list, err := rawfile.ReadRawMeasurements()

			for _, fc := range Config.Resolutions {
				s := Aggregate(list, fc.Resolution, 63)
				b := s.Bucketize(fc.Size)

				for start, summary := range b {
					date := FormatTime(start, true, fc.DateFmt)
					filename := fmt.Sprintf("data/%s/%s/%s/%s.gpr", fc.Name, date, node, metric)
					file, err := CreateOrOpenFile(filename, summary.ValueType, start, fc.Resolution, summary.Functions)
					if err != nil {
						log.Fatalln("fatal: Failed to open file", err)
					}
					defer file.Close()

					file.WriteSummary(s)
				}
			}
		}

	}
}

// TODO: this doesn't check for any files that were written to before it starts
// up.
func Watch(dir string) {
	ticker := time.NewTicker(10 * 1000000000)
	last := int64(0)
	for {
		select {
		case t := <-ticker.C:
			if last > 0 {
				filepath.Walk(dir, visitor(last), nil)
			}
			last = t
		}
	}
}
