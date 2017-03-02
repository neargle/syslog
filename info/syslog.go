package info

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

var BlackList []string

//获取复数的windows安全日志
var psTools = `&{$reslist=Get-WinEvent -FilterHashtable @{'ProviderName'='Microsoft-Windows-Security-Auditing';Id=4625%s};` +
	`If($reslist.length){For ($index=0;$index -le $reslist.length-1;++$index){Write-Host $reslist[$index].toxml()}}Else{Write-Host $res.toxml();}}`

//获取单个的windows安全日志
//var psTools_One = `& {$res = Get-WinEvent -FilterHashtable @{'ProviderName'` +
//	`='Microsoft-Windows-Security-Auditing';Id=4776%s} -MaxEvents 1;Write-Host $res.toxml();}`

//windows下的正则
var regEx = regexp.MustCompile(`<TimeCreated SystemTime='(?P<time>[\w\-\:]+)\.\w+'\/>.*<Data Name='TargetUserName'>(?P<username>[^<]+)</Data>.*<Data Name='TargetDomainName'>(?P<hostname>([^<]*))</Data><Data Name='Status'>(?P<status>\w+)</Data>.*<Data Name='IpAddress'>(?P<ip>[^<]+)</Data>`)

//linux的正则
var regExLinux = regexp.MustCompile(
	`(?P<username>\w+) +[^ ]+ +(?P<ip>[\d\.]+) +[A-Z][a-z]{2} (?P<time>[A-Z][a-z]{2} {1,2}\d{1,2} \d{2}:\d{2})`)

var HOSTNAME, _ = os.Hostname()

const TimeFormat = "Jan 2 15:04"
const TimeFormat2 = "01/02-15:04"
const TimeFormat3 = "2006-01-02T15:04:05"

func GetSysLog(system string, starttime string) []map[string]string {
	// system : windows or linux
	// starttime : all or MM/dd-hh:ss
	var log_list []map[string]string
	if system == "windows" {
		if PsExists() { // 如果存在powershell
			var res *exec.Cmd
			pstime := psDate(starttime)
			ps := fmt.Sprintf(psTools, pstime)
			res = exec.Command("PowerShell", "-Command", ps)
			//fmt.Println(ps)
			out, _ := res.Output()
			xmlstr := string(out)
			lines := strings.Split(xmlstr, "\n")
			for _, v := range lines {
				if str := strings.TrimSpace(v); str != "" {
					logMap := xml2logMap(str)
					if len(logMap) > 0 {
						log_list = append(log_list, logMap)
					}
				}
			}
			return log_list
		} else { // 如果不存在powershell

		}
	} else { // linux 获取日志
		lastlist := linuxLog("last", starttime)
		lastblist := linuxLog("lastb", starttime)
		log_list = append(log_list, lastlist...)
		log_list = append(log_list, lastblist...)
	}
	return log_list
}

func linuxLog(cmd string, starttime string) []map[string]string {
	var reslist []map[string]string
	if starttime == "all" {
		starttime = "01/01-00:00"
	}
	out := Cmdexec(cmd, "linux")
	lines := strings.Split(out, "\n")
	for _, v := range lines {
		res := last2logMap(v)
		if len(res) > 0 {
			if len(reslist) > 0 {
				t1, _ := time.Parse(TimeFormat2, res["time"])
				t2, _ := time.Parse(TimeFormat2, reslist[len(reslist)-1]["time"])
				starttime, _ := time.Parse(TimeFormat2, starttime)
				if !t1.Before(t2) || !starttime.Before(t1) {
					return reslist
				}
			}
			if cmd == "last" {
				res["status"] = "true"
			} else {
				res["status"] = "false"
			}
			reslist = append(reslist, res)
		}
	}
	return reslist
}

func outofdate(time string, loginfo map[string]string) bool {
	return false
}

func CurrentPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func psDate(time string) string {
	if time == "all" {
		return ""
	} else {
		//		if m, _ := regexp.MatchString(`^\d{4}\/\d{2}\/\d{2}\-\d{2}\:\d{2}\:\d{2}$`, time); !m {
		//			return ""
		//		}
		res := fmt.Sprintf(";starttime=[datetime]::ParseExact('%s','MM/dd-HH:mm',$null)", time)
		return res
	}
}

func xml2logMap(xml string) map[string]string {
	match := regEx.FindStringSubmatch(xml)
	result := make(map[string]string)
	subName := regEx.SubexpNames()
	if len(subName) == len(match) {
		for i, name := range subName {
			if i != 0 && name != "" {
				if name == "status" {
					if match[i] == "0x0" {
						result[name] = "true"
					} else {
						result[name] = "false"
					}
				} else if name == "time" {
					loc, _ := time.LoadLocation("Local")
					t1, _ := time.Parse(TimeFormat3, match[i])
					ntime := t1.In(loc)
					result[name] = ntime.Format(TimeFormat2)
				} else if name == "ip" {
					if stringInSlice(match[i], BlackList) {
						return make(map[string]string)
					}
					result[name] = match[i]
				} else {
					result[name] = match[i]
				}
			}
		}
	}
	return result
}

func last2logMap(lastOut string) map[string]string {
	match := regExLinux.FindStringSubmatch(lastOut)
	res := make(map[string]string)
	subName := regExLinux.SubexpNames()
	if len(subName) == len(match) {
		for i, name := range subName {
			if name != "" {
				res[name] = match[i]
				if name == "time" {
					//fmt.Println(match[i])
					t, _ := time.Parse(TimeFormat, match[i])
					res[name] = t.Format(TimeFormat2)
					//fmt.Println(t.Format(TimeFormat2))
				}
			}
		}
	}
	return res
}

func PsExists() bool {
	res := Cmdexec("powershell -Command {true}", "windows")
	flag := strings.TrimSpace(res)
	if flag == "true" {
		return true
	}
	return false
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
