package adclib

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"time"
)

// HttpForward: forward tcp: HttpForwardPort tcp:guestPort, then forward the http request serially.
// An avd can only have one forwarding port.
// If there are multiple guest ports, you should implement port forwarding inside the emulator yourself.
//
// - guestPort. guest port inside the emulator
//
// - method, request method, only support GET and POST.
//
// - paramUrl, e.g. (no '/')hello?param1=1&param2=2
//
// - jsonData, GET: nil, POST: content in json format
//
// - timeoutMS: timeout(millisecond)
//
// return status code, result, error. The error type during GET/POST is *url.Error
func (adbs *AdbS) HttpForward(guestPort string,  method string,
	paramUrl string, jsonData []byte, timeoutMS int) (int, []byte, error) {
	adbs.httpMutex.Lock()
	defer adbs.httpMutex.Unlock()
	// Make sure the forwarding service is alive
	_, _, err, _ := adbs.Exec("forward tcp:"+adbs.HttpForwardPort+" tcp:"+guestPort, 1000)
	if err != nil {
		return -1, nil, errors.New("forward error: " + err.Error())
	}
	// http request
	var resp *http.Response = nil
	client := http.Client {
		Timeout: time.Duration(timeoutMS) * time.Millisecond,
	}
	if method == "GET" { // GET
		resp, err = client.Get("http://" + adbs.Ip + ":" + adbs.HttpForwardPort + "/" + paramUrl)
		if err != nil {
			return -1, nil, err
		}
	} else if method == "POST" { // POST
		resp, err = client.Post("http://"+ adbs.Ip + ":" + adbs.HttpForwardPort + "/" + paramUrl, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			return -1, nil, err
		}
	} else {
		return -1, nil, errors.New("unknown http method")
	}
	ans, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return -1, nil, errors.New("read body error: " + err.Error())
	}
	return resp.StatusCode, ans, nil
}