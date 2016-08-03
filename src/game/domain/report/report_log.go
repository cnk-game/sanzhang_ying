package report

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
	"util"
)

type ReportLog struct {
	UserId       string    `bson:"userId"`
	ReportUserId string    `bson:"reportUserId"`
	Time         time.Time `bson:"time"`
}

const (
	reportLogC = "report_log"
)

func SaveReportLog(l *ReportLog) error {
	l.Time = time.Now()
	return util.WithLogCollection(reportLogC, func(c *mgo.Collection) error {
		_, err := c.Upsert(bson.M{"userId": l.UserId, "reportUserId": l.ReportUserId}, l)
		return err
	})
}
