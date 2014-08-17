package models

type PostDocument struct {
	Id 		string `bson:"_id,omitempty"`
	Title 	string
	ContentHtml string
	ContentMd string
}

func NewPost(id, title, contentHtml, contentMd string) *PostDocument{
	return &PostDocument{id, title, contentHtml, contentMd}
}
