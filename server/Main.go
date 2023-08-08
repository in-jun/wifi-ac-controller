package main

import (
	_ "embed"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
)

//go:embed HTML/error.html
var errorPage string

var ipMap = make(map[string][]string)
var dataFile = "data.txt"

func main() {
	loadDataFromFile()

	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/ip", handleStringRequest)
	http.HandleFunc("/error", handleError)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("서버 시작 실패: %v", err)
	}
}

func handleError(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%s", errorPage)
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	paths := strings.Split(r.URL.Path, "/")
	if len(paths) >= 3 {
		number, err := strconv.Atoi(paths[1])
		if err != nil {
			log.Printf("URL 경로에 잘못된 번호: %v", err)
			http.Error(w, "잘못된 번호", http.StatusBadRequest)
			return
		}
		publicIP := getClientPublicIP(r)
		redirectURL := getRedirectURL(publicIP, number, paths[2])

		log.Printf("요청: 클라이언트 IP 주소=%s, 경로=%s, 번호=%d, 리다이렉트 URL=%s", publicIP, r.URL.Path, number, redirectURL)

		http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
	} else {
		publicIP := getClientPublicIP(r)
		redirectURL := getRedirectURL(publicIP, 0, paths[1])

		log.Printf("요청: 클라이언트 IP 주소=%s, 경로=%s, 번호=기본값(0), 리다이렉트 URL=%s", publicIP, r.URL.Path, redirectURL)

		http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
	}
}

func getRedirectURL(publicIP string, number int, path string) string {
	privateIPs, exists := ipMap[publicIP]
	if !exists || number < 0 || number >= len(privateIPs) || privateIPs[number] == "" {
		log.Printf("잘못된 번호 또는 맵에 해당 공용 IP 주소가 없습니다. 공용 IP: %s, 번호: %d", publicIP, number)
		return "/error"
	}
	return "http://" + privateIPs[number] + "/" + path
}

func handleStringRequest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("허용되지 않는 메서드: %s", r.Method)
		http.Error(w, "허용되지 않는 메서드", http.StatusMethodNotAllowed)
		return
	}

	numberStr := r.FormValue("number")
	privateIP := r.FormValue("internal_ip")

	number, err := strconv.Atoi(numberStr)
	if err != nil {
		log.Printf("추가: 잘못된 번호: %v", err)
		http.Error(w, "잘못된 번호", http.StatusBadRequest)
		return
	}

	if number < 0 || number >= 3 {
		log.Printf("추가: 잘못된 번호: %v", number)
		http.Error(w, "잘못된 번호", http.StatusBadRequest)
		return
	}

	if !isValidIP(privateIP) {
		log.Printf("잘못된 내부 IP 주소: %s", privateIP)
		http.Error(w, "잘못된 내부 IP 주소", http.StatusBadRequest)
		return
	}

	publicIP := getClientPublicIP(r)
	addPrivateIP(publicIP, number, privateIP)

	log.Printf("추가: 공용 IP 주소=%s, 번호=%d, 내부 IP 주소=%s", publicIP, number, privateIP)
	fmt.Fprintf(w, "공용 IP 주소 %s에 대한 개인 IP 주소 %s를 번호 %d로 추가하였습니다", publicIP, privateIP, number)
}

func addPrivateIP(publicIP string, number int, privateIP string) {
	if _, exists := ipMap[publicIP]; !exists {
		ipMap[publicIP] = make([]string, 3)
	}
	ipMap[publicIP][number] = privateIP

	saveDataToFile()
}

func getClientPublicIP(r *http.Request) string {
	forwarded := r.Header.Get("X-FORWARDED-FOR")
	if forwarded != "" {
		ips := strings.Split(forwarded, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}
	return ip
}

func isValidIP(ip string) bool {
	return net.ParseIP(ip) != nil
}

func saveDataToFile() {
	content := ""

	for publicIP, privateIPs := range ipMap {
		content += publicIP + "," + strings.Join(privateIPs, ",") + "\n"
	}

	err := os.WriteFile(dataFile, []byte(content), 0644)
	if err != nil {
		log.Printf("데이터 파일 저장 실패: %v", err)
	}
}

func loadDataFromFile() {
	if _, err := os.Stat(dataFile); os.IsNotExist(err) {
		return
	}

	content, err := os.ReadFile(dataFile)
	if err != nil {
		log.Printf("데이터 파일 읽기 실패: %v", err)
		return
	}

	lines := strings.Split(string(content), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, ",")
		if len(parts) != 4 {
			continue
		}

		publicIP := parts[0]
		privateIPs := parts[1:]

		ipMap[publicIP] = privateIPs
	}
}
