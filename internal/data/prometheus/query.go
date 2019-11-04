package prometheus

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// QueryProm fetches metrics from prometheus
func QueryProm(query string, start, end int64, n int) (RangeResult, error) {
	res := RangeResult{}
	c := http.DefaultClient

	v := url.Values{}
	v.Set("query", query)
	v.Set("start", fmt.Sprintf("%d", start))
	v.Set("end", fmt.Sprintf("%d", end))
	v.Set("step", fmt.Sprintf("%d", (end-start)/int64(n)))
	u, err := url.Parse(fmt.Sprintf("http://prometheus:9090/api/v1/query_range?%s", v.Encode()))
	if err != nil {
		return res, err
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return res, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return res, err
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&res)
	return res, err
}
