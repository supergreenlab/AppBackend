package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

// curl 'http://localhost:9090/api/v1/query_range?query=g_BOX_0_SHT1X_TEMP_C%7Bid%3D%22a5a524ceee3a7d80%22%7D&start=1552927230&end=1552934445&step=15'

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

func (r RangeResult) toFloat64(min, max float64) [][]float64 {
	res := [][]float64{}
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

func queryProm(query string, start, end int64, res *RangeResult) error {
	c := http.DefaultClient

	v := url.Values{}
	v.Set("query", query)
	v.Set("start", fmt.Sprintf("%d", start))
	v.Set("end", fmt.Sprintf("%d", end))
	v.Set("step", fmt.Sprintf("%d", (end-start)/200))
	u, err := url.Parse(fmt.Sprintf("http://prometheus:9090/api/v1/query_range?%s", v.Encode()))
	if err != nil {
		return err
	}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	/*b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s", b)*/

	err = json.NewDecoder(resp.Body).Decode(&res)
	return err
}

type ServedResult struct {
	Humi [][]float64
	Temp [][]float64
}

func serveRange(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Access-Control-Allow-Origin")
	if r.Method == "OPTIONS" {
		w.WriteHeader(200)
		return
	} else if r.Method != "GET" {
		w.WriteHeader(404)
		return
	}
	box, err := strconv.Atoi(r.URL.Query().Get("box"))
	if err != nil {
		log.Println(err)
		w.WriteHeader(404)
		return
	}

	controller := r.URL.Query().Get("controller")
	if controller == "" {
		w.WriteHeader(404)
		return
	}

	humi := RangeResult{}
	queryProm(fmt.Sprintf(`g_BOX_%d_SHT1X_HUMI{id="%s"}`, box, controller), time.Now().Unix()-60*60*24*3, time.Now().Unix(), &humi)

	if humi.Status != "success" {
		w.WriteHeader(404)
		return
	}

	temp := RangeResult{}
	queryProm(fmt.Sprintf(`g_BOX_%d_SHT1X_TEMP_C{id="%s"}`, box, controller), time.Now().Unix()-60*60*24*3, time.Now().Unix(), &temp)

	if temp.Status != "success" {
		w.WriteHeader(404)
		return
	}

	if len(humi.Data.Result) < 1 || len(temp.Data.Result) < 1 {
		w.WriteHeader(200)
		return
	}

	sr := ServedResult{Humi: humi.toFloat64(0, 90), Temp: temp.toFloat64(0, 60)}

	js, err := json.Marshal(sr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func main() {
	http.HandleFunc("/", serveRange)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
