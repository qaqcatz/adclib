package adclib

import (
	"github.com/qaqcatz/nanoshlib"
	"log"
	"testing"
)

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

// emulator -avd test -port 5554 first
func TestAdbSExecNormal(t *testing.T) {
	adbs := NewAdbS("/home/hzy/Android/Sdk/platform-tools/adb", "emulator-5554","127.0.0.1", "0")
	outStream, errStream, err, rec := adbs.Exec("shell echo hello", 1000)
	if err != nil {
		t.Fatal("[err]\n" + err.Error() +
			"\n[out stream]\n" + string(outStream) + "\n[err stream]\n" + string(errStream))
	} else if rec {
		t.Fatal("why reconnect?")
	} else {
		t.Log("[out stream]\n" + string(outStream))
		t.Log("[err stream]\n" + string(errStream))
	}
}

// emulator -avd test -port 5554 first
func TestAdbSExecTimeout(t *testing.T) {
	adbs := NewAdbS("/home/hzy/Android/Sdk/platform-tools/adb", "emulator-5554","127.0.0.1","0")
	_, _, err, rec := adbs.Exec("shell sleep 10s", 1000)
	if err != nil {
		if rec {
			t.Fatal("why reconnect?")
		} else {
			switch err.(type) {
			case *nanoshlib.TimeoutError:
				t.Log(err.Error())
			default:
				t.Fatal("adb shell sleep 10s must timeout")
			}
		}
	} else {
		t.Fatal("adb shell must fail")
	}
}

// close avd first
func TestAdbSExecErrorAndReconnect(t *testing.T) {
	adbs := NewAdbS("/home/hzy/Android/Sdk/platform-tools/adb", "emulator-5554", "127.0.0.1","0")
	outStream, errStream, err, rec := adbs.Exec("shell echo hello", 1000)
	if err != nil && rec {
		t.Log("[err]\n" + err.Error() +
			"\n[out stream]\n" + string(outStream) + "\n[err stream]\n" + string(errStream))
	} else if !rec {
		t.Fatal("must reconnect")
	} else {
		t.Fatal("must fail")
	}
}