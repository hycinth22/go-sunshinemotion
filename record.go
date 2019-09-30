package ssmt

import (
	"encoding/json"
	"sort"
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

// Implements the interface sort.Interface
type RecordSetForSort []Record

func (r *RecordSetForSort) Len() int {
	return len(*r)
}
func (r *RecordSetForSort) Less(i int, j int) bool {
	return (*r)[i].BeginTime.Before((*r)[j].BeginTime)
}
func (r *RecordSetForSort) Swap(i int, j int) {
	(*r)[i], (*r)[j] = (*r)[j], (*r)[i]
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

func generateRandomTimeDuration(mpkmLimit Float64Range, distance float64) time.Duration {
	var mpkm float64 // MinutePerKM
	if mpkmLimit.Min != mpkmLimit.Max {
		mpkm = randRangeFloat(mpkmLimit.Min+1, mpkmLimit.Max)
	} else {
		mpkm = mpkmLimit.Min
	}
	println("minutePerKM:", mpkm)
	minute := distance * mpkm
	println("minute", minute)
	return time.Duration(minute*time.Minute.Seconds()) * time.Second
}

func nextTimeRangeASC(mpkmLimit Float64Range, distance float64, lastBeginTime time.Time, lastEndTime time.Time) (beginTime time.Time, endTime time.Time) {
	runDuration := generateRandomTimeDuration(mpkmLimit, distance)
	println("runDuration:", runDuration.String())

	idleDuration := time.Duration(randRange(1, 30))*time.Minute + time.Duration(randRange(1, 60))*time.Second
	beginTime = lastEndTime.Add(idleDuration)
	endTime = beginTime.Add(runDuration).Round(time.Second)
	return
}

func nextTimeRangeDESC(mpkmLimit Float64Range, distance float64, lastBeginTime time.Time, lastEndTime time.Time) (beginTime time.Time, endTime time.Time) {
	runDuration := generateRandomTimeDuration(mpkmLimit, distance)
	println("runDuration:", runDuration.String())

	idleDuration := time.Duration(randRange(1, 30))*time.Minute + time.Duration(randRange(1, 60))*time.Second
	endTime = lastBeginTime.Add(-idleDuration)
	beginTime = endTime.Add(-runDuration)
	return
}

func smartCreateRecords(schoolID int64, userID int64, limitParams LimitParams, remain float64, timePoint time.Time,
	timeRangeGenerator func(mpkmLimit Float64Range, distance float64, lastBeginTime time.Time, lastEndTime time.Time) (beginTime time.Time, endTime time.Time)) (records []Record) {
	println("remain", remain)

	var tmp RecordSetForSort = make([]Record, 0, int(remain/3))

	var (
		beginTime = timePoint
		endTime   = timePoint
	)
	sum := 0.0
	for remain >= limitParams.LimitSingleDistance.Min-EpsilonDistance {
		singleDistance := smartCreateDistance(limitParams, remain)
		if singleDistance < limitParams.LimitSingleDistance.Min-EpsilonDistance {
			break
		}
		beginTime, endTime = timeRangeGenerator(limitParams.MinutePerKM, singleDistance, beginTime, endTime)

		records = append(records, createRecord(userID, schoolID, NormalizeDistance(singleDistance), beginTime, endTime))

		remain -= singleDistance
		sum += singleDistance
		if sum > limitParams.LimitTotalMaxDistance+EpsilonDistance {
			if remain > 0.0+EpsilonDistance {
				println("drop remain:", remain)
				break
			}
			break
		}
	}
	sort.Sort(&tmp)
	return tmp
}

func SmartCreateRecordsAfter(schoolID int64, userID int64, limitParams LimitParams, remain float64, afterTime time.Time) (records []Record) {
	return smartCreateRecords(schoolID, userID, limitParams, remain, afterTime, nextTimeRangeASC)
}

func SmartCreateRecordsBefore(schoolID int64, userID int64, limitParams LimitParams, remain float64, beforeTime time.Time) []Record {
	return smartCreateRecords(schoolID, userID, limitParams, remain, beforeTime, nextTimeRangeDESC)
}

func CreateRecord(userID int64, schoolID int64, distance float64, endTime time.Time, duration time.Duration) Record {
	return createRecord(userID, schoolID, NormalizeDistance(distance), endTime.Add(-duration), endTime)
}

func createRecord(userID int64, schoolID int64, distance float64, beginTime, endTime time.Time) Record {
	r := Record{
		UserID:    userID,
		SchoolID:  schoolID,
		Distance:  distance,
		BeginTime: beginTime,
		EndTime:   endTime,
		IsValid:   true,
	}
	return r
}
