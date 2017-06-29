package forum

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

/*uid, pid, forum, title, description, topics, topic_count, post_count," +
"last_topic, last_post, tstamp, crdate, deleted, hidden from tx_typo3forum_domain_model_fo*/

type forumModel struct {
	uid              int
	pid              int
	parentForumID    int //id from parent forum (forum is parent forum object)
	title            string
	description      string
	topicsCnt        int // same as topic_count
	topicCount       int // same as topics
	postCount        int
	rawLastTopic     int
	rawLastPost      int
	lastTopic        *topicModel
	lastPost         *postModel
	rawTstamp        int64 // raw Timestamp
	rawCrdate        int64 // raw Timestamp
	tstamp           time.Time
	crdate           time.Time
	deleted          bool
	hidden           bool
	forum            *forumModel
	topics           []topicModel
	lastPostUpdated  bool
	lastTopicUpdated bool
}

//go tF.loadTopics(f.dBConn)
func (f *forumModel) loadTopics(d *sql.DB) {
	rows, err := d.Query(
		"select uid, pid, forum, subject, posts, post_count, last_post, last_post_crdate,"+
			"tstamp, crdate, deleted, hidden from tx_typo3forum_domain_model_forum_topic where forum = ?", f.uid)

	if err != nil {
		panic(err.Error())
	}

	for rows.Next() {
		var t topicModel
		if err := rows.Scan(
			&t.uid,
			&t.pid,
			&t.rawForum,
			&t.subject,
			&t.postsCnt, &t.postCount,
			&t.rawLastPost, &t.rawLastPostCrdate,
			&t.rawTstamp, &t.rawCrdate, &t.deleted, &t.hidden); err != nil {
			log.Fatal(err)
		}

		t.crdate = time.Unix(t.rawCrdate, 0)
		t.tstamp = time.Unix(t.rawTstamp, 0)

		t.loadPosts(d)
		if found, updated := t.loadLastPost(); updated == true {
			if found == false {
				log.Printf("Topic Check: last Post for Topic (%d) not found and updated!", t.uid)
			} else {
				log.Printf("Topic Check: last Post (%d) for Topic (%d) updated!", t.rawLastPost, t.uid)
			}
		}
		f.topics = append(f.topics, t)
	}
}

func (f *forumModel) getTopicByID(uid int) (bool, *topicModel) {
	for i := 0; i < len(f.topics); i++ {
		t := &f.topics[i]
		if t.uid == uid {
			return true, t
		}
	}
	return false, nil
}

func (f *forumModel) getLastTopic() *topicModel {

	post := f.getLastPost()

	if post == nil {
		return nil
	}

	return post.topic

	/*var highestID int
	var myLastTopic *topicModel

	for i := 0; i < len(f.topics); i++ {
		t := &f.topics[i]
		if t.uid > highestID && t.hidden == false && t.deleted == false {
			highestID = t.uid
			myLastTopic = t
		}
	}

	return myLastTopic
	*/
}

func (f *forumModel) getLastPost() *postModel {
	var highestID int
	var myLastPost *postModel

	for i := 0; i < len(f.topics); i++ {
		t := &f.topics[i]
		for i2 := 0; i2 < len(t.posts); i2++ {
			p := t.posts[i2]
			if p.uid > highestID && p.hidden == false && p.deleted == false && p.topic.deleted == false {
				highestID = p.uid
				myLastPost = &p
			}
		}

	}

	return myLastPost
}

func (f *forumModel) loadLastTopic() (found bool, updated bool) {

	myLastTopic := f.getLastTopic()
	if myLastTopic != nil {
		if myLastTopic.uid != f.rawLastTopic {
			f.rawLastTopic = myLastTopic.uid
			f.lastTopicUpdated = true
		}
		f.lastTopic = myLastTopic
		return true, f.lastTopicUpdated
	}

	if f.rawLastTopic > 0 {
		f.lastTopicUpdated = true
	}
	f.rawLastTopic = -1
	return false, f.lastTopicUpdated

}

func (f *forumModel) loadLastPost() (found bool, updated bool) {

	myLastPost := f.getLastPost()
	if myLastPost != nil {
		if myLastPost.uid != f.rawLastPost {
			f.rawLastPost = myLastPost.uid
			f.lastPostUpdated = true
		}
		f.lastPost = myLastPost
		return true, f.lastPostUpdated
	}

	if f.rawLastPost > 0 {
		f.lastPostUpdated = true
	}
	f.rawLastPost = -1
	return false, f.lastPostUpdated
}

//UpdateSQL returns query String to update table entry in database
func (f *forumModel) UpdateSQL() []string {
	var queries []string
	query := fmt.Sprintf("update tx_typo3forum_domain_model_forum_forum set last_topic = %d, last_post = %d where uid = %d;",
		f.rawLastTopic,
		f.rawLastPost,
		f.uid)
	log.Println("Database Update:", query)

	queries = append(queries, query)
	for _, t := range f.topics {
		if t.lastPostUpdated == true && t.deleted == false {
			queries = append(queries, t.UpdateSQL())
		}
	}
	return queries
}
