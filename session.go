package sunshinemotion

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/inkedawn/go-sunshinemotion/crypto"
)

type Session struct {
	device               *Device
	token                *Token
	userSportResultCache userSportResultCache
	httpClient           *http.Client
}

const (
	AppPackageName = "com.ccxyct.sunshinemotion" // the app package name to emulate
	AppVersion     = "2.2.2"                     // the app version to emulate
)

const ServiceSuccessStatus = 1

func CreateSession(device *Device, token *Token) *Session {
	return &Session{
		device: device,
		token:  token,
		userSportResultCache: userSportResultCache{
			ExpireDuration: 1 * time.Hour,
			FetchFunction: func() (updated UserSportResult, err error) {
				return UserSportResult{}, nil
			},
		},
		httpClient: &http.Client{
			Transport: &sessionTransport{
				device: device,
			},
		},
	}
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
	req.Header["UserID"] = []string{"0"}
	req.Header["crack"] = []string{"0"}
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
		// UserInfo           UserInfo
	}
	err = json.NewDecoder(resp.Body).Decode(&loginResult)
	if err != nil {
		return nil, fmt.Errorf("reslove Failed. %s", err.Error())
	}

	if loginResult.Status != ServiceSuccessStatus {
		return nil, serviceError{
			status:  loginResult.Status,
			message: loginResult.ErrorMessage,
		}
	}

	return &Token{
		TokenID:    loginResult.TokenID,
		UserID:     loginResult.UserID,
		SchoolID:   schoolID,
		ExpireTime: time.Unix(loginResult.UserExpirationTime, 0),
	}, nil
}

func (session *Session) UploadSportRecord(record Record) (e error) {
	if !session.token.Valid() {
		return ErrTokenExpired
	}
	xtcode := record.XTcode()
	bz := record.Remark.String()
	li := crypto.CalcLi(xtcode, bz)
	req, err := http.NewRequest(http.MethodPost, uploadSportRecordURL, strings.NewReader(url.Values{
		"results":   {toRPCDistanceStr(record.Distance)},
		"beginTime": {toRPCTimeStr(record.BeginTime)},
		"endTime":   {toRPCTimeStr(record.EndTime)},
		"isValid":   {"1"},
		"schoolId":  {strconv.FormatUint(session.token.SchoolID, 10)},
		"xtCode":    {xtcode},
		"bz":        {crypto.EncryptBZ(bz)},
		"li":        {li},
	}.Encode()))
	if err != nil {
		return errors.New("HTTP Create Request Failed." + err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header["UserID"] = []string{strconv.FormatUint(uint64(session.token.UserID), 10)}
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

	if uploadResult.Status != ServiceSuccessStatus {
		return serviceError{
			status:  uploadResult.Status,
			message: uploadResult.ErrorMessage,
		}
	}
	return nil
}

func (session *Session) UploadTestRecord(record Record) (e error) {
	if !session.token.Valid() {
		return ErrTokenExpired
	}
	xtcode := record.XTcode()
	bz := record.Remark.String()
	li := crypto.CalcLi("", bz)
	useTime := int(math.Floor(record.EndTime.Sub(record.BeginTime).Seconds()))
	req, err := http.NewRequest(http.MethodPost, uploadTestRecordURL, strings.NewReader(url.Values{
		"results":   {toRPCDistanceStr(record.Distance)},
		"beginTime": {toRPCTimeStr(record.BeginTime)},
		"endTime":   {toRPCTimeStr(record.EndTime)},
		"isValid":   {"1"},
		"schoolId":  {strconv.FormatUint(session.token.SchoolID, 10)},
		"xtCode":    {xtcode},
		"bz":        {crypto.EncryptBZ(bz)},
		"test_time": {strconv.Itoa(useTime)},
		"li":        {li},
	}.Encode()))
	if err != nil {
		return errors.New("HTTP Create Request Failed." + err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header["UserID"] = []string{strconv.FormatUint(uint64(session.token.UserID), 10)}
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

	if uploadResult.Status != ServiceSuccessStatus {
		return serviceError{
			status:  uploadResult.Status,
			message: uploadResult.ErrorMessage,
		}
	}
	return nil
}

func (session *Session) GetUserSportResult() (UserSportResult, error) {
	if !session.token.Valid() {
		return UserSportResult{}, ErrTokenExpired
	}
	req, err := http.NewRequest(http.MethodPost, getSportResultURL, strings.NewReader("flag=0"))
	if err != nil {
		return UserSportResult{}, fmt.Errorf("HTTP Create Request Failed. %s", err.Error())
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header["UserID"] = []string{strconv.FormatUint(uint64(session.token.UserID), 10)}
	req.Header["TokenID"] = []string{session.token.TokenID}
	req.Header["crack"] = []string{"0"}
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
	if responseResult.Status != ServiceSuccessStatus {
		return UserSportResult{}, serviceError{
			status:  responseResult.Status,
			message: responseResult.ErrorMessage,
		}
	}
	lastTime, err := fromRPCTimeStr(responseResult.LastTime)
	if err != nil {
		return UserSportResult{}, fmt.Errorf("reslove Failed. %s", err.Error())
	}
	return UserSportResult{
		UserID:            responseResult.UserID,
		Year:              responseResult.Year,
		Term:              responseResult.Term,
		QualifiedDistance: responseResult.Qualified,
		RunDistance:       responseResult.Result,
		LastTime:          lastTime,
	}, nil
}
