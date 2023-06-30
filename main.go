package main

import (
	"fmt"
	"github.com/karmdip-mi/go-fitz"
	"image/jpeg"
	//_ "net/http/pprof"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// читает файлы из папка pdf (создать в корне проекта) и преобразует их в картинки,
// сохраняет их в img (создать в корне проекта).
func runner(wg *sync.WaitGroup) {
	defer wg.Done()
	var files []string

	root := "pdf/"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if filepath.Ext(path) == ".pdf" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		doc, err := fitz.New(file)
		if err != nil {
			panic(err)
		}
		folder := strings.TrimSuffix(path.Base(file), filepath.Ext(path.Base(file)))

		// Extract pages as images
		for n := 0; n < doc.NumPage(); n++ {
			img, err := doc.ImageDPI(n, 150.0)
			if err != nil {
				panic(err)
			}
			err = os.MkdirAll("img/"+folder, 0755)
			if err != nil {
				panic(err)
			}

			f, err := os.Create(filepath.Join("img/"+folder+"/", fmt.Sprintf("image-%05d.jpg", n)))
			if err != nil {
				panic(err)
			}

			err = jpeg.Encode(f, img, &jpeg.Options{Quality: jpeg.DefaultQuality})
			if err != nil {
				panic(err)
			}

			f.Close()

		}
	}
}

func main() {
	fmt.Println(time.Now())

	//для запуска профайлера
	//go func() {
	//	fmt.Println(http.ListenAndServe("localhost:6060", nil))
	//}()
	var wg sync.WaitGroup
	wg.Add(1)
	//wg.Add(1)
	go runner(&wg)
	wg.Wait()
	fmt.Println(time.Now())
}
