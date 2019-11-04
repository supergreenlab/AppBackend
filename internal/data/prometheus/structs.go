package prometheus

import "strconv"

// RangeResult query result from prometheus
type RangeResult struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric struct {
				Name     string `json:"__name__"`
				ID       string `json:"id"`
				Instance string `json:"instance"`
				Job      string `json:"job"`
				Module   string `json:"module"`
			} `json:"metric"`
			Values [][]interface{} `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

// ToFloat64 returns the RangeResult as an array of [timestamp, value]
func (r RangeResult) ToFloat64(min, max float64) [][]float64 {
	res := [][]float64{}
	if len(r.Data.Result) < 1 {
		return res
	}
	var lasti float64
	for _, v := range r.Data.Result[0].Values {
		i, err := strconv.ParseFloat(v[1].(string), 64)
		if err != nil || i < min || i > max {
			i = lasti
		} else {
			lasti = i
		}
		res = append(res, []float64{
			v[0].(float64),
			i,
		})
	}
	return res
}
