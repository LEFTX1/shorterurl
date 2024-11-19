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

func newTLinkNetworkStat(db *gorm.DB, opts ...gen.DOOption) tLinkNetworkStat {
	_tLinkNetworkStat := tLinkNetworkStat{}

	_tLinkNetworkStat.tLinkNetworkStatDo.UseDB(db, opts...)
	_tLinkNetworkStat.tLinkNetworkStatDo.UseModel(&model.TLinkNetworkStat{})

	tableName := _tLinkNetworkStat.tLinkNetworkStatDo.TableName()
	_tLinkNetworkStat.ALL = field.NewAsterisk(tableName)
	_tLinkNetworkStat.ID = field.NewInt64(tableName, "id")
	_tLinkNetworkStat.FullShortURL = field.NewString(tableName, "full_short_url")
	_tLinkNetworkStat.Date = field.NewTime(tableName, "date")
	_tLinkNetworkStat.Cnt = field.NewInt32(tableName, "cnt")
	_tLinkNetworkStat.Network = field.NewString(tableName, "network")
	_tLinkNetworkStat.CreateTime = field.NewTime(tableName, "create_time")
	_tLinkNetworkStat.UpdateTime = field.NewTime(tableName, "update_time")
	_tLinkNetworkStat.DelFlag = field.NewBool(tableName, "del_flag")

	_tLinkNetworkStat.fillFieldMap()

	return _tLinkNetworkStat
}

type tLinkNetworkStat struct {
	tLinkNetworkStatDo

	ALL          field.Asterisk
	ID           field.Int64  // ID
	FullShortURL field.String // 完整短链接
	Date         field.Time   // 日期
	Cnt          field.Int32  // 访问量
	Network      field.String // 访问网络
	CreateTime   field.Time   // 创建时间
	UpdateTime   field.Time   // 修改时间
	DelFlag      field.Bool   // 删除标识 0：未删除 1：已删除

	fieldMap map[string]field.Expr
}

func (t tLinkNetworkStat) Table(newTableName string) *tLinkNetworkStat {
	t.tLinkNetworkStatDo.UseTable(newTableName)
	return t.updateTableName(newTableName)
}

func (t tLinkNetworkStat) As(alias string) *tLinkNetworkStat {
	t.tLinkNetworkStatDo.DO = *(t.tLinkNetworkStatDo.As(alias).(*gen.DO))
	return t.updateTableName(alias)
}

func (t *tLinkNetworkStat) updateTableName(table string) *tLinkNetworkStat {
	t.ALL = field.NewAsterisk(table)
	t.ID = field.NewInt64(table, "id")
	t.FullShortURL = field.NewString(table, "full_short_url")
	t.Date = field.NewTime(table, "date")
	t.Cnt = field.NewInt32(table, "cnt")
	t.Network = field.NewString(table, "network")
	t.CreateTime = field.NewTime(table, "create_time")
	t.UpdateTime = field.NewTime(table, "update_time")
	t.DelFlag = field.NewBool(table, "del_flag")

	t.fillFieldMap()

	return t
}

func (t *tLinkNetworkStat) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := t.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (t *tLinkNetworkStat) fillFieldMap() {
	t.fieldMap = make(map[string]field.Expr, 8)
	t.fieldMap["id"] = t.ID
	t.fieldMap["full_short_url"] = t.FullShortURL
	t.fieldMap["date"] = t.Date
	t.fieldMap["cnt"] = t.Cnt
	t.fieldMap["network"] = t.Network
	t.fieldMap["create_time"] = t.CreateTime
	t.fieldMap["update_time"] = t.UpdateTime
	t.fieldMap["del_flag"] = t.DelFlag
}

func (t tLinkNetworkStat) clone(db *gorm.DB) tLinkNetworkStat {
	t.tLinkNetworkStatDo.ReplaceConnPool(db.Statement.ConnPool)
	return t
}

func (t tLinkNetworkStat) replaceDB(db *gorm.DB) tLinkNetworkStat {
	t.tLinkNetworkStatDo.ReplaceDB(db)
	return t
}

type tLinkNetworkStatDo struct{ gen.DO }

type ITLinkNetworkStatDo interface {
	gen.SubQuery
	Debug() ITLinkNetworkStatDo
	WithContext(ctx context.Context) ITLinkNetworkStatDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() ITLinkNetworkStatDo
	WriteDB() ITLinkNetworkStatDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) ITLinkNetworkStatDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) ITLinkNetworkStatDo
	Not(conds ...gen.Condition) ITLinkNetworkStatDo
	Or(conds ...gen.Condition) ITLinkNetworkStatDo
	Select(conds ...field.Expr) ITLinkNetworkStatDo
	Where(conds ...gen.Condition) ITLinkNetworkStatDo
	Order(conds ...field.Expr) ITLinkNetworkStatDo
	Distinct(cols ...field.Expr) ITLinkNetworkStatDo
	Omit(cols ...field.Expr) ITLinkNetworkStatDo
	Join(table schema.Tabler, on ...field.Expr) ITLinkNetworkStatDo
	LeftJoin(table schema.Tabler, on ...field.Expr) ITLinkNetworkStatDo
	RightJoin(table schema.Tabler, on ...field.Expr) ITLinkNetworkStatDo
	Group(cols ...field.Expr) ITLinkNetworkStatDo
	Having(conds ...gen.Condition) ITLinkNetworkStatDo
	Limit(limit int) ITLinkNetworkStatDo
	Offset(offset int) ITLinkNetworkStatDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) ITLinkNetworkStatDo
	Unscoped() ITLinkNetworkStatDo
	Create(values ...*model.TLinkNetworkStat) error
	CreateInBatches(values []*model.TLinkNetworkStat, batchSize int) error
	Save(values ...*model.TLinkNetworkStat) error
	First() (*model.TLinkNetworkStat, error)
	Take() (*model.TLinkNetworkStat, error)
	Last() (*model.TLinkNetworkStat, error)
	Find() ([]*model.TLinkNetworkStat, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.TLinkNetworkStat, err error)
	FindInBatches(result *[]*model.TLinkNetworkStat, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.TLinkNetworkStat) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) ITLinkNetworkStatDo
	Assign(attrs ...field.AssignExpr) ITLinkNetworkStatDo
	Joins(fields ...field.RelationField) ITLinkNetworkStatDo
	Preload(fields ...field.RelationField) ITLinkNetworkStatDo
	FirstOrInit() (*model.TLinkNetworkStat, error)
	FirstOrCreate() (*model.TLinkNetworkStat, error)
	FindByPage(offset int, limit int) (result []*model.TLinkNetworkStat, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) ITLinkNetworkStatDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (t tLinkNetworkStatDo) Debug() ITLinkNetworkStatDo {
	return t.withDO(t.DO.Debug())
}

func (t tLinkNetworkStatDo) WithContext(ctx context.Context) ITLinkNetworkStatDo {
	return t.withDO(t.DO.WithContext(ctx))
}

func (t tLinkNetworkStatDo) ReadDB() ITLinkNetworkStatDo {
	return t.Clauses(dbresolver.Read)
}

func (t tLinkNetworkStatDo) WriteDB() ITLinkNetworkStatDo {
	return t.Clauses(dbresolver.Write)
}

func (t tLinkNetworkStatDo) Session(config *gorm.Session) ITLinkNetworkStatDo {
	return t.withDO(t.DO.Session(config))
}

func (t tLinkNetworkStatDo) Clauses(conds ...clause.Expression) ITLinkNetworkStatDo {
	return t.withDO(t.DO.Clauses(conds...))
}

func (t tLinkNetworkStatDo) Returning(value interface{}, columns ...string) ITLinkNetworkStatDo {
	return t.withDO(t.DO.Returning(value, columns...))
}

func (t tLinkNetworkStatDo) Not(conds ...gen.Condition) ITLinkNetworkStatDo {
	return t.withDO(t.DO.Not(conds...))
}

func (t tLinkNetworkStatDo) Or(conds ...gen.Condition) ITLinkNetworkStatDo {
	return t.withDO(t.DO.Or(conds...))
}

func (t tLinkNetworkStatDo) Select(conds ...field.Expr) ITLinkNetworkStatDo {
	return t.withDO(t.DO.Select(conds...))
}

func (t tLinkNetworkStatDo) Where(conds ...gen.Condition) ITLinkNetworkStatDo {
	return t.withDO(t.DO.Where(conds...))
}

func (t tLinkNetworkStatDo) Order(conds ...field.Expr) ITLinkNetworkStatDo {
	return t.withDO(t.DO.Order(conds...))
}

func (t tLinkNetworkStatDo) Distinct(cols ...field.Expr) ITLinkNetworkStatDo {
	return t.withDO(t.DO.Distinct(cols...))
}

func (t tLinkNetworkStatDo) Omit(cols ...field.Expr) ITLinkNetworkStatDo {
	return t.withDO(t.DO.Omit(cols...))
}

func (t tLinkNetworkStatDo) Join(table schema.Tabler, on ...field.Expr) ITLinkNetworkStatDo {
	return t.withDO(t.DO.Join(table, on...))
}

func (t tLinkNetworkStatDo) LeftJoin(table schema.Tabler, on ...field.Expr) ITLinkNetworkStatDo {
	return t.withDO(t.DO.LeftJoin(table, on...))
}

func (t tLinkNetworkStatDo) RightJoin(table schema.Tabler, on ...field.Expr) ITLinkNetworkStatDo {
	return t.withDO(t.DO.RightJoin(table, on...))
}

func (t tLinkNetworkStatDo) Group(cols ...field.Expr) ITLinkNetworkStatDo {
	return t.withDO(t.DO.Group(cols...))
}

func (t tLinkNetworkStatDo) Having(conds ...gen.Condition) ITLinkNetworkStatDo {
	return t.withDO(t.DO.Having(conds...))
}

func (t tLinkNetworkStatDo) Limit(limit int) ITLinkNetworkStatDo {
	return t.withDO(t.DO.Limit(limit))
}

func (t tLinkNetworkStatDo) Offset(offset int) ITLinkNetworkStatDo {
	return t.withDO(t.DO.Offset(offset))
}

func (t tLinkNetworkStatDo) Scopes(funcs ...func(gen.Dao) gen.Dao) ITLinkNetworkStatDo {
	return t.withDO(t.DO.Scopes(funcs...))
}

func (t tLinkNetworkStatDo) Unscoped() ITLinkNetworkStatDo {
	return t.withDO(t.DO.Unscoped())
}

func (t tLinkNetworkStatDo) Create(values ...*model.TLinkNetworkStat) error {
	if len(values) == 0 {
		return nil
	}
	return t.DO.Create(values)
}

func (t tLinkNetworkStatDo) CreateInBatches(values []*model.TLinkNetworkStat, batchSize int) error {
	return t.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (t tLinkNetworkStatDo) Save(values ...*model.TLinkNetworkStat) error {
	if len(values) == 0 {
		return nil
	}
	return t.DO.Save(values)
}

func (t tLinkNetworkStatDo) First() (*model.TLinkNetworkStat, error) {
	if result, err := t.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.TLinkNetworkStat), nil
	}
}

func (t tLinkNetworkStatDo) Take() (*model.TLinkNetworkStat, error) {
	if result, err := t.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.TLinkNetworkStat), nil
	}
}

func (t tLinkNetworkStatDo) Last() (*model.TLinkNetworkStat, error) {
	if result, err := t.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.TLinkNetworkStat), nil
	}
}

func (t tLinkNetworkStatDo) Find() ([]*model.TLinkNetworkStat, error) {
	result, err := t.DO.Find()
	return result.([]*model.TLinkNetworkStat), err
}

func (t tLinkNetworkStatDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.TLinkNetworkStat, err error) {
	buf := make([]*model.TLinkNetworkStat, 0, batchSize)
	err = t.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (t tLinkNetworkStatDo) FindInBatches(result *[]*model.TLinkNetworkStat, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return t.DO.FindInBatches(result, batchSize, fc)
}

func (t tLinkNetworkStatDo) Attrs(attrs ...field.AssignExpr) ITLinkNetworkStatDo {
	return t.withDO(t.DO.Attrs(attrs...))
}

func (t tLinkNetworkStatDo) Assign(attrs ...field.AssignExpr) ITLinkNetworkStatDo {
	return t.withDO(t.DO.Assign(attrs...))
}

func (t tLinkNetworkStatDo) Joins(fields ...field.RelationField) ITLinkNetworkStatDo {
	for _, _f := range fields {
		t = *t.withDO(t.DO.Joins(_f))
	}
	return &t
}

func (t tLinkNetworkStatDo) Preload(fields ...field.RelationField) ITLinkNetworkStatDo {
	for _, _f := range fields {
		t = *t.withDO(t.DO.Preload(_f))
	}
	return &t
}

func (t tLinkNetworkStatDo) FirstOrInit() (*model.TLinkNetworkStat, error) {
	if result, err := t.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.TLinkNetworkStat), nil
	}
}

func (t tLinkNetworkStatDo) FirstOrCreate() (*model.TLinkNetworkStat, error) {
	if result, err := t.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.TLinkNetworkStat), nil
	}
}

func (t tLinkNetworkStatDo) FindByPage(offset int, limit int) (result []*model.TLinkNetworkStat, count int64, err error) {
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

func (t tLinkNetworkStatDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = t.Count()
	if err != nil {
		return
	}

	err = t.Offset(offset).Limit(limit).Scan(result)
	return
}

func (t tLinkNetworkStatDo) Scan(result interface{}) (err error) {
	return t.DO.Scan(result)
}

func (t tLinkNetworkStatDo) Delete(models ...*model.TLinkNetworkStat) (result gen.ResultInfo, err error) {
	return t.DO.Delete(models)
}

func (t *tLinkNetworkStatDo) withDO(do gen.Dao) *tLinkNetworkStatDo {
	t.DO = *do.(*gen.DO)
	return t
}