package forum

import "github.com/mcbernie/t3-scheduler/db"

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

//Load Database
/*
select
	tx_typo3forum_domain_model_forum_post.*,
	tx_typo3forum_domain_model_forum_topic.subject,
	tx_typo3forum_domain_model_forum_topic.uid as topic_uid,
	tx_typo3forum_domain_model_forum_forum.title,
	tx_typo3forum_domain_model_forum_forum.uid as forum_uid
from tx_typo3forum_domain_model_forum_post
left join tx_typo3forum_domain_model_forum_topic on
	tx_typo3forum_domain_model_forum_topic.uid = tx_typo3forum_domain_model_forum_post.topic
left join tx_typo3forum_domain_model_forum_forum on
	tx_typo3forum_domain_model_forum_forum.uid = tx_typo3forum_domain_model_forum_topic.forum
where
	tx_typo3forum_domain_model_forum_topic.uid = 11
*/

// in topic -> last_post last_post_crdate
