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

func newTLinkGoto(db *gorm.DB, opts ...gen.DOOption) tLinkGoto {
	_tLinkGoto := tLinkGoto{}

	_tLinkGoto.tLinkGotoDo.UseDB(db, opts...)
	_tLinkGoto.tLinkGotoDo.UseModel(&model.TLinkGoto{})

	tableName := _tLinkGoto.tLinkGotoDo.TableName()
	_tLinkGoto.ALL = field.NewAsterisk(tableName)
	_tLinkGoto.ID = field.NewInt64(tableName, "id")
	_tLinkGoto.Gid = field.NewString(tableName, "gid")
	_tLinkGoto.FullShortURL = field.NewString(tableName, "full_short_url")

	_tLinkGoto.fillFieldMap()

	return _tLinkGoto
}

type tLinkGoto struct {
	tLinkGotoDo

	ALL          field.Asterisk
	ID           field.Int64  // ID
	Gid          field.String // 分组标识
	FullShortURL field.String // 完整短链接

	fieldMap map[string]field.Expr
}

func (t tLinkGoto) Table(newTableName string) *tLinkGoto {
	t.tLinkGotoDo.UseTable(newTableName)
	return t.updateTableName(newTableName)
}

func (t tLinkGoto) As(alias string) *tLinkGoto {
	t.tLinkGotoDo.DO = *(t.tLinkGotoDo.As(alias).(*gen.DO))
	return t.updateTableName(alias)
}

func (t *tLinkGoto) updateTableName(table string) *tLinkGoto {
	t.ALL = field.NewAsterisk(table)
	t.ID = field.NewInt64(table, "id")
	t.Gid = field.NewString(table, "gid")
	t.FullShortURL = field.NewString(table, "full_short_url")

	t.fillFieldMap()

	return t
}

func (t *tLinkGoto) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := t.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (t *tLinkGoto) fillFieldMap() {
	t.fieldMap = make(map[string]field.Expr, 3)
	t.fieldMap["id"] = t.ID
	t.fieldMap["gid"] = t.Gid
	t.fieldMap["full_short_url"] = t.FullShortURL
}

func (t tLinkGoto) clone(db *gorm.DB) tLinkGoto {
	t.tLinkGotoDo.ReplaceConnPool(db.Statement.ConnPool)
	return t
}

func (t tLinkGoto) replaceDB(db *gorm.DB) tLinkGoto {
	t.tLinkGotoDo.ReplaceDB(db)
	return t
}

type tLinkGotoDo struct{ gen.DO }

type ITLinkGotoDo interface {
	gen.SubQuery
	Debug() ITLinkGotoDo
	WithContext(ctx context.Context) ITLinkGotoDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() ITLinkGotoDo
	WriteDB() ITLinkGotoDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) ITLinkGotoDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) ITLinkGotoDo
	Not(conds ...gen.Condition) ITLinkGotoDo
	Or(conds ...gen.Condition) ITLinkGotoDo
	Select(conds ...field.Expr) ITLinkGotoDo
	Where(conds ...gen.Condition) ITLinkGotoDo
	Order(conds ...field.Expr) ITLinkGotoDo
	Distinct(cols ...field.Expr) ITLinkGotoDo
	Omit(cols ...field.Expr) ITLinkGotoDo
	Join(table schema.Tabler, on ...field.Expr) ITLinkGotoDo
	LeftJoin(table schema.Tabler, on ...field.Expr) ITLinkGotoDo
	RightJoin(table schema.Tabler, on ...field.Expr) ITLinkGotoDo
	Group(cols ...field.Expr) ITLinkGotoDo
	Having(conds ...gen.Condition) ITLinkGotoDo
	Limit(limit int) ITLinkGotoDo
	Offset(offset int) ITLinkGotoDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) ITLinkGotoDo
	Unscoped() ITLinkGotoDo
	Create(values ...*model.TLinkGoto) error
	CreateInBatches(values []*model.TLinkGoto, batchSize int) error
	Save(values ...*model.TLinkGoto) error
	First() (*model.TLinkGoto, error)
	Take() (*model.TLinkGoto, error)
	Last() (*model.TLinkGoto, error)
	Find() ([]*model.TLinkGoto, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.TLinkGoto, err error)
	FindInBatches(result *[]*model.TLinkGoto, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.TLinkGoto) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) ITLinkGotoDo
	Assign(attrs ...field.AssignExpr) ITLinkGotoDo
	Joins(fields ...field.RelationField) ITLinkGotoDo
	Preload(fields ...field.RelationField) ITLinkGotoDo
	FirstOrInit() (*model.TLinkGoto, error)
	FirstOrCreate() (*model.TLinkGoto, error)
	FindByPage(offset int, limit int) (result []*model.TLinkGoto, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) ITLinkGotoDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (t tLinkGotoDo) Debug() ITLinkGotoDo {
	return t.withDO(t.DO.Debug())
}

func (t tLinkGotoDo) WithContext(ctx context.Context) ITLinkGotoDo {
	return t.withDO(t.DO.WithContext(ctx))
}

func (t tLinkGotoDo) ReadDB() ITLinkGotoDo {
	return t.Clauses(dbresolver.Read)
}

func (t tLinkGotoDo) WriteDB() ITLinkGotoDo {
	return t.Clauses(dbresolver.Write)
}

func (t tLinkGotoDo) Session(config *gorm.Session) ITLinkGotoDo {
	return t.withDO(t.DO.Session(config))
}

func (t tLinkGotoDo) Clauses(conds ...clause.Expression) ITLinkGotoDo {
	return t.withDO(t.DO.Clauses(conds...))
}

func (t tLinkGotoDo) Returning(value interface{}, columns ...string) ITLinkGotoDo {
	return t.withDO(t.DO.Returning(value, columns...))
}

func (t tLinkGotoDo) Not(conds ...gen.Condition) ITLinkGotoDo {
	return t.withDO(t.DO.Not(conds...))
}

func (t tLinkGotoDo) Or(conds ...gen.Condition) ITLinkGotoDo {
	return t.withDO(t.DO.Or(conds...))
}

func (t tLinkGotoDo) Select(conds ...field.Expr) ITLinkGotoDo {
	return t.withDO(t.DO.Select(conds...))
}

func (t tLinkGotoDo) Where(conds ...gen.Condition) ITLinkGotoDo {
	return t.withDO(t.DO.Where(conds...))
}

func (t tLinkGotoDo) Order(conds ...field.Expr) ITLinkGotoDo {
	return t.withDO(t.DO.Order(conds...))
}

func (t tLinkGotoDo) Distinct(cols ...field.Expr) ITLinkGotoDo {
	return t.withDO(t.DO.Distinct(cols...))
}

func (t tLinkGotoDo) Omit(cols ...field.Expr) ITLinkGotoDo {
	return t.withDO(t.DO.Omit(cols...))
}

func (t tLinkGotoDo) Join(table schema.Tabler, on ...field.Expr) ITLinkGotoDo {
	return t.withDO(t.DO.Join(table, on...))
}

func (t tLinkGotoDo) LeftJoin(table schema.Tabler, on ...field.Expr) ITLinkGotoDo {
	return t.withDO(t.DO.LeftJoin(table, on...))
}

func (t tLinkGotoDo) RightJoin(table schema.Tabler, on ...field.Expr) ITLinkGotoDo {
	return t.withDO(t.DO.RightJoin(table, on...))
}

func (t tLinkGotoDo) Group(cols ...field.Expr) ITLinkGotoDo {
	return t.withDO(t.DO.Group(cols...))
}

func (t tLinkGotoDo) Having(conds ...gen.Condition) ITLinkGotoDo {
	return t.withDO(t.DO.Having(conds...))
}

func (t tLinkGotoDo) Limit(limit int) ITLinkGotoDo {
	return t.withDO(t.DO.Limit(limit))
}

func (t tLinkGotoDo) Offset(offset int) ITLinkGotoDo {
	return t.withDO(t.DO.Offset(offset))
}

func (t tLinkGotoDo) Scopes(funcs ...func(gen.Dao) gen.Dao) ITLinkGotoDo {
	return t.withDO(t.DO.Scopes(funcs...))
}

func (t tLinkGotoDo) Unscoped() ITLinkGotoDo {
	return t.withDO(t.DO.Unscoped())
}

func (t tLinkGotoDo) Create(values ...*model.TLinkGoto) error {
	if len(values) == 0 {
		return nil
	}
	return t.DO.Create(values)
}

func (t tLinkGotoDo) CreateInBatches(values []*model.TLinkGoto, batchSize int) error {
	return t.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (t tLinkGotoDo) Save(values ...*model.TLinkGoto) error {
	if len(values) == 0 {
		return nil
	}
	return t.DO.Save(values)
}

func (t tLinkGotoDo) First() (*model.TLinkGoto, error) {
	if result, err := t.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.TLinkGoto), nil
	}
}

func (t tLinkGotoDo) Take() (*model.TLinkGoto, error) {
	if result, err := t.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.TLinkGoto), nil
	}
}

func (t tLinkGotoDo) Last() (*model.TLinkGoto, error) {
	if result, err := t.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.TLinkGoto), nil
	}
}

func (t tLinkGotoDo) Find() ([]*model.TLinkGoto, error) {
	result, err := t.DO.Find()
	return result.([]*model.TLinkGoto), err
}

func (t tLinkGotoDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.TLinkGoto, err error) {
	buf := make([]*model.TLinkGoto, 0, batchSize)
	err = t.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (t tLinkGotoDo) FindInBatches(result *[]*model.TLinkGoto, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return t.DO.FindInBatches(result, batchSize, fc)
}

func (t tLinkGotoDo) Attrs(attrs ...field.AssignExpr) ITLinkGotoDo {
	return t.withDO(t.DO.Attrs(attrs...))
}

func (t tLinkGotoDo) Assign(attrs ...field.AssignExpr) ITLinkGotoDo {
	return t.withDO(t.DO.Assign(attrs...))
}

func (t tLinkGotoDo) Joins(fields ...field.RelationField) ITLinkGotoDo {
	for _, _f := range fields {
		t = *t.withDO(t.DO.Joins(_f))
	}
	return &t
}

func (t tLinkGotoDo) Preload(fields ...field.RelationField) ITLinkGotoDo {
	for _, _f := range fields {
		t = *t.withDO(t.DO.Preload(_f))
	}
	return &t
}

func (t tLinkGotoDo) FirstOrInit() (*model.TLinkGoto, error) {
	if result, err := t.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.TLinkGoto), nil
	}
}

func (t tLinkGotoDo) FirstOrCreate() (*model.TLinkGoto, error) {
	if result, err := t.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.TLinkGoto), nil
	}
}

func (t tLinkGotoDo) FindByPage(offset int, limit int) (result []*model.TLinkGoto, count int64, err error) {
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

func (t tLinkGotoDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = t.Count()
	if err != nil {
		return
	}

	err = t.Offset(offset).Limit(limit).Scan(result)
	return
}

func (t tLinkGotoDo) Scan(result interface{}) (err error) {
	return t.DO.Scan(result)
}

func (t tLinkGotoDo) Delete(models ...*model.TLinkGoto) (result gen.ResultInfo, err error) {
	return t.DO.Delete(models)
}

func (t *tLinkGotoDo) withDO(do gen.Dao) *tLinkGotoDo {
	t.DO = *do.(*gen.DO)
	return t
}
