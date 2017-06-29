package forum

import (
	"time"
)

type postModel struct {
	uid       int
	pid       int
	rawTopic  int
	topic     *topicModel
	text      string
	rawTstamp int64
	tstamp    time.Time
	rawCrdate int64
	crdate    time.Time
	deleted   bool
	hidden    bool
}
