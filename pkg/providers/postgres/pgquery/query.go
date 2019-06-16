package pgquery

import (
	"fmt"

	"github.com/lfittl/pg_query_go"
	nodes "github.com/lfittl/pg_query_go/nodes"
	"go.indent.com/apis/pkg/access/v1"
)

type AnalyzeResult struct {
	Actions []string
	Columns []string
	Tables  []string
}

type AnalyzeAccessResult struct {
	Actions   v1.Actions
	Resources v1.Resources
}

func print(str string, i ...interface{}) {
	return
	fmt.Printf(str, i...)
}

func AnalyzeAccess(q string) (AnalyzeAccessResult, error) {
	var aar AnalyzeAccessResult

	ar, err := Analyze(q)

	for _, act := range ar.Actions {
		action := "postgres:actions::sql:" + act
		aar.Actions = append(aar.Actions, action)
	}

	for _, col := range ar.Columns {
		resource := "postgres:resources::columns:" + col
		aar.Resources = append(aar.Resources, resource)
	}

	for _, tbl := range ar.Tables {
		resource := "postgres:resources::tables:" + tbl
		aar.Resources = append(aar.Resources, resource)
	}

	return aar, err
}

func Analyze(q string) (AnalyzeResult, error) {
	str, _ := pg_query.ParseToJSON(q)
	print("pgquery: Analyze: Query = %s\n", str)

	tree, err := pg_query.Parse(q)
	if err != nil {
		return AnalyzeResult{}, err
	}

	return traverseList(tree.Statements), err
}

func traverseNode(unode nodes.Node) (res AnalyzeResult) {
	switch node := unode.(type) {
	case nodes.RawStmt:
		print("pgquery: traverseNode: raw stmt\n")
		if node.Stmt != nil {
			res = combineResults(res, traverseNode(node.Stmt))
		}
	case nodes.SelectStmt:
		print("pgquery: traverseNode: select\n")
		res.Actions = append(res.Actions, "select")

		if len(node.FromClause.Items) > 0 {
			res = combineResults(res, traverseList(node.FromClause.Items))
		}
		if len(node.TargetList.Items) > 0 {
			res = combineResults(res, traverseList(node.TargetList.Items))
		}
	case nodes.ResTarget:
		print("pgquery: traverseNode: ResTarget\n")

		switch innerNode := node.Val.(type) {
		case nodes.ColumnRef:
			for _, ufield := range innerNode.Fields.Items {
				switch field := ufield.(type) {
				case nodes.A_Star:
					res.Columns = append(res.Columns, "*")
				case nodes.String:
					res.Columns = append(res.Columns, field.Str)
				}
			}
		}
	case nodes.RangeVar:
		print("pgquery: traverseNode: RangeVar = %s\n", *node.Relname)
		res.Tables = append(res.Tables, *node.Relname)
	case nodes.FromExpr:
		// print("pgquery: traverseNode: FromExpr = %v\n", node)
	case nodes.JoinExpr:
		// print("pgquery: traverseNode: JoinExpr = %v\n", node)
	default:
		// print("pgquery: traverseNode: unknown = %v\n", node)
	}

	return
}

func traverseList(list []nodes.Node) AnalyzeResult {
	var set []AnalyzeResult

	for _, unode := range list {
		set = append(set, traverseNode(unode))
	}

	return concatResultSet(set)
}

func combineResults(results ...AnalyzeResult) AnalyzeResult {
	return concatResultSet(results)
}

func concatResultSet(set []AnalyzeResult) AnalyzeResult {
	var actions, columns, tables []string

	for _, item := range set {
		if len(item.Actions) != 0 {
			actions = append(actions, item.Actions...)
		}
		if len(item.Columns) != 0 {
			columns = append(columns, item.Columns...)
		}
		if len(item.Tables) != 0 {
			tables = append(tables, item.Tables...)
		}
	}

	return AnalyzeResult{
		Actions: actions,
		Columns: columns,
		Tables:  tables,
	}
}
