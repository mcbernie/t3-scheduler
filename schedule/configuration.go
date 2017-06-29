package schedule

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/go-ini/ini"
)

type feuser struct {
	id       int
	username string
	mail     string
}

type watch struct {
	pageID      int
	title       string
	users       []feuser
	currentHash uint32
	lastHash    uint32
	newWatch    bool
	hidden      bool
	pageDeleted bool
}

func (w *watch) String() string {
	ret := fmt.Sprintf("watch for updates for %s(id:%d) and notifiy %d users:\n", w.title, w.pageID, len(w.users))
	for _, nu := range w.users {
		ret += fmt.Sprintf("\t%s\t%s\n", nu.username, nu.mail)
	}
	ret += fmt.Sprintf("\tcurrentHash: %d\t lastHash:%d\n", w.currentHash, w.lastHash)
	return ret
}

func (s *Schedule) getUserFromUsername(username string) feuser {
	user := feuser{}
	s.databaseOperation(func(d *sql.DB) interface{} {

		rows, err := d.Query("select uid, username, email from fe_users where username = ?", username)
		if err != nil {
			panic(err.Error())
		}

		for rows.Next() {
			if err := rows.Scan(&user.id, &user.username, &user.mail); err != nil {
				log.Fatal(err)
			}
		}
		return true
	})

	return user

}

func (s *Schedule) loadLastHashes() {
	configurationFile := filepath.Join(s.fileadminPath, "lastHashes.ini")
	var cfg *ini.File
	noHashes := false
	if _, err := os.Stat(configurationFile); err != nil {
		fmt.Fprintf(os.Stderr, "no lastHashes in %s found !\n", configurationFile)
		noHashes = true
	} else {
		var loadError error
		cfg, loadError = ini.LoadSources(ini.LoadOptions{AllowBooleanKeys: true}, configurationFile)
		if loadError != nil {
			fmt.Fprintf(os.Stderr, "error on loading lastHashes from %s!\n(%s)\n", configurationFile, loadError.Error())
			noHashes = true
		}
	}

	for i := 0; i < len(s.watches); i++ {
		w := &s.watches[i]
		if noHashes == true {
			w.newWatch = true
		} else {
			section, err := cfg.GetSection(fmt.Sprintf("%d", w.pageID))
			if err != nil {
				w.newWatch = true
				continue
			}
			var hash64 uint64
			hash64, err = section.Key("hash").Uint64()
			if err != nil {
				fmt.Fprintf(os.Stderr, "go error on loading hash for page %d: %s\n", w.pageID, err.Error())
				w.newWatch = true
			} else {
				hash := uint32(hash64)
				w.lastHash = hash
			}

		}

	}
}

func (s *Schedule) loadConfiguration() {
	watchSection, err := s.cfgIni.GetSection("watch")
	receiversSection, errReceiver := s.cfgIni.GetSection("receivers")
	if err != nil {
		panic("no watch configuration found!")
	}
	if errReceiver != nil {
		panic("no receivers configuration found!")
	}
	pageIDString := watchSection.Key("pages").String()

	pids := strings.Split(pageIDString, ",")

	for _, pidStr := range pids {

		pid, _ := strconv.Atoi(pidStr)
		w := watch{
			pageID: pid,
		}
		receiver := receiversSection.Key(pidStr).String()
		usernames := strings.Split(receiver, ",")
		for _, username := range usernames {
			user := s.getUserFromUsername(username)
			if user.mail != "" {
				w.users = append(w.users, user)
			}
		}

		dbPage := s.getPage(pid)
		w.title = dbPage.title
		w.currentHash = dbPage.Hash()
		w.hidden = dbPage.hidden
		w.pageDeleted = dbPage.deleted

		s.watches = append(s.watches, w)
	}

}
