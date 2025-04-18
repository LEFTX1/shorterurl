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

func newTLink(db *gorm.DB, opts ...gen.DOOption) tLink {
	_tLink := tLink{}

	_tLink.tLinkDo.UseDB(db, opts...)
	_tLink.tLinkDo.UseModel(&model.TLink{})

	tableName := _tLink.tLinkDo.TableName()
	_tLink.ALL = field.NewAsterisk(tableName)
	_tLink.ID = field.NewInt64(tableName, "id")
	_tLink.Domain = field.NewString(tableName, "domain")
	_tLink.ShortURI = field.NewString(tableName, "short_uri")
	_tLink.FullShortURL = field.NewString(tableName, "full_short_url")
	_tLink.OriginURL = field.NewString(tableName, "origin_url")
	_tLink.ClickNum = field.NewInt32(tableName, "click_num")
	_tLink.Gid = field.NewString(tableName, "gid")
	_tLink.Favicon = field.NewString(tableName, "favicon")
	_tLink.EnableStatus = field.NewBool(tableName, "enable_status")
	_tLink.CreatedType = field.NewBool(tableName, "created_type")
	_tLink.ValidDateType = field.NewBool(tableName, "valid_date_type")
	_tLink.ValidDate = field.NewTime(tableName, "valid_date")
	_tLink.Describe = field.NewString(tableName, "describe")
	_tLink.TotalPv = field.NewInt32(tableName, "total_pv")
	_tLink.TotalUv = field.NewInt32(tableName, "total_uv")
	_tLink.TotalUip = field.NewInt32(tableName, "total_uip")
	_tLink.CreateTime = field.NewTime(tableName, "create_time")
	_tLink.UpdateTime = field.NewTime(tableName, "update_time")
	_tLink.DelTime = field.NewInt64(tableName, "del_time")
	_tLink.DelFlag = field.NewBool(tableName, "del_flag")

	_tLink.fillFieldMap()

	return _tLink
}

type tLink struct {
	tLinkDo

	ALL           field.Asterisk
	ID            field.Int64  // ID
	Domain        field.String // 域名
	ShortURI      field.String // 短链接
	FullShortURL  field.String // 完整短链接
	OriginURL     field.String // 原始链接
	ClickNum      field.Int32  // 点击量
	Gid           field.String // 分组标识
	Favicon       field.String // 网站图标
	EnableStatus  field.Bool   // 启用标识 0：启用 1：未启用
	CreatedType   field.Bool   // 创建类型 0：接口创建 1：控制台创建
	ValidDateType field.Bool   // 有效期类型 0：永久有效 1：自定义
	ValidDate     field.Time   // 有效期
	Describe      field.String // 描述
	TotalPv       field.Int32  // 历史PV
	TotalUv       field.Int32  // 历史UV
	TotalUip      field.Int32  // 历史UIP
	CreateTime    field.Time   // 创建时间
	UpdateTime    field.Time   // 修改时间
	DelTime       field.Int64  // 删除时间戳
	DelFlag       field.Bool   // 删除标识 0：未删除 1：已删除

	fieldMap map[string]field.Expr
}

func (t tLink) Table(newTableName string) *tLink {
	t.tLinkDo.UseTable(newTableName)
	return t.updateTableName(newTableName)
}

func (t tLink) As(alias string) *tLink {
	t.tLinkDo.DO = *(t.tLinkDo.As(alias).(*gen.DO))
	return t.updateTableName(alias)
}

func (t *tLink) updateTableName(table string) *tLink {
	t.ALL = field.NewAsterisk(table)
	t.ID = field.NewInt64(table, "id")
	t.Domain = field.NewString(table, "domain")
	t.ShortURI = field.NewString(table, "short_uri")
	t.FullShortURL = field.NewString(table, "full_short_url")
	t.OriginURL = field.NewString(table, "origin_url")
	t.ClickNum = field.NewInt32(table, "click_num")
	t.Gid = field.NewString(table, "gid")
	t.Favicon = field.NewString(table, "favicon")
	t.EnableStatus = field.NewBool(table, "enable_status")
	t.CreatedType = field.NewBool(table, "created_type")
	t.ValidDateType = field.NewBool(table, "valid_date_type")
	t.ValidDate = field.NewTime(table, "valid_date")
	t.Describe = field.NewString(table, "describe")
	t.TotalPv = field.NewInt32(table, "total_pv")
	t.TotalUv = field.NewInt32(table, "total_uv")
	t.TotalUip = field.NewInt32(table, "total_uip")
	t.CreateTime = field.NewTime(table, "create_time")
	t.UpdateTime = field.NewTime(table, "update_time")
	t.DelTime = field.NewInt64(table, "del_time")
	t.DelFlag = field.NewBool(table, "del_flag")

	t.fillFieldMap()

	return t
}

func (t *tLink) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := t.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (t *tLink) fillFieldMap() {
	t.fieldMap = make(map[string]field.Expr, 20)
	t.fieldMap["id"] = t.ID
	t.fieldMap["domain"] = t.Domain
	t.fieldMap["short_uri"] = t.ShortURI
	t.fieldMap["full_short_url"] = t.FullShortURL
	t.fieldMap["origin_url"] = t.OriginURL
	t.fieldMap["click_num"] = t.ClickNum
	t.fieldMap["gid"] = t.Gid
	t.fieldMap["favicon"] = t.Favicon
	t.fieldMap["enable_status"] = t.EnableStatus
	t.fieldMap["created_type"] = t.CreatedType
	t.fieldMap["valid_date_type"] = t.ValidDateType
	t.fieldMap["valid_date"] = t.ValidDate
	t.fieldMap["describe"] = t.Describe
	t.fieldMap["total_pv"] = t.TotalPv
	t.fieldMap["total_uv"] = t.TotalUv
	t.fieldMap["total_uip"] = t.TotalUip
	t.fieldMap["create_time"] = t.CreateTime
	t.fieldMap["update_time"] = t.UpdateTime
	t.fieldMap["del_time"] = t.DelTime
	t.fieldMap["del_flag"] = t.DelFlag
}

func (t tLink) clone(db *gorm.DB) tLink {
	t.tLinkDo.ReplaceConnPool(db.Statement.ConnPool)
	return t
}

func (t tLink) replaceDB(db *gorm.DB) tLink {
	t.tLinkDo.ReplaceDB(db)
	return t
}

type tLinkDo struct{ gen.DO }

type ITLinkDo interface {
	gen.SubQuery
	Debug() ITLinkDo
	WithContext(ctx context.Context) ITLinkDo
	WithResult(fc func(tx gen.Dao)) gen.ResultInfo
	ReplaceDB(db *gorm.DB)
	ReadDB() ITLinkDo
	WriteDB() ITLinkDo
	As(alias string) gen.Dao
	Session(config *gorm.Session) ITLinkDo
	Columns(cols ...field.Expr) gen.Columns
	Clauses(conds ...clause.Expression) ITLinkDo
	Not(conds ...gen.Condition) ITLinkDo
	Or(conds ...gen.Condition) ITLinkDo
	Select(conds ...field.Expr) ITLinkDo
	Where(conds ...gen.Condition) ITLinkDo
	Order(conds ...field.Expr) ITLinkDo
	Distinct(cols ...field.Expr) ITLinkDo
	Omit(cols ...field.Expr) ITLinkDo
	Join(table schema.Tabler, on ...field.Expr) ITLinkDo
	LeftJoin(table schema.Tabler, on ...field.Expr) ITLinkDo
	RightJoin(table schema.Tabler, on ...field.Expr) ITLinkDo
	Group(cols ...field.Expr) ITLinkDo
	Having(conds ...gen.Condition) ITLinkDo
	Limit(limit int) ITLinkDo
	Offset(offset int) ITLinkDo
	Count() (count int64, err error)
	Scopes(funcs ...func(gen.Dao) gen.Dao) ITLinkDo
	Unscoped() ITLinkDo
	Create(values ...*model.TLink) error
	CreateInBatches(values []*model.TLink, batchSize int) error
	Save(values ...*model.TLink) error
	First() (*model.TLink, error)
	Take() (*model.TLink, error)
	Last() (*model.TLink, error)
	Find() ([]*model.TLink, error)
	FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.TLink, err error)
	FindInBatches(result *[]*model.TLink, batchSize int, fc func(tx gen.Dao, batch int) error) error
	Pluck(column field.Expr, dest interface{}) error
	Delete(...*model.TLink) (info gen.ResultInfo, err error)
	Update(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	Updates(value interface{}) (info gen.ResultInfo, err error)
	UpdateColumn(column field.Expr, value interface{}) (info gen.ResultInfo, err error)
	UpdateColumnSimple(columns ...field.AssignExpr) (info gen.ResultInfo, err error)
	UpdateColumns(value interface{}) (info gen.ResultInfo, err error)
	UpdateFrom(q gen.SubQuery) gen.Dao
	Attrs(attrs ...field.AssignExpr) ITLinkDo
	Assign(attrs ...field.AssignExpr) ITLinkDo
	Joins(fields ...field.RelationField) ITLinkDo
	Preload(fields ...field.RelationField) ITLinkDo
	FirstOrInit() (*model.TLink, error)
	FirstOrCreate() (*model.TLink, error)
	FindByPage(offset int, limit int) (result []*model.TLink, count int64, err error)
	ScanByPage(result interface{}, offset int, limit int) (count int64, err error)
	Scan(result interface{}) (err error)
	Returning(value interface{}, columns ...string) ITLinkDo
	UnderlyingDB() *gorm.DB
	schema.Tabler
}

func (t tLinkDo) Debug() ITLinkDo {
	return t.withDO(t.DO.Debug())
}

func (t tLinkDo) WithContext(ctx context.Context) ITLinkDo {
	return t.withDO(t.DO.WithContext(ctx))
}

func (t tLinkDo) ReadDB() ITLinkDo {
	return t.Clauses(dbresolver.Read)
}

func (t tLinkDo) WriteDB() ITLinkDo {
	return t.Clauses(dbresolver.Write)
}

func (t tLinkDo) Session(config *gorm.Session) ITLinkDo {
	return t.withDO(t.DO.Session(config))
}

func (t tLinkDo) Clauses(conds ...clause.Expression) ITLinkDo {
	return t.withDO(t.DO.Clauses(conds...))
}

func (t tLinkDo) Returning(value interface{}, columns ...string) ITLinkDo {
	return t.withDO(t.DO.Returning(value, columns...))
}

func (t tLinkDo) Not(conds ...gen.Condition) ITLinkDo {
	return t.withDO(t.DO.Not(conds...))
}

func (t tLinkDo) Or(conds ...gen.Condition) ITLinkDo {
	return t.withDO(t.DO.Or(conds...))
}

func (t tLinkDo) Select(conds ...field.Expr) ITLinkDo {
	return t.withDO(t.DO.Select(conds...))
}

func (t tLinkDo) Where(conds ...gen.Condition) ITLinkDo {
	return t.withDO(t.DO.Where(conds...))
}

func (t tLinkDo) Order(conds ...field.Expr) ITLinkDo {
	return t.withDO(t.DO.Order(conds...))
}

func (t tLinkDo) Distinct(cols ...field.Expr) ITLinkDo {
	return t.withDO(t.DO.Distinct(cols...))
}

func (t tLinkDo) Omit(cols ...field.Expr) ITLinkDo {
	return t.withDO(t.DO.Omit(cols...))
}

func (t tLinkDo) Join(table schema.Tabler, on ...field.Expr) ITLinkDo {
	return t.withDO(t.DO.Join(table, on...))
}

func (t tLinkDo) LeftJoin(table schema.Tabler, on ...field.Expr) ITLinkDo {
	return t.withDO(t.DO.LeftJoin(table, on...))
}

func (t tLinkDo) RightJoin(table schema.Tabler, on ...field.Expr) ITLinkDo {
	return t.withDO(t.DO.RightJoin(table, on...))
}

func (t tLinkDo) Group(cols ...field.Expr) ITLinkDo {
	return t.withDO(t.DO.Group(cols...))
}

func (t tLinkDo) Having(conds ...gen.Condition) ITLinkDo {
	return t.withDO(t.DO.Having(conds...))
}

func (t tLinkDo) Limit(limit int) ITLinkDo {
	return t.withDO(t.DO.Limit(limit))
}

func (t tLinkDo) Offset(offset int) ITLinkDo {
	return t.withDO(t.DO.Offset(offset))
}

func (t tLinkDo) Scopes(funcs ...func(gen.Dao) gen.Dao) ITLinkDo {
	return t.withDO(t.DO.Scopes(funcs...))
}

func (t tLinkDo) Unscoped() ITLinkDo {
	return t.withDO(t.DO.Unscoped())
}

func (t tLinkDo) Create(values ...*model.TLink) error {
	if len(values) == 0 {
		return nil
	}
	return t.DO.Create(values)
}

func (t tLinkDo) CreateInBatches(values []*model.TLink, batchSize int) error {
	return t.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (t tLinkDo) Save(values ...*model.TLink) error {
	if len(values) == 0 {
		return nil
	}
	return t.DO.Save(values)
}

func (t tLinkDo) First() (*model.TLink, error) {
	if result, err := t.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*model.TLink), nil
	}
}

func (t tLinkDo) Take() (*model.TLink, error) {
	if result, err := t.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*model.TLink), nil
	}
}

func (t tLinkDo) Last() (*model.TLink, error) {
	if result, err := t.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*model.TLink), nil
	}
}

func (t tLinkDo) Find() ([]*model.TLink, error) {
	result, err := t.DO.Find()
	return result.([]*model.TLink), err
}

func (t tLinkDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*model.TLink, err error) {
	buf := make([]*model.TLink, 0, batchSize)
	err = t.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (t tLinkDo) FindInBatches(result *[]*model.TLink, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return t.DO.FindInBatches(result, batchSize, fc)
}

func (t tLinkDo) Attrs(attrs ...field.AssignExpr) ITLinkDo {
	return t.withDO(t.DO.Attrs(attrs...))
}

func (t tLinkDo) Assign(attrs ...field.AssignExpr) ITLinkDo {
	return t.withDO(t.DO.Assign(attrs...))
}

func (t tLinkDo) Joins(fields ...field.RelationField) ITLinkDo {
	for _, _f := range fields {
		t = *t.withDO(t.DO.Joins(_f))
	}
	return &t
}

func (t tLinkDo) Preload(fields ...field.RelationField) ITLinkDo {
	for _, _f := range fields {
		t = *t.withDO(t.DO.Preload(_f))
	}
	return &t
}

func (t tLinkDo) FirstOrInit() (*model.TLink, error) {
	if result, err := t.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*model.TLink), nil
	}
}

func (t tLinkDo) FirstOrCreate() (*model.TLink, error) {
	if result, err := t.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*model.TLink), nil
	}
}

func (t tLinkDo) FindByPage(offset int, limit int) (result []*model.TLink, count int64, err error) {
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

func (t tLinkDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = t.Count()
	if err != nil {
		return
	}

	err = t.Offset(offset).Limit(limit).Scan(result)
	return
}

func (t tLinkDo) Scan(result interface{}) (err error) {
	return t.DO.Scan(result)
}

func (t tLinkDo) Delete(models ...*model.TLink) (result gen.ResultInfo, err error) {
	return t.DO.Delete(models)
}

func (t *tLinkDo) withDO(do gen.Dao) *tLinkDo {
	t.DO = *do.(*gen.DO)
	return t
}
