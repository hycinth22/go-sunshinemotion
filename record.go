package ssmt

import (
	"encoding/json"
	"time"
)

type Record struct {
	UserID    int64
	SchoolID  int64
	Distance  float64
	BeginTime time.Time
	EndTime   time.Time
	IsValid   bool
}

type IXTJsonSportData interface {
	ToJSON() string
	GetStrpa() string
}

type XTJsonSportData struct {
	Result       string `json:"results"`
	StartTimeStr string `json:"beginTime"`
	EndTimeStr   string `json:"endTime"`
	IsValid      int    `json:"isValid"`
	SchoolID     int64  `json:"schoolId"`
	BZ           string `json:"bz"`
	XTCode       string `json:"xtCode"`
}

func XTJsonSportDataFromRecord(r Record, bz string) XTJsonSportData {
	d := XTJsonSportData{
		Result:       toServiceStdDistance(r.Distance),
		StartTimeStr: toServiceStdTime(r.BeginTime),
		EndTimeStr:   toServiceStdTime(r.EndTime),
		IsValid:      1,
		BZ:           bz,
		SchoolID:     r.SchoolID,
	}
	if !r.IsValid {
		d.IsValid = 0
	}
	d.XTCode = CalcXTcode(r.UserID, d.StartTimeStr, d.Result)
	return d
}
func (r XTJsonSportData) ToJSON() string {
	j, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return string(j)
}
func (r XTJsonSportData) GetStrpa() string {
	return EncodeString(r.ToJSON())
}

type XTJsonSportTestData struct {
	Result       string `json:"results"`
	StartTimeStr string `json:"beginTime"`
	EndTimeStr   string `json:"endTime"`
	IsValid      int    `json:"isValid"`
	SchoolID     int64  `json:"schoolId"`
	BZ           string `json:"bz"`
	XTCode       string `json:"xtCode"`
	TestTime     int64  `json:"test_time"`
}

func XTJsonSportTestDataFromRecord(r Record, bz string) XTJsonSportTestData {
	d := XTJsonSportTestData{
		Result:       toServiceStdDistance(r.Distance),
		StartTimeStr: toServiceStdTime(r.BeginTime),
		EndTimeStr:   toServiceStdTime(r.EndTime),
		IsValid:      1,
		BZ:           bz,
		SchoolID:     r.SchoolID,
		TestTime:     int64(r.EndTime.Sub(r.BeginTime).Seconds()),
	}
	if !r.IsValid {
		d.IsValid = 0
	}
	d.XTCode = CalcXTcode(r.UserID, d.StartTimeStr, d.Result)
	return d
}
func (r XTJsonSportTestData) ToJSON() string {
	j, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return string(j)
}

func (r XTJsonSportTestData) GetStrpa() string {
	return EncodeString(r.ToJSON())
}

func SmartCreateRecordsAfter(schoolID int64, userID int64, limitParams LimitParams, distance float64, afterTime time.Time) []Record {
	records := make([]Record, 0, int(distance/3))
	remain := distance
	lastEndTime := afterTime
	println("distance", distance)
	sum := 0.0
	for remain >= limitParams.LimitSingleDistance.Min-EpsilonDistance {
		singleDistance := smartCreateDistance(limitParams, remain)
		if singleDistance < limitParams.LimitSingleDistance.Min-EpsilonDistance {
			break
		}
		// 时间间隔随机化
		minutePerKM := randRangeFloat(limitParams.MinutePerKM.Min, limitParams.MinutePerKM.Max)
		randomDuration := time.Duration(distance * minutePerKM * float64(time.Minute))
		println("randomDuration:", randomDuration.String())

		beginTime := lastEndTime.Add(time.Duration(randRange(1, 30))*time.Minute + time.Duration(randRange(1, 60))*time.Second)
		endTime := beginTime.Add(randomDuration)

		records = append(records, createRecord(userID, schoolID, NormalizeDistance(singleDistance), endTime, randomDuration))

		remain -= singleDistance
		lastEndTime = endTime
		sum += singleDistance
		if sum > limitParams.LimitTotalMaxDistance+EpsilonDistance {
			if remain > 0.0+EpsilonDistance {
				println("drop remain:", remain)
				break
			}
			break
		}
	}
	nRecord := len(records)
	reverse := make([]Record, nRecord)
	for i := 0; i < nRecord; i++ {
		reverse[i] = records[nRecord-i-1]
	}
	return reverse
}

func SmartCreateRecordsBefore(schoolID int64, userID int64, limitParams LimitParams, distance float64, beforeTime time.Time) []Record {
	records := make([]Record, 0, int(distance/3))
	remain := distance
	lastBeginTime := beforeTime
	println("distance", distance)
	sum := 0.0
	for remain >= limitParams.LimitSingleDistance.Min-EpsilonDistance {
		singleDistance := smartCreateDistance(limitParams, remain)
		println("singleDistance:", singleDistance)
		if singleDistance < limitParams.LimitSingleDistance.Min-EpsilonDistance {
			break
		}
		// 时间间隔随机化
		minutePerKM := randRangeFloat(limitParams.MinutePerKM.Min, limitParams.MinutePerKM.Max)
		randomDuration := time.Duration(distance * minutePerKM * float64(time.Minute))
		println("randomDuration:", randomDuration.String())

		endTime := lastBeginTime.Add(-time.Duration(randRange(1, 30))*time.Minute - time.Duration(randRange(1, 60))*time.Second)
		beginTime := endTime.Add(-randomDuration)

		records = append(records, createRecord(userID, schoolID, NormalizeDistance(singleDistance), endTime, randomDuration))

		remain -= singleDistance
		println("remain:", remain)
		lastBeginTime = beginTime
		sum += singleDistance
		if sum > limitParams.LimitTotalMaxDistance+EpsilonDistance {
			if remain > 0.0+EpsilonDistance {
				println("drop remain:", remain)
				break
			}
			break
		}
	}
	nRecord := len(records)
	reverse := make([]Record, nRecord)
	for i := 0; i < nRecord; i++ {
		reverse[i] = records[nRecord-i-1]
	}
	return reverse
}

func CreateRecord(userID int64, schoolID int64, distance float64, endTime time.Time, duration time.Duration) Record {
	return createRecord(userID, schoolID, NormalizeDistance(distance), endTime, duration)
}

func createRecord(userID int64, schoolID int64, distance float64, endTime time.Time, duration time.Duration) Record {
	r := Record{
		UserID:    userID,
		SchoolID:  schoolID,
		Distance:  distance,
		BeginTime: endTime.Add(-duration),
		EndTime:   endTime,
		IsValid:   true,
	}
	return r
}
