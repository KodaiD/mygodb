package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func prepareStatement(inputBuffer string, statement *Statement) PrepareResult {
	if inputBuffer[:6] == "insert" {
		statement.Type = STATEMENT_INSERT
		argsAssigned, err := fmt.Sscanf(inputBuffer, "insert %d %s %s", &statement.RowToInsert.ID,
			&statement.RowToInsert.UserName, &statement.RowToInsert.Email)
		if argsAssigned < 3 {
			log.Println(err)
			return PREPARE_SYNTAXERROR
		}
		return PREPARE_SUCCESS
	}
	if inputBuffer == "select" {
		statement.Type = STATEMENT_SELECT
		return PREPARE_SUCCESS
	}
	return PREPARE_UNRECOGNIZED_STATEMENT
}

func executeInsert(statement *Statement, table *Table) ExecuteResult {
	if table.NumRows >= TABLE_MAX_ROWS {
		return EXECUTE_TABLE_FULL
	}
	row := &statement.RowToInsert
	serializeRow(table, row)
	table.NumRows += 1
	return EXECUTE_SUCCESS
}

func executeSelect(statement *Statement, table *Table) ExecuteResult {
	row := &statement.RowToInsert
	for i := uint32(0); i < table.NumRows; i++ {
		deserializeRow(table, row, i)
		row.printRow()
	}
	return EXECUTE_SUCCESS
}

func executeStatement(statement *Statement, table *Table) ExecuteResult {
	switch statement.Type {
	case STATEMENT_INSERT:
		return executeInsert(statement, table)
	case STATEMENT_SELECT:
		return executeSelect(statement, table)
	default:
		fmt.Println("Warning...")
		return EXECUTE_SUCCESS
	}
}

func serializeRow(table *Table, r *Row) {
	data, err := json.Marshal(r)
	if err != nil {
		fmt.Println("---", err)
	}
	pn, rn := table.rowSlot(table.NumRows)
	table.Pager.Pages[pn].rows[rn] = data
}

func deserializeRow(table *Table, r *Row, i uint32) {
	pn, rn := table.rowSlot(i)
	data := table.Pager.Pages[pn].rows[rn]
	err := json.Unmarshal(data, r)
	if err != nil {
		fmt.Println("+++", err)
	}
}

func (r *Row) printRow() {
	fmt.Printf("(%d, %s, %s)\n", r.ID, r.UserName, r.Email)
}

func (table *Table) rowSlot(rowNum uint32) (uint32, uint32) {
	pageNum := rowNum / ROWS_PER_PAGE
	rowOffset := rowNum % ROWS_PER_PAGE
	return pageNum, rowOffset
}