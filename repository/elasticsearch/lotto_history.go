package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"
	"watcharis/go-migrate-lotto-history-els/models"
	"watcharis/go-migrate-lotto-history-els/repository/rest"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/update"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/updatebyquery"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/scriptlanguage"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
	"go.uber.org/zap/buffer"
)

type lottoHistoryElasticSearch struct {
	elasticClient *elasticsearch.Client
}

type LottoHistoryElasticSearch interface {
	InsertLottoHistory(ctx context.Context, index string, doc models.LottosHistoryElasticsearch) error
	GetLottoHistory(ctx context.Context, index string, uuid string, offsets int, recordPerpage int) ([]models.ElastictHitsHits, error)
	UpdateLottoHistoryRewardStatus(ctx context.Context, index string, id string, doc map[string]interface{}) error
	UpdateLottoHistoryRewardStatusByElasticQuery(ctx context.Context, index string, uuid string, _id int64, source string) error
	ExistsIndex(ctx context.Context, index string) (bool, error)
	CreateIndexWithMapping(ctx context.Context, index string, mapping string) error
	BulkData(ctx context.Context, index string) error
	BulkUpdateData(ctx context.Context, index string, ids []int, rewardStatus string) error
	GetPITtoken(ctx context.Context, index string, keepAlive string) (string, error)
	PointInTimeSearch(ctx context.Context, pitToken string, searchAfter []int64) ([]models.ElastictHitsHits, []int64, error)
	DeletePITtoken(ctx context.Context, pitToken string) error
}

func NewLottoHistoryElasticSearch(elasticClient *elasticsearch.Client) LottoHistoryElasticSearch {
	return &lottoHistoryElasticSearch{
		elasticClient: elasticClient,
	}
}

func (e *lottoHistoryElasticSearch) InsertLottoHistory(ctx context.Context, index string, doc models.LottosHistoryElasticsearch) error {

	data, _ := json.Marshal(doc)
	resp, err := e.elasticClient.Index(index, bytes.NewReader(data), func(request *esapi.IndexRequest) {
		request.DocumentID = strconv.FormatInt(doc.ID, 10)
	})

	if err != nil {
		log.Printf("error elastic insert failed index : %s , error is : %v, uuid: %v", index, err, doc.BuyerUUID)
		return err
	}

	if resp.IsError() {
		log.Printf("error elastic insert failed index : %s , error is : %v, uuid: %v", index, resp, doc.BuyerUUID)
		return err
	}

	log.Printf("insert successfully : %v ", resp)
	return nil
}

func (e *lottoHistoryElasticSearch) GetLottoHistory(ctx context.Context, index string, uuid string, offsets int, recordPerpage int) ([]models.ElastictHitsHits, error) {

	// fistPeriod := "3008-05-02"
	fistPeriod := time.Date(3008, 05, 02, 0, 0, 0, 0, time.Local).Format(time.DateOnly)
	// currentPeriod := "3010-06-01"
	currentPeriod := time.Date(3010, 06, 01, 0, 0, 0, 0, time.Local).Format(time.DateOnly)
	boolQuery := types.NewBoolQuery()

	// v := "ร้าน*"
	mustQueries := []types.Query{
		newTermQuery("buyer_uuid", uuid),
		newTermsQuery("lotto_status", []types.FieldValue{"sold", "transfered", "suspend", "ban"}),
		newDateRangeQuery("lotto_price_due", fistPeriod, currentPeriod),
		// {
		// 	Wildcard: map[string]types.WildcardQuery{
		// 		"retailer_yr_merchant": {
		// 			Value: &v,
		// 		},
		// 	},
		// },
		// {
		// 	Range: map[string]types.RangeQuery{
		// 		"lotto_price_due": types.DateRangeQuery{
		// 			Gte: &fistPeriod,
		// 			Lt:  &currentPeriod,
		// 		},
		// 	},
		// },
	}

	lottoTypeSource := `
		def lottoType = doc['lotto_type.keyword'].value;
		if (lottoType == '01') return 1;
		else if (lottoType == '02') return 2;
		else if (lottoType == 'P80') return 3;
		else if (lottoType == 'M') return 4;
		else return 0;
	`

	sortQueries := []types.SortCombinations{
		newFieldSort("lotto_price_due", &sortorder.Desc),
		newFieldScriptSortLottoType(lottoTypeSource, &sortorder.Asc),
		newFieldSort("lotto_number.keyword", &sortorder.Asc),
		newFieldSort("lotto_round.keyword", &sortorder.Asc),
		newFieldSort("lotto_set.keyword", &sortorder.Asc),
	}

	boolQuery.Must = mustQueries

	q := types.NewQuery()
	q.Bool = boolQuery

	r := search.NewRequest()
	r.Query = q
	r.Sort = sortQueries

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(r); err != nil {
		log.Fatalf("Error encoding the raw query: %s", err)
		return nil, err
	}
	log.Printf("lotto_history query els : %+v\n", buf.String())

	resp, err := e.elasticClient.Search(
		e.elasticClient.Search.WithContext(ctx),
		e.elasticClient.Search.WithIndex(index),
		e.elasticClient.Search.WithBody(&buf),
		e.elasticClient.Search.WithTrackTotalHits(true),
		e.elasticClient.Search.WithFrom(offsets),
		e.elasticClient.Search.WithSize(recordPerpage),
	)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var result models.ElasticResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Hits.Hits, nil
}

func (e *lottoHistoryElasticSearch) UpdateLottoHistoryRewardStatusByElasticQuery(ctx context.Context, index string, uuid string, _id int64, source string) error {

	// query terms
	// termsQuery := types.NewTermsQuery()
	// termsQuery.TermsQuery = map[string]types.TermsQueryField{
	// 	"_id": []types.FieldValue{
	// 		_id,
	// 	},
	// }

	inlineScript := types.NewInlineScript()
	inlineScript.Lang = &scriptlanguage.Painless
	inlineScript.Source = source

	// script query
	// scriptQuery := types.NewScriptQuery()
	// scriptQuery.Script = inlineScript

	q := types.NewQuery()
	q.Terms = newTermsQueryRootLayer("_id", []types.FieldValue{_id})
	// q.Script = scriptQuery

	r := updatebyquery.NewRequest()
	r.Query = q
	r.Script = inlineScript
	// r.Script = types.InlineScript{
	// 	Lang:   &scriptlanguage.Painless,
	// 	Source: source,
	// }

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(r); err != nil {
		log.Fatalf("Error encoding the raw query: %s", err)
		return err
	}
	log.Printf("lotto_history query els : %+v\n", buf.String())

	resp, err := e.elasticClient.UpdateByQuery(
		[]string{index},
		e.elasticClient.UpdateByQuery.WithContext(ctx),
		e.elasticClient.UpdateByQuery.WithBody(&buf),
	)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	fmt.Printf("resp update_by_query elastic : %+v\n", resp.Body)

	return nil
}

func (e *lottoHistoryElasticSearch) UpdateLottoHistoryRewardStatus(ctx context.Context, index string, id string, doc map[string]interface{}) error {

	_doc, err := json.Marshal(doc)
	if err != nil {
		return err
	}

	r := update.NewRequest()
	r.Doc = _doc

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(r); err != nil {
		log.Fatalf("Error encoding the raw query: %s", err)
		return err
	}
	log.Printf("lotto_history update query els : %+v\n", buf.String())

	resp, err := e.elasticClient.Update(index, id, &buf, e.elasticClient.Update.WithContext(ctx))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	fmt.Printf("resp update_by_query elastic : %+v\n", resp.Body)

	var result models.ElasticResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	fmt.Printf("result : %+v\n", result)

	return nil
}

func (e *lottoHistoryElasticSearch) ExistsIndex(ctx context.Context, index string) (bool, error) {

	resp, err := e.elasticClient.Indices.Exists([]string{index},
		e.elasticClient.Indices.Exists.WithContext(ctx),
	)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	log.Println("[ExistsIndex] response_status_code :", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		return false, nil
	}

	return true, nil
}

func (e *lottoHistoryElasticSearch) CreateIndexWithMapping(ctx context.Context, index string, mapping string) error {

	resp, err := e.elasticClient.Indices.Create(index,
		e.elasticClient.Indices.Create.WithBody(strings.NewReader(mapping)),
		e.elasticClient.Indices.Create.WithContext(ctx),
		e.elasticClient.Indices.Create.WithPretty(),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	log.Println("[CreateIndexWithMapping] response_status_code :", resp.StatusCode)
	if resp.StatusCode != 200 {
		return fmt.Errorf("create index failed invalid status_code : %d", resp.StatusCode)
	}

	return nil
}

func (e *lottoHistoryElasticSearch) BulkData(ctx context.Context, index string) error {

	query := []string{
		`{ "index": { "_index": "%s", "_id" : "3" } }`,
		`{ "name" : "mond", "password": "pwd" }`,
		`{ "update": { "_index": "%s", "_id" : "3" } }`,
		`{ "doc": { "name" : "mond-update-9X9", "password": "pwd-update-9X9", "email": "watcharis_9X9@test.com" } }`,
	}

	var buf bytes.Buffer
	for _, v := range query {

		v = strings.ReplaceAll(v, "%s", index)
		// fmt.Println("query :", query)

		buf.Write([]byte(v))
		buf.WriteByte('\n')
	}

	fmt.Printf("query bulk : %+v\n", buf.String())

	resp, err := e.elasticClient.Bulk(bytes.NewReader(buf.Bytes()), e.elasticClient.Bulk.WithContext(ctx))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	var response interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		slog.ErrorContext(ctx, "cannot json decode response body", slog.Any("err", err))
		return err
	}

	slog.InfoContext(ctx, "[BulkData] success", slog.Any("response", response))

	return nil
}

func (e *lottoHistoryElasticSearch) BulkUpdateData(ctx context.Context, index string, ids []int, rewardStatus string) error {

	if rewardStatus == "" {
		rewardStatus = "start_test"
	}

	var buf buffer.Buffer
	for _, lottoID := range ids {
		// meta := []byte(fmt.Sprintf(`{ "update" : { "_index" : "%s","_id" : "%d" } }%s{ "doc" : { "reward_status" : "%s","update_at": "%s" } }%s`,
		// 	index, lottoID, "\n", "dojodjdjdodojdo", time.Now().Format(time.DateTime), "\n"),
		// )
		meta := []byte(fmt.Sprintf(`{ "update" : { "_index" : "%s","_id" : "%d" } }%s`, index, lottoID, "\n"))
		doc := []byte(fmt.Sprintf(`{ "doc" : { "reward_status" : "%s","update_at": "%s" } }%s`, rewardStatus, time.Now().Format(time.DateTime), "\n"))
		buf.Write(meta)
		buf.Write(doc)
	}

	slog.InfoContext(ctx, "query bulk update data", slog.String("raw_query", buf.String()))

	resp, err := e.elasticClient.Bulk(bytes.NewReader(buf.Bytes()), e.elasticClient.Bulk.WithContext(ctx))
	if err != nil {
		slog.ErrorContext(ctx, "cannot bulk update lotto_history", slog.Any("error", err))
		return err
	}

	defer resp.Body.Close()

	if resp.IsError() {
		err := errors.New("error bulk update response is not true")
		slog.ErrorContext(ctx, "failed error bulk update data in elastic", slog.Any("error", err))
		return err
	}

	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		slog.ErrorContext(ctx, "cannot json decode response body", slog.Any("error", err))
		return err
	}
	slog.InfoContext(ctx, "encode bulk response", slog.Any("response", response))

	isError := response["errors"].(bool)
	if isError {
		slog.ErrorContext(ctx, "update lottos_history to elasticsearch with bulk", slog.Any("response", response))
		return errors.New("update data lottos_history with bulk failed")
	}

	slog.Info("[BulkUpdateData] success", slog.Any("response", response))

	return nil

}

func (e *lottoHistoryElasticSearch) SearchLottoHistory(ctx context.Context, req models.LottoSearch) (models.ElasticResult, error) {
	var (
		elsIndex string
		rs       models.ElasticResult
	)
	elsIndex = fmt.Sprintf("%s*", models.INDEX)

	q := types.NewQuery()
	q.Terms = &types.TermsQuery{
		TermsQuery: map[string]types.TermsQueryField{
			"_id": req.ID,
		},
	}

	r := search.NewRequest()
	r.Query = q
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(r); err != nil {
		return rs, err
	}

	resp, err := e.elasticClient.Search(
		e.elasticClient.Search.WithContext(ctx),
		e.elasticClient.Search.WithIndex(elsIndex),
		e.elasticClient.Search.WithBody(&buf),
	)
	if err != nil {
		return rs, err
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&rs); err != nil {
		return rs, err
	}

	return rs, nil
}

func (e *lottoHistoryElasticSearch) GetPITtoken(ctx context.Context, index string, keepAlive string) (string, error) {

	url := fmt.Sprintf("http://localhost:9200/%s/_pit?keep_alive=1m", index)

	client := rest.CreateHttpClient()

	// Create a new POST request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	slog.InfoContext(ctx, "response body", slog.Int("code", resp.StatusCode), slog.String("response", string(body)))

	var response map[string]interface{}
	if err := json.Unmarshal(body, &response); err != nil {
		return "", err
	}

	return response["id"].(string), nil
}

func (e *lottoHistoryElasticSearch) PointInTimeSearch(ctx context.Context, pitToken string, searchAfter []int64) ([]models.ElastictHitsHits, []int64, error) {

	var queryBody bytes.Buffer
	body := map[string]interface{}{
		"pit": map[string]interface{}{
			"id": pitToken,
		},
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
		// "_source":      false,
		"search_after": searchAfter,
		"sort": []map[string]interface{}{
			{
				"id": map[string]interface{}{
					"order": "asc",
				},
			},
		},
	}

	if err := json.NewEncoder(&queryBody).Encode(body); err != nil {
		return nil, nil, err
	}
	fmt.Printf("query body : %+v\n", queryBody.String())

	resp, err := e.elasticClient.Search(
		e.elasticClient.Search.WithContext(ctx),
		e.elasticClient.Search.WithBody(bytes.NewBufferString(queryBody.String())),
		e.elasticClient.Search.WithSize(models.ELSATIC_SIZE_1000),
		e.elasticClient.Search.WithPretty(),
	)
	if err != nil {
		return nil, nil, err
	}

	defer resp.Body.Close()

	var result models.ElasticResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, nil, err
	}

	if len(result.Hits.Hits) == 0 {
		return []models.ElastictHitsHits{}, []int64{}, nil
	}

	next := result.Hits.Hits[len(result.Hits.Hits)-1].Sort

	return result.Hits.Hits, next, nil
}

func (e *lottoHistoryElasticSearch) DeletePITtoken(ctx context.Context, pitToken string) error {

	body := map[string]interface{}{
		"id": pitToken,
	}

	var bodyBytes bytes.Buffer
	if err := json.NewEncoder(&bodyBytes).Encode(body); err != nil {
		return err
	}

	fmt.Printf("bodyBytes : %+v\n", bodyBytes.String())

	resp, err := e.elasticClient.ClosePointInTime(
		e.elasticClient.ClosePointInTime.WithContext(ctx),
		e.elasticClient.ClosePointInTime.WithBody(&bodyBytes),
		e.elasticClient.ClosePointInTime.WithPretty())

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("delete pit_token invalid response status_code : %d", resp.StatusCode)
	}

	var result models.ElasticResultPIT
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if !result.Succeeded {
		return fmt.Errorf("invalid result status for delete pit_token : %t", result.Succeeded)
	}

	return nil
}
