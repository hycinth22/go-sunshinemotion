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
	loginErr = session.LoginEx("031840607", "123", fmt.Sprintf("%x", md5.Sum([]byte("123456"))), 60)
}
func TestLoginEx(t *testing.T) {
	session = CreateSession()
	loginErr = session.LoginEx("031840607", "123", fmt.Sprintf("%x", md5.Sum([]byte("123456"))), 60)
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
	return
	records := SmartCreateRecords(session.UserID, session.UserInfo.InSchoolID, &LimitParams{
		RandDistance:        Float64Range{2.6, 4.0},
		LimitSingleDistance: Float64Range{2.0, 4.0},
		LimitTotalDistance:  Float64Range{2.0, 5.0},
		MinuteDuration:      IntRange{11, 20},
	}, 2, time.Now())

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
	return
	r := CreateRecord(session.UserID, session.UserInfo.InSchoolID, 3.211, time.Now(), 16*time.Minute+12*time.Second)
	t.Logf("%+v", r)
	err := session.UploadTestRecord(r)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}
