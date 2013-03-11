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
	"unsafe"
)

type ConnectionOracle struct {
	driver *DriverOracle

	hServer  *C.OCIServer
	hSession *C.OCISession
	//hError *C.OCIError
}

func (c *ConnectionOracle) NewStatement(sql string) (s *StatementOracle, err int) {

	s = new(StatementOracle)
	pStatement := unsafe.Pointer(&s.hStatement)
	res := C.OCIHandleAlloc(unsafe.Pointer(c.driver.hEnv), (*unsafe.Pointer)(pStatement), C.OCI_HTYPE_STMT, 0, nil)
	if C.OCI_SUCCESS != res /*&& C.OCI_STILL_EXECUTING != res*/ {
		fmt.Printf("Result hSession=%v\n", res)
		C.AuxOCIErrorGet(c.driver.hError)
		return nil, -1
	}

	cSql := C.CString(sql)
	defer C.free(unsafe.Pointer(cSql))
	var osql *C.OraText = C.Get_OraText(cSql)
	//fmt.Printf( "%v:%v\n", C.strlen(sql), C.GoString(sql))
	res = C.OCIStmtPrepare(s.hStatement, c.driver.hError, osql, C.ub4(C.strlen(cSql)), C.OCI_NTV_SYNTAX, C.OCI_DEFAULT)
	if C.OCI_SUCCESS != res /*&& C.OCI_STILL_EXECUTING != res*/ {
		fmt.Printf("Resultado OCIStmtPrepare=%v\n", res)
		C.AuxOCIErrorGet(c.driver.hError)
		return nil, -2
	}
	s.fetchDone = false
	s.con = c
	return s, 0
}
