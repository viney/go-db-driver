// Copyright 2011 Antonio Diaz Gil. All rights reserved.

package dboracle

/*
#cgo CFLAGS: -I/usr/lib/oracle/instantclient_11_2/sdk/include
#cgo LDFLAGS: -lclntsh -L/usr/lib/oracle/instantclient_11_2

#include <oci.h>
//#include <oratypes.h>
#include <stdlib.h>
#include <string.h>

extern void AuxOCIErrorGet( OCIError *hError);
extern OraText *Get_OraText( char* str);

*/
import "C"

import (
	//    "os"
	"fmt"

//    "unsafe"
)

type ColumnOracle struct {
	stmt *StatementOracle

	name       string
	size       int
	oracleType int

	valueInt int
	valueStr []C.char
}

/*
type ColumnStrOracle struct 
{
	StatementOracle *DriverOracle
	ColumnOracle

	value string
}
*/

func NewColumnChrOracle(stmt *StatementOracle, name string, size int, value []C.char) (col *ColumnOracle) {
	col = new(ColumnOracle)
	col.stmt = stmt
	col.name = name
	col.size = size
	col.valueStr = value
	col.oracleType = C.SQLT_CHR

	return col
}

func NewColumnNumberOracle(stmt *StatementOracle, name string, size int, value int) (col *ColumnOracle) {
	col = new(ColumnOracle)
	col.stmt = stmt
	col.name = name
	col.size = size
	col.valueInt = value
	col.oracleType = C.SQLT_NUM

	return col
}

func NewColumnUnknownOracle(stmt *StatementOracle, name string, size int, oracleType int) (col *ColumnOracle) {
	col = new(ColumnOracle)
	col.stmt = stmt
	col.name = name
	col.size = size
	col.oracleType = oracleType
	//col.valueStr = fmt.Sprintf( "¿%v?", oracleType)

	return col
}

func (col *ColumnOracle) GetValue() string {
	switch col.oracleType {
	case C.SQLT_CHR:
		return C.GoString(&col.valueStr[0])
	case C.SQLT_NUM:
		return fmt.Sprintf("%v", col.valueInt)
	}
	return "¿?"
}
