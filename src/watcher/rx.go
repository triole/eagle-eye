package watcher

import "regexp"

func (w *Watcher) matchRegex(rx string, content string) (result string) {
	// Initialize regexCache if nil (e.g., when Watcher is created without Init)
	if w.regexCache == nil {
		w.regexCache = make(map[string]*regexp.Regexp)
	}

	// Try to get compiled regex from cache (read lock)
	w.regexMu.RLock()
	compiled, cached := w.regexCache[rx]
	w.regexMu.RUnlock()

	if cached {
		return compiled.FindString(content)
	}

	// Compile and cache the regex (write lock)
	w.regexMu.Lock()
	// Double-check after acquiring write lock
	if compiled, cached = w.regexCache[rx]; cached {
		w.regexMu.Unlock()
		return compiled.FindString(content)
	}

	temp, err := regexp.Compile(rx)
	if err != nil {
		w.regexCache[rx] = nil
		w.regexMu.Unlock()
		return result
	}

	w.regexCache[rx] = temp
	w.regexMu.Unlock()

	return temp.FindString(content)
}
