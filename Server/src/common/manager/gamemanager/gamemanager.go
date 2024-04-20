package gamemanager

import (
	"common/connection"
	"common/manager/ordermanager"
	"common/model/game"
	"common/model/stadium"
	"fmt"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	Accept = "Accept"
	AcceptStr = "报名通过"
	Wait   = "Wait"
	WaitStr = "等待通过"
	Refuse = "Refuse"
	RefuseStr = "报名被拒绝"
	Cancel = "Cancel"
	CancelStr = "报名已取消"
)

func StatusTransition(status string) string {
	if status == Accept {
		return AcceptStr
	}
	if status == AcceptStr {
		return Accept
	}
	if status == Wait {
		return WaitStr
	}
	if status == WaitStr {
		return Wait
	}
	if status == Refuse {
		return RefuseStr
	}
	if status == RefuseStr {
		return Refuse
	}
	if status == Cancel {
		return CancelStr
	}
	return Cancel
}

// 创建比赛活动
func CreateGame(parameter *CreateParameter) error {
	// 开启事务
	err := connection.GetDB().Transaction(func(tx *gorm.DB) error {
		// 预约场馆
		id, err := ordermanager.OrderStadiumTransaction(tx, &ordermanager.OrderParameter{
			StadiumID: parameter.StadiumID,
			Number:    parameter.Number,
			Stadium:   parameter.Stadium,
			Date:      parameter.Date,
			Start:     parameter.Start,
			End:       parameter.End,
		})
		if err != nil {
			return err
		}
		// 预约成功 创建比赛活动
		g := game.Game {
			Title:      parameter.Title,
			Date:       parameter.Date,
			Stadium:    parameter.Stadium,
			StadiumID:  parameter.StadiumID,
			BookRecordID: id,
			Start:      parameter.Start,
			End:        parameter.End,
			UserNumber: parameter.Number,
			UserName:   parameter.UserName,
			Remark:     parameter.Remark,
			Count:      1,
			Maximum:    parameter.Maximum,
		}
		err = tx.Create(&g).Error
		if err != nil {
			return err
		}
		// 创建报名记录
		return tx.Create(&game.GameRecord {
			UserNumber: parameter.Number,
			UserName:   parameter.UserName,
			Stadium:    parameter.Stadium,
			GameID:     g.ID,
			Title:      parameter.Title,
			Date:       parameter.Date,
			Start:      parameter.Start,
			End:        parameter.End,
			HostNumber: parameter.Number,
			HostName:   parameter.UserName,
			Status:     Accept,
		}).Error
	})
	if err != nil {
		return err
	}
	// 删除缓存
	ordermanager.DeleteBookedCache(parameter.StadiumID, parameter.Date)
	return nil
}

// 取消比赛活动
func CancelGame(parameter *CancelParameter) error {
	var stadiumRecord stadium.StadiumRecord
	// 开启事务
	err := connection.GetDB().Transaction(func(tx *gorm.DB) error {
		var g game.Game
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&g, parameter.GameID).Error
		if err != nil {
			return err
		}
		err = tx.Unscoped().Delete(&game.Game{}, parameter.GameID).Error
		if err != nil {
			return err
		}
		// 取消体育场馆预约
		err = ordermanager.CancelOrderTransaction(tx, parameter.BookRecordID, &stadiumRecord)
		if err != nil {
			return err
		}
		// 取消所有报名申请
		return tx.Model(&game.GameRecord{}).Where("game_id = ?", parameter.GameID).Update("status", Cancel).Error
	})
	if err != nil {
		return err
	}
	ordermanager.DeleteBookedCache(stadiumRecord.StadiumId, stadiumRecord.Date)
	return nil
}


// 报名比赛活动
func ApplyGame(parameter *ApplyParameter) error {
	// 开启事务
	return connection.GetDB().Transaction(func(tx *gorm.DB) error {
		var g game.Game
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&g, parameter.GameID).Error
		if err != nil {
			return err
		}
		if g.Count >= g.Maximum {
			return fmt.Errorf("count maxed out")
		}
		var record game.GameRecord
		err = tx.Where(&game.GameRecord {
			GameID: parameter.GameID,
			UserNumber: parameter.Number,
		}).First(&record).Error
		if err == nil || err != gorm.ErrRecordNotFound {
			return fmt.Errorf("Applied")
		}
		// 报名
		return tx.Create(&game.GameRecord {
			UserNumber: parameter.Number,
			UserName: parameter.UserName,
			Stadium: g.Stadium,
			GameID: g.ID,
			Title: g.Title,
			Date: g.Date,
			Start: g.Start,
			End: g.End,
			HostNumber: g.UserNumber,
			HostName: g.UserName,
			Status: Wait,
		}).Error
	})
}

// 取消报名比赛
func CancelApplyGame(parameter *CancelApplyParameter) error {
	return connection.GetDB().Transaction(func(tx *gorm.DB) error {
		var record game.GameRecord
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&record, parameter.RecordID).Error
		if err != nil {
			return err
		}
		fmt.Println(record)
		// 取消报名
		if record.Status == Accept {
			err = tx.Model(&game.Game {Model: gorm.Model{ID: record.GameID}}).Update("count", gorm.Expr("count - ?", 1)).Error
			if err != nil {
				return err
			}
		}
		return tx.Unscoped().Delete(&record).Error
	})
}

// 通过报名
func AcceptApply(parameter *AcceptApplyParameter) error {
	return connection.GetDB().Transaction(func(tx *gorm.DB) error {
		var record game.GameRecord
		err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&record, parameter.RecordID).Error
		if err != nil {
			return err
		}
		// 通过报名
		err = tx.Model(&game.Game {Model: gorm.Model{ID: record.GameID}}).Where("count < maximum").Update("count", gorm.Expr("count + ?", 1)).Error
		if err != nil {
			return err
		}
		return tx.Model(&record).Update("status", Accept).Error
	})
}

// 拒绝报名
func RefuseApply(parameter *RefuseApplyParameter) error {
	return connection.GetDB().Model(&game.GameRecord{Model: gorm.Model{ID: parameter.RecordID}}).Update("status", Refuse).Error
}

// 返回比赛活动列表 当前
func GameList(parameter *GameListParameter) ([]game.Game, int, error) {
	var games []game.Game
	err := connection.GetDB().Where(&game.Game {
		UserNumber: parameter.Number,
		UserName: parameter.UserName,
		Date: parameter.Date,
		Stadium: parameter.Stadium,
		Title: parameter.Title,
	}).Find(&games).Error
	if err != nil {
		return nil, 0, err
	}
	ret := make([]game.Game, 0, len(games))
	now := time.Now().Local()
	for _, g := range games {
		t, err := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", g.Date, g.End), time.Local)
		if err !=  nil {
			continue
		}
		if parameter.CanApply && g.Count == g.Maximum {
			continue
		}
		if t.After(now) {
			ret = append(ret, g)
		}
	}
	start := parameter.PerPage * (parameter.Page - 1)
	end := start + parameter.PerPage
	if start >= len(ret) {
		start = len(ret)
	}
	if len(ret) < end {
		end = len(ret)
	}
	return ret[start:end], len(ret), nil
}

// 返回过期比赛活动
func OutdatedGameList(parameter *GameListParameter) ([]game.Game, int, error) {
	var games []game.Game
	err := connection.GetDB().Where(&game.Game {
		UserNumber: parameter.Number,
		UserName: parameter.UserName,
		Date: parameter.Date,
		Stadium: parameter.Stadium,
		Title: parameter.Title,
	}).Find(&games).Error
	if err != nil {
		return nil, 0, err
	}
	ret := make([]game.Game, 0, len(games))
	now := time.Now().Local()
	for _, g := range games {
		t, err := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", g.Date, g.End), time.Local)
		if err !=  nil {
			continue
		}
		if parameter.CanApply && g.Count == g.Maximum {
			continue
		}
		if t.Before(now) {
			ret = append(ret, g)
		}
	}
	start := parameter.PerPage * (parameter.Page - 1)
	end := start + parameter.PerPage
	if start >= len(ret) {
		start = len(ret)
	}
	if len(ret) < end {
		end = len(ret)
	}
	return ret[start:end], len(ret), nil
}

// 返回过期申请记录
func OutdatedRecord(parameter *ApplyRecordParameter) ([]GameRecord, int, error) {
	var records []game.GameRecord
	err := connection.GetDB().Where(&game.GameRecord{
		UserNumber: parameter.Number,
		Stadium: parameter.Stadium,
		Date: parameter.Date,
		HostNumber: parameter.HostNumber,
		Title: parameter.Title,
		Status: Accept,
	}).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}
	now := time.Now().Local()
	ret := make([]GameRecord, 0, len(records))
	for _, record := range records {
		t, err := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", record.Date, record.End), time.Local)
		if err != nil {
			continue
		}
		if t.Before(now) {
			ret = append(ret, GameRecord{
				ID: record.GameID,
				UserNumber: record.UserName,
				UserName: record.UserName,
				GameID: record.GameID,
				Title: record.Title,
				Stadium: record.Stadium,
				Date: record.Date,
				Start: record.Start,
				End: record.End,
				HostName: record.HostName,
				Status: record.Status,
				StatusStr: StatusTransition(record.Status),
			})
		}
	}
	start := parameter.PerPage * (parameter.Page - 1)
	end := start + parameter.PerPage
	if start >= len(ret) {
		start = len(ret)
	}
	if len(ret) < end {
		end = len(ret)
	}
	return ret[start:end], len(ret), nil
}

// 返回近期申请
func RecentRecords(parameter *ApplyRecordParameter) ([]GameRecord, int, error) {
	var records []game.GameRecord
	err := connection.GetDB().Where(&game.GameRecord{
		UserNumber: parameter.Number,
		UserName: parameter.UserName,
		Stadium: parameter.Stadium,
		Date: parameter.Date,
		HostNumber: parameter.HostNumber,
		Title: parameter.Title,
	}).Find(&records).Error
	if err != nil {
		return nil, 0, err
	}
	now := time.Now().Local()
	ret := make([]GameRecord, 0, len(records))
	for _, record := range records {
		t, err := time.ParseInLocation("2006-01-02 15:04", fmt.Sprintf("%s %s", record.Date, record.End), time.Local)
		if err != nil {
			continue
		}
		if t.After(now) {
			ret = append(ret, GameRecord{
				ID: record.ID,
				UserNumber: record.UserNumber,
				UserName: record.UserName,
				GameID: record.GameID,
				Title: record.Title,
				Stadium: record.Stadium,
				Date: record.Date,
				Start: record.Start,
				End: record.End,
				HostName: record.HostName,
				Status: record.Status,
				StatusStr: StatusTransition(record.Status),
			})
		}
	}
	start := parameter.PerPage * (parameter.Page - 1)
	end := start + parameter.PerPage
	if start >= len(ret) {
		start = len(ret)
	}
	if len(ret) < end {
		end = len(ret)
	}
	return ret[start:end], len(ret), nil
}
