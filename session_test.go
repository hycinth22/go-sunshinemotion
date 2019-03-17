package lib

import (
	"crypto/md5"
	"fmt"
	"testing"
	"time"
)

var session *Session
var loginErr error

func init() {
	session = CreateSession()
	loginErr = session.LoginEx("091840822", "123", fmt.Sprintf("%x", md5.Sum([]byte("123456"))), 60)
}
func TestLoginEx(t *testing.T) {
	session = CreateSession()
	loginErr = session.LoginEx("091840822", "123", fmt.Sprintf("%x", md5.Sum([]byte("123456"))), 60)
	if loginErr != nil {
		t.Log(loginErr.Error())
		t.Fatalf("%v", loginErr)
	}
	t.Logf("%+v", session)
}

func TestGetSportResult(t *testing.T) {
	r, err := session.GetSportResult()
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", r)
}

func TestSession_UploadData(t *testing.T) {
	return // Only Test If must required
	now := time.Now()
	beijing := time.FixedZone("Beijing Time", int((8 * time.Hour).Seconds()))
	endTime := time.Date(now.Year(), now.Month(), now.Day(), 8, 25, 0, 0, beijing)
	records := SmartCreateRecords(session.UserID, session.UserInfo.InSchoolID, &LimitParams{
		RandDistance:        Float64Range{3.6, 4.0},
		LimitSingleDistance: Float64Range{2.0, 4.0},
		LimitTotalDistance:  Float64Range{2.0, 5.0},
		MinuteDuration:      IntRange{35, 50},
	}, 2, endTime)

	for _, r := range records {
		t.Logf("%+v", r)
		err := session.UploadRecord(r)
		if err != nil {
			t.Error(err)
			t.Fail()
		}
	}
}

func TestSession_UploadTestRecord(t *testing.T) {
	return // Only Test If must required
	now := time.Now()
	beijing := time.FixedZone("Beijing Time", int((8 * time.Hour).Seconds()))
	endTime := time.Date(now.Year(), now.Month(), now.Day(), 8, 25, 0, 0, beijing)
	r := CreateRecord(session.UserID, session.UserInfo.InSchoolID, 3.211, endTime, 16*time.Minute+12*time.Second)
	t.Logf("%+v", r)
	err := session.UploadTestRecord(r)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}

func TestSession_GetAppInfo(t *testing.T) {
	info, err := session.GetAppInfo()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	t.Logf("info: %+v", info)
}
