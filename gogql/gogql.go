package gogql

import (
	"encoding/json"
	"fmt"
	"github.com/abeytom/goson"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"strings"
	"time"
)

type GraphqlClient struct {
	HttpClient    *http.Client
	Url           string
	Authorization string
}

type GraphQlRequest struct {
	ResultKeys []string `json:"resultKeys"`
	TotalKeys  []string `json:"totalKeys"`
	CountKeys  []string `json:"countKeys"`
	Graphql    string   `json:"graphql"`
	Limit      int      `json:"limit"`
	Max        int      `json:"max"`
	StartTime  string   `json:"startTime"`
	EndTime    string   `json:"endTime"`

	//state
	offset int `json:"offset"`
}

func (g *GraphqlClient) AsGqlPayload(gql string) string {
	jsonMap := make(map[string]interface{})
	jsonMap["variables"] = make(map[string]string)
	jsonMap["query"] = gql
	marshal, err := json.Marshal(jsonMap)
	if err != nil {
		panic(err)
	}
	return string(marshal)
}

func (g *GraphqlClient) ExecuteGraphqlIter(r GraphQlRequest, call func(*goson.ArrayNode) error) error {
	offset := r.offset
	index := 0
	total := 0
	for {
		log.Info().Msgf("running for offset %v", offset)
		payload := strings.ReplaceAll(g.AsGqlPayload(r.Graphql), "${OFFSET}", fmt.Sprintf("%v", offset))
		if r.Limit > 0 {
			payload = strings.ReplaceAll(payload, "${LIMIT}", fmt.Sprintf("%v", r.Limit))
		}
		if len(r.StartTime) > 0 {
			payload = strings.ReplaceAll(payload, "${START_TIME}", r.StartTime)
		}
		if len(r.EndTime) > 0 {
			payload = strings.ReplaceAll(payload, "${END_TIME}", r.EndTime)
		}
		b, err := g.ExecuteGraphQL(payload)
		if err != nil {
			return err
		}
		jsonNode, err := goson.ParseBytes(b)
		if err != nil {
			return errors.Wrapf(err, "cannot parse response %v", string(b))
		}
		mapNode := jsonNode.(*goson.MapNode)
		results := mapNode.GetArray(r.ResultKeys...)
		if results != nil && len(results.Objects) > 0 {
			err = call(results)
			if err != nil {
				return err
			}
		} else {
			errs := mapNode.GetArray("errors")
			if errs != nil && len(errs.Objects) > 0 {
				msg := ""
				for _, e := range errs.ItemsAsMap() {
					msg += e.GetString("message") + "."
				}
				return errors.Errorf("error exec http request with url=[%v] and payload=[%v]. The error is [%v]",
					g.Url, payload, msg)
			}
			log.Warn().Msgf("the results are null for %v", payload)
		}
		totalNode := mapNode.GetValue(r.TotalKeys...)
		if totalNode == nil {
			return errors.Wrapf(err, "cannot get the total from [%v]", r.TotalKeys)
		}
		countNode := mapNode.GetValue(r.CountKeys...)
		if countNode == nil {
			return errors.Wrapf(err, "cannot get the count from [%v]", r.CountKeys)
		}
		total = int(totalNode.Value().(float64))
		count := int(countNode.Value().(float64))
		offset += count
		index += count
		if index >= total {
			break
		}
		if r.Max > 0 && index >= r.Max {
			log.Info().Msgf("existing the pagination since the results [%v] crossed a max of [%v]. The total is  [%v]",
				index, r.Max, total)
			break
		}
	}
	log.Info().Msgf("completed the iteration, received %v of %v results", index, total)
	return nil
}

func (g *GraphqlClient) ExecuteGraphQL(payload string) ([]byte, error) {
	request, err := http.NewRequest("POST", g.Url, strings.NewReader(payload))
	if err != nil {
		return nil, errors.Wrapf(err, "cannot create new request for graphql")
	}
	if len(g.Authorization) > 0 {
		request.Header.Add("Authorization", g.Authorization)
	}
	log.Debug().Msgf("running the query \n%v", payload)
	resp, err := g.HttpClient.Do(request)
	if err != nil {
		return nil, errors.Wrapf(err,
			"error exec http request with url=[%v] and payload=[%v]", g.Url, payload)
	}
	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err,
			"error reading the response from url=[%v] and payload=[%v]", g.Url, payload)
	}
	if resp.StatusCode != 200 {
		return nil, errors.Wrapf(err,
			"bad http response from url=[%v] and payload=[%v]; status=[%v],response=[%v]",
			g.Url, payload, resp.Status, string(respBytes))
	}
	return respBytes, nil
}

func TimeRangeLastX(duration time.Duration) (string, string) {
	now := time.Now()
	startMillis := now.UnixMilli() - duration.Milliseconds()
	return time.UnixMilli(startMillis).Format(time.RFC3339), now.Format(time.RFC3339)
}
