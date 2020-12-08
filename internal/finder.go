package internal

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/ktr0731/go-fuzzyfinder"
)

type Finder struct {
	Depth int
}

type dirInfo struct {
	path  string
	name  string
	depth int
}

func (f *Finder) Run(dir string) {
	chDir := make(chan dirInfo)
	chWalk := make(chan dirInfo, 2_000_000)
	var dirs []string

	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		go func() {
			for {
				select {
				case d := <-chDir:
					// fmt.Println(d.path)
					dirs = append(dirs, d.path)
					select {
					case chWalk <- d:
					case <-time.After(3 * time.Second):
						fmt.Println("cannot send to chWalk")
					}
					// case <-time.After(5 * time.Second):
					//   fmt.Println("dir timeout...")
				}
			}
		}()
	}

	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		go func() {
			for {
				select {
				case d := <-chWalk:
					if err := f.walk(d.path, d.depth, chDir); err != nil {
						_, _ = fmt.Fprintln(os.Stderr, err)
					}
					// case <-time.After(5 * time.Second):
					//   fmt.Println("walk timeout...")
				}
			}
		}()
	}

	name := filepath.Base(dir)
	chDir <- dirInfo{path: dir, name: name, depth: 0}
	// wg.Wait()
	// select {}

	idx, err := fuzzyfinder.Find(&dirs, func(i int) string {
		return dirs[i]
	},
		fuzzyfinder.WithHotReload(),
		fuzzyfinder.WithPreviewWindow(func(i, width, height int) string {
			if i < 0 || i >= len(dirs) {
				return ""
			}
			path := dirs[i]
			files, err := ioutil.ReadDir(path)
			if err != nil {
				return ""
			}
			sb := &strings.Builder{}
			sb.WriteString(path + "\n")
			for i, f := range files {
				if i == len(files)-1 {
					sb.WriteString(`└── `)
				} else {
					sb.WriteString(`├── `)
				}
				sb.WriteString(f.Name())
				sb.WriteRune('\n')
			}
			return sb.String()
		}),
	)
	// fmt.Println(dirs[idx], err)

	if err == nil {
		fmt.Println(dirs[idx])
	}
}

func (f *Finder) walk(dir string, depth int, chDir chan<- dirInfo) error {
	// fmt.Println("walk", dir, depth)
	if depth > f.Depth {
		return nil
	}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err
	}
	for _, f := range files {
		if !f.IsDir() || strings.HasPrefix(f.Name(), ".") || f.Name() == "node_modules" {
			continue
		}
		chDir <- dirInfo{path: filepath.Join(dir, f.Name()), name: f.Name(), depth: depth + 1}
	}
	return nil
}
