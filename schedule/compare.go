package schedule

import (
	"fmt"
	"path/filepath"

	"github.com/go-ini/ini"
)

func (w *watch) changed() bool {
	if w.hidden == true || w.pageDeleted == true {
		return false
	}

	if w.newWatch == true {
		return false
	}
	if w.currentHash != w.lastHash {
		return true
	}
	return false
}

// compare hases and write new configuration for oldhashes
func (s *Schedule) compare() int {
	var count int

	for _, w := range s.watches {
		if w.changed() {
			count++
			w.notifiyUsers(s)
		}
	}

	s.generateNewlastHashes()
	return count
}

func (s *Schedule) generateNewlastHashes() {
	configurationFile := filepath.Join(s.fileadminPath, "lastHashes.ini")
	cfg := ini.Empty()
	for _, w := range s.watches {
		if w.hidden == true || w.pageDeleted == true {
			continue
		}
		section, err := cfg.NewSection(fmt.Sprintf("%d", w.pageID))
		if err != nil {
			panic("Error creating section!")
		}
		section.NewKey("hash", fmt.Sprintf("%d", w.currentHash))
	}

	err := cfg.SaveTo(configurationFile)

	if err != nil {
		panic(fmt.Sprintf("Error on saving lastHashes.ini:%s", err.Error()))
	}
}
