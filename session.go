package lib

import (
	"encoding/json"
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
	loginURL          = server + "/sunShine_Sports-server1/loginSport.action"
	uploadDataURL     = server + "/sunShine_Sports-server1/xtUploadData.action"
	postTestDataURL   = server + "/sunShine_Sports-server1/postTestData.action"
	getSportResultURL = server + "/sunShine_Sports-server1/xtGetSportResult.action"
	DefaultUserAgent  = "Dalvik/2.1.0 (Linux; U; Android 7.0)"

	defaultSchoolId = 60
	defaultIMSI     = "1234567890"
	defaultDevice   = "Android,25,7.1.2"
)

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
	req.Header.Set("User-Agent", s.UserAgent)
	req.Header.Set("UserID", "0")
	req.Header.Set("crack", "0")
	req.Header["UserID"] = []string{"0"}
	req.Header["crack"] = []string{"0"}

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
		UserID             int64
		TokenID            string
		UserExpirationTime int64
		UserInfo           UserInfo
	}
	err = json.Unmarshal(respBytes, &loginResult)
	if err != nil {
		return fmt.Errorf("reslove Failed. %s %s", err.Error(), string(respBytes))
	}
	if loginResult.Status != 1 {
		return fmt.Errorf("resp status not ok. %d", loginResult.Status)
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
	return s.UploadData(record.Distance, record.BeginTime, record.EndTime, record.xtcode)
}
func (s *Session) UploadTestRecord(record Record) (e error) {
	return s.uploadTestRecord(record.Distance, record.BeginTime, record.EndTime, record.xtcode, int64(record.EndTime.Sub(record.BeginTime).Seconds()))
}

func (s *Session) uploadTestRecord(distance float64, beginTime time.Time, endTime time.Time, xtCode string, useTime int64) (e error) {
	bz := "[" +
		strconv.FormatInt(endTime.Unix(), 10) + ", " +
		defaultDevice + ", " +
		s.PhoneIMEI + ", " +
		defaultIMSI +
		"]"
	li := GetLi("", bz)
	req, err := http.NewRequest(http.MethodPost, postTestDataURL, strings.NewReader(url.Values{
		"results":   {toExchangeDistanceStr(distance)},
		"beginTime": {toExchangeTimeStr(beginTime)},
		"endTime":   {toExchangeTimeStr(endTime)},
		"isValid":   {"1"},
		"schoolId":  {strconv.FormatInt(s.UserInfo.InSchoolID, 10)},
		"xtCode":    {xtCode},
		"bz":        {EncodeBZ(bz)},
		"test_time": {toExchangeInt64Str(useTime)},
		"li":        {li},
	}.Encode()))
	if err != nil {
		panic(httpError{"HTTP Create Request Failed.", err})
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", s.UserAgent)
	req.Header["UserID"] = []string{strconv.FormatInt(s.UserID, 10)}
	req.Header["TokenID"] = []string{s.TokenID}
	req.Header["app"] = []string{"com.ccxyct.sunshinemotion"}
	req.Header["ver"] = []string{"2.2.2"}
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
		Status       int
		ErrorMessage string
	}
	err = json.Unmarshal(respBytes, &uploadResult)
	if err != nil {
		panic(fmt.Errorf("reslove Failed. %s %s", err.Error(), string(respBytes)))
	}
	const successCode = 1
	if uploadResult.Status != successCode {
		return fmt.Errorf("server status %d , message: %s", uploadResult.Status, uploadResult.ErrorMessage)
	}
	return nil
}

func (s *Session) UploadData(distance float64, beginTime time.Time, endTime time.Time, xtCode string) (e error) {
	bz := "[" +
		strconv.FormatInt(endTime.Unix(), 10) + ", " +
		defaultDevice + ", " +
		s.PhoneIMEI + ", " +
		defaultIMSI +
		"]"
	li := GetLi(xtCode, bz)
	req, err := http.NewRequest(http.MethodPost, uploadDataURL, strings.NewReader(url.Values{
		"results":   {toExchangeDistanceStr(distance)},
		"beginTime": {toExchangeTimeStr(beginTime)},
		"endTime":   {toExchangeTimeStr(endTime)},
		"isValid":   {"1"},
		"schoolId":  {strconv.FormatInt(s.UserInfo.InSchoolID, 10)},
		"xtCode":    {xtCode},
		"bz":        {EncodeBZ(bz)},
		"li":        {li},
	}.Encode()))
	if err != nil {
		panic(httpError{"HTTP Create Request Failed.", err})
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", s.UserAgent)
	req.Header["UserID"] = []string{strconv.FormatInt(s.UserID, 10)}
	req.Header["TokenID"] = []string{s.TokenID}
	req.Header["app"] = []string{"com.ccxyct.sunshinemotion"}
	req.Header["ver"] = []string{"2.2.2"}
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
		Status       int
		ErrorMessage string
	}
	err = json.Unmarshal(respBytes, &uploadResult)
	if err != nil {
		panic(fmt.Errorf("reslove Failed. %s %s", err.Error(), string(respBytes)))
	}
	const successCode = 1
	if uploadResult.Status != successCode {
		return fmt.Errorf("server status %d , message: %s", uploadResult.Status, uploadResult.ErrorMessage)
	}
	return nil
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
	req.Header.Set("User-Agent", s.UserAgent)
	req.Header["UserID"] = []string{strconv.FormatInt(s.UserID, 10)}
	req.Header["TokenID"] = []string{s.TokenID}
	req.Header["crack"] = []string{"0"}

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
		Status       int
		ErrorMessage string
		LastTime     string  `json:"lastTime"`
		Qualified    float64 `json:"qualified"`
		Result       float64 `json:"result"`
		UserID       int64   `json:"userID"`
		Term         string  `json:"term"`
		Year         int     `json:"year"`
	}
	err = json.Unmarshal(respBytes, &httpSporstResult)
	if err != nil {
		return nil, fmt.Errorf("reslove Failed. %s %s", err.Error(), string(respBytes))
	}
	const successCode = 1
	if httpSporstResult.Status != successCode {
		return nil, fmt.Errorf("server status %d , message: %s", httpSporstResult.Status, httpSporstResult.ErrorMessage)
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
