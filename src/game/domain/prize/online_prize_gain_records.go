package prize

import (
	"github.com/golang/glog"
	"time"
	"util"
	"pb"
	"code.google.com/p/goprotobuf/proto"
)

type OnlinePrizeGainRecords struct {
	UserId  string
	Records map[int]*OnlinePrizeGainRecord
}

func NewOnlinePrizeGainRecords(userId string) *OnlinePrizeGainRecords {
	records := &OnlinePrizeGainRecords{}
	records.UserId = userId
	records.Records = make(map[int]*OnlinePrizeGainRecord)

	rs, err := FindOnlinePrizeGainRecords(userId)
	if err == nil {
		for _, r := range rs {
			glog.V(2).Info("===>在线奖励领取记录:", r)
			records.Records[r.PrizeId] = r
		}
	}

	return records
}

func (rs *OnlinePrizeGainRecords) Save() {
	for _, r := range rs.Records {
		SaveOnlinePrizeGainRecord(r)
	}
}

func (rs *OnlinePrizeGainRecords) IsGained(prizeId int) bool {
	r := rs.Records[prizeId]
	if r == nil {
		return false
	}

	return util.CompareDate(time.Now(), r.GainTime)
}

func (rs *OnlinePrizeGainRecords) SetGained(prizeId int, gainTime time.Time) {
	r := rs.Records[prizeId]
	if r == nil {
		r = &OnlinePrizeGainRecord{}
		r.UserId = rs.UserId
		r.PrizeId = prizeId
		rs.Records[r.PrizeId] = r
	}
	r.GainTime = gainTime
}

/*func (rs *OnlinePrizeGainRecords) BuildMessage() []int32 {
	result := []int32{}

	for _, r := range rs.Records {
		if !util.CompareDate(r.GainTime, time.Now()) {
			continue
		}
		result = append(result, int32(r.PrizeId))
	}

	return result
}*/

func (rs *OnlinePrizeGainRecords) BuildMessage() []*pb.PrizeOnline {
	result := []*pb.PrizeOnline{}
    prizes := GetOnlinePrizeManager().GetOnlineAllPrize()
    now := time.Now()
    for _, v := range prizes {
        beginTime := time.Date(now.Year(), now.Month(), now.Day(), int(v.BeginTime/60), int(v.BeginTime%60), 0, 0, time.Local)
        endTime := time.Date(now.Year(), now.Month(), now.Day(), int(v.EndTime/60), int(v.EndTime%60), 0, 0, time.Local)

        r := rs.Records[v.PrizeID]
        p := &pb.PrizeOnline{}
        if r != nil && util.CompareDate(r.GainTime, time.Now()) {
            p.State = proto.Int(1)
        } else if time.Since(beginTime).Seconds() >= 0 && time.Since(endTime).Seconds() <= 0 {
            p.State = proto.Int(2)
        } else {
            p.State = proto.Int(0)
        }
        p.PrizeId = proto.Int(v.PrizeID)
        p.PrizeTitle = proto.String(v.PrizeTitle)
        p.IcoRes = proto.String(v.IcoRes)
        p.PrizeGold = proto.Int(v.PrizeGold)
        result = append(result, p)
    }

	return result
}
