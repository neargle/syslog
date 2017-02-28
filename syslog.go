package info

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

//获取复数的windows安全日志
var psTools = `& {$reslist =  Get-WinEvent -FilterHashtable ` +
	`@{'ProviderName'='Microsoft-Windows-Security-Auditing';Id=4776%s} -MaxEvents %d;` +
	`For ($index = 0; $index -le $reslist.length-1; ++$index){Write-Host $reslist[$index].toxml()}}`

//获取单个的windows安全日志
var psTools_One = `& {$res = Get-WinEvent -FilterHashtable @{'ProviderName'` +
	`='Microsoft-Windows-Security-Auditing';Id=4776%s} -MaxEvents 1;Write-Host $res.toxml();}`

//获取部分日志信息的正则
var regEx = regexp.MustCompile(`<TimeCreated SystemTime='(?P<time>[\w\-\:\.]+)'\/>.*` +
	`<Data Name='TargetUserName'>(?P<username>[^<]+)</Data>.*<Data Name='Workstation'>(?P<ip>([\d\.\-]+))</Data><Data Name='Status'>(?P<status>\w+)</Data>`)

func GetSysLog(system string, logCount int32, time string) []map[string]string {
	var log_list []map[string]string
	if system == "windows" {
		if PsExists() { // 如果存在powershell
			var res *exec.Cmd
			pstime := psDate(time)
			if logCount > 1 {
				ps := fmt.Sprintf(psTools, pstime, logCount)
				res = exec.Command("PowerShell", "-Command", ps)
			} else {
				ps := fmt.Sprintf(psTools, pstime)
				res = exec.Command("cmd", "/c", "PowerShell", "-Command", ps)
			}
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
	} else {

	}
	return log_list
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
		if m, _ := regexp.MatchString(`^\d{4}\/\d{2}\/\d{2}\-\d{2}\:\d{2}\:\d{2}$`, time); !m {
			return ""
		}
		res := fmt.Sprintf(";starttime=[datetime]::ParseExact('%s','yyyy/MM/dd-HH:mm:ss',$null)", time)
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
				} else {
					result[name] = match[i]
				}
			}
		}
	}
	return result
}

func PsExists() bool {
	res := Cmdexec("powershell -Command {true}", "windows")
	flag := strings.TrimSpace(res)
	if flag == "true" {
		return true
	}
	return false
}
