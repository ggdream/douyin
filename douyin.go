package main

import (
	"os"
	"strings"
	"unsafe"

	"github.com/tidwall/gjson"
)

func DouYin(name, path string) error {
	data, err := os.ReadFile(name)
	if err != nil {
		return err
	}

	for _, v := range strings.Split(*(*string)(unsafe.Pointer(&data)), "\n") {
		if res := parseLink(v); res != nil {
			// fmt.Println(data)
			NewPersist(path).Run(res)
		}
	}

	return nil
}

func parseLink(data string) (medias []Media) {
	authors := gjson.Get(data, `aweme_list.#.author.nickname`).Array()
	descs := gjson.Get(data, `aweme_list.#.desc`).Array()
	covers := gjson.Get(data, `aweme_list.#.video.cover.url_list.1`).Array()
	videos := gjson.Get(data, `aweme_list.#.video.bit_rate.0.play_addr.url_list.0`).Array()

	if len(videos) == 0 || len(covers) == 0 {
		return nil
	}

	for i := 0; i < len(covers); i++ {
		medias = append(medias, Media{
			Author: authors[i].String(),
			Desc:   descs[i].String(),
			Cover:  covers[i].String(),
			Video:  videos[i].String(),
		})
	}

	return
}
