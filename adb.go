package adclib

import (
	"errors"
	"github.com/qaqcatz/nanoshlib"
	"strings"
	"sync"
	"time"
)

// AdbS: adb -s Device
type AdbS struct {
	AdbPath         string     // e.g. /home/hzy/Android/Sdk/platform-tools/adb
	Device          string     // e.g. emulator-5554
	adbMutex        sync.Mutex // execute commands serially on Device
	Ip              string     // avd's ip, usually 127.0.0.1
	HttpForwardPort string     // forward tcp:HttpForwardPort(host port) tcp:guest port. An avd can only have one forwarding port.
	httpMutex       sync.Mutex // forward http request serially on Device
}

// NewAdbS: create an AdbS object.
//
// - adbPath, e.g. /home/hzy/Android/Sdk/platform-tools/adb
//
// - device, e.g. emulator-5554
//
// - ip, avd's ip, usually 127.0.0.1
//
// - httpForwardPort, forward tcp:HttpForwardPort(host port) tcp:guest port.
// An avd can only have one forwarding port.
// If there are multiple guest ports, you should implement port forwarding within the emulator yourself.
// If you do not want to use Http, just set any value.
func NewAdbS(adbPath string, device string, ip string, httpForwardPort string) *AdbS {
	return &AdbS{AdbPath: adbPath, Device: device, Ip: ip, HttpForwardPort: httpForwardPort}
}

// Exec: execute adb -s Device cmdStr serially. wait for the result, or timeout,
// return out stream, error stream, and an error, which can be nil, normal error or *nanoshlib.TimeoutError.
// timeoutMS <= 0 means timeoutMS = inf.
//
// Exec will execute adb -s Device reconnect automatically according to the device status before executing cmdStr.
// If a reconnection has occurred, the last parameter will return true, otherwise false.
// May wait 0~5.5+timoutMs
func (adbs *AdbS) Exec(cmdStr string, timoutMs int) ([]byte, []byte, error, bool) {
	adbs.adbMutex.Lock()
	defer adbs.adbMutex.Unlock()
	isReconnect := false
	if !adbs.isAlive() {
		err := adbs.reconnect()
		isReconnect = true
		if err != nil {
			return nil, nil, errors.New("adb is not alive and reconnect error: " + err.Error()), isReconnect
		}
	}
	outStream, errStream, err := nanoshlib.Exec(adbs.AdbPath+" -s "+adbs.Device+" "+cmdStr, timoutMs)
	return outStream, errStream, err, isReconnect
}

// isAlive: use adb -s Device get-state to check if adb is alive.
// May wait 0~1s.
func (adbs *AdbS) isAlive() bool {
	outStream, _, err := nanoshlib.Exec(adbs.AdbPath+" -s "+adbs.Device+" get-state", 1000)
	if err != nil || !strings.HasSuffix(strings.TrimSpace(string(outStream)), "device") {
		// why use strings.HashSuffix?
		// sometimes the outStream may be:
		// * daemon not running; starting now at tcp:5037
		// * daemon started successfully
		// device
		return false
	}
	return true
}

// reconnect: use adb -s Device reconnect to reconnect adb.
// If the adb is alive after reconnection, then return nil, otherwise return a reconnect error.
// May wait 0~4.5s.
func (adbs *AdbS) reconnect() error {
	_, _, err := nanoshlib.Exec(adbs.AdbPath+" -s "+adbs.Device+" reconnect", 1000)
	if err != nil {
		// adb -s Device reconnect usually only prints the error stream,
		// but still exit 0, so this error detection usually does not work
		return err
	}
	// adb -s Device reconnect normally does not block, so we need to wait until the adb is alive,
	// or timeout
	for i := 0; i < 3; i++ {
		time.Sleep(500 * time.Millisecond)
		if adbs.isAlive() {
			return nil
		}
	}
	return errors.New("reconnect timeout")
}
