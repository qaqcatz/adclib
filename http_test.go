package adclib

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"testing"
)

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

// test dependency:
// nanohttpd, 8624
// GET: hello?param1=1&param2=2 -> hello get param1 param2
//      sleep -> sleep 10s
// POST: hello {"param1"=1,"param2"=2} -> hello post {"param1"=1,"param2"=2}
// otherwise: 404

func TestAdbSHttpForwardNormal(t *testing.T) {
	adbs := NewAdbS("/home/hzy/Android/Sdk/platform-tools/adb",
		"emulator-5554","127.0.0.1", "8080")
	// 200
	fmt.Println("200 GET==================================================")
	statusCode, res, err := adbs.HttpForward("8624", "GET",
		"hello?param1=1&param2=2", nil, 0)
	if err != nil {
		t.Fatal("[err]\n" + err.Error())
	} else {
		if statusCode != http.StatusOK {
			t.Fatal("statusCode != 200, res: " + string(res))
		} else {
			t.Log(string(res))
		}
	}
	fmt.Println("200 POST==================================================")
	statusCode, res, err = adbs.HttpForward("8624", "POST",
		"hello", []byte("{\"param1\"=1,\"param2\"=2}"), 0)
	if err != nil {
		t.Fatal("[err]\n" + err.Error())
	} else {
		if statusCode != http.StatusOK {
			t.Fatal("statusCode != 200, res: " + string(res))
		} else {
			t.Log(string(res))
		}
	}
	// 404
	fmt.Println("404==================================================")
	statusCode, res, err = adbs.HttpForward("8624", "GET",
		"xxx", nil, 0)
	if err != nil {
		t.Fatal("[err]\n" + err.Error())
	} else {
		if statusCode != http.StatusNotFound {
			t.Fatal("statusCode != 404, res: " + string(res))
		} else {
			t.Log(string(res))
		}
	}
}

func TestAdbSHttpForwardError(t *testing.T) {
	adbs := NewAdbS("/home/hzy/Android/Sdk/platform-tools/adb",
		"emulator-5554","127.0.0.1", "8080")
	_, _, err := adbs.HttpForward("8625", "GET",
		"hello?param1=1&param2=2", nil, 0)
	if err != nil {
		t.Log("[err]\n" + err.Error())
	} else {
		t.Fatal("must fail")
	}
}

func TestAdbSHttpForwardTimeout(t *testing.T) {
	adbs := NewAdbS("/home/hzy/Android/Sdk/platform-tools/adb",
		"emulator-5554","127.0.0.1", "8080")
	_, _, err := adbs.HttpForward("8624", "GET",
		"sleep", nil, 1000)
	if err != nil {
		switch err.(type) {
		case *url.Error:
			if err.(*url.Error).Timeout() {
				t.Log(err.Error())
			} else {
				t.Fatal("must timeout. error: " + err.Error())
			}
		default:
			t.Fatal("must timeout. error: " + err.Error())
		}
	} else {
		t.Fatal("must fail")
	}
}