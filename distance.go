package ssmt

import (
	"math"
)

const (
	EpsilonDistance    = 0.0004999999999999999
	WidthAfterDecPoint = 3 // 保留3位小数
)

var (
	widthFactor      = math.Pow10(WidthAfterDecPoint)
	distanceAccuracy = math.Pow10(-WidthAfterDecPoint)
)

// 我们限制区间[min, max) 不能到达右区间。
// 当最后结果四舍五入，可能造成到达右区间。
//
// 该函数对此情况做了修正，保证返回值在经过四舍五入之后也仍然在区间之内。
func randRangeDistance(min, max float64) float64 {
	// 对max修正的同时，也对min做了修正，这是为了尽力维持原概率分布。
	if math.Trunc(max*widthFactor)/widthFactor >= EpsilonDistance {
		min = math.Trunc(min*widthFactor)/widthFactor - EpsilonDistance
		max -= EpsilonDistance
	}
	return NormalizeDistance(randRangeFloat(min, max))
}

// 返回x的浮动区间[min, max)，保证浮动区间内四舍五入后仍然在该精度精度上等于x
func DistanceRangeAround(x float64, width int) (min, max float64) {
	wf := math.Pow10(width)
	t := math.Trunc(x * wf)
	return (t - 0.5) / wf, (t + 0.5) / wf
}

func NormalizeDistance(distance float64) (normalizeDistance float64) {
	return math.Round(distance*widthFactor) / widthFactor
}

// 结果
// 尽量位于[RandDistance.Min, RandDistance.Max)
// 一定位于[LimitSingleDistance.Min, LimitSingleDistance.Max)
// 可能返回0.0，代表在LimitParams下无法从remain生成距离
func smartCreateDistance(limitParams LimitParams, remain float64) (singleDistance float64) {
	// 参数检查
	if remain-limitParams.LimitSingleDistance.Min < EpsilonDistance {
		println("smartCreateDistance参数不正确", singleDistance)
		return 0.0
	}
	// 兜底检测，检测结果合法性
	defer func() {
		if singleDistance != 0.0 {
			if singleDistance-limitParams.LimitSingleDistance.Min < -EpsilonDistance {
				// 丢弃不合法距离
				println("检查算法正确性, Too little distance: ", singleDistance)
				singleDistance = 0.0
			}
			if singleDistance-limitParams.LimitSingleDistance.Max >= EpsilonDistance {
				println("检查算法正确性, Too much distance", singleDistance)
				singleDistance = 0.0
			}
		}
	}()
	if remain-limitParams.RandDistance.Min >= EpsilonDistance {
		// 剩余足够大，正常取随机值
		// Use RandDistance Params
		low := math.Max(limitParams.RandDistance.Min, limitParams.LimitSingleDistance.Min)
		high := math.Min(remain, math.Min(limitParams.RandDistance.Max, limitParams.LimitSingleDistance.Max))
		println("remain", remain)
		println("low", low)
		println("high", high)
		if low <= high {
			singleDistance = randRangeDistance(low, high)
			println("p1", singleDistance)
			return NormalizeDistance(singleDistance)
		} else {
			println("fail to inRand, downgrade")
		}
	}
	// Downgrade
	low := limitParams.LimitSingleDistance.Min
	high := math.Min(remain, limitParams.LimitSingleDistance.Max)
	println("remain", remain)
	println("low", low)
	println("high", high)
	if remain-limitParams.LimitSingleDistance.Min >= EpsilonDistance && low <= high {
		singleDistance = randRangeDistance(low, high)
		println("p2", singleDistance)
	} else {
		println("drop", singleDistance)
		return 0.0
	}
	return NormalizeDistance(singleDistance)
}
