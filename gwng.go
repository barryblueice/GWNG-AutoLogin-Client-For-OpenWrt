package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	mrand "math/rand"
	"net"
	"os"
	"os/exec"

	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func checkInternetConnection(interfaceName string) bool {
	cmd := exec.Command("ping", "-I", interfaceName, "-c", "1", "www.baidu.com")
	err := cmd.Run()
	if err != nil {
		fmt.Printf("(%s) Error executing ping command:%s", interfaceName, err)
		return false
	}
	log_handle(fmt.Sprintf("%s联网成功\n", interfaceName))
	return true
}

func import_config() (map[string]string, bool) {
	ConfigFile, err := os.Open("/etc/config/gwng_autologin")
	var emptyBytes map[string]string
	if err != nil {
		return emptyBytes, false
	}

	Config_ram := ""

	scanner := bufio.NewScanner(ConfigFile)

	for scanner.Scan() {
		line := scanner.Text()
		if line != "\n" {
			line = strings.ReplaceAll(line, "	", "")
			line = strings.ReplaceAll(line, "list ", "")
			line = strings.ReplaceAll(line, "option ", "")
			line = strings.ReplaceAll(line, "'", "")
			line = strings.ReplaceAll(line, "\n", "")
			line = strings.TrimLeft(line, " ")
			line += " "
			Config_ram += line
		}
	}
	Config_ram = strings.TrimLeft(Config_ram, " ")
	configList := strings.Split(Config_ram, " ")

	result := make(map[string]string)
	for i := 0; i < len(configList); i += 2 {
		key := configList[i]
		value := configList[i+1]
		result[key] = value
	}

	ConfigData, err := json.Marshal(result)

	if err != nil {
		return emptyBytes, false
	}

	var ConfigMap map[string]string
	err = json.Unmarshal(ConfigData, &ConfigMap)
	if err != nil {
		return emptyBytes, false
	}
	return ConfigMap, true
}

func getIP(device string) (string, error) {
	// networkPath := "/etc/config/network"
	// file, err := os.Open(networkPath)
	// if err != nil {
	// 	return nil, err
	// }
	// defer file.Close()

	// networkList := make(map[string]string)
	// scanner := bufio.NewScanner(file)
	// for scanner.Scan() {
	// 	line := strings.ReplaceAll(scanner.Text(), "\t", "")
	// 	line = strings.ReplaceAll(line, "'", "")
	// 	line = strings.TrimSpace(line)

	// 	if strings.HasPrefix(line, "config interface ") {
	// 		iface := strings.TrimPrefix(line, "config interface ")
	// 		iface += " "
	// 		networkList[iface[:len(iface)-1]] = ""
	// 	} else if strings.HasPrefix(line, "option ipaddr ") {
	// 		line = strings.ReplaceAll(line, " ", ":")
	// 		line = strings.TrimPrefix(line, "option:ipaddr:")
	// 		line += " "
	// 		lastKey := ""
	// 		for k := range networkList {
	// 			lastKey = k
	// 		}
	// 		if lastKey != "" {
	// 			networkList[lastKey] = strings.TrimSpace(line)
	// 		} else {
	// 			return nil, fmt.Errorf("found 'option ipaddr' line without preceding 'config interface' line")
	// 		}
	// 	}
	// }

	// if err := scanner.Err(); err != nil {
	// 	return nil, err
	// }

	// return networkList, nil
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Name == device {
			addrs, err := iface.Addrs()
			if err != nil {
				return "", err
			}

			for _, addr := range addrs {
				ipNet, ok := addr.(*net.IPNet)
				if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
					return ipNet.IP.String(), nil
				}
			}
		}
	}

	return "", fmt.Errorf("IPv4 address not found")
}

func getMac(iface string) (string, error) {
	cmd := exec.Command("ifconfig", iface)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	re := regexp.MustCompile(`([[:xdigit:]]{1,2}:){5}[[:xdigit:]]{1,2}`)
	match := re.FindString(out.String())
	if match == "" {
		return "", fmt.Errorf("no MAC address found for interface %s", iface)
	}

	return strings.ToLower(match), nil
}

func generateRandomMac(separator string) (string, error) {
	mac := make([]byte, 6)
	if _, err := rand.Read(mac); err != nil {
		return "", err
	}

	mac[0] |= 0x02
	mac[0] &^= 0x01

	var macStr string
	for i, b := range mac {
		if i > 0 {
			macStr += separator
		}
		macStr += fmt.Sprintf("%02x", b)
	}

	return macStr, nil
}

func generateRandomName(length int, r *mrand.Rand) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}

func generateRandomSeed(length int, r *mrand.Rand) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return strings.ToLower(string(b))
}

func gen_uuid() string {
	randomUUID := uuid.New()
	randomUUIDString := randomUUID.String()
	return (string(randomUUIDString))
}

func get_timestamp() string {
	currentTime := time.Now()
	timestampMillis := currentTime.UnixNano() / int64(time.Millisecond)
	timestampString := strconv.FormatInt(timestampMillis, 10)
	return (string(timestampString))
}

func curl(userid string, passwd string, hostname string, mac string, userip string) {
	randomGenerator := mrand.New(mrand.NewSource(time.Now().UnixNano()))
	generatedMac1, err := generateRandomMac(":")
	if err != nil {
		generatedMac1 = "as:as:as:as:as:as"
	}
	generatedMac2, err := generateRandomMac(":")
	if err != nil {
		generatedMac2 = "sa:sa:sa:sa:sa:sa"
	}
	generatedUUID := gen_uuid()
	generatedSeed := generateRandomSeed(14, randomGenerator)
	gentimestamp := get_timestamp()
	// fmt.Println(generatedMac1)
	// fmt.Println(generatedMac2)
	// fmt.Println(generatedUUID)
	// fmt.Println(gentimestamp)
	// generatedSeed := "6d3ef521eff4e8"
	cmd := exec.Command("curl",
		fmt.Sprintf("http://10.244.0.13/quickauth.do?userid=%s&passwd=%s&wlanuserip=%s&wlanacname=gwng&wlanacIp=10.244.0.1&ssid=&vlan=502&mac=%s&version=0&portalpageid=1&timestamp=%s&uuid=%s&portaltype=1&hostname=%s", userid, passwd, userip, mac, gentimestamp, generatedUUID, hostname),
		"-H", "Accept: application/json, text/javascript, */*; q=0.01",
		"-H", "Accept-Language: zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6",
		"-H", "Connection: keep-alive",
		"-H", fmt.Sprintf("Cookie: macAuth=%s|%s|%s; ABMS=19aa6d57-bd91-4aee-bba1-49dda5c9654c", generatedMac1, generatedMac2, mac),
		"-H", fmt.Sprintf("Referer: http://10.244.0.13/portal.do?wlanuserip=%s&wlanacname=gwng&mac=%s&vlan=502&hostname=%s&rand=%s&url=http://www.msftconnecttest.com/redirec", userip, mac, hostname, generatedSeed),
		"-H", "User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0",
		"-H", "X-Requested-With: XMLHttpRequest",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log_handle("Error executing curl command:" + err.Error())
		return
	}

	fmt.Println(string(output))
}

func curl_dual(userid string, passwd string, hostname1 string, hostname2 string, mac1 string, mac2 string, userip1 string, userip2 string) {

	randomGenerator := mrand.New(mrand.NewSource(time.Now().UnixNano()))
	generatedMac1, err := generateRandomMac(":")
	if err != nil {
		generatedMac1 = "as:as:as:as:as:as"
	}
	generatedUUID := gen_uuid()
	generatedSeed := generateRandomSeed(14, randomGenerator)
	gentimestamp := get_timestamp()
	// fmt.Println(generatedMac1)
	// fmt.Println(generatedMac2)
	// fmt.Println(generatedUUID)
	// fmt.Println(gentimestamp)
	// generatedSeed := "6d3ef521eff4e8"
	cmd := exec.Command("curl",
		fmt.Sprintf("http://10.244.0.13/quickauth.do?userid=%s&passwd=%s&wlanuserip=%s&wlanacname=gwng&wlanacIp=10.244.0.1&ssid=&vlan=502&mac=%s&version=0&portalpageid=1&timestamp=%s&uuid=%s&portaltype=1&hostname=%s", userid, passwd, userip1, mac1, gentimestamp, generatedUUID, hostname1),
		"-H", "Accept: application/json, text/javascript, */*; q=0.01",
		"-H", "Accept-Language: zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6",
		"-H", "Connection: keep-alive",
		"-H", fmt.Sprintf("Cookie: macAuth=%s|%s|%s; ABMS=19aa6d57-bd91-4aee-bba1-49dda5c9654c", generatedMac1, mac1, mac2),
		"-H", fmt.Sprintf("Referer: http://10.244.0.13/portal.do?wlanuserip=%s&wlanacname=gwng&mac=%s&vlan=502&hostname=%s&rand=%s&url=http://www.msftconnecttest.com/redirec", userip1, mac1, hostname1, generatedSeed),
		"-H", "User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0",
		"-H", "X-Requested-With: XMLHttpRequest",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log_handle("Error executing curl command:" + err.Error())
		return
	}

	log_handle(string(output))

	cmd = exec.Command("curl",
		fmt.Sprintf("http://10.244.0.13/quickauth.do?userid=%s&passwd=%s&wlanuserip=%s&wlanacname=gwng&wlanacIp=10.244.0.1&ssid=&vlan=502&mac=%s&version=0&portalpageid=1&timestamp=%s&uuid=%s&portaltype=1&hostname=%s", userid, passwd, userip2, mac2, gentimestamp, generatedUUID, hostname2),
		"-H", "Accept: application/json, text/javascript, */*; q=0.01",
		"-H", "Accept-Language: zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6",
		"-H", "Connection: keep-alive",
		"-H", fmt.Sprintf("Cookie: macAuth=%s|%s|%s; ABMS=19aa6d57-bd91-4aee-bba1-49dda5c9654c", generatedMac1, mac1, mac2),
		"-H", fmt.Sprintf("Referer: http://10.244.0.13/portal.do?wlanuserip=%s&wlanacname=gwng&mac=%s&vlan=502&hostname=%s&rand=%s&url=http://www.msftconnecttest.com/redirec", userip2, mac2, hostname2, generatedSeed),
		"-H", "User-Agent: Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0",
		"-H", "X-Requested-With: XMLHttpRequest",
	)

	output, err = cmd.CombinedOutput()
	if err != nil {
		log_handle("Error executing curl command:" + err.Error())
		return
	}

	log_handle(string(output))
}

func log_handle(message string) {
	logDir := "/var/log/gwng"
	logFile := "gwng.log"
	logFilePath := filepath.Join(logDir, logFile)

	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	f, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer f.Close()

	// 设置日志输出到文件
	log.SetOutput(f)

	log.Printf("%s\n", message)

	currentTime := time.Now().Format("2006-01-02 15:04:05 ")

	fmt.Println(currentTime + message)
}

func main() {

	// a, err := getIP()
	// if err != nil {
	// 	log_handle(err)
	// }
	// log_handle(a["lan"])
	randomGenerator := mrand.New(mrand.NewSource(time.Now().UnixNano()))

	log_handle("程序开始执行……\n")

	// for iface, ip := range networkList {
	// 	fmt.Printf("Interface: %s, IP: %s\n", iface, ip)
	// }

	for {
		for i := 0; i < 10; i++ {
			Config_Content, Ok := import_config()
			num := Config_Content["num"]
			if Ok {
				if num == "1" {
					interfaces1 := Config_Content["interfaces1"]
					if checkInternetConnection(interfaces1) {
						fmt.Println(interfaces1 + "联网成功")
					} else {
						for {
							log_handle("联网失败，正在自动联网……")
							interfaces1 = Config_Content["interfaces1"]
							userid := Config_Content["username"]
							passwd := Config_Content["password"]
							hostname := generateRandomName(12, randomGenerator)
							userip, err := getIP(interfaces1)
							if err != nil {
								userip = "127.0.0.1"
							}
							mac, err := getMac(interfaces1)
							if err != nil {
								mac = "88:88:88:88:88:88"
							}
							curl(userid, passwd, hostname, mac, userip)
							time.Sleep(10 * time.Second)
							log_handle("开始检测联网状态")
							if checkInternetConnection(interfaces1) {
								fmt.Println(interfaces1 + "联网成功\n")
								break
							}
						}
					}
				} else {
					interfaces1 := Config_Content["interfaces1"]
					interfaces2 := Config_Content["interfaces2"]
					ping_1 := checkInternetConnection(interfaces1)
					ping_2 := checkInternetConnection(interfaces2)
					if ping_1 && ping_2 {
						fmt.Println(interfaces1 + "联网成功\n")
						fmt.Println(interfaces2 + "联网成功\n")
					} else {
						userid := Config_Content["username"]
						passwd := Config_Content["password"]
						hostname1 := generateRandomName(12, randomGenerator)
						hostname2 := generateRandomName(12, randomGenerator)
						var userip1, userip2 string
						userip1, err1 := getIP(interfaces1)
						if err1 != nil {
							userip1 = "127.0.0.1"
						}
						userip2, err2 := getIP(interfaces2)
						if err2 != nil {
							userip2 = "127.0.0.1"
						}
						usermac1, err1 := getMac(interfaces1)
						usermac2, err2 := getMac(interfaces2)
						if err1 != nil {
							usermac1 = "88:88:88:88:88:88"
						}
						if err2 != nil {
							usermac2 = "ff:ff:ff:ff:ff:ff"
						}
						curl_dual(userid, passwd, hostname1, hostname2, usermac1, usermac2, userip1, userip2)
						time.Sleep(10 * time.Second)
						log_handle("开始检测联网状态")
						if checkInternetConnection(interfaces1) {
							fmt.Println(interfaces1 + "联网成功\n")
							if checkInternetConnection(interfaces2) {
								fmt.Println(interfaces2 + "联网成功\n")
								break
							}
						}
					}
				}
			}
		}
		time.Sleep(3 * time.Minute)
	}
}

// func main() {
// 	randomGenerator := mrand.New(mrand.NewSource(time.Now().UnixNano()))
// 	Config_Content, Ok := import_config()
// 	if Ok {
// 		interfaces1 := Config_Content["interfaces1"]
// 		userid := Config_Content["username"]
// 		passwd := Config_Content["password"]
// 		hostname := generateRandomName(12, randomGenerator)
// 		userip, err := getIP(interfaces1)
// 		if err != nil {
// 			userip = "127.0.0.1"
// 		}
// 		mac, err := getMac(interfaces1)
// 		if err != nil {
// 			mac = "88:88:88:88:88:88"
// 		}
// 		fmt.Println(interfaces1)
// 		fmt.Println(userid)
// 		fmt.Println(passwd)
// 		fmt.Println(hostname)
// 		fmt.Println(mac)
// 		curl(userid, passwd, hostname, mac, userip)
// 	}
// }
