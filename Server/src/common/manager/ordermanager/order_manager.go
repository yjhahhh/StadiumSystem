package ordermanager

import (
	// "encoding/json"
	"fmt"
	"time"

	"common/connection"
	"common/log"
	"common/model/stadium"
	"common/redis"
	"common/model/game"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// 时段表
var timeTable []string

func InitTimetable() {
	timeTable = []string{"07:00", "07:30", "08:00", "08:30", "09:00", "09:30", "10:00", "10:30", "11:00", "11:30", "12:00", "12:30", "13:00", "13:30"}
}

const (
	IsBoked   = "true"
	IsIdle    = "false"
	OrderFail = "order fail"
	layout    = "2006-01-02 15:04"
)

// 起始时间返回时间段表
func GetTimes(start, end string) []string {
	flag := false
	ret := make([]string, 0)
	for i := range timeTable {
		if timeTable[i] == start {
			flag = true
		}
		if timeTable[i] == end {
			break
		}
		if flag {
			ret = append(ret, timeTable[i])
		}
	}
	return ret
}

// 返回全部时段
func GetAllTimes(id uint, date string) stadium.BookedTime {
	ret := make(stadium.BookedTime)
	for i := range timeTable {
		ret[timeTable[i]] = IsIdle
	}
	booked, err := getStadiumBooked(id, date)
	log.Debugf("GetAllTimes id = %d, date = %s, booked : %+v\n", id, date, booked)
	if err != nil {
		log.Debug(err)
	}
	for k := range booked {
		ret[k] = IsBoked
	}
	return ret
}

// 返回体育馆当前日期的已预约时段   获取缓存
func getStadiumBooked(id uint, date string) (stadium.BookedTime, error) {
	exists, err := redis.KeyExists(getStadiumKey(id, date))
	if err != nil {
		log.Error(err)
		return loadIdlesCache(id, date)
	}
	if exists {
		all, err := redis.HGetAll(getStadiumKey(id, date))
		log.Errorf("getStadiumBooked err = %s\n", err)
		if err != nil {
			log.Debug("从DB获取")
			// 从DB获取
			return loadIdlesCache(id, date)
		}
		return all, nil
	}
	log.Debugf("key = %s not exists\n", getStadiumKey(id, date))
	return loadIdlesCache(id, date)
}

// 返回所有可预约时段
func GetAllowableTimes(query *IdleTimeQuery) ([]IdleTimes, int) {
	now := time.Now()
	nextDay := now.AddDate(0, 0, 1)
	today := now.Format(layout)
	tomorrow := nextDay.Format(layout)
	ret := make([]IdleTimes, 0)
	todayBooked := GetAllTimes(query.StadiumID, today)
	tomorrowBooked := GetAllTimes(query.StadiumID, tomorrow)
	ret = append(ret, allTimes2IdleTimes(query.StadiumID, today, todayBooked)...)
	ret = append(ret, allTimes2IdleTimes(query.StadiumID, tomorrow, tomorrowBooked)...)
	start := query.PerPage * (query.Page - 1)
	end := start + query.PerPage
	if start >= len(ret) {
		start = len(ret)
	}
	if end > len(ret) {
		end = len(ret)
	}
	return ret[start:end], len(ret)
}

// 转换可预约时段
func allTimes2IdleTimes(id uint, date string, all stadium.BookedTime) []IdleTimes {
	ret := make([]IdleTimes, 0)
	flag := false
	var start string
	for _, t := range timeTable {
		if all[t] == IsIdle {
			if !flag {
				start = t
				flag = true
			}
		} else {
			if flag {
				ret = append(ret, IdleTimes{
					Date:  date,
					Start: start,
					End:   t,
				})
				flag = false
			}
		}
	}
	end := timeTable[len(timeTable)-1]
	if flag && start != end {
		ret = append(ret, IdleTimes{
			Date:  date,
			Start: start,
			End:   end,
		})
	}
	return ret
}

// 加载缓存
func loadIdlesCache(id uint, date string) (stadium.BookedTime, error) {
	ret, err := getStadiumBookedFromDB(id, date)
	if err != nil {
		return nil, err
	}
	kv := make(map[string]interface{})
	for k, v := range ret {
		kv[k] = v
	}
	if len(kv) > 0 {
		err = redis.HMSet(getStadiumKey(id, date), kv)
	} else {
		err = redis.HSet(getStadiumKey(id, date), "", "")
	}
	return ret, err
}

// 删除redis缓存
func DeleteBookedCache(id uint, date string) error {
	return redis.Del(getStadiumKey(id, date))
}

// 从数据库获取体育馆当前日期的已预约时段
func getStadiumBookedFromDB(id uint, date string) (stadium.BookedTime, error) {
	var stadiumIdle stadium.StadiumBooked
	err := connection.GetDB().Where(&stadium.StadiumBooked{StadiumId: id, Date: date}).First(&stadiumIdle).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if err == gorm.ErrRecordNotFound {
		stadiumIdle.StadiumId = id
		stadiumIdle.Date = date
		stadiumIdle.BookedTimes = make(stadium.BookedTime)
		err = connection.GetDB().Create(&stadiumIdle).Error
	}
	return stadiumIdle.BookedTimes, err
}

// 获取redis键
func getStadiumKey(id uint, date string) string {
	return fmt.Sprintf("stadium:%d:%s", id, date)
}

// 预约场馆
func OrderStadium(parameter *OrderParameter) error {

	// 开启事务
	err := connection.GetDB().Transaction(func(tx *gorm.DB) error {
		_, err := OrderStadiumTransaction(tx, parameter)
		return err
	})
	if err != nil {
		return err
	}
	// 预约成功 删除缓存
	DeleteBookedCache(parameter.StadiumID, parameter.Date)
	return nil
}

// 预约场馆事务
func OrderStadiumTransaction(tx *gorm.DB, parameter *OrderParameter) (uint, error) {

	var booked stadium.StadiumBooked
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&stadium.StadiumBooked{StadiumId: parameter.StadiumID, Date: parameter.Date}).First(&booked).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return 1, err
	}
	if err == gorm.ErrRecordNotFound {
		booked.StadiumId = parameter.StadiumID
		booked.Date = parameter.Date
		booked.BookedTimes = make(stadium.BookedTime)
	}
	flag := true
	times := GetTimes(parameter.Start, parameter.End)
	for i := range times {
		_, exists := booked.BookedTimes[times[i]]
		if exists {
			flag = false
		}
	}
	if !flag {
		return 1, fmt.Errorf(OrderFail)
	}
	for i := range times {
		booked.BookedTimes[times[i]] = IsBoked
	}
	if err == gorm.ErrRecordNotFound {
		err = tx.Create(&booked).Error
	} else {
		err = tx.Save(&booked).Error
	}
	if err != nil {
		return 1, err
	}
	record := stadium.StadiumRecord {
		StadiumId: parameter.StadiumID,
		Date: parameter.Date,
		Stadium: parameter.Stadium,
		UserNo: parameter.Number,
		Start: parameter.Start,
		End: parameter.End,
	}
	return record.ID, tx.Create(&record).Error
}

// 取消预约
func CancelOrder(id uint) error {
	var record stadium.StadiumRecord
	err := connection.GetDB().Transaction(func(tx *gorm.DB) error {
		var g game.Game
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&game.Game{BookRecordID: id}).First(&g).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return err
		}
		if err == nil {
			err = tx.Unscoped().Delete(&g).Error
			if err != nil {
				return err
			}
			// 取消所有报名申请
			err = tx.Model(&game.GameRecord{}).Where("game_id = ?", g.ID).Update("status", "Cancel").Error
			if err != nil {
				return err
			}
		}
		return CancelOrderTransaction(tx, id, &record)
	})
	if err != nil {
		return err
	}
	// 取消预约成功 删除缓存
	DeleteBookedCache(id, record.Date)
	return nil
}

// 取消预约事务
func CancelOrderTransaction(tx *gorm.DB, id uint, record *stadium.StadiumRecord) error {
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(record, id).Error
	if err != nil {
		return err
	}
	err = tx.Unscoped().Delete(&stadium.StadiumRecord{}, id).Error
	if err != nil {
		return err
	}
	var booked stadium.StadiumBooked
	err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where(&stadium.StadiumBooked{StadiumId: record.StadiumId, Date: record.Date}).First(&booked).Error
	if err != nil {
		return err
	}
	times := GetTimes(record.Start, record.End)
	for i := range times {
		delete(booked.BookedTimes, times[i])
	}
	return tx.Save(&booked).Error
}

// 获取用户预约记录
func GetOrderRecord(query *Paramter, outdate bool) ([]stadium.StadiumRecord, int, error) {
	var ret []stadium.StadiumRecord
	err := connection.GetDB().Where(&stadium.StadiumRecord{
		UserNo:  query.Number,
		Stadium: query.Stadium,
		Date:    query.Date,
	}).Find(&ret).Error
	if err != nil {
		return nil, 0, err
	}
	now := time.Now().Local()
	records := make([]stadium.StadiumRecord, 0, len(ret))
	for _, record := range ret {
		t, err := time.ParseInLocation(layout, fmt.Sprintf("%s %s", record.Date, record.End), time.Local)
		if err != nil {
			return nil, 0, err
		}
		if outdate && t.Local().Before(now) {
			records = append(records, record)
		} else if !outdate && t.Local().After(now) {
			records = append(records, record)
		}
	}
	ret = records
	start := query.PerPage * (query.Page - 1)
	end := start + query.PerPage
	if start >= len(records) {
		start = len(records)
	}
	if len(records) < end {
		end = len(records)
	}
	return ret[start:end], len(records), nil
}