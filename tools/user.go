package tools

import (
	"database/sql"
	
	"github.com/infernalfire72/acache/log"
	"github.com/infernalfire72/acache/config"
)

func GetFriends(id int) []int {
	rows, err := config.DB.Query("SELECT user2 FROM users_relationships WHERE user1 = ?", id)
	var friends []int
	if err != nil {
		if err != sql.ErrNoRows {
			log.Error(err)
		}
		return friends
	}
	defer rows.Close()
	
	for rows.Next() {
		var user2 int
		err = rows.Scan(&user2)
		if err != nil {
			log.Error(err)
			continue
		}
		friends = append(friends, user2)
	}
	
	return friends
}