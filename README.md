# adclib

[![Go Reference](https://pkg.go.dev/badge/github.com/qaqcatz/adclib.svg)](https://pkg.go.dev/github.com/qaqcatz/adclib)

Provide stable communication interfaces with avd, which encapsulates adb and http request forwarding.

Avoid communication failures caused by device offline and forwarding shutdown as much as possible.

Support timeout.

# How to use

## import

```golang
// go.mod
require github.com/qaqcatz/adclib v1.0.0
// xxx.go
import "github.com/qaqcatz/adclib"
```

## AdbS

```golang
// AdbS: adb -s Device
type AdbS struct {
	AdbPath         string     // e.g. /home/hzy/Android/Sdk/platform-tools/adb
	Device          string     // e.g. emulator-5554
	adbMutex        sync.Mutex // execute commands serially on Device
	Ip              string     // avd's ip, usually 127.0.0.1
	HttpForwardPort string     // forward tcp:HttpForwardPort(host port) tcp:guest port. An avd can only have one forwarding port.
	httpMutex       sync.Mutex // forward http request serially on Device
}
```

### NewAdbS

```golang
func NewAdbS(adbPath string, device string, ip string, httpForwardPort string) *AdbS
```


NewAdbS: create an AdbS object.

- adbPath, e.g. /home/hzy/Android/Sdk/platform-tools/adb
- device, e.g. emulator-5554
- ip, avd's ip, usually 127.0.0.1
- httpForwardPort, forward tcp:HttpForwardPort(host port) tcp:guest port. An avd can only have one forwarding port. If there are multiple guest ports, you should implement port forwarding within the device yourself. If you do not want to use Http, just set any value.

### Exec

```golang
func (adbs *AdbS) Exec(cmdStr string, timoutMs int) ([]byte, []byte, error, bool)
```

Exec: execute adb -s Device cmdStr serially. wait for the result, or timeout, return out stream, error stream, and an error, which can be nil, normal error or *nanoshlib.TimeoutError. timeoutMS <= 0 means timeoutMS = inf.
Exec will execute adb -s Device reconnect automatically according to the device status before executing cmdStr. If a reconnection has occurred, the last parameter will return true, otherwise false. May wait 0~5.5+timoutMs

**Example**

```golang
func ExampleAdbS_Exec_normal() {
	// emulator -avd test -port 5554 first
	adbs := NewAdbS("/home/hzy/Android/Sdk/platform-tools/adb", "emulator-5554","127.0.0.1", "0")
	outStream, errStream, err, rec := adbs.Exec("shell echo hello", 1000)
	if err != nil {
		log.Fatal("[err]\n" + err.Error() +
			"\n[out stream]\n" + string(outStream) + "\n[err stream]\n" + string(errStream))
	} else if rec {
		log.Fatal("why reconnect?")
	} else {
		log.Println("[out stream]\n" + string(outStream))
		log.Println("[err stream]\n" + string(errStream))
	}
}

func ExampleAdbS_Exec_errorAndReconnect() {
	// close avd first
	adbs := NewAdbS("/home/hzy/Android/Sdk/platform-tools/adb", "emulator-5554","127.0.0.1", "0")
	outStream, errStream, err, rec := adbs.Exec("shell echo hello", 1000)
	if err != nil && rec {
		log.Println("[err]\n" + err.Error() +
			"\n[out stream]\n" + string(outStream) + "\n[err stream]\n" + string(errStream))
	} else if !rec {
		log.Fatal("must reconnect")
	} else {
		log.Fatal("must fail")
	}
}

func ExampleAdbS_Exec_timeout() {
	// emulator -avd test -port 5554 first
	adbs := NewAdbS("/home/hzy/Android/Sdk/platform-tools/adb", "emulator-5554","127.0.0.1","0")
	_, _, err, rec := adbs.Exec("shell sleep 10s", 1000)
	if err != nil {
		if rec {
			log.Fatal("why reconnect?")
		} else {
			switch err.(type) {
			case *nanoshlib.TimeoutError:
				log.Println(err.Error())
			default:
				log.Fatal("adb shell sleep 10s must timeout")
			}
		}
	} else {
		log.Fatal("adb shell must fail")
	}
}
```

### HttpForward

```golang
func (adbs *AdbS) HttpForward(guestPort string, method string, paramUrl string, jsonData []byte, timeoutMS int) (int, []byte, error)
```


HttpForward: forward tcp: HttpForwardPort tcp:guestPort, then forward the http request serially. An avd can only have one forwarding port. If there are multiple guest ports, you should implement port forwarding inside the device yourself.

- guestPort. guest port inside the device
- method, request method, only support GET and POST.
- paramUrl, e.g. (no '/')hello?param1=1&param2=2
- jsonData, GET: nil, POST: content in json format
- timeoutMS: timeout(millisecond)
return status code, result, error. The error type during GET/POST is *url.Error

**Example**

```golang
func ExampleAdbS_HttpForward_normal() {
	// test dependency:
	// nanohttpd, 8624
	// GET: hello?param1=1&param2=2 -> hello get param1 param2
	// POST: hello {"param1"=1,"param2"=2} -> hello post {"param1"=1,"param2"=2}
	// otherwise: 404
	adbs := NewAdbS("/home/hzy/Android/Sdk/platform-tools/adb",
		"emulator-5554","127.0.0.1", "8080")
	// 200
	fmt.Println("200 GET==================================================")
	statusCode, res, err := adbs.HttpForward("8624", "GET",
		"hello?param1=1&param2=2", nil, 0)
	if err != nil {
		log.Fatal("[err]\n" + err.Error())
	} else {
		if statusCode != http.StatusOK {
			log.Fatal("statusCode != 200, res: " + string(res))
		} else {
			log.Println(string(res))
		}
	}
	fmt.Println("200 POST==================================================")
	statusCode, res, err = adbs.HttpForward("8624", "POST",
		"hello", []byte("{\"param1\"=1,\"param2\"=2}"), 0)
	if err != nil {
		log.Fatal("[err]\n" + err.Error())
	} else {
		if statusCode != http.StatusOK {
			log.Fatal("statusCode != 200, res: " + string(res))
		} else {
			log.Println(string(res))
		}
	}
	// 404
	fmt.Println("404==================================================")
	statusCode, res, err = adbs.HttpForward("8624", "GET",
		"xxx", nil, 0)
	if err != nil {
		log.Fatal("[err]\n" + err.Error())
	} else {
		if statusCode != http.StatusNotFound {
			log.Fatal("statusCode != 404, res: " + string(res))
		} else {
			log.Println(string(res))
		}
	}
}

func ExampleAdbS_HttpForward_timeout() {
	// test dependency:
	// nanohttpd, 8624
	// GET: sleep -> sleep 10s
	adbs := NewAdbS("/home/hzy/Android/Sdk/platform-tools/adb",
		"emulator-5554","127.0.0.1", "8080")
	_, _, err := adbs.HttpForward("8624", "GET",
		"sleep", nil, 1000)
	if err != nil {
		switch err.(type) {
		case *url.Error:
			if err.(*url.Error).Timeout() {
				log.Println(err.Error())
			} else {
				log.Fatal("must timeout. error: " + err.Error())
			}
		default:
			log.Fatal("must timeout. error: " + err.Error())
		}
	} else {
		log.Fatal("must fail")
	}
}
```



