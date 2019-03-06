package sunshinemotion

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Session struct {
	device           *Device
	token            *Token
	httpClient       *http.Client
	userInfo         UserInfo
	sportResultCache sportResultCache
}

const (
	AppPackageName = "com.ccxyct.sunshinemotion" // the app package name to emulate
	AppVersion     = "2.2.6"                     // the app version to emulate
)

func CreateSession(device *Device, token *Token) *Session {
	s := &Session{
		device: device,
		token:  token,
	}
	httpClient := &http.Client{
		Transport: &sessionTransport{
			device: device,
			s:      s,
		},
	}
	sportResultCache := sportResultCache{
		ExpireDuration: 15 * time.Second,
	}
	sportResultCache.Update = func() (err error) {
		_, err = s.GetUserSportResult()
		return errors.New("failed to get results for cache update: " + err.Error())
	}
	s.httpClient = httpClient
	s.sportResultCache = sportResultCache
	return s
}

// request a token and use it in the session automatically.
func (session *Session) UpdateToken(schoolID uint64, username string, passwordHash string, phone string) error {
	token, err := session.RequestToken(schoolID, username, passwordHash, phone)
	if err == nil {
		session.token = token
	}
	return err
}

// only return a token but does not use it automatically.
// if you expect it to use token automatically, see session.UpdateToken().
//
// can use when token is nil.
func (session *Session) RequestToken(schoolID uint64, username string, passwordHash string, phone string) (*Token, error) {
	req, err := http.NewRequest(http.MethodPost, loginURL, strings.NewReader(url.Values{
		"stuNum":   {username},
		"phoneNum": {phone},
		"passWd":   {passwordHash},
		"schoolId": {strconv.FormatUint(schoolID, 10)},
		"stuId":    {"1"},
		"token":    {""},
	}.Encode()))
	if err != nil {
		return nil, errors.New("HTTP Create Request Failed! " + err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := session.httpClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, errors.New("HTTP Send Request Failed! " + err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP Response Status: %d(%s)", resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	var loginResult struct {
		Status             int64
		ErrorMessage       string
		UserID             uint
		TokenID            string
		UserExpirationTime int64
		UserInfo           UserInfo
	}
	err = json.NewDecoder(resp.Body).Decode(&loginResult)
	if err != nil {
		return nil, fmt.Errorf("reslove Failed. %s", err.Error())
	}
	err = serviceCodeToGoError(loginResult.Status, loginResult.ErrorMessage)
	if err != nil {
		return nil, err
	}
	session.userInfo = loginResult.UserInfo
	return &Token{
		TokenID:    loginResult.TokenID,
		UserID:     loginResult.UserID,
		SchoolID:   schoolID,
		ExpireTime: time.Unix(loginResult.UserExpirationTime, 0),
	}, nil
}

func (session *Session) UploadSportRecord(record Record) (e error) {
	if !session.token.ValidFormat() {
		return ErrTokenInvalid
	}
	if session.token.Expired() {
		return ErrTokenExpired
	}
	queryParams := url.Values{
		"item_param": []string{record.EncryptedJSON()},
	}
	req, err := http.NewRequest(http.MethodPost, uploadSportRecordURL+"?"+queryParams.Encode(), nil)
	if err != nil {
		return errors.New("HTTP Create Request Failed." + err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header["TokenID"] = []string{session.token.TokenID}
	req.Header["app"] = []string{AppPackageName}
	req.Header["ver"] = []string{AppVersion}
	req.Header["device"] = []string{session.device.Name}
	req.Header["model"] = []string{session.device.Model}
	req.Header["screen"] = []string{session.device.Screen}
	req.Header["imei"] = []string{session.device.IMEI}
	req.Header["imsi"] = []string{session.device.IMSI}
	req.Header["crack"] = []string{"0"}
	req.Header["latitude"] = []string{session.device.Latitude}
	req.Header["longitude"] = []string{session.device.Longitude}
	resp, err := session.httpClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return fmt.Errorf("HTTP Send Request Failed! %s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP Response Status: %d(%s)", resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	var uploadResult struct {
		Status       int64
		ErrorMessage string
	}
	err = json.NewDecoder(resp.Body).Decode(&uploadResult)
	if err != nil {
		return fmt.Errorf("reslove Failed. %s", err.Error())
	}
	err = serviceCodeToGoError(uploadResult.Status, uploadResult.ErrorMessage)
	if err != nil {
		return err
	}
	return nil
}

func (session *Session) UploadTestRecord(record Record) (e error) {
	if !session.token.ValidFormat() {
		return ErrTokenInvalid
	}
	if session.token.Expired() {
		return ErrTokenExpired
	}
	queryParams := url.Values{
		"item_param": []string{record.EncryptedJSON()},
	}
	req, err := http.NewRequest(http.MethodPost, uploadTestRecordURL+"?"+queryParams.Encode(), nil)
	if err != nil {
		return errors.New("HTTP Create Request Failed." + err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header["app"] = []string{AppPackageName}
	req.Header["ver"] = []string{AppVersion}
	req.Header["device"] = []string{session.device.Name}
	req.Header["model"] = []string{session.device.Model}
	req.Header["screen"] = []string{session.device.Screen}
	req.Header["imei"] = []string{session.device.IMEI}
	req.Header["imsi"] = []string{session.device.IMSI}
	req.Header["latitude"] = []string{session.device.Latitude}
	req.Header["longitude"] = []string{session.device.Longitude}
	resp, err := session.httpClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return fmt.Errorf("HTTP Send Request Failed! %s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HTTP Response Status: %d(%s)", resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	var uploadResult struct {
		Status       int64
		ErrorMessage string
	}
	err = json.NewDecoder(resp.Body).Decode(&uploadResult)
	if err != nil {
		return fmt.Errorf("reslove Failed. %s", err.Error())
	}
	err = serviceCodeToGoError(uploadResult.Status, uploadResult.ErrorMessage)
	if err != nil {
		return err
	}
	return nil
}

func (session *Session) GetUserSportResult() (UserSportResult, error) {
	if !session.token.ValidFormat() {
		return UserSportResult{}, ErrTokenInvalid
	}
	if session.token.Expired() {
		return UserSportResult{}, ErrTokenExpired
	}
	req, err := http.NewRequest(http.MethodPost, getSportResultURL, strings.NewReader("flag=0"))
	if err != nil {
		return UserSportResult{}, fmt.Errorf("HTTP Create Request Failed. %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := session.httpClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return UserSportResult{}, fmt.Errorf("HTTP Send Request Failed! %s", err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return UserSportResult{}, fmt.Errorf("HTTP Response Status: %d(%s)", resp.StatusCode, http.StatusText(resp.StatusCode))
	}
	var responseResult struct {
		Status       int64   `json:"Status"`
		ErrorMessage string  `json:"ErrorMessage"`
		LastTime     string  `json:"lastTime"`
		Qualified    float64 `json:"qualified"`
		Result       float64 `json:"result"`
		UserID       uint    `json:"userID"`
		Term         string  `json:"term"`
		Year         int     `json:"year"`
	}
	err = json.NewDecoder(resp.Body).Decode(&responseResult)
	if err != nil {
		return UserSportResult{}, fmt.Errorf("reslove Failed. %s", err.Error())
	}
	err = serviceCodeToGoError(responseResult.Status, responseResult.ErrorMessage)
	if err != nil {
		return UserSportResult{}, err
	}
	lastTime, err := fromRPCTimeStr(responseResult.LastTime)
	if err != nil {
		return UserSportResult{}, fmt.Errorf("reslove Failed. %s", err.Error())
	}
	r := UserSportResult{
		UserID:            responseResult.UserID,
		Year:              responseResult.Year,
		Term:              responseResult.Term,
		QualifiedDistance: responseResult.Qualified,
		RunDistance:       responseResult.Result,
		LastTime:          lastTime,
	}
	// put UserSportResult into cache
	{
		var cacheTime time.Time
		respTimeStr := resp.Header.Get("Date")
		if respTimeStr != "" {
			// try to use response date as cacheTime
			cacheTime, _, err = parseHTTPDate(respTimeStr)
			if err != nil {
				// use now as cacheTime
				cacheTime = time.Now()
			}
		} else {
			// use now as cacheTime
			cacheTime = time.Now()
		}
		session.sportResultCache.Put(r, cacheTime)
	}
	return r, nil
}

// userInfo is only updated when UpdateToken/RequestToken.
// since userInfo may be obsoleted..
func (session *Session) GetUserInfo() UserInfo {
	return session.userInfo
}
