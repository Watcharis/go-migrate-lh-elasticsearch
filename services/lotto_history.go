package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
	"strconv"
	"strings"
	"time"
	"watcharis/go-migrate-lotto-history-els/models"
	"watcharis/go-migrate-lotto-history-els/repository/cache"
	"watcharis/go-migrate-lotto-history-els/repository/db"
	lottoHistoryELS "watcharis/go-migrate-lotto-history-els/repository/elasticsearch"

	"watcharis/go-migrate-lotto-history-els/util/errorr"
)

type lottoHistoryServices struct {
	slogger                   *slog.Logger
	cache                     cache.Cache
	lottoHistoryRepository    db.LottoHistoryRepository
	lottoHistoryElasticSearch lottoHistoryELS.LottoHistoryElasticSearch
	timetableRepository       db.TimeTableRepository
	lottoRepository           db.LottoRepository
	rewardRepository          db.RewardRepository
}

func NewLottoHistoryServices(
	slogger *slog.Logger,
	cache cache.Cache,
	lottoHistoryRepository db.LottoHistoryRepository,
	lottoHistoryElasticSearch lottoHistoryELS.LottoHistoryElasticSearch,
	timetableRepository db.TimeTableRepository,
	lottoRepository db.LottoRepository,
	rewardRepository db.RewardRepository,
) LottoHistoryServices {
	return &lottoHistoryServices{
		slogger:                   slogger,
		cache:                     cache,
		lottoHistoryRepository:    lottoHistoryRepository,
		lottoHistoryElasticSearch: lottoHistoryElasticSearch,
		timetableRepository:       timetableRepository,
		lottoRepository:           lottoRepository,
		rewardRepository:          rewardRepository,
	}
}

func (s *lottoHistoryServices) MigrateDataLottoHistoryToElastic(ctx context.Context) error {

	currentPeriod, err := s.cache.Get(ctx, models.REDISKEY_CURRENTPERIOD)
	if err != nil {
		log.Printf("cannot get current_period in redis [error]: %+v\n", err)
		return err
	}
	fmt.Println("currentPeriod :", currentPeriod)

	currentPeriodTime, err := time.Parse(models.DATETIME_FORMAT_PRICE_DUE_IN_REDIS, currentPeriod)
	if err != nil {
		log.Printf("cannot time.parse string to time.time [error]: %+v\n", err)
		return err
	}

	firstPeriodStr := currentPeriodTime.AddDate(-2, -1, 0).Format(models.DATETIME_FORMAT_PRICE_DUE_IN_DB)

	currentPeriodTimeStr := currentPeriodTime.Format(models.DATETIME_FORMAT_PRICE_DUE_IN_DB)

	lottoPriceDueList, err := s.timetableRepository.GetTimetable(ctx, firstPeriodStr, currentPeriodTimeStr)
	if err != nil {
		log.Printf("cannot query round_date from timetable [error] : %+v", err)
		return err
	}

	lottoMapTimeIndex := make(map[string]string, 0)
	lottoRoundDB := []string{}
	for _, lpd := range lottoPriceDueList {

		ltRoundDB := lpd.RoundDate.Format(models.DATETIME_FORMAT_PRICE_DUE_IN_DB)

		if _, ok := lottoMapTimeIndex[ltRoundDB]; !ok {
			lottoMapTimeIndex[ltRoundDB] = lpd.RoundDate.Format(models.DATETIME_FORMAT_PRICE_DUE_NO_DAT)
		}

		lottoRoundDB = append(lottoRoundDB, ltRoundDB)
	}

	for _, currentPreiod := range lottoRoundDB {

		lottoHistory, err := s.lottoHistoryRepository.GetLottoHistoryByLottoPriceDue(ctx, currentPreiod)
		if err != nil {
			log.Printf("cannot query lotto_history [error]: %+v", err)
			return err
		}
		fmt.Printf("lottoHistory lenght : %+v\n", len(lottoHistory))

		if len(lottoHistory) > 5000 {
			lottoHistory = lottoHistory[0:3000]
		}

		for _, lh := range lottoHistory {
			if lh.WinDesc == "" {
				lh.WinDesc = `{"description": null, "reward":null ,"is_first_reward": null}`
			}
			// fmt.Printf("lottoHistory : %+v\n", lh)

			lottoPriceDue := lh.LottoPriceDue.Format(models.DATETIME_FORMAT_PRICE_DUE)

			lottoHistoryEls := models.LottoHistoryElasticSearchOld{
				ElasticId:     strconv.Itoa(lh.ID),
				ID:            lh.ID,
				LottoNumber:   lh.LottoNumber,
				LottoRound:    lh.LottoRound,
				LottoSet:      lh.LottoSet,
				LottoYear:     lh.LottoYear,
				LottoItem:     lh.LottoItem,
				LottoPriceDue: lottoPriceDue,
				LottoPrice:    lh.LottoPrice,
				ScannerUUID:   lh.ScannerUUID,
				BuyerUUID:     lh.BuyerUUID,
				LottoURL:      lh.LottoURL,
				LottoBcRef:    lh.LottoBcRef,
				LottoBcHash:   lh.LottoBcHash,
				LottoType:     lh.LottoType,
				LottoStatus:   lh.LottoStatus,
				CreateAt:      time.Now().Format(models.DATETIME_FORMAT_PRICE_DUE_IN_DB),
				CreateBy:      "",
				UpdateAt:      time.Now().Format(models.DATETIME_FORMAT_PRICE_DUE_IN_DB),
				UpdateBy:      "",
				PaymentID:     lh.PaymentID,
				LottoUUID:     lh.LottoUUID,
				Tags:          lh.Tags,
				IsWin:         lh.IsWin,
				WinDesc:       []models.WinDescObject{},
			}
			fmt.Printf("lottoHistoryEls : %+v\n", lottoHistoryEls)

			lottoHistoryIndex := fmt.Sprintf(models.INDEX, lottoMapTimeIndex[lh.LottoPriceDue.Format(models.DATETIME_FORMAT_PRICE_DUE_IN_DB)])
			fmt.Printf("lottoHistoryIndex : %+v\n", lottoHistoryIndex)

			// if err := s.lottoHistoryElasticSearch.InsertLottoHistory(ctx, lottoHistoryIndex, lottoHistoryEls); err != nil {
			// 	log.Printf("cannot insert lotto_history els [error]: %+v\n", err)
			// 	return err
			// }
		}
	}

	return nil
}

func (s *lottoHistoryServices) GetLottoHistory(ctx context.Context) error {

	// lottoHistoryIndex := fmt.Sprintf(models.INDEX, "30100503")
	lottoHistoryIndex := "lotto_history.*"
	fmt.Println("lotto_history_index :", lottoHistoryIndex)

	// buyerUuiid := "97a136db1d2bf20a6588fbf5702e48c292c5bd517e1b7da840d9c76915fd39f1"
	buyerUuiid := "3fbf6121a3dfb7e2dac7076c13250b57eae306ab3b6f485a010e23eed5436fc5"

	// requestPage :=

	recordPerpage := 50
	pageNo := 1

	dupLottoPriceDue := map[string]bool{}
	listLottoPriceDue := []string{}
	lottoHistoryELS := []models.LottosHistoryElasticsearch{}

	for {
		if pageNo < 1 {
			break
		}

		offsets := recordPerpage * (pageNo - 1)
		fmt.Println("pageNo :", pageNo)
		fmt.Println("offsets :", offsets)

		lottoHistory, err := s.lottoHistoryElasticSearch.GetLottoHistory(ctx, lottoHistoryIndex, buyerUuiid, offsets, recordPerpage)
		if err != nil {
			log.Printf("get lotto_history failed [error]: %+v\n", err)
			break
		}

		fmt.Println("lotto_history lenght : ")
		if len(lottoHistory) == 0 {
			if pageNo > 1 {
				log.Printf("found all lotto_history")
				break
			} else {
				log.Printf("lotto_history not found")
				break
			}
		}

		lottoPriceDueTime, err := time.ParseInLocation(time.DateTime, lottoHistory[0].Source.LottoPriceDue, time.Local)
		if err != nil {
			log.Printf("cannot convert lotto_price_due to time.time [error] : %+v", err)
			return err
		}

		lastLottoPriceDue := lottoPriceDueTime.Format(time.DateOnly)

		if _, ok := dupLottoPriceDue[lastLottoPriceDue]; !ok {
			dupLottoPriceDue[lastLottoPriceDue] = true
			listLottoPriceDue = append(listLottoPriceDue, lastLottoPriceDue)
		}

		// if pageNo == requestPage {
		// for _, lth := range lottoHistory {
		// 	source := lth.Source
		// 	fmt.Printf("lottoHistoryELS : %+v\n", source)
		// 	lottoPriceDue := source.LottoPriceDue.Format(time.DateTime)
		// 	fmt.Println("lottoPriceDue :", lottoPriceDue)

		// 	// lottoHistoryELS = append(lottoHistoryELS, lth.Source)
		// }

		for _, lh := range lottoHistory {
			// if index > 10 {
			// 	break
			// }
			sizeByte, err := json.Marshal(lh.Source)
			if err != nil {
				return err
			}
			fmt.Printf("lotto_history size : %d bytes\n", len(sizeByte))

			// switch len(sizeByte) {
			// case 850:
			// 	fmt.Printf("lh 850 : %+v\n", string(sizeByte))
			// case 950:
			// 	fmt.Printf("lh 950 : %+v\n", string(sizeByte))
			// case 1082:
			// 	fmt.Printf("lh 1082 : %+v\n", string(sizeByte))
			// default:
			// 	fmt.Printf("lh default : %+v\n", string(sizeByte))
			// }
		}
		// }

		pageNo += 1
	}
	// fmt.Println("listLottoPriceDue :", listLottoPriceDue)
	fmt.Printf("lottoHistoryELS lenght : %+v\n", len(lottoHistoryELS))
	fmt.Printf("lottoHistoryELS : %+v\n", lottoHistoryELS)
	fmt.Println("listLottoPriceDue :", listLottoPriceDue)
	return nil
}

func (s *lottoHistoryServices) TestDBUseTransaction(ctx context.Context) error {

	uuid := "4e43640c08a809f52dd95b55d4f2304d154421185268bf4b6a5dd7c381fee291"

	tx := s.lottoHistoryRepository.Begin(ctx)

	lhTx, err := s.lottoHistoryRepository.GetLottoHistoryTransaction(ctx, tx, uuid)
	if err != nil {
		log.Printf("cannot get lotto_history_tx [error]: %+v\n", err)
		tx.Rollback()
		return err
	}
	fmt.Printf("lotto_history_tx : %+v\n", len(lhTx))

	if err := s.lottoHistoryRepository.Commit(ctx, tx); err != nil {
		log.Printf("cannot commit transaction lotto_history_tx [error]: %+v\n", err)
		return err
	}

	lhWithoutTx, err := s.lottoHistoryRepository.GetLottoHistoryWithoutTransaction(ctx, uuid)
	if err != nil {
		log.Printf("cannot get lotto_history_tx [error]: %+v\n", err)
		return err
	}
	fmt.Printf("lhWithoutTx : %+v\n", len(lhWithoutTx))

	return nil
}

func (s *lottoHistoryServices) MigrateSpeacificLotto(ctx context.Context) error {

	uuid := "4e43640c08a809f52dd95b55d4f2304d154421185268bf4b6a5dd7c381fee291"

	lottos, err := s.lottoRepository.GetLottoWithReward(ctx, uuid)
	if err != nil {
		log.Printf("cannot commit transaction lotto_history_tx [error]: %+v\n", err)
		return err
	}
	// fmt.Printf("lottos : %+v\n", lottos)

	lottoPriceDueIndex := ""
	for i, lt := range lottos {
		fmt.Println("lt :", lt.Tags)

		if i == 0 {
			lottoPriceDueIndex = lt.LottoPriceDue.Format(time.DateOnly)
		}

		if i%5 == 2 {
			lt.Tags = "[1]"
		}
		if i == len(lottos)-2 {
			lt.Tags = "[1, 2]"
		}

		var dt1, dt2, dt3, dt4, dt5, dt6 int
		if lt.LottoNumber != "" {
			lottoNumber := strings.Split(lt.LottoNumber, "")
			dt1, _ = strconv.Atoi(lottoNumber[0])
			dt2, _ = strconv.Atoi(lottoNumber[1])
			dt3, _ = strconv.Atoi(lottoNumber[2])
			dt4, _ = strconv.Atoi(lottoNumber[3])
			dt5, _ = strconv.Atoi(lottoNumber[4])
			dt6, _ = strconv.Atoi(lottoNumber[5])
		}

		tags := []int{}

		if !IsEmptyString(lt.Tags) && lt.Tags != "[]" {
			tagStr := lt.Tags
			fmt.Println("tagStr :", tagStr)
			tagStr = strings.TrimPrefix(tagStr, "[")
			tagStr = strings.TrimSuffix(tagStr, "]")
			tagStr = strings.ReplaceAll(tagStr, " ", "")

			tagList := strings.Split(tagStr, ",")

			for _, tag := range tagList {
				tagInt, err := strconv.Atoi(tag)
				if err != nil {
					log.Printf("cannot convert tag string to int [error]: %+v\n", err)
					return err
				}

				tags = append(tags, tagInt)
			}
		}

		var winDesc []models.WinDesc
		if lt.WinDesc != "" {
			winDescByte, err := json.Marshal(lt.WinDesc)
			if err != nil {
				log.Printf("cannot json.Mashal lotto.win_desc [error] : %+v\n", err)
				return err
			}

			if err := json.Unmarshal(winDescByte, &winDesc); err != nil {
				log.Printf("cannot json.Unmarshal lotto.win_desc to array [error] : %+v\n", err)
				return err
			}
		} else {
			winDesc = []models.WinDesc{}
		}

		lottoHistory := models.LottosHistoryElasticsearch{
			ID:            int64(lt.ID),
			LottoNumber:   lt.LottoNumber,
			LottoDt1:      dt1,
			LottoDt2:      dt2,
			LottoDt3:      dt3,
			LottoDt4:      dt4,
			LottoDt5:      dt5,
			LottoDt6:      dt6,
			LottoRound:    lt.LottoRound,
			LottoSet:      lt.LottoSet,
			LottoYear:     lt.LottoYear,
			LottoItem:     lt.LottoItem,
			LottoPriceDue: lt.LottoPriceDue.Format(time.DateOnly),
			LottoPrice:    lt.LottoPrice,
			ScannerUUID:   lt.ScannerUUID,
			BuyerUUID:     lt.BuyerUUID,
			LottoURL:      lt.LottoURL,
			LottoBCRef:    lt.LottoBCRef,
			LottoBCHash:   lt.LottoBCHash,
			LottoType:     lt.LottoType,
			LottoStatus:   lt.LottoStatus,
			CreateAt:      lt.CreateAt.Format(time.DateTime),
			CreateBy:      lt.CreateBy,
			UpdateAt:      lt.UpdateAt.Format(time.DateTime),
			UpdateBy:      lt.UpdateBy,
			PaymentID:     lt.PaymentID,
			LottoUUID:     lt.LottoUUID,
			Tags:          tags,
			IsWin:         lt.IsWin,
			WinDesc:       winDesc,
			RewardStatus:  lt.RewardStatus,
		}
		fmt.Printf("lottoHistory : %+v\n", lottoHistory)

		if lt.LottoPriceDue.Format(time.DateOnly) != lottoPriceDueIndex {
			fmt.Println("diff lotto_price_due :", lt.LottoPriceDue.Format(time.DateOnly))
			lottoPriceDueIndex = lt.LottoPriceDue.Format(time.DateOnly)
		}

		index := fmt.Sprintf(models.INDEX_POV, lottoPriceDueIndex)
		fmt.Println("index :", index)

		if err := s.lottoHistoryElasticSearch.InsertLottoHistory(ctx, index, lottoHistory); err != nil {
			log.Printf("cannot commit transaction lotto_history_tx [error]: %+v\n", err)
			return err
		}

	}

	return nil
}

func (s *lottoHistoryServices) MigrateSoldLottos(ctx context.Context) error {

	uuid := "4e43640c08a809f52dd95b55d4f2304d154421185268bf4b6a5dd7c381fee291"
	// lottoPriceDue := "3010-08-16 00:00:00"
	lottoPriceDue := time.Date(3010, 8, 16, 0, 0, 0, 0, time.Local)

	lottos, err := s.lottoRepository.GetLottoSoldWithRewardSpeacificRoundDate(ctx, uuid, lottoPriceDue.Format(time.DateTime))
	if err != nil {
		log.Printf("cannot get lotto_with_reward [error]: %+v\n", err)
		return err
	}
	fmt.Printf("lotto : %+v\n", len(lottos))

	for _, lt := range lottos {

		fmt.Printf("lt : %+v\n", lt)

		var dt1, dt2, dt3, dt4, dt5, dt6 int
		if lt.LottoNumber != "" {
			lottoNumber := strings.Split(lt.LottoNumber, "")
			dt1, _ = strconv.Atoi(lottoNumber[0])
			dt2, _ = strconv.Atoi(lottoNumber[1])
			dt3, _ = strconv.Atoi(lottoNumber[2])
			dt4, _ = strconv.Atoi(lottoNumber[3])
			dt5, _ = strconv.Atoi(lottoNumber[4])
			dt6, _ = strconv.Atoi(lottoNumber[5])
		}

		tags := []int{}
		if !IsEmptyString(lt.Tags) {
			if err := json.Unmarshal([]byte(lt.Tags), &tags); err != nil {
				log.Printf("cannot json.Unmarshal lotto.win_desc to array [error] : %+v\n", err)
				return err
			}
		}

		var winDesc []models.WinDesc
		if lt.WinDesc != "" {
			fmt.Println("win_desc not empty")
			if err := json.Unmarshal([]byte(lt.WinDesc), &winDesc); err != nil {
				log.Printf("cannot json.Unmarshal lotto.win_desc to array [error] : %+v\n", err)
				return err
			}
		} else {
			winDesc = []models.WinDesc{}
		}

		lottoHistory := models.LottosHistoryElasticsearch{
			ID:            int64(lt.ID),
			LottoNumber:   lt.LottoNumber,
			LottoDt1:      dt1,
			LottoDt2:      dt2,
			LottoDt3:      dt3,
			LottoDt4:      dt4,
			LottoDt5:      dt5,
			LottoDt6:      dt6,
			LottoRound:    lt.LottoRound,
			LottoSet:      lt.LottoSet,
			LottoYear:     lt.LottoYear,
			LottoItem:     lt.LottoItem,
			LottoPriceDue: lt.LottoPriceDue.Format(time.DateOnly),
			LottoPrice:    lt.LottoPrice,
			ScannerUUID:   lt.ScannerUUID,
			BuyerUUID:     lt.BuyerUUID,
			LottoURL:      lt.LottoURL,
			LottoBCRef:    lt.LottoBCRef,
			LottoBCHash:   lt.LottoBCHash,
			LottoType:     lt.LottoType,
			LottoStatus:   lt.LottoStatus,
			CreateAt:      lt.CreateAt.Format(time.DateTime),
			CreateBy:      lt.CreateBy,
			UpdateAt:      lt.UpdateAt.Format(time.DateTime),
			UpdateBy:      lt.UpdateBy,
			PaymentID:     lt.PaymentID,
			LottoUUID:     lt.LottoUUID,
			Tags:          tags,
			IsWin:         lt.IsWin,
			WinDesc:       winDesc,
			RewardStatus:  lt.RewardStatus,
		}
		fmt.Printf("lottoHistory : %+v\n", lottoHistory)

		index := fmt.Sprintf(models.INDEX_POV, lottoPriceDue.Format(time.DateOnly))
		fmt.Println("index :", index)

		if err := s.lottoHistoryElasticSearch.InsertLottoHistory(ctx, index, lottoHistory); err != nil {
			log.Printf("cannot commit transaction lotto_history_tx [error]: %+v\n", err)
			return err
		}

	}

	return nil
}

// service : create index with mapping for query elastic
func (s *lottoHistoryServices) CreateIndexIfNotExists(ctx context.Context) error {

	currentDate := time.Now().Format(time.DateOnly)

	index := fmt.Sprintf(models.INDEX_POV, currentDate)

	exists, err := s.lottoHistoryElasticSearch.ExistsIndex(ctx, index)
	if (err != nil && !exists) || (err == nil && !exists) {
		// log.Printf("Notfound index [EXISTS_STATUS] is %v [ERROR] : %+v\n", exists, err)
		s.slogger.ErrorContext(ctx, "Notfound index in elasticsearch", slog.Bool("exists", exists), slog.Any("error", err))

		lottoHistoryIndexMapping := LottoHistoryElasticSearchMapping()

		log.Println("start create index with mapping .....")
		if err := s.lottoHistoryElasticSearch.CreateIndexWithMapping(ctx, index, lottoHistoryIndexMapping); err != nil {
			// log.Printf("cannot create index with mapping [error]: %+v\n", err)
			s.slogger.ErrorContext(ctx, "cannot create index with mapping", slog.Any("error", err))
			return err
		}

		// log.Printf("create index in elasticsearch success index_name : %+v\n", index)
		s.slogger.InfoContext(ctx, "create index in elasticsearch success", slog.String("index", index))
		return nil
	}

	// log.Printf("elastic Index is exists : %+v\n", exists)
	s.slogger.InfoContext(ctx, "elastic Index is exists", slog.Bool("exists", exists))
	return nil
}

func (s *lottoHistoryServices) GetMultipleLottoHistoryAndMigrateToElasticsearch(ctx context.Context) error {

	errorLottoHistory := new(errorr.Error)
	lottoPriceDue := "3011-02-01"

	if lottoPriceDue == "" {
		err := errors.New("lotto_price_due is empty")
		slog.ErrorContext(ctx, "lotto_price_due is empty", slog.Any("error", err))
		return errorLottoHistory.Error(err)
	}

	count, err := s.lottoHistoryRepository.CountLottoHistoryByLottoPriceDue(ctx, lottoPriceDue)
	if err != nil {
		slog.ErrorContext(ctx, "cannot count lotto_history", slog.Any("error", err))
		return errorLottoHistory.Error(err)
	}
	slog.InfoContext(ctx, "count all lotto_history in round date success", slog.String("lotto_price_due", lottoPriceDue), slog.Int("count", count))

	limitSize := models.LIMIT_LOTTO_HISTORY_QUERY_DB
	n := 0
	for n < count {
		offset := n
		lottoHistory, err := s.lottoHistoryRepository.GetLottoHistoryWithLimitAndOffset(ctx, lottoPriceDue, limitSize, offset)
		if err != nil {
			slog.ErrorContext(ctx, "cannot get lotto_history by limit/offset", slog.Any("error", err))
			return nil
		}

		LottoIDs := make([]int, 0)
		for _, lh := range lottoHistory {
			lottoID := lh.ID
			LottoIDs = append(LottoIDs, lottoID)
		}

		rewards, err := s.rewardRepository.GetRewardByListLottoID(ctx, LottoIDs)
		if err != nil {
			slog.ErrorContext(ctx, "cannot get list reward", slog.Any("error", err))
			return errorLottoHistory.Error(err)
		}

		mapRewards := make(map[int]models.Rewards, 0)
		for _, rw := range rewards {
			lottoID := rw.LottoID
			mapRewards[lottoID] = rw
		}

		sequent := 0
		lastSequent := len(lottoHistory) - 1
		for {

			if sequent > lastSequent {
				break
			}

			lottoID := LottoIDs[sequent]
			lh := lottoHistory[sequent]
			reward := mapRewards[lottoID]

			if lottoID != lh.ID {
				slog.InfoContext(ctx, "map lotto failed different lotto_id", slog.Int("lotto_1_id", lottoID), slog.Int("lotto_2_id", lh.ID))
				break
			}

			var dt1, dt2, dt3, dt4, dt5, dt6 int
			if lh.LottoNumber != "" {
				lottoNumber := strings.Split(lh.LottoNumber, "")
				dt1, _ = strconv.Atoi(lottoNumber[0])
				dt2, _ = strconv.Atoi(lottoNumber[1])
				dt3, _ = strconv.Atoi(lottoNumber[2])
				dt4, _ = strconv.Atoi(lottoNumber[3])
				dt5, _ = strconv.Atoi(lottoNumber[4])
				dt6, _ = strconv.Atoi(lottoNumber[5])
			}

			tags := []int{}
			if !IsEmptyString(lh.Tags) {
				if err := json.Unmarshal([]byte(lh.Tags), &tags); err != nil {
					slog.ErrorContext(ctx, "cannot json.Unmarshal lh.tags to array", slog.Any("error", err))
					return errorLottoHistory.Error(err)
				}
			}

			var winDesc []models.WinDesc
			if reward.WinDesc != "" {
				if err := json.Unmarshal([]byte(reward.WinDesc), &winDesc); err != nil {
					slog.ErrorContext(ctx, "cannot json.Unmarshal lh.win_desc to array", slog.Any("error", err))
					return errorLottoHistory.Error(err)
				}
			} else {
				winDesc = []models.WinDesc{}
			}

			lottoHistoryElasticsearch := models.LottosHistoryElasticsearch{
				ID:            int64(lh.ID),
				LottoNumber:   lh.LottoNumber,
				LottoDt1:      dt1,
				LottoDt2:      dt2,
				LottoDt3:      dt3,
				LottoDt4:      dt4,
				LottoDt5:      dt5,
				LottoDt6:      dt6,
				LottoRound:    lh.LottoRound,
				LottoSet:      lh.LottoSet,
				LottoYear:     lh.LottoYear,
				LottoItem:     lh.LottoItem,
				LottoPriceDue: lh.LottoPriceDue.Format(time.DateOnly),
				LottoPrice:    lh.LottoPrice,
				ScannerUUID:   lh.ScannerUUID,
				BuyerUUID:     lh.BuyerUUID,
				LottoURL:      lh.LottoURL,
				LottoBCRef:    "xxx",
				LottoBCHash:   "xxx",
				LottoType:     lh.LottoType,
				LottoStatus:   lh.LottoStatus,
				CreateAt:      lh.CreateAt.Format(time.DateTime),
				CreateBy:      lh.CreateBy,
				UpdateAt:      lh.UpdateAt.Format(time.DateTime),
				UpdateBy:      lh.UpdateBy,
				PaymentID:     int64(lh.PaymentID),
				LottoUUID:     lh.LottoUUID,
				Tags:          tags,
				IsWin:         reward.IsWin,
				WinDesc:       winDesc,
				RewardStatus:  reward.RewardStatus,
			}

			lottoHistoryElasticsearchIndex := fmt.Sprintf(models.INDEX_POV, lottoPriceDue)

			if err := s.lottoHistoryElasticSearch.InsertLottoHistory(ctx, lottoHistoryElasticsearchIndex, lottoHistoryElasticsearch); err != nil {
				slog.ErrorContext(ctx, "cannot insert lotto_history to elastic search",
					slog.String("index", lottoHistoryElasticsearchIndex),
					slog.Any("doc", lottoHistoryElasticsearch),
					slog.Any("error", err))
				return errorLottoHistory.Error(err)
			}

			sequent++
		}

		n += limitSize
	}

	return nil
}

func (s *lottoHistoryServices) GetLottoHistoryBlABlA(ctx context.Context) error {
	return nil
}

func (s *lottoHistoryServices) PocElasticsearchPIT(ctx context.Context) error {

	startTime := time.Now()
	var (
		index = "lotto_history_pov.3011-02-01"
	)

	pitToken, err := s.lottoHistoryElasticSearch.GetPITtoken(ctx, index, "5m")
	if err != nil {
		slog.ErrorContext(ctx, "cannot get pit_id from elasticsearch", slog.Any("error", err))
		return err
	}
	slog.InfoContext(ctx, "get point_in_time_token success")

	i := 0
	var searchAfter []int64
	for {
		if i == 0 {
			searchAfter = []int64{0, 0}
		}

		resultPIT, sort, err := s.lottoHistoryElasticSearch.PointInTimeSearch(ctx, pitToken, searchAfter)
		if err != nil {
			slog.ErrorContext(ctx, "cannot use point_in_time", slog.Any("error", err))
			return err
		}

		if len(resultPIT) == 0 {
			slog.InfoContext(ctx, "end of process point_in_time")
			break
		}

		lottoIDs := []int{}
		for _, lh := range resultPIT {
			ltID := int(lh.Source.ID)
			lottoIDs = append(lottoIDs, ltID)
		}

		if err := s.lottoHistoryElasticSearch.BulkUpdateData(ctx, index, lottoIDs, "success"); err != nil {
			slog.ErrorContext(ctx, "cannot bulk update data to elastic", slog.Any("error", err))
			return err
		}

		searchAfter = sort
		i++
	}

	if err := s.lottoHistoryElasticSearch.DeletePITtoken(ctx, pitToken); err != nil {
		slog.ErrorContext(ctx, "cannot delet pit_token", slog.Any("error", err))
		return err
	}
	slog.InfoContext(ctx, "process del pit_token success")

	endtime := time.Since(startTime)

	slog.InfoContext(ctx, "usage time for finish process",
		slog.Any("start_time", startTime),
		slog.Any("end_time", startTime.Add(endtime)),
		slog.Any("usage", endtime.Seconds()),
	)

	return nil
}
