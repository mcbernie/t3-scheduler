package forum

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type topicModel struct {
	uid               int
	pid               int
	rawForum          int
	forum             *forumModel
	subject           string
	postsCnt          int // posts count
	postCount         int // posts count
	rawLastPost       int
	lastPost          *postModel
	rawLastPostCrdate int64
	lastPostCrdate    time.Time
	rawTstamp         int64
	tstamp            time.Time
	rawCrdate         int64
	crdate            time.Time
	deleted           bool
	hidden            bool
	posts             []postModel
	lastPostUpdated   bool
}

//go tF.loadTopics(f.dBConn)
func (t *topicModel) loadPosts(d *sql.DB) {
	rows, err := d.Query(
		"select uid, pid, topic, text, tstamp, crdate, deleted, hidden from tx_typo3forum_domain_model_forum_post where topic = ?", t.uid)

	if err != nil {
		panic(err.Error())
	}

	for rows.Next() {
		var p postModel
		if err := rows.Scan(
			&p.uid,
			&p.pid,
			&p.rawTopic,
			&p.text,
			&p.rawTstamp, &p.rawCrdate, &p.deleted, &p.hidden); err != nil {
			log.Fatal(err)
		}

		p.crdate = time.Unix(p.rawCrdate, 0)
		p.tstamp = time.Unix(p.rawTstamp, 0)
		p.topic = t

		t.posts = append(t.posts, p)
	}
}

func (t *topicModel) getPostByID(uid int) (bool, *postModel) {
	for i := 0; i < len(t.posts); i++ {
		p := &t.posts[i]
		if p.uid == uid {
			return true, p
		}
	}
	return false, nil
}

func (t *topicModel) getLastPost() *postModel {
	var newestPost *postModel
	var highestID int
	for i := 0; i < len(t.posts); i++ {
		if t.posts[i].uid > highestID && t.posts[i].hidden == false && t.posts[i].deleted == false {
			highestID = t.posts[i].uid
			newestPost = &t.posts[i]
		}
	}

	return newestPost
}

func (t *topicModel) loadLastPost() (found bool, updated bool) {

	myLastPost := t.getLastPost()

	if myLastPost != nil {

		if t.rawLastPost != myLastPost.uid {
			t.rawLastPost = myLastPost.uid
			t.rawLastPostCrdate = myLastPost.rawCrdate
			t.lastPostCrdate = myLastPost.crdate
			t.lastPostUpdated = true
		}

		t.lastPost = myLastPost
		return true, t.lastPostUpdated
	}

	if t.rawLastPost > 0 {
		t.lastPostUpdated = true
	}
	t.rawLastPost = -1
	return false, t.lastPostUpdated

}

//UpdateSQL returns query String to update table entry in database
func (t *topicModel) UpdateSQL() string {
	query := fmt.Sprintf("update tx_typo3forum_domain_model_forum_topic set last_post = %d, last_post_crdate = %d where uid = %d;",
		t.rawLastPost,
		t.rawLastPostCrdate,
		t.uid)
	log.Println("Database Update:", query)
	return query
}
