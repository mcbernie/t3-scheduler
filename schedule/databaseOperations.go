package schedule

import (
	"github.com/mcbernie/t3-scheduler/db"
)

func (s *Schedule) databaseOperation(fn db.Query) interface{} {
	d := db.Create(s.cfgTypo3.DatabaseConfiguration.String())
	ret := d.DatabaseOperation(fn)
	return ret
}
