package stats

import (
	"errors"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/devusSs/crosshairs/database"
	"github.com/devusSs/crosshairs/storage"
	"github.com/devusSs/crosshairs/updater"
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

type SystemInfo struct {
	BuildInfo struct {
		BuildVersion string `json:"build_version"`
		BuildDate    string `json:"build_date"`
		BuildOS      string `json:"build_os"`
		BuildArch    string `json:"build_arch"`
		GoVersion    string `json:"go_version"`
	} `json:"build_info"`
	Systeminfo struct {
		CPUCount        int    `json:"cpu_count"`
		CGOCalls        int64  `json:"cgo_calls"`
		GoRoutinesCount int    `json:"goroutines_count"`
		Pagesize        int    `json:"pagesize"`
		ProcessID       int    `json:"process_id"`
		PathInfo        string `json:"path_info"`
		HostInfo        string `json:"host_info"`
		ResolvedAddr    bool   `json:"resolved_addr"`
	} `json:"system_info"`
	Integration struct {
		PostgresVersion string `json:"postgres_version"`
		RedisVersion    string `json:"redis_version"`
		MinioVersion    string `json:"minio_version"`
	} `json:"integration"`
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

func CollectSystemStats(svc database.Service, minioSvc *storage.Service) (*SystemInfo, error) {
	pPath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}

	dnsWorks, err := testDNS()
	if err != nil {
		return nil, err
	}

	postgresVersion, err := svc.GetPostgresVersion()
	if err != nil {
		return nil, err
	}

	// TODO: does this work on docker only machines without redis installed?
	redisVersion, err := exec.Command("redis-server", "--version").Output()
	if err != nil {
		return nil, err
	}

	minioVersion, err := minioSvc.CheckMinioVersion()
	if err != nil {
		return nil, err
	}

	var info SystemInfo

	info.BuildInfo.BuildVersion = updater.BuildVersion
	info.BuildInfo.BuildDate = updater.BuildDate
	info.BuildInfo.BuildOS = updater.BuildOS
	info.BuildInfo.BuildArch = updater.BuildARCH
	info.BuildInfo.GoVersion = updater.BuildGo

	info.Systeminfo.CPUCount = runtime.NumCPU()
	info.Systeminfo.CGOCalls = runtime.NumCgoCall()
	info.Systeminfo.GoRoutinesCount = runtime.NumGoroutine()
	info.Systeminfo.Pagesize = os.Getpagesize()
	info.Systeminfo.ProcessID = os.Getpid()
	info.Systeminfo.PathInfo = pPath
	info.Systeminfo.HostInfo = hostname
	info.Systeminfo.ResolvedAddr = dnsWorks

	info.Integration.PostgresVersion = postgresVersion
	info.Integration.RedisVersion = strings.TrimSpace(string(redisVersion))
	info.Integration.MinioVersion = minioVersion

	return &info, nil
}

func testDNS() (bool, error) {
	ips, err := net.LookupHost("github.com")
	if err != nil {
		return false, err
	}

	if len(ips) == 0 {
		return false, errors.New("no ip address found for test host")
	}

	return true, nil
}
