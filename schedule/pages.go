package schedule

import (
	"database/sql"
	"fmt"
	"hash/fnv"
	"log"
	"time"
)

type pageResults struct {
	uid                int
	tstamp             int64
	datetime           time.Time
	title              string
	lastchange         int64
	lastchangeDatetime time.Time
	content            []pageContent
	hidden             bool
	deleted            bool
}
type pageContent struct {
	uid            int
	tstamp         int64
	datetime       time.Time
	crdate         int64
	createDatetime time.Time
	bodytext       string
}

//String generate a string for hash generation
func (pr *pageResults) String() string {
	page := fmt.Sprintf("%d%d%s%d", pr.uid, pr.tstamp, pr.title, pr.lastchange)

	for _, c := range pr.content {
		page += fmt.Sprintf("%d%d%d%s", c.uid, c.tstamp, c.crdate, c.bodytext)
	}

	return page
}

//Hash generates a uniq Hash from current page with all contents with help of String() function
func (pr *pageResults) Hash() uint32 {
	h := fnv.New32a()
	h.Write([]byte(pr.String()))
	return h.Sum32()
}

//getPages(uid) get specified page with content from Database
func (s *Schedule) getPage(uid int) pageResults {
	page := pageResults{}
	s.databaseOperation(func(d *sql.DB) interface{} {

		rows, err := d.Query("select uid, tstamp, title, SYS_LASTCHANGED, hidden from pages where uid = ? and deleted = 0", uid)
		if err != nil {
			panic(err.Error())
		}
		for rows.Next() {
			if err := rows.Scan(&page.uid, &page.tstamp, &page.title, &page.lastchange, &page.hidden); err != nil {
				log.Fatal(err)
			}

			page.datetime = time.Unix(page.tstamp, 0)
			page.lastchangeDatetime = time.Unix(page.lastchange, 0)
			s.getContent(&page)

		}
		return true
	})

	if page.uid < 1 {
		page.deleted = true
	}

	return page
}

func (s *Schedule) getContent(page *pageResults) {
	s.databaseOperation(func(d *sql.DB) interface{} {

		rows, err := d.Query("select uid, tstamp, crdate from tt_content where pid = ? and hidden = 0 and deleted = 0", page.uid)
		if err != nil {
			panic(err.Error())
		}

		for rows.Next() {

			content := pageContent{}
			if err := rows.Scan(&content.uid, &content.tstamp, &content.crdate); err != nil {
				log.Printf("error on scan getContent: %s", err.Error())
				continue
			}

			content.datetime = time.Unix(content.tstamp, 0)
			content.createDatetime = time.Unix(content.crdate, 0)

			page.content = append(page.content, content)
		}
		return true
	})

}
