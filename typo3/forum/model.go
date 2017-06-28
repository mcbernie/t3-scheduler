package forum

type forumModel struct {
	uid         int
	title       string
	description string
	forum       *forumModel
}

type topicModel struct {
	uid            int
	forum          forumModel
	subject        string
	posts          int // posts count
	postCount      int // posts count
	lastPost       int
	lastPostCrdate string
}

type postModel struct {
	uid     int
	topic   topicModel
	text    string
	deleted bool
	hidden  bool
}
