package main

import (
	"encoding/binary"
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
	deserializeRow(table)
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
	/*
	how to store data:
	データ長(4B) データ
	 */
	startIdx := table.Pager.pages[table.Pager.point].point

	data, err := json.Marshal(r)
	if err != nil {
		log.Println("---", err)
	}
	rowLength := uint32(len(data))
	// pageが一杯なら他のpageに書く
	if table.Pager.point + uint32(rowLength) + uint32(1) > PAGE_SIZE {
		table.Pager.point++
	}
	buf := table.Pager.pages[table.Pager.point].buf
	binary.BigEndian.PutUint32(buf[startIdx:startIdx+4], rowLength)
	copy(buf[startIdx+4:], data)
	table.Pager.pages[table.Pager.point].point += 4 + rowLength
	table.NumRows++
}

func deserializeRow(table *Table) {
	var r Row
	var buf []byte
	var rowLength uint32
	var data []byte
	for i := uint32(0); i <= table.Pager.point; i++ {
		for j := uint32(0); j < table.Pager.pages[i].point; {
			buf = table.Pager.pages[i].buf
			rowLength = binary.BigEndian.Uint32(buf[j:j+4])
			//if rowLength <= 0 {
			//	break
			//}
			data = buf[j+4:j+4+rowLength]
			err := json.Unmarshal(data, &r)
			if err != nil {
				log.Println("+++", err)
			}
			r.printRow()
			j += rowLength + 4
		}
	}

}

func (r *Row) printRow() {
	fmt.Printf("(%d, %s, %s)\n", r.ID, r.UserName, r.Email)
}