package lib

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Session struct {
	UserID             int64
	TokenID            string
	UserExpirationTime time.Time
	UserInfo           UserInfo
	UserAgent          string
	LimitParams        *LimitParams
	PhoneIMEI          string
	PhoneModel         string
}
type httpError struct {
	msg     string
	httpErr error
}

func (e httpError) Error() string {
	return e.msg + "\n" + e.httpErr.Error()
}

const (
	server            = "http://www.ccxyct.com:8080"
	loginURL          = server + "/sunShine_Sports/loginSport.action"
	uploadDataURL     = server + "/sunShine_Sports/xtUploadData.action"
	postTestDataURL   = server + "/sunShine_Sports/postTestData.action"
	getSportResultURL = server + "/sunShine_Sports/xtGetSportResult.action"
	getAppInfoURL     = server + "/sunShine_Sports/xtGetAppInfo.action"
	DefaultUserAgent  = "Dalvik/2.1.0 (Linux; U; Android 7.0)"

	defaultSchoolId = 60
	defaultIMSI     = "1234567890"
	defaultDevice   = "Android,25,7.1.2"
	appVersion      = "2.2.6"
	appVersionID    = 13
)

var (
	serviceErrorTable = make(map[int64]error)
	ErrTokenExpired   = errors.New("超时，请重新登录")
)

func init() {
	serviceErrorTable[1] = nil
	serviceErrorTable[2] = ErrTokenExpired
}

func translateServiceError(statusCode int64, statusMessage string) error {
	err, exist := serviceErrorTable[statusCode]
	if exist {
		return err
	}
	return fmt.Errorf("response status %d , message: %s", statusCode, statusMessage)
}

func CreateSession() *Session {
	return &Session{UserID: 0, TokenID: "", UserAgent: DefaultUserAgent, PhoneIMEI: GenerateIMEI(), PhoneModel: RandModel()}
}

// DEPRECATED: use LoginEx instead
func (s *Session) Login(stuNum string, phoneNum string, passwordHash string) (e error) {
	return s.LoginEx(stuNum, phoneNum, passwordHash, defaultSchoolId)
}
func (s *Session) LoginEx(stuNum string, phoneNum string, passwordHash string, schoolID int64) (e error) {
	req, err := http.NewRequest(http.MethodPost, loginURL, strings.NewReader(url.Values{
		"stuNum":   {stuNum},
		"phoneNum": {phoneNum},
		"passWd":   {passwordHash},
		"schoolId": {strconv.FormatInt(schoolID, 10)},
		"stuId":    {"1"},
		"token":    {""},
	}.Encode()))
	if err != nil {
		return httpError{"HTTP Create Request Failed.", err}
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header["UserID"] = []string{"0"}
	req.Header["TokenID"] = []string{""}
	req.Header["app"] = []string{"com.ccxyct.sunshinemotion"}
	req.Header["ver"] = []string{appVersion}
	req.Header["device"] = []string{defaultDevice}
	req.Header["model"] = []string{s.PhoneModel}
	req.Header["screen"] = []string{"1080x1920"}
	req.Header["imei"] = []string{s.PhoneIMEI}
	req.Header["imsi"] = []string{defaultIMSI}
	req.Header["crack"] = []string{"0"}
	req.Header["latitude"] = []string{"0.0"}
	req.Header["longitude"] = []string{"0.0"}
	req.Header["crack"] = []string{"0"}
	req.Header.Set("User-Agent", s.UserAgent)
	resp, err := http.DefaultClient.Do(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return httpError{"HTTP Send Request Failed! ", err}
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP Response Status: %d(%s)", resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("HTTP Read Resp Failed! %s", err.Error())
	}

	var loginResult struct {
		Status             int64
		ErrorMessage       string
		UserID             int64
		TokenID            string
		UserExpirationTime int64
		UserInfo           UserInfo
	}
	err = json.Unmarshal(respBytes, &loginResult)
	if err != nil {
		return fmt.Errorf("reslove Failed. %s %s", err.Error(), string(respBytes))
	}

	err = translateServiceError(loginResult.Status, loginResult.ErrorMessage)
	if err != nil {
		return err
	}
	s.UserID, s.TokenID, s.UserExpirationTime, s.UserInfo = loginResult.UserID, loginResult.TokenID, time.Unix(loginResult.UserExpirationTime/1000, 0), loginResult.UserInfo
	s.UpdateLimitParams()
	return nil
}

func (s *Session) UpdateLimitParams() {
	// 参数设定：
	// MinuteDuration: min>minDis*3, max<maxDis*10
	switch s.UserInfo.Sex {
	case "F":
		s.LimitParams = &LimitParams{
			RandDistance:        Float64Range{2.0, 3.0},
			LimitSingleDistance: Float64Range{1.0, 3.0},
			LimitTotalDistance:  Float64Range{1.0, 3.0},
			MinuteDuration:      IntRange{11, 20},
		}
	case "M":
		s.LimitParams = &LimitParams{
			RandDistance:        Float64Range{2.6, 4.0},
			LimitSingleDistance: Float64Range{2.0, 4.0},
			LimitTotalDistance:  Float64Range{2.0, 5.0},
			MinuteDuration:      IntRange{14, 25},
		}

	default:
		panic("Unknown Sex" + s.UserInfo.Sex)
	}
}
func (s *Session) UploadRecord(record Record) (e error) {
	return s.UploadData(record.Distance, record.BeginTime, record.EndTime, record.xtcode, record.SchoolID)
}
func (s *Session) UploadTestRecord(record Record) (e error) {
	return s.uploadTestRecord(record.Distance, record.BeginTime, record.EndTime, record.xtcode, int64(record.EndTime.Sub(record.BeginTime).Seconds()), record.SchoolID)
}

func (s *Session) uploadTestRecord(distance float64, beginTime time.Time, endTime time.Time, xtCode string, useTime int64, schoolID int64) (e error) {
	bz := "[" +
		strconv.FormatInt(endTime.Unix(), 10) + ", " +
		defaultDevice + ", " +
		s.PhoneIMEI + ", " +
		defaultIMSI +
		"]"
	j := XTJsonSportTestData{
		Result:       toExchangeDistanceStr(distance),
		StartTimeStr: toExchangeTimeStr(beginTime),
		EndTimeStr:   toExchangeTimeStr(endTime),
		IsValid:      1,
		BZ:           bz,
		XTCode:       xtCode,
		SchoolID:     schoolID,
		TestTime:     useTime,
	}.ToJSON()
	fmt.Println("j:", j)
	strpa := EncodeSportData(j)
	fmt.Println("strpa:", strpa)
	query := url.Values{"item_param": []string{strpa}}.Encode()
	fmt.Println("query:", query)
	req, err := http.NewRequest(http.MethodPost, postTestDataURL+"?"+query, nil)
	if err != nil {
		panic(httpError{"HTTP Create Request Failed.", err})
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header["UserID"] = []string{strconv.FormatInt(s.UserID, 10)}
	req.Header["TokenID"] = []string{s.TokenID}
	req.Header["app"] = []string{"com.ccxyct.sunshinemotion"}
	req.Header["ver"] = []string{appVersion}
	req.Header["device"] = []string{defaultDevice}
	req.Header["model"] = []string{s.PhoneModel}
	req.Header["screen"] = []string{"1080x1920"}
	req.Header["imei"] = []string{s.PhoneIMEI}
	req.Header["imsi"] = []string{defaultIMSI}
	req.Header["crack"] = []string{"0"}
	req.Header["latitude"] = []string{"0.0"}
	req.Header["longitude"] = []string{"0.0"}
	req.Header.Set("User-Agent", s.UserAgent)

	resp, err := http.DefaultClient.Do(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		panic(fmt.Errorf("HTTP Send Request Failed! %s", err.Error()))
	}
	if resp.StatusCode != http.StatusOK {
		panic(fmt.Errorf("HTTP Get Failed Resp! %s", http.StatusText(resp.StatusCode)))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Errorf("HTTP Read Resp Failed! %s", err.Error()))
	}
	var uploadResult struct {
		Status       int64
		ErrorMessage string
	}
	err = json.Unmarshal(respBytes, &uploadResult)
	if err != nil {
		panic(fmt.Errorf("reslove Failed. %s %s", err.Error(), string(respBytes)))
	}

	return translateServiceError(uploadResult.Status, uploadResult.ErrorMessage)
}

func (s *Session) UploadData(distance float64, beginTime time.Time, endTime time.Time, xtCode string, schoolId int64) (e error) {
	bz := "[" +
		strconv.FormatInt(endTime.Unix(), 10) + ", " +
		defaultDevice + ", " +
		s.PhoneIMEI + ", " +
		defaultIMSI +
		"]"
	j := XTJsonSportData{
		Result:       toExchangeDistanceStr(distance),
		StartTimeStr: toExchangeTimeStr(beginTime),
		EndTimeStr:   toExchangeTimeStr(endTime),
		IsValid:      1,
		BZ:           bz,
		XTCode:       xtCode,
		SchoolID:     schoolId,
	}.ToJSON()
	fmt.Println("j:", j)
	strpa := EncodeSportData(j)
	fmt.Println("strpa:", strpa)
	query := url.Values{"item_param": []string{strpa}}.Encode()
	fmt.Println("query:", query)
	req, err := http.NewRequest(http.MethodPost, uploadDataURL+"?"+query, nil)
	if err != nil {
		panic(httpError{"HTTP Create Request Failed.", err})
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header["UserID"] = []string{strconv.FormatInt(s.UserID, 10)}
	req.Header["TokenID"] = []string{s.TokenID}
	req.Header["app"] = []string{"com.ccxyct.sunshinemotion"}
	req.Header["ver"] = []string{appVersion}
	req.Header["device"] = []string{defaultDevice}
	req.Header["model"] = []string{s.PhoneModel}
	req.Header["screen"] = []string{"1080x1920"}
	req.Header["imei"] = []string{s.PhoneIMEI}
	req.Header["imsi"] = []string{defaultIMSI}
	req.Header["crack"] = []string{"0"}
	req.Header["latitude"] = []string{"0.0"}
	req.Header["longitude"] = []string{"0.0"}
	req.Header.Set("User-Agent", s.UserAgent)

	resp, err := http.DefaultClient.Do(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		panic(fmt.Errorf("HTTP Send Request Failed! %s", err.Error()))
	}
	if resp.StatusCode != http.StatusOK {
		panic(fmt.Errorf("HTTP Get Failed Resp! %s", http.StatusText(resp.StatusCode)))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(fmt.Errorf("HTTP Read Resp Failed! %s", err.Error()))
	}

	var uploadResult struct {
		Status       int64
		ErrorMessage string
	}
	err = json.Unmarshal(respBytes, &uploadResult)
	if err != nil {
		panic(fmt.Errorf("reslove Failed. %s %s", err.Error(), string(respBytes)))
	}

	return translateServiceError(uploadResult.Status, uploadResult.ErrorMessage)
}

type SportResult struct {
	LastTime  time.Time
	Qualified float64
	Distance  float64
}

func (s *Session) GetSportResult() (r *SportResult, e error) {
	req, err := http.NewRequest(http.MethodPost, getSportResultURL, strings.NewReader("flag=0"))
	if err != nil {
		return nil, fmt.Errorf("HTTP Create Request Failed. %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header["UserID"] = []string{strconv.FormatInt(s.UserID, 10)}
	req.Header["TokenID"] = []string{s.TokenID}
	req.Header["app"] = []string{"com.ccxyct.sunshinemotion"}
	req.Header["ver"] = []string{appVersion}
	req.Header["device"] = []string{defaultDevice}
	req.Header["model"] = []string{s.PhoneModel}
	req.Header["screen"] = []string{"1080x1920"}
	req.Header["imei"] = []string{s.PhoneIMEI}
	req.Header["imsi"] = []string{defaultIMSI}
	req.Header["crack"] = []string{"0"}
	req.Header["latitude"] = []string{"0.0"}
	req.Header["longitude"] = []string{"0.0"}
	resp, err := http.DefaultClient.Do(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, fmt.Errorf("HTTP Send Request Failed! %s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP Get Failed Resp! %s", http.StatusText(resp.StatusCode))
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("HTTP Read Resp Failed! %s", err.Error())
	}
	var httpSporstResult struct {
		Status       int64
		ErrorMessage string
		LastTime     string  `json:"lastTime"`
		Qualified    float64 `json:"qualified"`
		Result       float64 `json:"Result"`
		UserID       int64   `json:"userID"`
		Term         string  `json:"term"`
		Year         int     `json:"year"`
	}
	err = json.Unmarshal(respBytes, &httpSporstResult)
	if err != nil {
		return nil, fmt.Errorf("reslove Failed. %s %s", err.Error(), string(respBytes))
	}

	err = translateServiceError(httpSporstResult.Status, httpSporstResult.ErrorMessage)
	if err != nil {
		return nil, err
	}

	r = new(SportResult)
	if httpSporstResult.LastTime != "" {
		r.LastTime, err = fromExchangeTimeStr(httpSporstResult.LastTime)
	} else {
		r.LastTime = time.Now()
	}
	if err != nil {
		log.Println(string(respBytes))
		panic(err)
	}
	r.Qualified = httpSporstResult.Qualified
	r.Distance = httpSporstResult.Result
	return r, nil
}

type AppInfo struct {
	ID        int    `json:"iID"`
	VerNumber int    `json:"iVerNumber"`
	Url       string `json:"strUrl"`
	Ver       string `json:"strVer"`
}

// Fetch the latest app info
func (s *Session) GetAppInfo() (r AppInfo, e error) {
	req, err := http.NewRequest(http.MethodPost, getAppInfoURL, nil)
	if err != nil {
		e = fmt.Errorf("HTTP Create Request Failed. %s", err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header["UserID"] = []string{strconv.FormatInt(s.UserID, 10)}
	req.Header["TokenID"] = []string{s.TokenID}
	req.Header["app"] = []string{"com.ccxyct.sunshinemotion"}
	req.Header["ver"] = []string{appVersion}
	req.Header["device"] = []string{defaultDevice}
	req.Header["model"] = []string{s.PhoneModel}
	req.Header["screen"] = []string{"1080x1920"}
	req.Header["imei"] = []string{s.PhoneIMEI}
	req.Header["imsi"] = []string{defaultIMSI}
	req.Header["crack"] = []string{"0"}
	req.Header["latitude"] = []string{"0.0"}
	req.Header["longitude"] = []string{"0.0"}
	resp, err := http.DefaultClient.Do(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		e = fmt.Errorf("HTTP Send Request Failed! %s", err.Error())
		return
	}
	if resp.StatusCode != http.StatusOK {
		e = fmt.Errorf("HTTP Get Failed Resp! %s", http.StatusText(resp.StatusCode))
		return
	}
	respBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		e = fmt.Errorf("HTTP Read Resp Failed! %s", err.Error())
		return
	}
	var httpResult struct {
		Status       int64   `json:"Status"`
		ErrorMessage string  `json:"ErrorMessage"`
		AppInfo      AppInfo `json:"AppInfo"`
	}
	err = json.Unmarshal(respBytes, &httpResult)
	if err != nil {
		e = fmt.Errorf("reslove Failed. %s %s", err.Error(), string(respBytes))
		return
	}

	err = translateServiceError(httpResult.Status, httpResult.ErrorMessage)
	if err != nil {
		e = err
		return
	}

	return httpResult.AppInfo, nil
}
