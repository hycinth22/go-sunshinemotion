package lib

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"
)

type Record struct {
	UserID    int64
	SchoolID  int64
	Distance  float64
	BeginTime time.Time
	EndTime   time.Time
	xtcode    string
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

func (r XTJsonSportData) ToJSON() string {
	j, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return string(j)
}

func (r XTJsonSportData) GetStrpa() string {
	return EncodeSportData(r.ToJSON())
}

func (r XTJsonSportTestData) ToJSON() string {
	j, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	return string(j)
}

func (r XTJsonSportTestData) GetStrpa() string {
	return EncodeSportData(r.ToJSON())
}

func SmartCreateRecords(userID int64, schoolID int64, limitParams *LimitParams, distance float64, beforeTime time.Time) []Record {
	records := make([]Record, 0, int(distance/3))
	remain := distance
	lastBeginTime := beforeTime
	for remain > 0 {
		var singleDistance float64
		// 范围取随机
		if remain > limitParams.RandDistance.Max {
			// 检查是否下一条可能丢弃较大的距离
			// 防止：剩下比较多，但却不满足最小限制距离，不能生成下一条记录
			if remain-limitParams.RandDistance.Max > limitParams.LimitSingleDistance.Min {
				// 正常取随机值
				singleDistance = float64(randRange(int(limitParams.RandDistance.Min*1000), int(limitParams.RandDistance.Max*1000))) / 1000
			} else {
				// 随机选择本条为最小限制距离，或者为下一条预留最小限制距离
				singleDistance = []float64{limitParams.LimitSingleDistance.Min, remain - limitParams.LimitSingleDistance.Min}[randRange(0, 1)]
			}
		} else if remain >= limitParams.LimitSingleDistance.Min && remain <= limitParams.LimitSingleDistance.Max {
			// 剩余的符合限制区间，直接使用剩余的生成最后一条记录
			singleDistance = remain
		} else {
			// 剩余较多，但不符合限制区间无法再生成一条合法记录，输出提醒
			if remain > 0.5 {
				fmt.Println("提醒：由于随机原则与区间限制的冲突，丢弃了较大的距离", remain, "公里，考虑重新设定距离值。")
			}
			break
		}

		// 小数部分随机化 -0.09 ~ 0.09
		tinyPart := float64(randRange(0, 99999)) / 1000000
		switch r := singleDistance + tinyPart; {
		case r < limitParams.LimitSingleDistance.Min:
			singleDistance = limitParams.LimitSingleDistance.Min
			/*case r > userInfo.LimitParams.LimitSingleDistance.Max:
			singleDistance = userInfo.LimitParams.LimitSingleDistance.Max
			*/
		default:
			singleDistance += tinyPart
		}

		// 检测结果合法性，由于TinyPart允许上下浮动0.1
		if singleDistance < limitParams.LimitSingleDistance.Min-0.1 || singleDistance > limitParams.LimitSingleDistance.Max+0.1 {
			// 丢弃不合法距离
			log.Println("Drop distance: ", singleDistance)
			continue
		}

		var randomDuration time.Duration
		// 时间间隔随机化
		randomDuration = time.Duration(randRange(limitParams.MinuteDuration.Min, limitParams.MinuteDuration.Max)) * time.Minute
		randomDuration += time.Duration(randRange(0, 60)) * time.Second // 时间间隔秒级随机化

		endTime := lastBeginTime.Add(-time.Duration(randRange(1, 30)) * time.Minute)
		endTime = endTime.Add(-time.Duration(randRange(1, 60)) * time.Second)
		beginTime := endTime.Add(-randomDuration)

		records = append(records, Record{
			UserID:    userID,
			SchoolID:  schoolID,
			Distance:  singleDistance,
			BeginTime: beginTime,
			EndTime:   endTime,
			xtcode:    GetXTcodeV3(userID, toExchangeTimeStr(beginTime), toExchangeDistanceStr(singleDistance)),
		})

		remain -= singleDistance - tinyPart
		lastBeginTime = beginTime
	}
	nRecord := len(records)
	reverse := make([]Record, nRecord)
	for i := 0; i < nRecord; i++ {
		reverse[i] = records[nRecord-i-1]
	}
	return reverse
}
func CreateRecord(userID int64, schoolID int64, distance float64, beforeTime time.Time, duration time.Duration) Record {
	r := Record{
		UserID:    userID,
		SchoolID:  schoolID,
		Distance:  distance,
		BeginTime: beforeTime.Add(-duration),
		EndTime:   beforeTime,
	}
	r.xtcode = GetXTcodeV3(userID, toExchangeTimeStr(r.BeginTime), toExchangeDistanceStr(r.Distance))
	return r
}

func GetXTcode(userId int64, beginTime string) string {
	key := MD5String(strconv.FormatInt(userId, 10) + beginTime + "stlchang")
	var xtCode bytes.Buffer
	xtCode.WriteByte(key[7])
	xtCode.WriteByte(key[3])
	xtCode.WriteByte(key[15])
	xtCode.WriteByte(key[24])
	xtCode.WriteByte(key[9])
	xtCode.WriteByte(key[17])
	xtCode.WriteByte(key[29])
	xtCode.WriteByte(key[23])
	return xtCode.String()
}

func GetXTcodeV2(userId int64, beginTime string, distance string) string {
	phrase := strconv.FormatInt(userId, 10) + beginTime + distance + "stlchang"
	key := MD5String(phrase)
	log.Println(phrase, key)
	var xtCode bytes.Buffer
	xtCode.WriteByte(key[7])
	xtCode.WriteByte(key[3])
	xtCode.WriteByte(key[15])
	xtCode.WriteByte(key[24])
	xtCode.WriteByte(key[9])
	xtCode.WriteByte(key[17])
	xtCode.WriteByte(key[29])
	xtCode.WriteByte(key[23])
	return xtCode.String()
}
