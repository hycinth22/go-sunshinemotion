package ssmt

import "math"

type LimitParams struct {
	// 随机区间（若还有足够的欲生成的距离，生成单条记录的的随机距离区间）
	RandDistance Float64Range
	// 限制区间（生成单条记录的的硬性限制）
	LimitSingleDistance Float64Range
	// 所有记录总和的最大值
	LimitTotalMaxDistance float64
	// 每公里的分钟系数区间
	MinutePerKM Float64Range
}

const (
	EpsilonDistance  = 0.0004999_9999_9999_9999
	DistanceAccuracy = 0.001
)

func GetDefaultLimitParams(sex string) LimitParams {
	// 参数设定：
	// MinuteDuration: min>minDis*3, max<maxDis*10
	switch sex {
	case "F":
		return LimitParams{
			RandDistance:          Float64Range{2.995, 3.0},
			LimitSingleDistance:   Float64Range{1.0, 3.0},
			LimitTotalMaxDistance: 3.0,
			MinutePerKM:           Float64Range{7.0, 13.0},
		}
	case "M":
		return LimitParams{
			RandDistance:          Float64Range{4.995, 5.0},
			LimitSingleDistance:   Float64Range{2.0, 5.0},
			LimitTotalMaxDistance: 5.0,
			MinutePerKM:           Float64Range{8.0, 14.0},
		}
	}
	return LimitParams{}
}

// 区间合法、RandDistance与LimitSingleDistance是否有区间交集
func (p LimitParams) IsValid() bool {
	return p.RandDistance.Min <= p.RandDistance.Max && // valid range
		p.LimitSingleDistance.Min <= p.LimitSingleDistance.Max && // valid range
		math.Max(p.RandDistance.Min, p.LimitSingleDistance.Min) <= math.Min(p.RandDistance.Max, p.LimitSingleDistance.Max) // range overlap
}
