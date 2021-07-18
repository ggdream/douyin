package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

type Persist struct {
	dest string
	wg   sync.WaitGroup

	cover string
	video string
}

func NewPersist(path string) *Persist {
	return &Persist{
		dest: filepath.Join(path, "./dist"),
		wg:   sync.WaitGroup{},

		cover: "cover",
		video: "video",
	}
}

func (p *Persist) Run(data []Media) {
	p.SaveFiles(data)
}

// 创建用户目录
func (p *Persist) Mkdir(user string) error {
	if _, err := os.Stat(filepath.Join(p.dest, user)); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Join(p.dest, user, p.cover), 0755); err != nil {
			return err
		}
		if err := os.MkdirAll(filepath.Join(p.dest, user, p.video), 0755); err != nil {
			return err
		}
	}

	return nil
}

func (p *Persist) save(media Media) {
	defer p.wg.Done()

	if err := p.Mkdir(media.Author); err != nil {
		panic(err)
	}

	// fmt.Printf("开始下载封面：%s", media.Cover)
	if err := p.request(media.Cover, media.Desc, media.Author, true); err != nil {
		panic(err)
	}
	// fmt.Printf("下载完成封面：%s", media.Cover)

	// fmt.Printf("开始下载视频：%s", media.Video)
	if err := p.request(media.Video, media.Desc, media.Author, false); err != nil {
		panic(err)
	}
	// fmt.Printf("下载完成视频：%s", media.Video)
}

func (p *Persist) request(url, dest, user string, isCover bool) error {
	client := http.Client{}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("host", "v26-web.douyinvod.com")
	req.Header.Set("referer", url)
	req.Header.Set("range", "bytes=0-")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36 Edg/91.0.864.70")

	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	file, err := os.OpenFile(filepath.Join(p.dest, user, func() string {
		if isCover {
			return p.cover + "/" + dest + ".jpeg"
		}
		return p.video + "/" + dest + ".mp4"
	}()), os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, res.Body)
	return err
}

func (p *Persist) SaveFiles(medias []Media) {
	for _, v := range medias {
		p.wg.Add(1)
		go p.save(v)
	}

	p.wg.Wait()
}
