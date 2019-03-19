package ssmt

type LimitParams struct {
	// 随机区间（若还有足够的欲生成的距离， 生成单条记录的的随机距离区间）
	RandDistance Float64Range
	// 限制区间（生成单条记录的的硬性限制）
	LimitSingleDistance Float64Range
	// 限制区间（所有记录总和的限制区间）
	LimitTotalDistance Float64Range
	// 每条记录的时间区间
	MinuteDuration IntRange
}

func GetDefaultLimitParams(sex string) LimitParams {
	// 参数设定：
	// MinuteDuration: min>minDis*3, max<maxDis*10
	switch sex {
	case "F":
		return LimitParams{
			RandDistance:        Float64Range{2.5, 3.0},
			LimitSingleDistance: Float64Range{1.0, 3.0},
			LimitTotalDistance:  Float64Range{1.0, 3.0},
			MinuteDuration:      IntRange{25, 35},
		}
	case "M":
		return LimitParams{
			RandDistance:        Float64Range{4.5, 5.0},
			LimitSingleDistance: Float64Range{2.0, 5.0},
			LimitTotalDistance:  Float64Range{2.0, 5.0},
			MinuteDuration:      IntRange{45, 65},
		}
	}
	return LimitParams{}
}
