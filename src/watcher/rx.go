package watcher

import "regexp"

func (w Watcher) rxfind(rx string, content string) (r string) {
	temp, _ := regexp.Compile(rx)
	r = temp.FindString(content)
	return
}
