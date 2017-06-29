package forum

import (
	"database/sql"
	"log"
	"time"

	"github.com/mcbernie/t3-scheduler/db"
)

//Forum - Typo3Forum
type Forum struct {
	forums []forumModel
	dbConn db.Database
}

//Create a new Forum
func Create(dbConnectionString string) Forum {
	return Forum{
		dbConn: db.Create(dbConnectionString),
	}
}

//Run loading data and check database
func (f *Forum) Run() {
	log.Println("Load all Typo3 Forum data from Database")
	f.getForums()

	log.Println("Generate and Make database updates")
	affRows := f.makeUpdates()

	log.Printf("Datbase update Complete! %d affected rows.", affRows)

}

func (f *Forum) makeUpdates() int64 {
	var affectedRows int64
	for _, tF := range f.forums {
		if tF.lastPostUpdated == true || tF.lastTopicUpdated == true {
			for _, updateQuery := range tF.UpdateSQL() {
				affectedRows += f.executeUpdate(updateQuery)
			}
		}
	}

	return affectedRows
}

func (f *Forum) executeUpdate(query string) int64 {
	var affectedRows int64
	f.dbConn.DatabaseOperation(func(d *sql.DB) interface{} {
		res, err := d.Exec(query)

		if err != nil {
			panic(err.Error())
		}
		affectedRows, err = res.RowsAffected()
		return true
	})
	return affectedRows
}

func (f *Forum) getForums() {
	f.dbConn.DatabaseOperation(func(d *sql.DB) interface{} {
		rows, err := d.Query(
			"select uid, pid, forum, title, description, topics, topic_count, post_count," +
				"last_topic, last_post, tstamp, crdate, deleted, hidden from tx_typo3forum_domain_model_forum_forum")

		if err != nil {
			panic(err.Error())
		}

		for rows.Next() {
			var tF forumModel
			if err := rows.Scan(
				&tF.uid,
				&tF.pid,
				&tF.parentForumID,
				&tF.title, &tF.description,
				&tF.topicsCnt, &tF.topicCount, &tF.postCount,
				&tF.rawLastTopic, &tF.rawLastPost,
				&tF.rawTstamp, &tF.rawCrdate, &tF.deleted, &tF.hidden); err != nil {
				log.Fatal(err)
			}

			tF.crdate = time.Unix(tF.rawCrdate, 0)
			tF.tstamp = time.Unix(tF.rawTstamp, 0)

			tF.loadTopics(d)

			if found, updated := tF.loadLastPost(); updated == true {
				if found == false {
					log.Printf("Forum Check: last Post for Forum (%d) not found and updated!", tF.uid)
				} else {
					log.Printf("Forum Check: last Post (%d) for Forum (%d) updated!", tF.rawLastPost, tF.uid)
				}
			}

			if found, updated := tF.loadLastTopic(); updated == true {
				if found == false {
					log.Printf("Forum Check: last Topic for Forum (%d) not found and updated!", tF.uid)
				} else {
					log.Printf("Forum Check: last Topic (%d) for Forum (%d) updated!", tF.rawLastTopic, tF.uid)
				}
			}

			f.forums = append(f.forums, tF)
		}
		return true
	})

}
