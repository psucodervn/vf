package internal

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
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
	chDir := make(chan dirInfo, 2_000)
	chWalk := make(chan dirInfo, 2_000)
	var dirs struct {
		sync.Mutex
		arr []string
	}

	for i := 0; i < runtime.GOMAXPROCS(0)*2; i++ {
		go func() {
			for {
				select {
				case d := <-chDir:
					dirs.Lock()
					dirs.arr = append(dirs.arr, d.path)
					dirs.Unlock()
					go func() {
						select {
						case chWalk <- d:
						case <-time.After(3 * time.Second):
							fmt.Println("cannot send to chWalk")
						}
					}()
				}
			}
		}()
	}

	for i := 0; i < runtime.GOMAXPROCS(0)*8; i++ {
		go func() {
			for d := range chWalk {
				if err := f.walk(d.path, d.depth, chDir); err != nil {
					// _, _ = fmt.Fprintln(os.Stderr, err)
				}
			}
		}()
	}

	name := filepath.Base(dir)
	chDir <- dirInfo{path: dir, name: name, depth: 0}

	idx, err := fuzzyfinder.Find(&dirs.arr, func(i int) string {
		return dirs.arr[i]
	},
		fuzzyfinder.WithHotReload(),
		fuzzyfinder.WithPreviewWindow(func(i, width, height int) string {
			if i < 0 || i >= len(dirs.arr) {
				return ""
			}
			path := dirs.arr[i]
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

	if err == nil {
		fmt.Println(dirs.arr[idx])
	}
}

func (f *Finder) walk(dir string, depth int, chDir chan<- dirInfo) error {
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
