package services

import "context"

type LottoHistoryServices interface {
	MigrateDataLottoHistoryToElastic(ctx context.Context) error
	GetLottoHistory(ctx context.Context) error
	TestDBUseTransaction(ctx context.Context) error
	MigrateSpeacificLotto(ctx context.Context) error
	UpdateRewardStatusByElasticQuery(ctx context.Context) error
	UpdateRewardStatus(ctx context.Context) error
	MigrateSoldLottos(ctx context.Context) error
	CreateIndexIfNotExists(ctx context.Context) error
	PocBulkData(ctx context.Context) error
	GetMultipleLottoHistoryAndMigrateToElasticsearch(ctx context.Context) error
	PocElasticsearchPIT(ctx context.Context) error
}
