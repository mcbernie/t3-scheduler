package schedule

import (
	"fmt"

	"github.com/go-ini/ini"
	"github.com/mcbernie/t3-scheduler/typo3"
)

type databaseConfiguration struct {
	database string
	username string
	password string
	host     string
}

//Schedule simple Schedule
type Schedule struct {
	cfgIni        *ini.File
	cfgTypo3      typo3.Typo3
	typo3Path     string
	fileadminPath string
	watches       []watch
}

//Create Create a scheduler
func Create(cfg *ini.File, t3Path string, faPath string) Schedule {
	return Schedule{
		cfgIni:        cfg,
		cfgTypo3:      typo3.Load(t3Path),
		typo3Path:     t3Path,
		fileadminPath: faPath,
	}
}

//Run start the schedule
func (s *Schedule) Run() (int, error) {

	//load all data from ini
	s.loadConfiguration()

	//load last hashes from hash filepath
	s.loadLastHashes()

	/*for _, w := range s.watches {
		fmt.Println(w.String())
	}*/

	notifications := s.compare()
	return notifications, fmt.Errorf("send %d notifications for %d pages", notifications, len(s.watches))
}
