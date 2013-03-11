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
	"unsafe"
)

type StatementOracle struct {
	con     *ConnectionOracle
	columns []*ColumnOracle

	hStatement *C.OCIStmt
	hDefine    *C.OCIDefine

	fetchDone  bool
	numColumns int
}

func (s *StatementOracle) Next() (err int) {

	if !s.fetchDone {
		s.hDefine = new(C.OCIDefine)
		res := C.OCIStmtExecute(s.con.driver.hService, s.hStatement, s.con.driver.hError, 0, 0, nil, nil, C.OCI_DEFAULT)
		if C.OCI_SUCCESS != res && C.OCI_STILL_EXECUTING != res {
			fmt.Printf("Resultado OCIStmtExecute=%v\n", res)
			C.AuxOCIErrorGet(s.con.driver.hError)
			return -1
		}
		s.fetchDone = true
		/*		
		    var idValue int
		    res = C.OCIDefineByPos( s.hStatement, &s.hDefine, s.con.driver.hError, 1 /*pos* /, unsafe.Pointer(&idValue), C.sb4(unsafe.Sizeof(idValue)), C.SQLT_INT, nil, nil, nil, C.OCI_DEFAULT);
		    if (C.OCI_SUCCESS != res /*&& C.OCI_STILL_EXECUTING != res* /) {
		 		fmt.Printf( "Resultado OCIDefineByPos=%v\n", res)
		     	C.AuxOCIErrorGet( s.con.driver.hError);
		     	return -2
		    }
		    var surnameValue [40]C.char
		    res = C.OCIDefineByPos( s.hStatement, &s.hDefine, s.con.driver.hError, 2 /*pos* /, unsafe.Pointer(&surnameValue[0]), C.sb4(unsafe.Sizeof(surnameValue)), C.SQLT_STR, nil, nil, nil, C.OCI_DEFAULT);
		    if (C.OCI_SUCCESS != res /*&& C.OCI_STILL_EXECUTING != res* /) {
		 		fmt.Printf( "Resultado OCIDefineByPos=%v\n", res)
		     	C.AuxOCIErrorGet( s.con.driver.hError);
		     	return -3
		    }
		*/
		res = C.OCIAttrGet(unsafe.Pointer(s.hStatement), C.OCI_HTYPE_STMT, unsafe.Pointer(&s.numColumns), nil, C.OCI_ATTR_PARAM_COUNT, s.con.driver.hError)
		if C.OCI_SUCCESS != res /*&& C.OCI_STILL_EXECUTING != res*/ {
			fmt.Printf("Resultado OCIAttrGet=%v\n", res)
			C.AuxOCIErrorGet(s.con.driver.hError)
			return -3
		}
		fmt.Printf("numColumns=%v\n", s.numColumns)

		s.columns = make([]*ColumnOracle, 0, 50)
		var hColumn C.ub4
		pColumn := unsafe.Pointer(&hColumn)
		for i := 1; i <= s.numColumns; i++ {
			/* get parameter for column i */
			res = C.OCIParamGet(unsafe.Pointer(s.hStatement), C.OCI_HTYPE_STMT, s.con.driver.hError, &pColumn, C.ub4(i))
			if C.OCI_SUCCESS != res /*&& C.OCI_STILL_EXECUTING != res*/ {
				fmt.Printf("Resultado OCIParamGet=%v\n", res)
				C.AuxOCIErrorGet(s.con.driver.hError)
				return -4
			}

			/* get data-type of column i */
			var columnType C.ub2
			pType := unsafe.Pointer(&columnType)
			res = C.OCIAttrGet(pColumn, C.OCI_DTYPE_PARAM, pType, nil, C.OCI_ATTR_DATA_TYPE, s.con.driver.hError)
			if C.OCI_SUCCESS != res /*&& C.OCI_STILL_EXECUTING != res*/ {
				fmt.Printf("Resultado OCIParamGet=%v\n", res)
				C.AuxOCIErrorGet(s.con.driver.hError)
				return -4
			}

			size := C.ub4(100)
			var column_name_tmp *C.char
			res = C.OCIAttrGet(pColumn, C.OCI_DTYPE_PARAM, unsafe.Pointer(&column_name_tmp), nil, C.OCI_ATTR_NAME, s.con.driver.hError)
			fmt.Printf("coumnName (size=%v) = %v\n", size, column_name_tmp)
			//column_name_tmp[size] = 0
			columnName := C.GoString(column_name_tmp)
			var columnSize int
			res = C.OCIAttrGet(pColumn, C.OCI_DTYPE_PARAM, unsafe.Pointer(&columnSize), nil, C.OCI_ATTR_DATA_SIZE, s.con.driver.hError)

			fmt.Printf("Column:%v Name:%v Size:%v Type:%v\n", i, columnName, columnSize, columnType)

			switch columnType {
			case C.SQLT_CHR:
				var columnValue = make([]C.char, columnSize+1)

				res = C.OCIDefineByPos(s.hStatement, &s.hDefine, s.con.driver.hError, C.ub4(i), /*pos*/
					unsafe.Pointer(&columnValue[0]), C.sb4(columnSize), columnType, nil, nil, nil, C.OCI_DEFAULT)
				if C.OCI_SUCCESS != res /*&& C.OCI_STILL_EXECUTING != res*/ {
					fmt.Printf("Resultado OCIDefineByPos=%v\n", res)
					C.AuxOCIErrorGet(s.con.driver.hError)
					return -2
				}
				s.columns = append(s.columns, NewColumnChrOracle(s, columnName, int(columnSize), columnValue))

			case C.SQLT_NUM:
				var columnValue int

				res = C.OCIDefineByPos(s.hStatement, &s.hDefine, s.con.driver.hError, C.ub4(i), /*pos*/
					unsafe.Pointer(&columnValue), C.sb4(4 /*sizeof(columnValue)*/), columnType, nil, nil, nil, C.OCI_DEFAULT)
				if C.OCI_SUCCESS != res /*&& C.OCI_STILL_EXECUTING != res*/ {
					fmt.Printf("Resultado OCIDefineByPos=%v\n", res)
					C.AuxOCIErrorGet(s.con.driver.hError)
					return -2
				}
				s.columns = append(s.columns, NewColumnNumberOracle(s, columnName, int(columnSize), columnValue))

			//case C.SQLT_LNG:
			//case C.SQLT_RID:
			//case C.SQLT_DAT:
			//case C.SQLT_BIN:
			//case C.SQLT_LBI:
			//case C.SQLT_AFC:	
			//case C.SQLT_REF:	
			//case C.SQLT_CLOB:	
			//case C.SQLT_BLOB:	
			//case C.SQLT_BFILEE:	
			//case C.SQLT_TIMESTAMP:	
			//case C.SQLT_TIME_TZ:	
			//case C.SQLT_INTERVAL_YM:	
			//case C.SQLT_INTERVAL_DS:	
			//case C.SQLT_TIMESTAMP_LTZ:	
			//case C.SQLT_RDD:	
			default:
				s.columns = append(s.columns, NewColumnUnknownOracle(s, columnName, int(columnSize), int(columnType)))
			}
			fmt.Printf("Next len(stmt.columns)=%v\n", len(s.columns))
		}
		//return 0
	}

	res := C.OCIStmtFetch2(s.hStatement, s.con.driver.hError, 1 /*nrows*/, C.OCI_FETCH_NEXT, 0, C.OCI_DEFAULT)
	if C.OCI_SUCCESS != res && C.OCI_NO_DATA != res {
		fmt.Printf("Resultado OCIStmtFetch2=%v\n", res)
		C.AuxOCIErrorGet(s.con.driver.hError)
		return -4
	} else if C.OCI_NO_DATA == res {
		//fmt.Printf( "The result is:%v-%v\n", idValue, C.GoString(&surnameValue[0]))
		return 0
	}
	return 1
}

func (stmt *StatementOracle) GetByPos(pos int) (col *ColumnOracle, err int) {
	// if stmt.state != FETCHED .... ERROR
	// if pos > len(stmt.columns) .... ERROR
	fmt.Printf("GetByPos(%v/%v)\n", pos, len(stmt.columns))
	return stmt.columns[pos], 0
}
