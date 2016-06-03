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

func (this *GMrepModel) InitbyGID(gid string) {
	mapforGroupName := map[string]string{"19624531": "d9c55c1ed0d017b645908fac5f", "15529154": "90af1423a7b97665968ad4bcdd", "12000977": "72693bab1250b6b353d66947f1", "14806448": "7080dbdf1ca3f3e0f2fd67f16d", "21088731": "a38b2c9bc6c7f0f40de1d1d841"}
	this.Bot_id = mapforGroupName[gid]
	return
}
func (this *GMrepModel) SetText(text string) {
	this.Text = text
	return
}
