package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"
)

const UNISWAP_V3_API_URL = "https://api.thegraph.com/subgraphs/name/uniswap/uniswap-v3"

const START_DAY = "2022-01-01T00:00:00Z"

const END_DAY = "2022-02-28T00:00:00Z"

// maximum pagination size
const PAGE_SIZE = 1000

// query template
const POOL_DAY_DATA_QUERY_TEMPLATE = `{
	poolDayDatas(first: {{.First}}, where: {id_gt: "{{.IdGt}}", date_gte: {{.StartDate}}, date_lte: {{.EndDate}}}) {
	  id
	  date
	  pool {
		id
	  }
	  feesUSD
	  tvlUSD
	}
  }`

type Params struct {
	First     int64
	IdGt      string
	StartDate int64
	EndDate   int64
}

type PoolDayDataResponse struct {
	Data struct {
		PoolDayDatas []PoolDayData `json:"poolDayDatas"`
	} `json:"data"`
}

type PoolDayData struct {
	Id   string `json:"id"`
	Date int64  `json:"date"`
	Pool struct {
		Id string `json:"id"`
	} `json:"pool"`
	FeesUSD string `json:"feesUSD"`
	TvlUSD  string `json:"tvlUSD"`
}

// Convert datetime string into unix timestamp
func ConvertDatetimeToUnixTimestamp(datetime string) (int64, error) {
	time, err := time.Parse(time.RFC3339, datetime)
	if err != nil {
		return 0, err
	}
	return time.Unix(), nil
}

func BuildQuery(tmpl *template.Template, params Params) (string, error) {
	var query bytes.Buffer
	if err := tmpl.Execute(&query, params); err != nil {
		return "", err
	}
	return query.String(), nil
}

func FetchPoolDayDatas(query string) ([]PoolDayData, error) {
	payload := struct {
		Query string `json:"query"`
	}{
		query,
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(UNISWAP_V3_API_URL, "application/json", bytes.NewReader(payloadBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response PoolDayDataResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}
	return response.Data.PoolDayDatas, nil
}

func Largest(earningRates map[string]float64) (string, float64) {
	var largestPoolId string
	var largestEarningRate float64
	for poolId, earningRate := range earningRates {
		if earningRate > largestEarningRate {
			largestPoolId = poolId
			largestEarningRate = earningRate
		}
	}
	return largestPoolId, largestEarningRate
}

func main() {
	tmpl := template.Must(template.New("query").Parse(POOL_DAY_DATA_QUERY_TEMPLATE))
	startDate, err := ConvertDatetimeToUnixTimestamp(START_DAY)
	if err != nil {
		log.Fatal(err)
	}
	endDate, err := ConvertDatetimeToUnixTimestamp(END_DAY)
	if err != nil {
		log.Fatal(err)
	}

	// efficient pagination leveaging the last-id pattern
	// see https://thegraph.com/docs/en/developer/graphql-api/#example-4
	var lastId string

	// poolId -> earningRate mapping
	earningRates := map[string]float64{}

	// Number of records fetched
	nFetched := 0

	// Fetch poolDayDatas until there is no more
	for {
		query, err := BuildQuery(tmpl, Params{
			First:     PAGE_SIZE,
			IdGt:      lastId,
			StartDate: startDate,
			EndDate:   endDate,
		})
		if err != nil {
			log.Fatal(err)
		}
		datas, err := FetchPoolDayDatas(query)
		if err != nil {
			log.Fatal(err)
		}

		// reach the end of pagination
		if len(datas) == 0 {
			log.Println("done")
			break
		}

		nFetched += len(datas)
		log.Println("fetched", nFetched, "datas")

		for _, data := range datas {
			// update lastId
			lastId = data.Id

			// calculate the earning rate of the day
			feesFloat, err := strconv.ParseFloat(data.FeesUSD, 64)
			if err != nil {
				log.Fatal(err)
			}
			tvlFloat, err := strconv.ParseFloat(data.TvlUSD, 64)
			if err != nil {
				log.Fatal(err)
			}

			// skip when the pool has no swap
			if feesFloat == 0 || tvlFloat == 0 {
				continue
			}
			earningRateDelta := feesFloat / tvlFloat

			// update earning rate
			if _, ok := earningRates[data.Pool.Id]; !ok {
				earningRates[data.Pool.Id] = earningRateDelta
			} else {
				earningRates[data.Pool.Id] += earningRateDelta
			}
		}
	}

	largestPoolId, largestEarningRate := Largest(earningRates)
	fmt.Println(largestPoolId)
	fmt.Println(largestEarningRate)
}
