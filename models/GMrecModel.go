/**{
  "attachments": [],
  "avatar_url": "http://i.groupme.com/123456789",
  "created_at": 1302623328,
  "group_id": "1234567890",
  "id": "1234567890",
  "name": "John",
  "sender_id": "12345",
  "sender_type": "user",
  "source_guid": "GUID",
  "system": false,
  "text": "Hello world ☃☃",
  "user_id": "1234567890"
}**/
package models

type GMrecModel struct {
	Attachments interface{}
	Avatar_url  string
	Created_at  int
	Group_id    string
	Id          string
	Name        string
	Sender_id   string
	Sender_type string
	Source_guid string
	System      bool
	Text        string
	User_id     string
}
