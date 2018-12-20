package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"syscall"
)

// get local hostname
func get_fqdn() string {
	fqdn, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	return fqdn
}

// get type of service by hostnmae
func get_service(fqdn string) string {

	var rVCC = regexp.MustCompile(`^(qa|tag|pm)+(core|api|apiewbe|apiewfe|apife|cp|dsb|env|ftp|envws|rpt|bil)+[0-9]{3}[a-z]{0,1}\..*$`)
	var rFDM = regexp.MustCompile(`^(qa|tag|pm)+(dsbfrm|frmapp|frmchat|frmstream|frmsup|frmvcc|frmweb|frmauth|frmstat)+[\d]{3}[a-z]{0,1}\..*$`)
	var rSCC = regexp.MustCompile(`^(qa|tag|pm)+(sccnas|sccrpt|sccsvc|sccsync|sccws|qasocial)+[\d]{3}[a-z]{0,1}\..*$`)
//	var rBLK = regexp.MustCompile(`^(qa|tag|pm)+(vccrmq|db|cplb|apilb|fdsdb|db|frmaerosup|frmrmq|frmaero|lb|sccdb|loginbe|loginlb|nlu|rules|nss|sip|sfgw)+[\d]{3}[a-z]{0,1}\..*$`)

	// what service is placed on this host type
	switch {
	case rVCC.MatchString(fqdn):
		return "VCC"
	case rFDM.MatchString(fqdn):
		return "FDM"
	case rSCC.MatchString(fqdn):
		return "SCC"
	}
	return "UNKNOWN"
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// get version by type of service
func get_version(host_service string) string {

	// get version regarding to service type
	switch host_service {
	case "VCC":
		out, err := exec.Command("/bin/rpm", "-qa").Output()
		check(err)
		file := string(out)
		temp := strings.Split(file, "\n")
		var rVERSION_VCC = regexp.MustCompile(`^vcc-.*noarch`)

		for _, item := range temp {
			switch {
			case rVERSION_VCC.MatchString(item):
				return strings.Split(item, "-")[1]
			}
		}
		return "VCC_unknown_version"

	case "FDM":
		// FDM storages its version here
		data, err := ioutil.ReadFile("/myuser/tmp/build.info")
		check(err)
		file := string(data)
		temp := strings.Split(file, "\n")
		var rVERSION_FDM = regexp.MustCompile(`^FREEDOM_version=.*`)

		// find string which contains version
		for _, item := range temp {
			switch {
			case rVERSION_FDM.MatchString(item):
				// shrink unwanted part of string
				re := regexp.MustCompile("^.*=")
				s := ""
				return re.ReplaceAllString(item, s)
			}
		}

		return "FDM_unknown_version"

	case "SCC":
		// SCC version gathering
		data, err := ioutil.ReadFile("/myuser/tmp/build.info")
		check(err)
		file := string(data)
		re := regexp.MustCompile(`^(SCC|SCC_MAINSTREAM_DAILY)\.`)
		s := ""
		return re.ReplaceAllString(strings.Split(strings.Split(file, "=")[1], "-")[0], s)

	}
	return "UNKNOWN"
}

func get_status(host_service string) string {

	switch host_service {
	// get status VCC by request of service in init.d
	case "VCC", "FDM", "SCC":

	    cmd := exec.Command("/etc/init.d/"+strings.ToLower(host_service), "status")
	    // var waitStatus syscall.WaitStatus

	    if err := cmd.Run(); err != nil {
	        if err != nil {
	            os.Stderr.WriteString(fmt.Sprintf("Error: %s\n", err.Error()))
	            }
        // FAIL
            if exitError, ok := err.(*exec.ExitError); ok {
                _ = exitError.Sys().(syscall.WaitStatus)
         //       fmt.Printf("Output: %s\n", []byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
            }
         } else {
         // SUCCESS
            _ = cmd.ProcessState.Sys().(syscall.WaitStatus)
         //   fmt.Printf("Output: %s\n", []byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
            return "True"
         }


    }
	return "False"
}

func main() {
	url := "http://cmt001.infra.mycompany.com/status/"
    host_fqdn    := get_fqdn()
    host_service := get_service(host_fqdn)
    host_version := get_version(host_service)
    host_status  := get_status(host_service)

	fmt.Println("URL     :>", url)
	fmt.Println("FQDN    :>", host_fqdn)
	fmt.Println("SERVICE :>", host_service)
	fmt.Println("VERSION :>", host_version)
	fmt.Println("STATUS  :>", host_status)

	// preparation of json for sending
	var jsonStr = []byte(`
    {
        "server" : "` + host_fqdn + `",
        "service": "` + host_service + `",
        "status" : "` + host_status + `",
        "version": "` + host_version + `"
        }
        `)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

	// adding headers
	req.Header.Add("User-Agent", "monclient-v1")
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}

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
