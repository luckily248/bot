/**
{
  "bot_id"  : "j5abcdefg",
  "text"    : "Hello world"
}

token iN18M7dmlgEMhUaKwyuWED2nHBDkTGSskj7Iv4KC
**/

package models

type GMrepModel struct {
	Text   string `json:"text"`
	Bot_id string `json:"bot_id"`
}

func (this *GMrepModel) Init() {
	this.Bot_id = "90af1423a7b97665968ad4bcdd" //init botid
	return
}
func (this *GMrepModel) SetText(text string) {
	this.Text = text
	return
}
