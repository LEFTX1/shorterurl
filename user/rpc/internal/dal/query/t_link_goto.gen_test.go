// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"
	"fmt"
	"testing"

	"shorterurl/user/rpc/internal/dal/model"

	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm/clause"
)

func init() {
	InitializeDB()
	err := _gen_test_db.AutoMigrate(&model.TLinkGoto{})
	if err != nil {
		fmt.Printf("Error: AutoMigrate(&model.TLinkGoto{}) fail: %s", err)
	}
}

func Test_tLinkGotoQuery(t *testing.T) {
	tLinkGoto := newTLinkGoto(_gen_test_db)
	tLinkGoto = *tLinkGoto.As(tLinkGoto.TableName())
	_do := tLinkGoto.WithContext(context.Background()).Debug()

	primaryKey := field.NewString(tLinkGoto.TableName(), clause.PrimaryKey)
	_, err := _do.Unscoped().Where(primaryKey.IsNotNull()).Delete()
	if err != nil {
		t.Error("clean table <t_link_goto> fail:", err)
		return
	}

	_, ok := tLinkGoto.GetFieldByName("")
	if ok {
		t.Error("GetFieldByName(\"\") from tLinkGoto success")
	}

	err = _do.Create(&model.TLinkGoto{})
	if err != nil {
		t.Error("create item in table <t_link_goto> fail:", err)
	}

	err = _do.Save(&model.TLinkGoto{})
	if err != nil {
		t.Error("create item in table <t_link_goto> fail:", err)
	}

	err = _do.CreateInBatches([]*model.TLinkGoto{{}, {}}, 10)
	if err != nil {
		t.Error("create item in table <t_link_goto> fail:", err)
	}

	_, err = _do.Select(tLinkGoto.ALL).Take()
	if err != nil {
		t.Error("Take() on table <t_link_goto> fail:", err)
	}

	_, err = _do.First()
	if err != nil {
		t.Error("First() on table <t_link_goto> fail:", err)
	}

	_, err = _do.Last()
	if err != nil {
		t.Error("First() on table <t_link_goto> fail:", err)
	}

	_, err = _do.Where(primaryKey.IsNotNull()).FindInBatch(10, func(tx gen.Dao, batch int) error { return nil })
	if err != nil {
		t.Error("FindInBatch() on table <t_link_goto> fail:", err)
	}

	err = _do.Where(primaryKey.IsNotNull()).FindInBatches(&[]*model.TLinkGoto{}, 10, func(tx gen.Dao, batch int) error { return nil })
	if err != nil {
		t.Error("FindInBatches() on table <t_link_goto> fail:", err)
	}

	_, err = _do.Select(tLinkGoto.ALL).Where(primaryKey.IsNotNull()).Order(primaryKey.Desc()).Find()
	if err != nil {
		t.Error("Find() on table <t_link_goto> fail:", err)
	}

	_, err = _do.Distinct(primaryKey).Take()
	if err != nil {
		t.Error("select Distinct() on table <t_link_goto> fail:", err)
	}

	_, err = _do.Select(tLinkGoto.ALL).Omit(primaryKey).Take()
	if err != nil {
		t.Error("Omit() on table <t_link_goto> fail:", err)
	}

	_, err = _do.Group(primaryKey).Find()
	if err != nil {
		t.Error("Group() on table <t_link_goto> fail:", err)
	}

	_, err = _do.Scopes(func(dao gen.Dao) gen.Dao { return dao.Where(primaryKey.IsNotNull()) }).Find()
	if err != nil {
		t.Error("Scopes() on table <t_link_goto> fail:", err)
	}

	_, _, err = _do.FindByPage(0, 1)
	if err != nil {
		t.Error("FindByPage() on table <t_link_goto> fail:", err)
	}

	_, err = _do.ScanByPage(&model.TLinkGoto{}, 0, 1)
	if err != nil {
		t.Error("ScanByPage() on table <t_link_goto> fail:", err)
	}

	_, err = _do.Attrs(primaryKey).Assign(primaryKey).FirstOrInit()
	if err != nil {
		t.Error("FirstOrInit() on table <t_link_goto> fail:", err)
	}

	_, err = _do.Attrs(primaryKey).Assign(primaryKey).FirstOrCreate()
	if err != nil {
		t.Error("FirstOrCreate() on table <t_link_goto> fail:", err)
	}

	var _a _another
	var _aPK = field.NewString(_a.TableName(), "id")

	err = _do.Join(&_a, primaryKey.EqCol(_aPK)).Scan(map[string]interface{}{})
	if err != nil {
		t.Error("Join() on table <t_link_goto> fail:", err)
	}

	err = _do.LeftJoin(&_a, primaryKey.EqCol(_aPK)).Scan(map[string]interface{}{})
	if err != nil {
		t.Error("LeftJoin() on table <t_link_goto> fail:", err)
	}

	_, err = _do.Not().Or().Clauses().Take()
	if err != nil {
		t.Error("Not/Or/Clauses on table <t_link_goto> fail:", err)
	}
}
