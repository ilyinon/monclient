package main

import (
	"bytes"
	"fmt"
	"net/http"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"strings"
//	"syscall"
)


func check(e error) {
	if e != nil {
		panic(e)
	}
}

// get local hostname
func get_node_name() string {
	fqdn, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	return fqdn
}

// get winenv name
func get_winenv_name() string {
		data, err := ioutil.ReadFile("c:\\salt\\conf\\minion.d\\grain.conf")
		if err != nil {
		    return "no_data"
		   }
		file := string(data)
		re := regexp.MustCompile(`^    node_env\.`)
		s := ""
		return strings.TrimSpace(re.ReplaceAllString(strings.Split(file, ":")[2], s))
}

func get_java_version() string {
        out, err := exec.Command("c:\\ProgramData\\Oracle\\Java\\javapath\\java.exe", "-version").CombinedOutput()
		check(err)
		temp := string(out)
		re := regexp.MustCompile(`^java version\.`)
        s := ""
		return re.ReplaceAllString(strings.Split(strings.Split(temp, " ")[2], "\"")[1], s)
}

func get_chrome_version() string {
        out, err := exec.Command("wmic", "datafile", "where", `name="C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe"`, "get", "Version", "/value").Output()
        check(err)
        temp := string(out)
        re := regexp.MustCompile(`^Version=\.`)
		s := ""
		return strings.TrimSpace(re.ReplaceAllString(strings.Split(temp, "=")[1], s))
}

func get_chromedriver_version() string {
        out, err := exec.Command("c:\\selenium\\chromedriver.exe", "/version").Output()
        check(err)
        temp := string(out)
        re := regexp.MustCompile(`^ChromeDriver\.`)
		s := ""
		return re.ReplaceAllString(strings.Split(temp, " ")[1], s)
}

func get_firefox_version() string {
        out, err := exec.Command("C:\\Program Files\\Mozilla Firefox\\firefox.exe", "-v", "|more").Output()
        check(err)
        temp := string(out)
        re := regexp.MustCompile(`^Mozilla Firefox\.`)
		s := ""
		return strings.TrimSpace(re.ReplaceAllString(strings.Split(temp, " ")[2], s))
}

func get_gecko_version() string {
        out, err := exec.Command("c:\\selenium\\geckodriver.exe", "--version").Output()
        check(err)
        temp := string(out)
        re := regexp.MustCompile(`^geckodriver\.`)
		s := ""
		return re.ReplaceAllString(strings.Split(strings.Split(temp, " ")[1], "\n")[0], s)
}

func get_python_version() string {
        out, err := exec.Command("C:\\Python27\\python.exe", "--version").CombinedOutput()
        check(err)
        temp := string(out)
        re := regexp.MustCompile(`^Python\.`)
		s := ""
    return strings.TrimSpace(re.ReplaceAllString(strings.Split(temp, " ")[1], s))
}

func get_selenium_version() string {
        out, err := exec.Command("C:\\Python27\\python.exe", "-m", "pip", "freeze").Output()
        check(err)
        re := regexp.MustCompile(".*selenium.*")
        s := ""
        match := re.FindStringSubmatch(string(out))
    return strings.TrimSpace(re.ReplaceAllString(strings.Split(match[0], "==")[1], s))
}


func get_windows_activated() string {
        out, err := exec.Command("cscript.exe" ,"//nologo", "C:\\Windows\\system32\\slmgr.vbs", "/dli").Output()
        check(err)
        re := regexp.MustCompile(".*License Status.*")
        s := ""
        match := re.FindStringSubmatch(string(out))
    return strings.TrimSpace(re.ReplaceAllString(strings.Split(match[0], ":")[1], s))
}

func get_windows_version() string {
        out, err := exec.Command("wmic" ,"os", "get", "Caption", "/value").Output()
        check(err)
        temp := string(out)
        re := regexp.MustCompile(`^Caption\.`)
		s := ""
		return strings.TrimSpace(re.ReplaceAllString(strings.Split(strings.Split(temp, "=")[1], "\n")[0], s))
}

func main() {
	url := "http://cmt001.infra.mycompany.com:80/winstatus/"
    my_node_name    := get_node_name()
    my_winenv       := get_winenv_name()
    my_java         := get_java_version()
    my_chrome       := get_chrome_version()
    my_chromedriver := get_chromedriver_version()
    my_firefox      := get_firefox_version()
    my_gecko        := get_gecko_version()
    my_python       := get_python_version()
    my_selenium     := get_selenium_version()
    my_activated    := get_windows_activated()
    my_winver       := get_windows_version()

	fmt.Println("URL                   :>", url)
	fmt.Println("node_name             :>", my_node_name)
	fmt.Println("winenv                :>", my_winenv)
	fmt.Println("java_version          :>", my_java)
	fmt.Println("chrome_version        :>", my_chrome)
	fmt.Println("chromedriver_version  :>", my_chromedriver)
	fmt.Println("firefox_version       :>", my_firefox)
	fmt.Println("gecko_version         :>", my_gecko)
	fmt.Println("python_version        :>", my_python)
	fmt.Println("selenium_version      :>", my_selenium)
	fmt.Println("windows_activated     :>", my_activated)
	fmt.Println("windows_version       :>", my_winver)

	// preparation of json for sending
	var jsonStr = []byte(`
{"node_name"       : "` + my_node_name + `",
"winenv"          : "` + my_winenv    + `",
"java_version"    : "` + my_java      + `",
"chrome_version"  : "` + my_chrome    + `",
"chromedriver_version": "` + my_chromedriver + `",
"firefox_version" : "` + my_firefox   + `",
"gecko_version"   : "` + my_gecko     + `",
"python_version"  : "` + my_python    + `",
"selenium_version": "` + my_selenium  + `",
"windows_activated": "` + my_activated  + `",
"windows_version": "` + my_winver + `"}`)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

	// adding headers
	req.Header.Add("User-Agent", "winclient-v1")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}

    //
    // fmt.Println("request", req)
	// send request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	// close connection
	defer resp.Body.Close()

	// show answer,  code and header
	fmt.Println("response Status:", resp.Status)
	//fmt.Println("response Headers:", resp.Header)

}
