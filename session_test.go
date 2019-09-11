package ssmt

import (
	"crypto/md5"
	"fmt"
	"testing"
	"time"
)

var session *Session
var loginErr error
var info UserInfo
var beijingZone = time.FixedZone("Beijing Time", int((8 * time.Hour).Seconds()))

func init() {
	session = CreateSession()
	info, loginErr = session.Login(60, "061841216", "17390940000", fmt.Sprintf("%x", md5.Sum([]byte("LJW66666"))))
}

func TestGetSchoolList(t *testing.T) {
	r, err := CreateSession().GetSchoolList()
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", r)
}

func TestLoginEx(t *testing.T) {
	if loginErr != nil {
		t.Log(loginErr.Error())
		t.Fatalf("%v", loginErr)
	}
	t.Logf("%+v", session)
}

func TestGetSportResult(t *testing.T) {
	t.Logf("%+v %+v", info, *session)
	r, err := session.GetSportResult()
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", r)
}

func TestSession_UploadSingleData(t *testing.T) {
	return // Only Test If must required
	r := Record{
		UserID:    session.User.UserID,
		SchoolID:  session.User.SchoolID,
		Distance:  4.871,
		BeginTime: time.Date(2019, 5, 16, 20, 30, 33, 795044373, beijingZone),
		EndTime:   time.Date(2019, 5, 16, 21, 15, 26, 863371645, beijingZone),
		IsValid:   true,
	}
	err := session.UploadRecord(r)
	if err != nil {
		t.Error(err)
		t.Fail()
	}
}
func TestSession_UploadData(t *testing.T) {
	return // Only Test If must required
	if !t.Run("GetRandRoute", TestGetRandRoute) {
		t.Skip()
		return
	}
	endTime := time.Now()
	records := SmartCreateRecordsBefore(session.User.SchoolID, session.User.UserID, GetDefaultLimitParams(info.Sex), 3, endTime)
	for _, r := range records {
		t.Logf("%+v", r)
		//var err error = nil
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
	endTime := time.Date(now.Year(), now.Month(), now.Day(), 8, 25, 0, 0, TimeZoneCST)
	r := CreateRecord(session.User.UserID, session.User.SchoolID, 3.211, endTime, 16*time.Minute+12*time.Second)
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

func TestGetRandRoute(t *testing.T) {
	r, err := session.GetRandRoute()
	if err != nil {
		t.Fatalf("%v", err)
	}
	t.Logf("%+v", r)
}

func Test_makeQuery(t *testing.T) {
	sportData := XTJsonSportData{
		Result:       "0.000",
		StartTimeStr: "2019-03-19 14:38:02",
		EndTimeStr:   "2019-03-19 14:38:07",
		IsValid:      0,
		BZ:           "C4F4DD3BDD73BAEEE4CF3655C868FCBE640C5135CD76567F79420B0AE2D490132542B908D88399BB7BA2DA8E13C88099892AC29ECA1C0992CB2FB893B1C8FA996401394E121F845BA72C70416184371163426924E4EF6CBFA81F231EC8EFFA5DD56404EC64776BFEAEED95EB48C808FC91D54607A2E3A2C58E2AA8CC1EE54B91",
		XTCode:       "08E10F3A5C53D1F2622198360D5833AB",
		SchoolID:     60,
	}
	query := makeQuery(sportData)
	if query != "item_param=894036009DE93B9BF92A69D0E15560D79BAFB2748486833B83AF2BE5E7A6A08FFAF781EFD491E2A42CA0888B48AF6681B62D1CA4CFBC5664F039524C432639F0FAF781EFD491E2A42CA0888B48AF66816A24CB7FBEB8F783B29C0ED944F6815FFAB38DAAC9065348888D19AA8CCD1502A18C6EA1A273EEBC71EAAFA91A65186E107A3E74630E7ADB7AEF99BFB5A511DAD8A7AA9F987AF474138436DCD75A54408D6D1555BC3CA3374BEE550DA1A711A3B4474FC2C3634B0E377A811745979E349030241D5E16366EFBE81AD6E102C7263882D234C72D4A7189A50959FE64EB4B8BFF162BE1E39969A0C4A2730E768EBADAB9B0539246B50CB68126155B4DC446036BDBB80199A42F1B3F616E3F3525A83CBA4E0B2124C7FEB97D56873DD55F3BF3418278763452E2B85B654E0BBBED16AD21CD4ECD0F8FD0F1C5F7FA7D160E4B84FFE2C8FCECF5B6EFDBF8E238BFBE4170DD7A4323FDC212EFF3930458439440393D07A51CEC4D9F3FB9367611474298CA2EB4E608DA91C6967EA1B9A2C36E07D50A34A8194CA2B4BF5826F71D145172443A21FA47B682C14121E94AC51FEEFF4C8AA0C8A20687AC41F89A104478C9AC" {
		t.Fail()
	}
	t.Logf("%s", query)
}
