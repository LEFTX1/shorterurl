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

	"shorterurl/user/rpc/internal/dal/model"
)

func newTGroup(db *gorm.DB, opts ...gen.DOOption) tGroup {
	_tGroup := tGroup{}

	_tGroup.tGroupDo.UseDB(db, opts...)
	_tGroup.tGroupDo.UseModel(&model.TGroup{})

	tableName := _tGroup.tGroupDo.TableName()
	_tGroup.ALL = field.NewAsterisk(tableName)
	_tGroup.ID = field.NewInt64(tableName, "id")
	_tGroup.Gid = field.NewString(tableName, "gid")
	_tGroup.Name = field.NewString(tableName, "name")
	_tGroup.Username = field.NewString(tableName, "username")
	_tGroup.SortOrder = field.NewInt32(tableName, "sort_order")
	_tGroup.CreateTime = field.NewTime(tableName, "create_time")
	_tGroup.UpdateTime = field.NewTime(tableName, "update_time")
	_tGroup.DelFlag = field.NewBool(tableName, "del_flag")

	_tGroup.fillFieldMap()

	return _tGroup
}

type tGroup struct {
	tGroupDo

	ALL        field.Asterisk
	ID         field.Int64  // ID
	Gid        field.String // 分组标识
	Name       field.String // 分组名称
	Username   field.String // 创建分组用户名
	SortOrder  field.Int32  // 分组排序
	CreateTime field.Time   // 创建时间
	UpdateTime field.Time   // 修改时间
	DelFlag    field.Bool   // 删除标识 0：未删除 1：已删除

	fieldMap map[string]field.Expr
}

func (t tGroup) Table(newTableName string) *tGroup {
	t.tGroupDo.UseTable(newTableName)
	return t.updateTableName(newTableName)
}

func (t tGroup) As(alias string) *tGroup {
	t.tGroupDo.DO = *(t.tGroupDo.As(alias).(*gen.DO))
	return t.updateTableName(alias)
}

func (t *tGroup) updateTableName(table string) *tGroup {
	t.ALL = field.NewAsterisk(table)
	t.ID = field.NewInt64(table, "id")
	t.Gid = field.NewString(table, "gid")
	t.Name = field.NewString(table, "name")
	t.Username = field.NewString(table, "username")
	t.SortOrder = field.NewInt32(table, "sort_order")
	t.CreateTime = field.NewTime(table, "create_time")
	t.UpdateTime = field.NewTime(table, "update_time")
	t.DelFlag = field.NewBool(table, "del_flag")

	t.fillFieldMap()

	return t
}

func (t *tGroup) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := t.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (t *tGroup) fillFieldMap() {
	t.fieldMap = make(map[string]field.Expr, 8)
	t.fieldMap["id"] = t.ID
	t.fieldMap["gid"] = t.Gid
	t.fieldMap["name"] = t.Name
	t.fieldMap["username"] = t.Username
	t.fieldMap["sort_order"] = t.SortOrder
	t.fieldMap["create_time"] = t.CreateTime
	t.fieldMap["update_time"] = t.UpdateTime
	t.fieldMap["del_flag"] = t.DelFlag
}

func (t tGroup) clone(db *gorm.DB) tGroup {
	t.tGroupDo.ReplaceConnPool(db.Statement.ConnPool)
	return t
}

func (t tGroup) replaceDB(db *gorm.DB) tGroup {
	t.tGroupDo.ReplaceDB(db)
	return t
}

type tGroupDo struct{ gen.DO }

type ITGroupDo interface {
	gen.SubQuery
	Debug() ITGroupDo
	WithContext(ctx context.Context) ITGroupDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() ITGroupDo
	WriteDB() ITGroupDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) ITGroupDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) ITGroupDo
	Not(conds ...gen.Condition) ITGroupDo
	Or(conds ...gen.Condition) ITGroupDo
	Select(conds ...field.Expr) ITGroupDo
	Where(conds ...gen.Condition) ITGroupDo
	Order(conds ...field.Expr) ITGroupDo
	Distinct(cols ...field.Expr) ITGroupDo
	Omit(cols ...field.Expr) ITGroupDo
	Join(table schema.Tabler, on ...field.Expr) ITGroupDo
	LeftJoin(table schema.Tabler, on ...field.Expr) ITGroupDo
	RightJoin(table schema.Tabler, on ...field.Expr) ITGroupDo
	Group(cols ...field.Expr) ITGroupDo
	Having(conds ...gen.Condition) ITGroupDo
	Limit(limit int) ITGroupDo
	Offset(offset int) ITGroupDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) ITGroupDo
	Unscoped() ITGroupDo
	Create(values ...*model.TGroup) error
	CreateInBatches(values []*model.TGroup, batchSize int) error
	Save(values ...*model.TGroup) error
	First() (*model.TGroup, error)
	Take() (*model.TGroup, error)
	Last() (*model.TGroup, error)
	Find() ([]*model.TGroup, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.TGroup, err error)
	FindInBatches(result *[]*model.TGroup, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.TGroup) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) ITGroupDo
	Assign(attrs ...field.AssignExpr) ITGroupDo
	Joins(fields ...field.RelationField) ITGroupDo
	Preload(fields ...field.RelationField) ITGroupDo
	FirstOrInit() (*model.TGroup, error)
	FirstOrCreate() (*model.TGroup, error)
	FindByPage(offset int, limit int) (result []*model.TGroup, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) ITGroupDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (t tGroupDo) Debug() ITGroupDo {
	return t.withDO(t.DO.Debug())
}

func (t tGroupDo) WithContext(ctx context.Context) ITGroupDo {
	return t.withDO(t.DO.WithContext(ctx))
}

func (t tGroupDo) ReadDB() ITGroupDo {
	return t.Clauses(dbresolver.Read)
}

func (t tGroupDo) WriteDB() ITGroupDo {
	return t.Clauses(dbresolver.Write)
}

func (t tGroupDo) Session(config *gorm.Session) ITGroupDo {
	return t.withDO(t.DO.Session(config))
}

func (t tGroupDo) Clauses(conds ...clause.Expression) ITGroupDo {
	return t.withDO(t.DO.Clauses(conds...))
}

func (t tGroupDo) Returning(value interface{}, columns ...string) ITGroupDo {
	return t.withDO(t.DO.Returning(value, columns...))
}

func (t tGroupDo) Not(conds ...gen.Condition) ITGroupDo {
	return t.withDO(t.DO.Not(conds...))
}

func (t tGroupDo) Or(conds ...gen.Condition) ITGroupDo {
	return t.withDO(t.DO.Or(conds...))
}

func (t tGroupDo) Select(conds ...field.Expr) ITGroupDo {
	return t.withDO(t.DO.Select(conds...))
}

func (t tGroupDo) Where(conds ...gen.Condition) ITGroupDo {
	return t.withDO(t.DO.Where(conds...))
}

func (t tGroupDo) Order(conds ...field.Expr) ITGroupDo {
	return t.withDO(t.DO.Order(conds...))
}

func (t tGroupDo) Distinct(cols ...field.Expr) ITGroupDo {
	return t.withDO(t.DO.Distinct(cols...))
}

func (t tGroupDo) Omit(cols ...field.Expr) ITGroupDo {
	return t.withDO(t.DO.Omit(cols...))
}

func (t tGroupDo) Join(table schema.Tabler, on ...field.Expr) ITGroupDo {
	return t.withDO(t.DO.Join(table, on...))
}

func (t tGroupDo) LeftJoin(table schema.Tabler, on ...field.Expr) ITGroupDo {
	return t.withDO(t.DO.LeftJoin(table, on...))
}

func (t tGroupDo) RightJoin(table schema.Tabler, on ...field.Expr) ITGroupDo {
	return t.withDO(t.DO.RightJoin(table, on...))
}

func (t tGroupDo) Group(cols ...field.Expr) ITGroupDo {
	return t.withDO(t.DO.Group(cols...))
}

func (t tGroupDo) Having(conds ...gen.Condition) ITGroupDo {
	return t.withDO(t.DO.Having(conds...))
}

func (t tGroupDo) Limit(limit int) ITGroupDo {
	return t.withDO(t.DO.Limit(limit))
}

func (t tGroupDo) Offset(offset int) ITGroupDo {
	return t.withDO(t.DO.Offset(offset))
}

func (t tGroupDo) Scopes(funcs ...func(gen.Dao) gen.Dao) ITGroupDo {
	return t.withDO(t.DO.Scopes(funcs...))
}

func (t tGroupDo) Unscoped() ITGroupDo {
	return t.withDO(t.DO.Unscoped())
}

func (t tGroupDo) Create(values ...*model.TGroup) error {
	if len(values) == 0 {
		return nil
	}
	return t.DO.Create(values)
}

func (t tGroupDo) CreateInBatches(values []*model.TGroup, batchSize int) error {
	return t.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (t tGroupDo) Save(values ...*model.TGroup) error {
	if len(values) == 0 {
		return nil
	}
	return t.DO.Save(values)
}

func (t tGroupDo) First() (*model.TGroup, error) {
	if result, err := t.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.TGroup), nil
	}
}

func (t tGroupDo) Take() (*model.TGroup, error) {
	if result, err := t.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.TGroup), nil
	}
}

func (t tGroupDo) Last() (*model.TGroup, error) {
	if result, err := t.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.TGroup), nil
	}
}

func (t tGroupDo) Find() ([]*model.TGroup, error) {
	result, err := t.DO.Find()
	return result.([]*model.TGroup), err
}

func (t tGroupDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.TGroup, err error) {
	buf := make([]*model.TGroup, 0, batchSize)
	err = t.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (t tGroupDo) FindInBatches(result *[]*model.TGroup, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return t.DO.FindInBatches(result, batchSize, fc)
}

func (t tGroupDo) Attrs(attrs ...field.AssignExpr) ITGroupDo {
	return t.withDO(t.DO.Attrs(attrs...))
}

func (t tGroupDo) Assign(attrs ...field.AssignExpr) ITGroupDo {
	return t.withDO(t.DO.Assign(attrs...))
}

func (t tGroupDo) Joins(fields ...field.RelationField) ITGroupDo {
	for _, _f := range fields {
		t = *t.withDO(t.DO.Joins(_f))
	}
	return &t
}

func (t tGroupDo) Preload(fields ...field.RelationField) ITGroupDo {
	for _, _f := range fields {
		t = *t.withDO(t.DO.Preload(_f))
	}
	return &t
}

func (t tGroupDo) FirstOrInit() (*model.TGroup, error) {
	if result, err := t.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.TGroup), nil
	}
}

func (t tGroupDo) FirstOrCreate() (*model.TGroup, error) {
	if result, err := t.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.TGroup), nil
	}
}

func (t tGroupDo) FindByPage(offset int, limit int) (result []*model.TGroup, count int64, err error) {
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

func (t tGroupDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = t.Count()
	if err != nil {
		return
	}

	err = t.Offset(offset).Limit(limit).Scan(result)
	return
}

func (t tGroupDo) Scan(result interface{}) (err error) {
	return t.DO.Scan(result)
}

func (t tGroupDo) Delete(models ...*model.TGroup) (result gen.ResultInfo, err error) {
	return t.DO.Delete(models)
}

func (t *tGroupDo) withDO(do gen.Dao) *tGroupDo {
	t.DO = *do.(*gen.DO)
	return t
}
