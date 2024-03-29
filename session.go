package ssmt

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Session struct {
	User   *UserIdentify
	Token  *UserToken
	Device *Device
}

type UserIdentify struct {
	SchoolID int64
	StuNum   string
	UserID   int64
}

type UserToken struct {
	TokenID        string
	ExpirationTime time.Time
}

var (
	ErrNotLogin = errors.New("this operation need to login first")
)

const (
	serverAPIRoot     = "http://www.ccxyct.com:8080"
	loginURL          = serverAPIRoot + "/sunShine_Sports/loginSport.action"
	uploadDataURL     = serverAPIRoot + "/sunShine_Sports/xtUploadData.action"
	postTestDataURL   = serverAPIRoot + "/sunShine_Sports/postTestData.action"
	getSportResultURL = serverAPIRoot + "/sunShine_Sports/xtGetSportResult.action"
	getAppInfoURL     = serverAPIRoot + "/sunShine_Sports/xtGetAppInfo.action"
	getSchoolURL      = serverAPIRoot + "/sunShine_Sports/xtGetSchool.action"
	getRandRouteURL   = serverAPIRoot + "/sunShine_Sports/xtGetRandRoute.action"
	getTestRuleURL    = serverAPIRoot + "/sunShine_Sports/getTestRule.action"

	AppPackageName = "com.ccxyct.sunshinemotion"
	AppVersion     = "3.1.12"
	AppVersionID   = 18
)

func (s *Session) setRandomDevice() {
	s.Device = GenerateDevice()
}

func CreateSession() *Session {
	return &Session{}
}

func (s *Session) setHTTPHeader(req *http.Request) {
	if !s.login() {
		s.setHTTPHeaderWithoutLogin(req)
		return
	}
	if s.Device == nil {
		s.setRandomDevice()
	}
	req.Header["Content-Type"] = []string{"application/x-www-form-urlencoded"}
	req.Header["UserID"] = []string{strconv.FormatInt(s.User.UserID, 10)}
	req.Header["TokenID"] = []string{s.Token.TokenID}
	req.Header["app"] = []string{AppPackageName}
	req.Header["Ver"] = []string{AppVersion}
	req.Header["Device"] = []string{s.Device.DeviceName}
	req.Header["Model"] = []string{s.Device.ModelType}
	req.Header["Screen"] = []string{s.Device.Screen}
	req.Header["IMEI"] = []string{s.Device.IMEI}
	req.Header["IMSI"] = []string{s.Device.IMSI}
	req.Header["crack"] = []string{"0"}
	req.Header["latitude"] = []string{"0.0"}
	req.Header["longitude"] = []string{"0.0"}
	req.Header["VerID"] = []string{strconv.FormatInt(AppVersionID, 10)}
	req.Header.Set("User-Agent", s.Device.UserAgent)
}
func (s *Session) setHTTPHeaderWithoutLogin(req *http.Request) {
	if s.Device == nil {
		s.setRandomDevice()
	}
	req.Header["Content-Type"] = []string{"application/x-www-form-urlencoded"}
	req.Header["UserID"] = []string{"0"}
	req.Header["TokenID"] = []string{""}
	req.Header["app"] = []string{AppPackageName}
	req.Header["Ver"] = []string{AppVersion}
	req.Header["Device"] = []string{s.Device.DeviceName}
	req.Header["Model"] = []string{s.Device.ModelType}
	req.Header["Screen"] = []string{s.Device.Screen}
	req.Header["IMEI"] = []string{s.Device.IMEI}
	req.Header["IMSI"] = []string{s.Device.IMSI}
	req.Header["crack"] = []string{"0"}
	req.Header["latitude"] = []string{"0.0"}
	req.Header["longitude"] = []string{"0.0"}
	req.Header["VerID"] = []string{strconv.FormatInt(AppVersionID, 10)}
	req.Header.Set("User-Agent", s.Device.UserAgent)
}

type UserInfo struct {
	ClassID       int64  `json:"inClassID"`
	ClassName     string `json:"inClassName"`
	CollegeID     int64  `json:"inCollegeID"`
	CollegeName   string `json:"inCollegeName"`
	SchoolID      int64  `json:"inSchoolID"`
	SchoolName    string `json:"inSchoolName"`
	SchoolNumber  string `json:"inSchoolNumber"`
	NickName      string `json:"nickName"`
	StudentName   string `json:"studentName"`
	StudentNumber string `json:"studentNumber"`
	IsTeacher     int    `json:"isTeacher"`
	Sex           string `json:"sex"`
	PhoneNumber   string `json:"phoneNumber"`
	UserRoleID    int    `json:"UserRoleID"`
}

func (s *Session) login() bool {
	return s.Token != nil
}
func (s *Session) Login(schoolID int64, stuNum string, phoneNum string, passwordHash string) (info UserInfo, e error) {
	if s.Device == nil {
		s.setRandomDevice()
	}
	req, err := http.NewRequest(http.MethodPost, loginURL, strings.NewReader(url.Values{
		"stuNum":   {stuNum},
		"phoneNum": {phoneNum},
		"passWd":   {passwordHash},
		"schoolId": {strconv.FormatInt(schoolID, 10)},
		"stuId":    {"1"},
		"token":    {""},
	}.Encode()))
	if err != nil {
		return UserInfo{}, httpError{"HTTP Create Request Failed.", err}
	}
	s.User, s.Token = new(UserIdentify), new(UserToken)
	s.setHTTPHeaderWithoutLogin(req)
	resp, err := http.DefaultClient.Do(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return UserInfo{}, httpError{"HTTP Send Request Failed! ", err}
	}
	if resp.StatusCode != http.StatusOK {
		return UserInfo{}, fmt.Errorf("HTTP Response Status: %d(%s)", resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	var loginResult struct {
		Status             int64
		ErrorMessage       string
		Date               string
		UserID             int64
		TokenID            string
		UserExpirationTime int64
		UserInfo           UserInfo
	}
	err = json.NewDecoder(resp.Body).Decode(&loginResult)
	if err != nil {
		return UserInfo{}, fmt.Errorf("HTTP Response reslove Failed. %s", err.Error())
	}
	err = translateServiceError(loginResult.Status, loginResult.ErrorMessage)
	if err != nil {
		return UserInfo{}, err
	}
	s.User.UserID, s.User.SchoolID, s.User.StuNum = loginResult.UserID, schoolID, stuNum
	s.Token.TokenID, s.Token.ExpirationTime = loginResult.TokenID, time.Unix(0, loginResult.UserExpirationTime*1000000)
	return loginResult.UserInfo, nil
}

func (s *Session) UploadRecord(record Record) (e error) {
	if !s.login() {
		return ErrNotLogin
	}
	if s.Device == nil {
		s.setRandomDevice()
	}
	bz := EncodeString("[ccxyct:" +
		strconv.FormatInt(record.EndTime.UnixNano()/1000000, 10) + ", " +
		s.Device.DeviceName + ", " +
		s.Device.IMEI + ", " +
		s.Device.IMSI +
		"]")
	sportData := XTJsonSportDataFromRecord(record, bz)
	fmt.Println("sportData:", sportData)
	query := makeQuery(sportData)
	req, err := http.NewRequest(http.MethodPost, uploadDataURL+"?"+query, nil)
	if err != nil {
		return httpError{"HTTP Create Request Failed.", err}
	}
	s.setHTTPHeader(req)
	fmt.Println("Upload Request Headers: ", DumpStruct(req.Header))
	fmt.Println("Upload Request Query: ", query)
	beginT := time.Now()
	resp, err := http.DefaultClient.Do(req)
	endT := time.Now()
	fmt.Println("Upload Time", beginT, endT, endT.Sub(beginT))
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return httpError{"HTTP Send Request Failed! %s", err}
	}
	fmt.Println("Upload Response Headers: ", DumpStruct(resp.Header))
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP Status: %s", resp.Status)
	}
	var uploadResult struct {
		Status       int64
		ErrorMessage string
	}
	err = json.NewDecoder(resp.Body).Decode(&uploadResult)
	fmt.Println("Upload Response Result: ", DumpStruct(uploadResult))
	if err != nil {
		return fmt.Errorf("HTTP Response reslove Failed. %s", err.Error())
	}
	return translateServiceError(uploadResult.Status, uploadResult.ErrorMessage)
}
func (s *Session) UploadTestRecord(record Record) (e error) {
	if !s.login() {
		return ErrNotLogin
	}
	if s.Device == nil {
		s.setRandomDevice()
	}
	bz := EncodeString("[ccxyct:" +
		strconv.FormatInt(record.EndTime.UnixNano()/1000000, 10) +
		"]")
	sportData := XTJsonSportTestDataFromRecord(record, bz)
	fmt.Println("sportData:", sportData)
	query := makeQuery(sportData)
	fmt.Println("query:", query)
	req, err := http.NewRequest(http.MethodPost, postTestDataURL+"?"+query, nil)
	if err != nil {
		return httpError{"HTTP Create Request Failed.", err}
	}
	s.setHTTPHeader(req)
	resp, err := http.DefaultClient.Do(req)
	if resp != nil && resp.Body != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return fmt.Errorf("HTTP Send Request Failed! %s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP Get Failed Resp! %s", http.StatusText(resp.StatusCode))
	}
	var uploadResult struct {
		Status       int64
		ErrorMessage string
	}
	err = json.NewDecoder(resp.Body).Decode(&uploadResult)
	if err != nil {
		return fmt.Errorf("reslove Failed. %s", err.Error())
	}

	return translateServiceError(uploadResult.Status, uploadResult.ErrorMessage)
}
func makeQuery(d IXTJsonSportData) string {
	j := d.ToJSON()
	fmt.Println("json:", j)
	pa := EncodeString(j)
	fmt.Println("pa:", pa)
	return url.Values{"item_param": []string{pa}}.Encode()
}

type SportResult struct {
	UserID            int64     // 用户ID，作废字段，目前暂时恒为0
	Year              int       // 年度
	Term              string    // 学期
	QualifiedDistance float64   // 达标距离
	ActualDistance    float64   // 已计距离
	LastTime          time.Time // 上次跑步时间
}

func (s *Session) GetSportResult() (r *SportResult, e error) {
	if !s.login() {
		return nil, ErrNotLogin
	}
	if s.Device == nil {
		s.setRandomDevice()
	}
	req, err := http.NewRequest(http.MethodPost, getSportResultURL, strings.NewReader("flag=0"))
	if err != nil {
		return nil, fmt.Errorf("HTTP Create Request Failed. %s", err.Error())
	}
	s.setHTTPHeader(req)
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
	var httpSportResult struct {
		Status       int64
		ErrorMessage string
		LastTime     string  `json:"lastTime"`
		Qualified    float64 `json:"qualified"`
		Result       float64 `json:"Result"`
		UserID       int64   `json:"userID"`
		Term         string  `json:"term"`
		Year         int     `json:"year"`
	}
	err = json.NewDecoder(resp.Body).Decode(&httpSportResult)
	if err != nil {
		return nil, fmt.Errorf("reslove Failed. %s", err.Error())
	}
	err = translateServiceError(httpSportResult.Status, httpSportResult.ErrorMessage)
	if err != nil {
		return nil, err
	}

	r = &SportResult{
		UserID:            httpSportResult.UserID,
		Year:              httpSportResult.Year,
		Term:              httpSportResult.Term,
		QualifiedDistance: httpSportResult.Qualified,
		ActualDistance:    httpSportResult.Result,
	}
	if httpSportResult.LastTime != "" {
		r.LastTime, err = fromServiceStdTime(httpSportResult.LastTime)
		if err != nil {
			log.Println(err)
			return nil, errors.New("response parsing error: " + err.Error())
		}
	}
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
	if s.Device == nil {
		s.setRandomDevice()
	}
	req, err := http.NewRequest(http.MethodPost, getAppInfoURL, nil)
	if err != nil {
		e = fmt.Errorf("HTTP Create Request Failed. %s", err.Error())
		return
	}
	s.setHTTPHeader(req)
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
	var httpResult struct {
		Status       int64   `json:"Status"`
		ErrorMessage string  `json:"ErrorMessage"`
		AppInfo      AppInfo `json:"AppInfo"`
	}
	err = json.NewDecoder(resp.Body).Decode(&httpResult)
	if err != nil {
		e = fmt.Errorf("reslove Failed. %s", err.Error())
		return
	}

	err = translateServiceError(httpResult.Status, httpResult.ErrorMessage)
	if err != nil {
		e = err
		return
	}

	return httpResult.AppInfo, nil
}

type RandRoute struct {
	Dots []RandRouteDot `json:"dotArray"`
	ID   int            `json:"routeID"`
	Name string         `json:"routeName"`
}
type RandRouteDot struct {
	ID        int    `json:"dotID"`
	Name      string `json:"dotName"`
	IsKey     int    `json:"isKey"`
	Latitude  string `json:"latitude"`
	Longitude string `json:"longitude"`
}

func (s *Session) GetRandRoute() (r []RandRoute, e error) {
	if !s.login() {
		return []RandRoute{}, ErrNotLogin
	}
	if s.Device == nil {
		s.setRandomDevice()
	}
	req, err := http.NewRequest(http.MethodPost, getRandRouteURL, strings.NewReader(url.Values{
		"schoolId": {strconv.FormatInt(s.User.SchoolID, 10)},
	}.Encode()))
	if err != nil {
		e = fmt.Errorf("HTTP Create Request Failed. %s", err.Error())
		return
	}
	s.setHTTPHeader(req)
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
	var httpResult struct {
		Status       int64  `json:"Status"`
		ErrorMessage string `json:"ErrorMessage"`

		// Warning: misleading fileds names, only use if you really understand it.
		// these fileds is confusing in target system, the following notes give some explanations.
		Qualified float64 `json:"Qualified"` // Confirmed: Qualified Distance
		MinSpeed  float64 `json:"MinSpeed"`  // Unconfirmed
		DayLimit  float64 `json:"DayLimit"`  // Inconsistent Name&Usage: Minimal Distance One Record
		MaxSpeed  float64 `json:"MaxSpeed"`  // Inconsistent Name&Usage: Maximal Distance One Day
		WeekLimit float64 `json:"WeekLimit"` // Unconfirmed

		RouteArray []RandRoute `json:"RouteArray"`
	}
	err = json.NewDecoder(resp.Body).Decode(&httpResult)
	if err != nil {
		e = fmt.Errorf("reslove Failed. %s", err.Error())
		return
	}
	println(DumpStructValue(httpResult))
	err = translateServiceError(httpResult.Status, httpResult.ErrorMessage)
	if err != nil {
		e = err
		return
	}

	return httpResult.RouteArray, nil
}

type SchoolList []School
type School struct {
	SchoolId     int64  `json:"schoolId"`
	SchoolName   string `json:"schoolName"`
	SchoolNumber string `json:"schoolNumber"`
}

func (s *Session) GetSchoolList() (r SchoolList, e error) {
	if s.Device == nil {
		s.setRandomDevice()
	}
	req, err := http.NewRequest(http.MethodPost, getSchoolURL, nil)
	if err != nil {
		e = fmt.Errorf("HTTP Create Request Failed. %s", err.Error())
		return
	}
	s.setHTTPHeaderWithoutLogin(req)
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
	var httpResult struct {
		Status       int64      `json:"Status"`
		ErrorMessage string     `json:"ErrorMessage"`
		SchoolList   SchoolList `json:"ShoolArray"`
	}
	err = json.NewDecoder(resp.Body).Decode(&httpResult)
	if err != nil {
		e = fmt.Errorf("reslove Failed. %s", err.Error())
		return
	}

	err = translateServiceError(httpResult.Status, httpResult.ErrorMessage)
	if err != nil {
		e = err
		return
	}

	return httpResult.SchoolList, nil
}

func (s *Session) GetTestRule() (r TestRule, e error) {
	if s.Device == nil {
		s.setRandomDevice()
	}
	req, err := http.NewRequest(http.MethodPost, getTestRuleURL+"?"+url.Values{
		"school_id": {strconv.FormatInt(s.User.SchoolID, 10)},
	}.Encode(), nil)
	if err != nil {
		e = fmt.Errorf("HTTP Create Request Failed. %s", err.Error())
		return
	}
	s.setHTTPHeader(req)
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
	var httpResult struct {
		Status       int64    `json:"Status"`
		ErrorMessage string   `json:"ErrorMessage"`
		Rule         TestRule `json:"TestRule"`
	}
	err = json.NewDecoder(resp.Body).Decode(&httpResult)
	if err != nil {
		e = fmt.Errorf("reslove Failed. %s", err.Error())
		return
	}

	err = translateServiceError(httpResult.Status, httpResult.ErrorMessage)
	if err != nil {
		e = err
		return
	}

	return httpResult.Rule, nil
}
