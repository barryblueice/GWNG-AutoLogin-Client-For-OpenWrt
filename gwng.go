package main

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"log"
	mrand "math/rand"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"

	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func checkInternetConnection() bool {
	host := "www.baidu.com"
	port := 80
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return false
	}
	defer conn.Close()
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

func getIP() (map[string]string, error) {
	networkPath := "/etc/config/network"
	file, err := os.Open(networkPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	networkList := make(map[string]string)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.ReplaceAll(scanner.Text(), "\t", "")
		line = strings.ReplaceAll(line, "'", "")
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "config interface ") {
			iface := strings.TrimPrefix(line, "config interface ")
			iface += " "
			networkList[iface[:len(iface)-1]] = ""
		} else if strings.HasPrefix(line, "option ipaddr ") {
			line = strings.ReplaceAll(line, " ", ":")
			line = strings.TrimPrefix(line, "option:ipaddr:")
			line += " "
			lastKey := ""
			for k := range networkList {
				lastKey = k
			}
			if lastKey != "" {
				networkList[lastKey] = strings.TrimSpace(line)
			} else {
				return nil, fmt.Errorf("found 'option ipaddr' line without preceding 'config interface' line")
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return networkList, nil
}

func getMac(iface string) (string, error) {
	// 执行 ifconfig 命令
	cmd := exec.Command("ifconfig", iface)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	// 解析 ifconfig 的输出以找到 MAC 地址
	re := regexp.MustCompile(`([[:xdigit:]]{1,2}:){5}[[:xdigit:]]{1,2}`)
	match := re.FindString(out.String())
	if match == "" {
		return "", fmt.Errorf("no MAC address found for interface %s", iface)
	}

	// 转换为小写并返回
	return strings.ToLower(match), nil
}

func generateRandomMac(separator string) (string, error) {
	mac := make([]byte, 6)
	if _, err := rand.Read(mac); err != nil {
		return "", err
	}

	// Set the locally administered bit and clear the multicast bit
	mac[0] |= 0x02
	mac[0] &^= 0x01

	// Convert the MAC address bytes to a formatted string
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
	requestURL := "http://10.244.0.13/quickauth.do"
	params := url.Values{}
	params.Add("userid", userid)
	params.Add("passwd", passwd)
	params.Add("wlanuserip", userip)
	params.Add("wlanacname", "gwng")
	params.Add("wlanacIp", "10.244.0.1")
	params.Add("ssid", "")
	params.Add("vlan", "502")
	params.Add("mac", mac)
	params.Add("version", "0")
	params.Add("portalpageid", "5")
	params.Add("timestamp", gentimestamp)
	params.Add("uuid", generatedUUID)
	params.Add("portaltype", "1")
	params.Add("hostname", hostname)
	headers := map[string]string{
		"Accept":           "application/json, text/javascript, */*; q=0.01",
		"Accept-Language":  "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6",
		"Connection":       "keep-alive",
		"Cookie":           fmt.Sprintf("%s|%s|%s; ABMS=c0c701f6-ea15-448e-a48b-6230c9117317", generatedMac1, generatedMac2, mac),
		"Referer":          fmt.Sprintf("http://10.244.0.13/portal.do?wlanuserip=%s&wlanacname=gwng&mac=%s&vlan=502&hostname=%s&rand=%s&url=http://www.msftconnecttest.com/redirec", userip, mac, hostname, generatedSeed),
		"User-Agent":       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36 Edg/117.0.2045.31",
		"X-Requested-With": "XMLHttpRequest",
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", requestURL+"?"+params.Encode(), nil)
	if err != nil {
		log.Fatal(err)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log_handle(string(body))
}

func curl_dual(userid string, passwd string, hostname string, mac1 string, mac2 string, userip1 string, userip2 string) {

	randomGenerator := mrand.New(mrand.NewSource(time.Now().UnixNano()))
	generatedMac, err := generateRandomMac(":")
	if err != nil {
		generatedMac = "as:as:as:as:as:as"
	}
	generatedUUID := gen_uuid()
	generatedSeed := generateRandomSeed(14, randomGenerator)
	gentimestamp := get_timestamp()
	requestURL := "http://10.244.0.13/quickauth.do"
	params := url.Values{}
	params.Add("userid", userid)
	params.Add("passwd", passwd)
	params.Add("wlanuserip", userip1)
	params.Add("wlanacname", "gwng")
	params.Add("wlanacIp", "10.244.0.1")
	params.Add("ssid", "")
	params.Add("vlan", "502")
	params.Add("mac", mac1)
	params.Add("version", "0")
	params.Add("portalpageid", "5")
	params.Add("timestamp", gentimestamp)
	params.Add("uuid", generatedUUID)
	params.Add("portaltype", "1")
	params.Add("hostname", hostname)
	headers := map[string]string{
		"Accept":           "application/json, text/javascript, */*; q=0.01",
		"Accept-Language":  "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6",
		"Connection":       "keep-alive",
		"Cookie":           fmt.Sprintf("%s|%s|%s; ABMS=c0c701f6-ea15-448e-a48b-6230c9117317", generatedMac, mac1, mac2),
		"Referer":          fmt.Sprintf("http://10.244.0.13/portal.do?wlanuserip=%s&wlanacname=gwng&mac=%s&vlan=502&hostname=%s&rand=%s&url=http://www.msftconnecttest.com/redirec", userip1, mac1, hostname, generatedSeed),
		"User-Agent":       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36 Edg/117.0.2045.31",
		"X-Requested-With": "XMLHttpRequest",
	}
	client := &http.Client{}
	req, err := http.NewRequest("GET", requestURL+"?"+params.Encode(), nil)
	if err != nil {
		log.Fatal(err)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	log_handle(string(body))

	generatedMac, err = generateRandomMac(":")
	if err != nil {
		generatedMac = "as:as:as:as:as:as"
	}
	generatedUUID = gen_uuid()
	generatedSeed = generateRandomSeed(14, randomGenerator)
	gentimestamp = get_timestamp()
	requestURL = "http://10.244.0.13/quickauth.do"
	params = url.Values{}
	params.Add("userid", userid)
	params.Add("passwd", passwd)
	params.Add("wlanuserip", userip2)
	params.Add("wlanacname", "gwng")
	params.Add("wlanacIp", "10.244.0.1")
	params.Add("ssid", "")
	params.Add("vlan", "502")
	params.Add("mac", mac2)
	params.Add("version", "0")
	params.Add("portalpageid", "5")
	params.Add("timestamp", gentimestamp)
	params.Add("uuid", generatedUUID)
	params.Add("portaltype", "1")
	params.Add("hostname", hostname)
	headers = map[string]string{
		"Accept":           "application/json, text/javascript, */*; q=0.01",
		"Accept-Language":  "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6",
		"Connection":       "keep-alive",
		"Cookie":           fmt.Sprintf("%s|%s|%s; ABMS=c0c701f6-ea15-448e-a48b-6230c9117317", generatedMac, mac1, mac2),
		"Referer":          fmt.Sprintf("http://10.244.0.13/portal.do?wlanuserip=%s&wlanacname=gwng&mac=%s&vlan=502&hostname=%s&rand=%s&url=http://www.msftconnecttest.com/redirec", userip2, mac2, hostname, generatedSeed),
		"User-Agent":       "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/117.0.0.0 Safari/537.36 Edg/117.0.2045.31",
		"X-Requested-With": "XMLHttpRequest",
	}
	client = &http.Client{}
	req, err = http.NewRequest("GET", requestURL+"?"+params.Encode(), nil)
	if err != nil {
		log.Fatal(err)
	}
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	resp, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	log_handle(string(body))
}

func log_handle(message string) {
	logDir := "/var/log/gwng"
	logFile := "gwng.log"
	logFilePath := filepath.Join(logDir, logFile)

	// 创建日志目录（如果不存在）
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	// 打开（或创建）日志文件
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
			if checkInternetConnection() {
				log_handle("联网成功")
			} else {
				for {
					log_handle("联网失败，正在自动联网……")
					Config_Content, Ok := import_config()
					if Ok {
						num := Config_Content["num"]
						if num == "1" {
							interfaces1 := Config_Content["interfaces1"]
							userid := Config_Content["username"]
							passwd := Config_Content["password"]
							hostname := generateRandomName(12, randomGenerator)
							IP_list, err := getIP()
							var userip string
							if err != nil {
								userip = "127.0.0.1"
							} else {
								userip = IP_list[Config_Content["interfaces1"]]
							}
							mac, err := getMac(interfaces1)
							if err != nil {
								mac = "88:88:88:88:88:88"
							}
							curl(userid, passwd, hostname, mac, userip)
						} else {
							interfaces1 := Config_Content["interfaces1"]
							interfaces2 := Config_Content["interfaces2"]
							userid := Config_Content["username"]
							passwd := Config_Content["password"]
							hostname1 := generateRandomName(12, randomGenerator)
							hostname2 := generateRandomName(12, randomGenerator)
							var userip1, userip2 string
							IP_list, err := getIP()
							if err != nil {
								userip1 = "127.0.0.1"
								userip2 = "127.0.0.1"
							} else {
								userip1 = IP_list[Config_Content["interfaces1"]]
								userip2 = IP_list[Config_Content["interfaces2"]]
							}
							usermac1, err1 := getMac(interfaces1)
							usermac2, err2 := getMac(interfaces2)
							if err1 != nil {
								usermac1 = "88:88:88:88:88:88"
							}
							if err2 != nil {
								usermac2 = "ff:ff:ff:ff:ff:ff"
							}
							curl_dual(userid, passwd, hostname1, usermac1, usermac2, userip1, userip2)
							curl_dual(userid, passwd, hostname2, usermac1, usermac2, userip1, userip2)
						}
					}
					if checkInternetConnection() {
						log_handle("联网成功")
						break
					}
				}
			}
		}
		time.Sleep(3 * time.Minute)
	}
}
