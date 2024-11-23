package testutil

import (
	"github.com/longbridgeapp/sqlparser"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

// TestReplaceTablePrefixInWhere tests the replaceTablePrefixInWhere function
func TestReplaceTablePrefixInWhere(t *testing.T) {
	oldTableName := "orders"
	newTableName := "orders_1"

	tests := []struct {
		name          string
		query         string
		expectedWhere string
		description   string
	}{
		{
			name:          "Simple WHERE condition with table prefix",
			query:         "SELECT * FROM orders WHERE orders.user_id = 1",
			expectedWhere: "orders_1.user_id = 1",
			description:   "简单 WHERE 子句，带表前缀的字段应替换为新表名前缀",
		},
		{
			name:          "Complex WHERE condition with multiple table prefixes",
			query:         "SELECT * FROM orders WHERE orders.user_id = 1 AND orders.status = 'active'",
			expectedWhere: "orders_1.user_id = 1 AND orders_1.status = 'active'",
			description:   "复杂 WHERE 子句，多个字段均替换表前缀",
		},
		{
			name:          "WHERE condition with parentheses",
			query:         "SELECT * FROM orders WHERE (orders.user_id = 1 OR orders.status = 'active')",
			expectedWhere: "(orders_1.user_id = 1 OR orders_1.status = 'active')",
			description:   "括号嵌套条件，字段表前缀替换保持逻辑结构",
		},
		{
			name:          "No table prefix in WHERE condition",
			query:         "SELECT * FROM orders WHERE status = 'active'",
			expectedWhere: "status = 'active'",
			description:   "没有表前缀的字段应保持不变",
		},
		{
			name:          "Nested parentheses with multiple table prefixes",
			query:         "SELECT * FROM orders WHERE ((orders.user_id = 1 AND orders.status = 'active') OR (orders.price > 100))",
			expectedWhere: "((orders_1.user_id = 1 AND orders_1.status = 'active') OR (orders_1.price > 100))",
			description:   "嵌套括号中多个字段的表前缀应全部替换",
		},
		{
			name:          "OR conditions without table prefix",
			query:         "SELECT * FROM orders WHERE status = 'active' OR orders.user_id = 1",
			expectedWhere: "status = 'active' OR orders_1.user_id = 1",
			description:   "部分字段没有表前缀，保持不变，其他字段替换",
		},
		{
			name:          "Field comparison between two tables",
			query:         "SELECT * FROM orders WHERE orders.user_id = customers.id",
			expectedWhere: "orders_1.user_id = customers.id",
			description:   "跨表字段比较，替换指定表前缀",
		},
		{
			name:          "IN clause with table prefix",
			query:         "SELECT * FROM orders WHERE orders.user_id IN (1, 2, 3)",
			expectedWhere: "orders_1.user_id IN (1, 2, 3)",
			description:   "IN 子句中字段表前缀应替换",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Parse query into sqlparser.Statement
			parser := sqlparser.NewParser(strings.NewReader(tt.query))
			expr, err := parser.ParseStatement()
			assert.NoError(t, err, "解析 SQL 失败")

			// Extract and modify WHERE condition
			stmt, ok := expr.(*sqlparser.SelectStatement)
			assert.True(t, ok, "Expected a SelectStatement")

			condition := stmt.Condition
			updatedCondition := replaceTablePrefixInWhere(condition, oldTableName, newTableName)
			stmt.Condition = updatedCondition

			// Validate only the WHERE condition part
			assert.Equal(t, tt.expectedWhere, stmt.Condition.String(), tt.description)
		})
	}
}

// replaceTablePrefixInWhere replaces table prefixes in a WHERE condition
func replaceTablePrefixInWhere(condition sqlparser.Expr, oldTableName, newTableName string) sqlparser.Expr {
	if condition == nil {
		return nil
	}

	switch cond := condition.(type) {
	case *sqlparser.BinaryExpr:
		// 二元表达式（如 a = b 或 a > b），递归处理左右表达式
		cond.X = replaceTablePrefixInWhere(cond.X, oldTableName, newTableName)
		cond.Y = replaceTablePrefixInWhere(cond.Y, oldTableName, newTableName)
		return cond

	case *sqlparser.ParenExpr:
		// 括号表达式，递归处理内部表达式
		cond.X = replaceTablePrefixInWhere(cond.X, oldTableName, newTableName)
		return cond

	case *sqlparser.QualifiedRef:
		// 如果是带表前缀的字段，检查是否需要替换表名
		if cond.Table != nil && cond.Table.Name == oldTableName {
			cond.Table.Name = newTableName
		}
		return cond

	default:
		// 对于其他类型的表达式，直接返回
		return condition
	}
}
