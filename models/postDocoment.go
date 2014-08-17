package models

type PostDocument struct {
	Id 		string `bson:"_id,omitempty"`
	Title 	string
	Content string
}
