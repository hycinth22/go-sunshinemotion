package ssmt

import (
	"encoding/json"
	"log"
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
	for remain > 0 {
		singleDistance := smartCreateDistance(schoolID, userID, limitParams, remain)
		if singleDistance <= 0.0+EPSILON_Distance {
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
	for remain > 0 {
		singleDistance := smartCreateDistance(schoolID, userID, limitParams, remain)
		if singleDistance <= 0.0+EPSILON_Distance {
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
		lastBeginTime = beginTime
	}
	nRecord := len(records)
	reverse := make([]Record, nRecord)
	for i := 0; i < nRecord; i++ {
		reverse[i] = records[nRecord-i-1]
	}
	return reverse
}

func smartCreateDistance(schoolID int64, userID int64, limitParams LimitParams, remain float64) (singleDistance float64) {
	const tinyPartLimit = 0.10
	// 范围取随机
	// 会检查是否下一条可能丢弃较大的距离，防止：剩下比较多，但却不满足最小限制距离，不能生成下一条记录
	if remain >= 2*limitParams.LimitSingleDistance.Max {
		// 剩余足够大，正常取随机值
		singleDistance = randRangeFloat(limitParams.RandDistance.Min, limitParams.RandDistance.Max-tinyPartLimit)
		println("p1", singleDistance)
	} else if remain >= 2*limitParams.LimitSingleDistance.Min {
		// 即将耗尽，首先尝试放入一条记录内，否则为下一条预留
		if remain <= limitParams.LimitSingleDistance.Max {
			singleDistance = remain
		} else {
			// 为下一条预留最小限制距离
			singleDistance = remain - limitParams.LimitSingleDistance.Min
		}
		println("p2", singleDistance)
	} else if remain >= limitParams.LimitSingleDistance.Min {
		// 剩余的符合最小限制距离，直接使用剩余的生成最后一条记录
		singleDistance = remain
		println("p3", singleDistance)
	} else if remain > EPSILON_Distance {
		println("检查算法正确性", remain)
		return 0.0
	} else {
		// 最后的零星距离，可以直接丢弃
		return 0.0
	}

	// 小数部分随机化 -0.09 ~ 0.09
	tinyPart := randRangeFloat(0.0, tinyPartLimit)
	switch r := singleDistance + tinyPart; {
	case r < limitParams.LimitSingleDistance.Min:
		singleDistance = limitParams.LimitSingleDistance.Min
		/*case r > userInfo.LimitParams.LimitSingleDistance.Max:
		singleDistance = userInfo.LimitParams.LimitSingleDistance.Max
		*/
	default:
		singleDistance += tinyPart
	}

	// 检测结果合法性
	if singleDistance < limitParams.LimitSingleDistance.Min-EPSILON_Distance {
		// 丢弃不合法距离
		log.Println("Drop distance: ", singleDistance)
		return 0.0
	}
	if singleDistance > limitParams.LimitSingleDistance.Max {
		singleDistance = limitParams.LimitSingleDistance.Max - EPSILON_Distance
	}
	return NormalizeDistance(singleDistance)
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
