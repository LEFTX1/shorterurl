// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"go-zero-shorterurl/admin/internal/dal/model"
)

func newTLinkStatsToday(db *gorm.DB, opts ...gen.DOOption) tLinkStatsToday {
	_tLinkStatsToday := tLinkStatsToday{}

	_tLinkStatsToday.tLinkStatsTodayDo.UseDB(db, opts...)
	_tLinkStatsToday.tLinkStatsTodayDo.UseModel(&model.TLinkStatsToday{})

	tableName := _tLinkStatsToday.tLinkStatsTodayDo.TableName()
	_tLinkStatsToday.ALL = field.NewAsterisk(tableName)
	_tLinkStatsToday.ID = field.NewInt64(tableName, "id")
	_tLinkStatsToday.FullShortURL = field.NewString(tableName, "full_short_url")
	_tLinkStatsToday.Date = field.NewTime(tableName, "date")
	_tLinkStatsToday.TodayPv = field.NewInt32(tableName, "today_pv")
	_tLinkStatsToday.TodayUv = field.NewInt32(tableName, "today_uv")
	_tLinkStatsToday.TodayUip = field.NewInt32(tableName, "today_uip")
	_tLinkStatsToday.CreateTime = field.NewTime(tableName, "create_time")
	_tLinkStatsToday.UpdateTime = field.NewTime(tableName, "update_time")
	_tLinkStatsToday.DelFlag = field.NewBool(tableName, "del_flag")

	_tLinkStatsToday.fillFieldMap()

	return _tLinkStatsToday
}

type tLinkStatsToday struct {
	tLinkStatsTodayDo

	ALL          field.Asterisk
	ID           field.Int64  // ID
	FullShortURL field.String // 短链接
	Date         field.Time   // 日期
	TodayPv      field.Int32  // 今日PV
	TodayUv      field.Int32  // 今日UV
	TodayUip     field.Int32  // 今日IP数
	CreateTime   field.Time   // 创建时间
	UpdateTime   field.Time   // 修改时间
	DelFlag      field.Bool   // 删除标识 0：未删除 1：已删除

	fieldMap map[string]field.Expr
}

func (t tLinkStatsToday) Table(newTableName string) *tLinkStatsToday {
	t.tLinkStatsTodayDo.UseTable(newTableName)
	return t.updateTableName(newTableName)
}

func (t tLinkStatsToday) As(alias string) *tLinkStatsToday {
	t.tLinkStatsTodayDo.DO = *(t.tLinkStatsTodayDo.As(alias).(*gen.DO))
	return t.updateTableName(alias)
}

func (t *tLinkStatsToday) updateTableName(table string) *tLinkStatsToday {
	t.ALL = field.NewAsterisk(table)
	t.ID = field.NewInt64(table, "id")
	t.FullShortURL = field.NewString(table, "full_short_url")
	t.Date = field.NewTime(table, "date")
	t.TodayPv = field.NewInt32(table, "today_pv")
	t.TodayUv = field.NewInt32(table, "today_uv")
	t.TodayUip = field.NewInt32(table, "today_uip")
	t.CreateTime = field.NewTime(table, "create_time")
	t.UpdateTime = field.NewTime(table, "update_time")
	t.DelFlag = field.NewBool(table, "del_flag")

	t.fillFieldMap()

	return t
}

func (t *tLinkStatsToday) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := t.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (t *tLinkStatsToday) fillFieldMap() {
	t.fieldMap = make(map[string]field.Expr, 9)
	t.fieldMap["id"] = t.ID
	t.fieldMap["full_short_url"] = t.FullShortURL
	t.fieldMap["date"] = t.Date
	t.fieldMap["today_pv"] = t.TodayPv
	t.fieldMap["today_uv"] = t.TodayUv
	t.fieldMap["today_uip"] = t.TodayUip
	t.fieldMap["create_time"] = t.CreateTime
	t.fieldMap["update_time"] = t.UpdateTime
	t.fieldMap["del_flag"] = t.DelFlag
}

func (t tLinkStatsToday) clone(db *gorm.DB) tLinkStatsToday {
	t.tLinkStatsTodayDo.ReplaceConnPool(db.Statement.ConnPool)
	return t
}

func (t tLinkStatsToday) replaceDB(db *gorm.DB) tLinkStatsToday {
	t.tLinkStatsTodayDo.ReplaceDB(db)
	return t
}

type tLinkStatsTodayDo struct{ gen.DO }

type ITLinkStatsTodayDo interface {
	gen.SubQuery
	Debug() ITLinkStatsTodayDo
	WithContext(ctx context.Context) ITLinkStatsTodayDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() ITLinkStatsTodayDo
	WriteDB() ITLinkStatsTodayDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) ITLinkStatsTodayDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) ITLinkStatsTodayDo
	Not(conds ...gen.Condition) ITLinkStatsTodayDo
	Or(conds ...gen.Condition) ITLinkStatsTodayDo
	Select(conds ...field.Expr) ITLinkStatsTodayDo
	Where(conds ...gen.Condition) ITLinkStatsTodayDo
	Order(conds ...field.Expr) ITLinkStatsTodayDo
	Distinct(cols ...field.Expr) ITLinkStatsTodayDo
	Omit(cols ...field.Expr) ITLinkStatsTodayDo
	Join(table schema.Tabler, on ...field.Expr) ITLinkStatsTodayDo
	LeftJoin(table schema.Tabler, on ...field.Expr) ITLinkStatsTodayDo
	RightJoin(table schema.Tabler, on ...field.Expr) ITLinkStatsTodayDo
	Group(cols ...field.Expr) ITLinkStatsTodayDo
	Having(conds ...gen.Condition) ITLinkStatsTodayDo
	Limit(limit int) ITLinkStatsTodayDo
	Offset(offset int) ITLinkStatsTodayDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) ITLinkStatsTodayDo
	Unscoped() ITLinkStatsTodayDo
	Create(values ...*model.TLinkStatsToday) error
	CreateInBatches(values []*model.TLinkStatsToday, batchSize int) error
	Save(values ...*model.TLinkStatsToday) error
	First() (*model.TLinkStatsToday, error)
	Take() (*model.TLinkStatsToday, error)
	Last() (*model.TLinkStatsToday, error)
	Find() ([]*model.TLinkStatsToday, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.TLinkStatsToday, err error)
	FindInBatches(result *[]*model.TLinkStatsToday, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.TLinkStatsToday) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) ITLinkStatsTodayDo
	Assign(attrs ...field.AssignExpr) ITLinkStatsTodayDo
	Joins(fields ...field.RelationField) ITLinkStatsTodayDo
	Preload(fields ...field.RelationField) ITLinkStatsTodayDo
	FirstOrInit() (*model.TLinkStatsToday, error)
	FirstOrCreate() (*model.TLinkStatsToday, error)
	FindByPage(offset int, limit int) (result []*model.TLinkStatsToday, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) ITLinkStatsTodayDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (t tLinkStatsTodayDo) Debug() ITLinkStatsTodayDo {
	return t.withDO(t.DO.Debug())
}

func (t tLinkStatsTodayDo) WithContext(ctx context.Context) ITLinkStatsTodayDo {
	return t.withDO(t.DO.WithContext(ctx))
}

func (t tLinkStatsTodayDo) ReadDB() ITLinkStatsTodayDo {
	return t.Clauses(dbresolver.Read)
}

func (t tLinkStatsTodayDo) WriteDB() ITLinkStatsTodayDo {
	return t.Clauses(dbresolver.Write)
}

func (t tLinkStatsTodayDo) Session(config *gorm.Session) ITLinkStatsTodayDo {
	return t.withDO(t.DO.Session(config))
}

func (t tLinkStatsTodayDo) Clauses(conds ...clause.Expression) ITLinkStatsTodayDo {
	return t.withDO(t.DO.Clauses(conds...))
}

func (t tLinkStatsTodayDo) Returning(value interface{}, columns ...string) ITLinkStatsTodayDo {
	return t.withDO(t.DO.Returning(value, columns...))
}

func (t tLinkStatsTodayDo) Not(conds ...gen.Condition) ITLinkStatsTodayDo {
	return t.withDO(t.DO.Not(conds...))
}

func (t tLinkStatsTodayDo) Or(conds ...gen.Condition) ITLinkStatsTodayDo {
	return t.withDO(t.DO.Or(conds...))
}

func (t tLinkStatsTodayDo) Select(conds ...field.Expr) ITLinkStatsTodayDo {
	return t.withDO(t.DO.Select(conds...))
}

func (t tLinkStatsTodayDo) Where(conds ...gen.Condition) ITLinkStatsTodayDo {
	return t.withDO(t.DO.Where(conds...))
}

func (t tLinkStatsTodayDo) Order(conds ...field.Expr) ITLinkStatsTodayDo {
	return t.withDO(t.DO.Order(conds...))
}

func (t tLinkStatsTodayDo) Distinct(cols ...field.Expr) ITLinkStatsTodayDo {
	return t.withDO(t.DO.Distinct(cols...))
}

func (t tLinkStatsTodayDo) Omit(cols ...field.Expr) ITLinkStatsTodayDo {
	return t.withDO(t.DO.Omit(cols...))
}

func (t tLinkStatsTodayDo) Join(table schema.Tabler, on ...field.Expr) ITLinkStatsTodayDo {
	return t.withDO(t.DO.Join(table, on...))
}

func (t tLinkStatsTodayDo) LeftJoin(table schema.Tabler, on ...field.Expr) ITLinkStatsTodayDo {
	return t.withDO(t.DO.LeftJoin(table, on...))
}

func (t tLinkStatsTodayDo) RightJoin(table schema.Tabler, on ...field.Expr) ITLinkStatsTodayDo {
	return t.withDO(t.DO.RightJoin(table, on...))
}

func (t tLinkStatsTodayDo) Group(cols ...field.Expr) ITLinkStatsTodayDo {
	return t.withDO(t.DO.Group(cols...))
}

func (t tLinkStatsTodayDo) Having(conds ...gen.Condition) ITLinkStatsTodayDo {
	return t.withDO(t.DO.Having(conds...))
}

func (t tLinkStatsTodayDo) Limit(limit int) ITLinkStatsTodayDo {
	return t.withDO(t.DO.Limit(limit))
}

func (t tLinkStatsTodayDo) Offset(offset int) ITLinkStatsTodayDo {
	return t.withDO(t.DO.Offset(offset))
}

func (t tLinkStatsTodayDo) Scopes(funcs ...func(gen.Dao) gen.Dao) ITLinkStatsTodayDo {
	return t.withDO(t.DO.Scopes(funcs...))
}

func (t tLinkStatsTodayDo) Unscoped() ITLinkStatsTodayDo {
	return t.withDO(t.DO.Unscoped())
}

func (t tLinkStatsTodayDo) Create(values ...*model.TLinkStatsToday) error {
	if len(values) == 0 {
		return nil
	}
	return t.DO.Create(values)
}

func (t tLinkStatsTodayDo) CreateInBatches(values []*model.TLinkStatsToday, batchSize int) error {
	return t.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (t tLinkStatsTodayDo) Save(values ...*model.TLinkStatsToday) error {
	if len(values) == 0 {
		return nil
	}
	return t.DO.Save(values)
}

func (t tLinkStatsTodayDo) First() (*model.TLinkStatsToday, error) {
	if result, err := t.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.TLinkStatsToday), nil
	}
}

func (t tLinkStatsTodayDo) Take() (*model.TLinkStatsToday, error) {
	if result, err := t.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.TLinkStatsToday), nil
	}
}

func (t tLinkStatsTodayDo) Last() (*model.TLinkStatsToday, error) {
	if result, err := t.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.TLinkStatsToday), nil
	}
}

func (t tLinkStatsTodayDo) Find() ([]*model.TLinkStatsToday, error) {
	result, err := t.DO.Find()
	return result.([]*model.TLinkStatsToday), err
}

func (t tLinkStatsTodayDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.TLinkStatsToday, err error) {
	buf := make([]*model.TLinkStatsToday, 0, batchSize)
	err = t.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (t tLinkStatsTodayDo) FindInBatches(result *[]*model.TLinkStatsToday, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return t.DO.FindInBatches(result, batchSize, fc)
}

func (t tLinkStatsTodayDo) Attrs(attrs ...field.AssignExpr) ITLinkStatsTodayDo {
	return t.withDO(t.DO.Attrs(attrs...))
}

func (t tLinkStatsTodayDo) Assign(attrs ...field.AssignExpr) ITLinkStatsTodayDo {
	return t.withDO(t.DO.Assign(attrs...))
}

func (t tLinkStatsTodayDo) Joins(fields ...field.RelationField) ITLinkStatsTodayDo {
	for _, _f := range fields {
		t = *t.withDO(t.DO.Joins(_f))
	}
	return &t
}

func (t tLinkStatsTodayDo) Preload(fields ...field.RelationField) ITLinkStatsTodayDo {
	for _, _f := range fields {
		t = *t.withDO(t.DO.Preload(_f))
	}
	return &t
}

func (t tLinkStatsTodayDo) FirstOrInit() (*model.TLinkStatsToday, error) {
	if result, err := t.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.TLinkStatsToday), nil
	}
}

func (t tLinkStatsTodayDo) FirstOrCreate() (*model.TLinkStatsToday, error) {
	if result, err := t.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.TLinkStatsToday), nil
	}
}

func (t tLinkStatsTodayDo) FindByPage(offset int, limit int) (result []*model.TLinkStatsToday, count int64, err error) {
	result, err = t.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = t.Offset(-1).Limit(-1).Count()
	return
}

func (t tLinkStatsTodayDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = t.Count()
	if err != nil {
		return
	}

	err = t.Offset(offset).Limit(limit).Scan(result)
	return
}

func (t tLinkStatsTodayDo) Scan(result interface{}) (err error) {
	return t.DO.Scan(result)
}

func (t tLinkStatsTodayDo) Delete(models ...*model.TLinkStatsToday) (result gen.ResultInfo, err error) {
	return t.DO.Delete(models)
}

func (t *tLinkStatsTodayDo) withDO(do gen.Dao) *tLinkStatsTodayDo {
	t.DO = *do.(*gen.DO)
	return t
}
