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
	t.Run("Test_makeQuery subTest", func(t *testing.T) {
		t.Parallel()
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
		if query != "item_param=BBD7C2A1A48BF5184498E36D6D2889892439F8A60DF895E5538082266B7AD4BCD2CFA18933972621D4D32125C6D9EE6CBECDD7DAB9410A6A2ABC5559CE433857D2CFA18933972621D4D32125C6D9EE6C1E6A992627A2F6C839F66FD21AE0C79D8FE3B87B86999C5D0C6FCCFAE4D159E98B31340A8B595EB917C50C7F9D0D4E0CEE8C4FF8833572011DFA1FBB947F351B9F07DF14CF9E9FE36DC698CF3D9D54FC49B8B06DC3124202BD953A4C5A72DB7A0A8B0C6A496C6808710ADD2E1749380EB9050163CC6A9219A8CD95E6E1B132052C2C71A32B65FD4CBC6A252342DF045B8B9FA1B7AF1AFBC33BD0AD5972294C973EC466F85A53CDEC99862206ABC272E6508BEB826CBDAA66655B8646B775AAB8BA67A62E5EE9B42D8DBD31C57E07D02B5A542FFF46C5A676B952DBD3E9B2441CB74C8D1EE4D09A25AEE7D67AA55B0C4812750D5FC845331F4A25EAF96BCCE35B99F14A48033A6F87ACD615983D9C3ADE0D3814AE4D71B2830BA22722E8064450CF91BEFFA684DA4ED9072B3F86BCC3FD3135437193DAF14A256747793CF8410F33C4156A4F31EAE5D8F354630B664C6BF7FD3E55F607F312CEE32A23E7AE9C24" {
			t.Fail()
		}
		t.Logf("%s", query)
	})
}
