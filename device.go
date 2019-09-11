package ssmt

type Device struct {
	DeviceName string
	ModelType  string
	Screen     string
	IMEI       string
	IMSI       string
	UserAgent  string
}

const (
	defaultDevice    = "Android,26,8.0 release-keys"
	defaultScreen    = "1080x1920"
	DefaultUserAgent = "Dalvik/2.2.0 (Linux; U; Android 8.0)"
)

func GenerateDevice() *Device {
	imei := GenerateIMEI()
	return &Device{
		DeviceName: defaultDevice,
		ModelType:  RandModel(),
		Screen:     defaultScreen,
		IMEI:       imei,
		IMSI:       imei, // same as imei when couldn't get
		UserAgent:  DefaultUserAgent,
	}
}
