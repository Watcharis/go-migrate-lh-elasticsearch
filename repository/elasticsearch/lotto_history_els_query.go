package elasticsearch

import (
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/scriptlanguage"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/scriptsorttype"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
)

const GetLottoHistoryRawQuery = `{
	"query": {
		"bool": {
			"must": [
			{
				"term": {
					"buyer_uuid": {
						"value": %s
					}
				}
			},
			{
				"terms": {
					"lotto_status": ["sold","transfered","suspend","ban"]
				}
			}
			]
		}
	},
	"sort": [
		{
			"lotto_price_due": {
				"order": "desc"
			}
		},
		{
			"lotto_number.keyword": {
				"order": "asc"
			}
		},
		{
			"lotto_round.keyword": {
				"order": "asc"
			}
		},
		{
			"lotto_set.keyword": {
				"order": "asc"
			}
		}
	]
}`

func newTermQuery(k string, v types.FieldValue) types.Query {
	q := types.NewQuery()
	q.Term[k] = types.TermQuery{Value: v}

	return *q
}

func newMatchQuery(k string, v string) types.Query {
	q := types.NewQuery()
	q.Match[k] = types.MatchQuery{Query: v}

	return *q
}

func newFieldSort(k string, o *sortorder.SortOrder) types.SortOptions {
	so := types.NewSortOptions()
	so.SortOptions[k] = types.FieldSort{Order: o}

	return *so
}

func newTermsQuery(k string, v []types.FieldValue) types.Query {
	q := types.NewQuery()
	q.Terms = &types.TermsQuery{TermsQuery: map[string]types.TermsQueryField{k: v}}

	return *q
}

func newDateRangeQuery(k string, firstPeriod, secondPeriod string) types.Query {
	q := types.NewQuery()
	q.Range = map[string]types.RangeQuery{
		k: types.DateRangeQuery{
			Gte: &firstPeriod,
			Lt:  &secondPeriod,
		},
	}
	return *q
}

func newFieldScriptSortLottoType(source string, o *sortorder.SortOrder) types.SortOptions {
	so := types.NewSortOptions()

	inlineScript := types.NewInlineScript()
	inlineScript.Lang = &scriptlanguage.Painless
	inlineScript.Source = source

	_script := types.NewScriptSort()
	_script.Type = &scriptsorttype.Number
	_script.Script = inlineScript
	_script.Order = o

	so.Script_ = _script

	return *so
}

func newFieldWirldCard(k string, v string) types.Query {
	q := types.NewQuery()
	q.Wildcard = map[string]types.WildcardQuery{
		k: {Value: &v},
	}
	return *q
}

func newTermsQueryRootLayer(k string, v []types.FieldValue) *types.TermsQuery {
	terms := types.NewTermsQuery()
	terms.TermsQuery = map[string]types.TermsQueryField{
		k: v,
	}

	return terms
}
