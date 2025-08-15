package database

import (
	"testing"
	"time"

	"github.com/k0kubun/pp"
)

func TestUser(t *testing.T) {
	d := CreateRepository("user.db")
	u, err := d.CreateUser(
		&User{
			UserID:     "mocho",
			Name:       "Momo Asakura",
			PW:         "sadasifh92q3hfuawef",
			DebutDate:  time.Date(1994, 6, 25, 4, 17, 0, 0, time.Local),
			PlayerType: 3,
			Stance:     3.4,
			Club:       "MusicRayn",
			Major:      "昭和歌謡曲",
			Profile:    "麻倉ももです！イントロクイズはじめました！",
		},
	)
	if err != nil {
		t.Errorf("user creation error: %v", err)
	}
	pp.Println(u)

	u2, err := d.GetUser("mocho")
	if err != nil {
		t.Errorf("user get error: %v", err)
	}
	pp.Println(u2)
}
