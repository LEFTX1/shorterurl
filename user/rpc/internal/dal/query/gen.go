// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"
	"database/sql"

	"gorm.io/gorm"

	"gorm.io/gen"

	"gorm.io/plugin/dbresolver"
)

var (
	Q                = new(Query)
	TGroup           *tGroup
	TGroupUnique     *tGroupUnique
	TLink            *tLink
	TLinkAccessLog   *tLinkAccessLog
	TLinkAccessStat  *tLinkAccessStat
	TLinkBrowserStat *tLinkBrowserStat
	TLinkDeviceStat  *tLinkDeviceStat
	TLinkGoto        *tLinkGoto
	TLinkLocaleStat  *tLinkLocaleStat
	TLinkNetworkStat *tLinkNetworkStat
	TLinkOsStat      *tLinkOsStat
	TLinkStatsToday  *tLinkStatsToday
	TUser            *tUser
)

func SetDefault(db *gorm.DB, opts ...gen.DOOption) {
	*Q = *Use(db, opts...)
	TGroup = &Q.TGroup
	TGroupUnique = &Q.TGroupUnique
	TLink = &Q.TLink
	TLinkAccessLog = &Q.TLinkAccessLog
	TLinkAccessStat = &Q.TLinkAccessStat
	TLinkBrowserStat = &Q.TLinkBrowserStat
	TLinkDeviceStat = &Q.TLinkDeviceStat
	TLinkGoto = &Q.TLinkGoto
	TLinkLocaleStat = &Q.TLinkLocaleStat
	TLinkNetworkStat = &Q.TLinkNetworkStat
	TLinkOsStat = &Q.TLinkOsStat
	TLinkStatsToday = &Q.TLinkStatsToday
	TUser = &Q.TUser
}

func Use(db *gorm.DB, opts ...gen.DOOption) *Query {
	return &Query{
		db:               db,
		TGroup:           newTGroup(db, opts...),
		TGroupUnique:     newTGroupUnique(db, opts...),
		TLink:            newTLink(db, opts...),
		TLinkAccessLog:   newTLinkAccessLog(db, opts...),
		TLinkAccessStat:  newTLinkAccessStat(db, opts...),
		TLinkBrowserStat: newTLinkBrowserStat(db, opts...),
		TLinkDeviceStat:  newTLinkDeviceStat(db, opts...),
		TLinkGoto:        newTLinkGoto(db, opts...),
		TLinkLocaleStat:  newTLinkLocaleStat(db, opts...),
		TLinkNetworkStat: newTLinkNetworkStat(db, opts...),
		TLinkOsStat:      newTLinkOsStat(db, opts...),
		TLinkStatsToday:  newTLinkStatsToday(db, opts...),
		TUser:            newTUser(db, opts...),
	}
}

type Query struct {
	db *gorm.DB

	TGroup           tGroup
	TGroupUnique     tGroupUnique
	TLink            tLink
	TLinkAccessLog   tLinkAccessLog
	TLinkAccessStat  tLinkAccessStat
	TLinkBrowserStat tLinkBrowserStat
	TLinkDeviceStat  tLinkDeviceStat
	TLinkGoto        tLinkGoto
	TLinkLocaleStat  tLinkLocaleStat
	TLinkNetworkStat tLinkNetworkStat
	TLinkOsStat      tLinkOsStat
	TLinkStatsToday  tLinkStatsToday
	TUser            tUser
}

func (q *Query) Available() bool { return q.db != nil }

func (q *Query) clone(db *gorm.DB) *Query {
	return &Query{
		db:               db,
		TGroup:           q.TGroup.clone(db),
		TGroupUnique:     q.TGroupUnique.clone(db),
		TLink:            q.TLink.clone(db),
		TLinkAccessLog:   q.TLinkAccessLog.clone(db),
		TLinkAccessStat:  q.TLinkAccessStat.clone(db),
		TLinkBrowserStat: q.TLinkBrowserStat.clone(db),
		TLinkDeviceStat:  q.TLinkDeviceStat.clone(db),
		TLinkGoto:        q.TLinkGoto.clone(db),
		TLinkLocaleStat:  q.TLinkLocaleStat.clone(db),
		TLinkNetworkStat: q.TLinkNetworkStat.clone(db),
		TLinkOsStat:      q.TLinkOsStat.clone(db),
		TLinkStatsToday:  q.TLinkStatsToday.clone(db),
		TUser:            q.TUser.clone(db),
	}
}

func (q *Query) ReadDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Read))
}

func (q *Query) WriteDB() *Query {
	return q.ReplaceDB(q.db.Clauses(dbresolver.Write))
}

func (q *Query) ReplaceDB(db *gorm.DB) *Query {
	return &Query{
		db:               db,
		TGroup:           q.TGroup.replaceDB(db),
		TGroupUnique:     q.TGroupUnique.replaceDB(db),
		TLink:            q.TLink.replaceDB(db),
		TLinkAccessLog:   q.TLinkAccessLog.replaceDB(db),
		TLinkAccessStat:  q.TLinkAccessStat.replaceDB(db),
		TLinkBrowserStat: q.TLinkBrowserStat.replaceDB(db),
		TLinkDeviceStat:  q.TLinkDeviceStat.replaceDB(db),
		TLinkGoto:        q.TLinkGoto.replaceDB(db),
		TLinkLocaleStat:  q.TLinkLocaleStat.replaceDB(db),
		TLinkNetworkStat: q.TLinkNetworkStat.replaceDB(db),
		TLinkOsStat:      q.TLinkOsStat.replaceDB(db),
		TLinkStatsToday:  q.TLinkStatsToday.replaceDB(db),
		TUser:            q.TUser.replaceDB(db),
	}
}

type queryCtx struct {
	TGroup           ITGroupDo
	TGroupUnique     ITGroupUniqueDo
	TLink            ITLinkDo
	TLinkAccessLog   ITLinkAccessLogDo
	TLinkAccessStat  ITLinkAccessStatDo
	TLinkBrowserStat ITLinkBrowserStatDo
	TLinkDeviceStat  ITLinkDeviceStatDo
	TLinkGoto        ITLinkGotoDo
	TLinkLocaleStat  ITLinkLocaleStatDo
	TLinkNetworkStat ITLinkNetworkStatDo
	TLinkOsStat      ITLinkOsStatDo
	TLinkStatsToday  ITLinkStatsTodayDo
	TUser            ITUserDo
}

func (q *Query) WithContext(ctx context.Context) *queryCtx {
	return &queryCtx{
		TGroup:           q.TGroup.WithContext(ctx),
		TGroupUnique:     q.TGroupUnique.WithContext(ctx),
		TLink:            q.TLink.WithContext(ctx),
		TLinkAccessLog:   q.TLinkAccessLog.WithContext(ctx),
		TLinkAccessStat:  q.TLinkAccessStat.WithContext(ctx),
		TLinkBrowserStat: q.TLinkBrowserStat.WithContext(ctx),
		TLinkDeviceStat:  q.TLinkDeviceStat.WithContext(ctx),
		TLinkGoto:        q.TLinkGoto.WithContext(ctx),
		TLinkLocaleStat:  q.TLinkLocaleStat.WithContext(ctx),
		TLinkNetworkStat: q.TLinkNetworkStat.WithContext(ctx),
		TLinkOsStat:      q.TLinkOsStat.WithContext(ctx),
		TLinkStatsToday:  q.TLinkStatsToday.WithContext(ctx),
		TUser:            q.TUser.WithContext(ctx),
	}
}

func (q *Query) Transaction(fc func(tx *Query) error, opts ...*sql.TxOptions) error {
	return q.db.Transaction(func(tx *gorm.DB) error { return fc(q.clone(tx)) }, opts...)
}

func (q *Query) Begin(opts ...*sql.TxOptions) *QueryTx {
	tx := q.db.Begin(opts...)
	return &QueryTx{Query: q.clone(tx), Error: tx.Error}
}

type QueryTx struct {
	*Query
	Error error
}

func (q *QueryTx) Commit() error {
	return q.db.Commit().Error
}

func (q *QueryTx) Rollback() error {
	return q.db.Rollback().Error
}

func (q *QueryTx) SavePoint(name string) error {
	return q.db.SavePoint(name).Error
}

func (q *QueryTx) RollbackTo(name string) error {
	return q.db.RollbackTo(name).Error
}