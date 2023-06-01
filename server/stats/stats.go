package stats

import (
	"time"

	"github.com/devusSs/crosshairs/database"
)

var (
	UsersRegisteredLast24Hours int
	UsersLoggedInLast24Hours   int
	RequestsInLast24Hours      int
)

type StatsAllTime struct {
	RegisteredUsers      int `json:"registered_users"`
	RegisteredCrosshairs int `json:"registered_crosshairs"`
}

type Stats24Hours struct {
	UsersRegistered int `json:"users_registered"`
	UserLogins      int `json:"user_logins"`
	Requests        int `json:"api_requests"`
}

func Reset24Statistics() {
	UsersLoggedInLast24Hours = 0
	UsersRegisteredLast24Hours = 0
	RequestsInLast24Hours = 0

	time.AfterFunc(CalculateTimeUntilMidnight(), Reset24Statistics)
}

func CalculateTimeUntilMidnight() time.Duration {
	currentTime := time.Now()
	midnight := time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day()+1, 0, 0, 0, 0, currentTime.Location())
	return midnight.Sub(currentTime)
}

func GetStatsAllTime(svc database.Service) (*StatsAllTime, error) {
	users, err := svc.GetAllUsers()
	if err != nil {
		return nil, err
	}

	chs, err := svc.GetAllCrosshairs()
	if err != nil {
		return nil, err
	}

	return &StatsAllTime{
		RegisteredUsers:      len(users),
		RegisteredCrosshairs: len(chs),
	}, nil
}

func GetStats24Hours() *Stats24Hours {
	return &Stats24Hours{
		UsersRegistered: UsersRegisteredLast24Hours,
		UserLogins:      UsersLoggedInLast24Hours,
		Requests:        RequestsInLast24Hours,
	}
}
