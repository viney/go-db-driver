// Copyright 2011 Antonio Diaz Gil. All rights reserved.

package dboracle

/*
#cgo CFLAGS: -I/usr/lib/oracle/instantclient_11_2/sdk/include
#cgo LDFLAGS: -lclntsh -L/usr/lib/oracle/instantclient_11_2

#include <oci.h>
//#include <oratypes.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

int Printw (const char* str) { 
	return puts (str); 
}

sword AuxOCIServerAttach( OCIServer *hServer, OCIError *hError, char *tns) 
{
	return OCIServerAttach( hServer, hError, (OraText*)tns, strlen(tns), OCI_DEFAULT);
}

void AuxOCIErrorGet( OCIError *hError) {
	int code;
	char str[200];
    OCIErrorGet( hError, 1, NULL, &code, (OraText*)str, 100, OCI_HTYPE_ERROR);
	printf( "%d: ", code);
	switch( code) {
	  	case 0: 
	  		printf ("OCI_SUCCESS");
	  		break;
	  	case -1: 
	  		printf ("OCI_ERROR");
	  		break;
	  	case -2: 
	  		printf ("OCI_INVALID_HANDLE");
	  		break;
	  	case -3123: 
	  		printf ("OCI_STILL_EXECUTING");
	  		break;
	  	case -24200: 
	  		printf ("OCI_CONTINUE");
	  		break;
	  	case 1: 
	  		printf ("OCI_SUCCESS_WITH_INFO");
	  		break;
	  	case 99: 
	  		printf ("OCI_NEED_DATA");
	  		break;
	  	case 100: 
	  		printf ("OCI_NO_DATA");
	  		break;
	  	default:
	  		printf ("???");
	  		break;
	}
	printf( " '%s'\n", str);
}

OraText *Get_OraText( char* str) {
	return (OraText *)str;
}
*/
import "C"

import (
	//    "os"
	"fmt"
	"unsafe"
	//"ConnectionOracle"
)

type DriverOracle struct {
	hEnv     *C.OCIEnv
	hService *C.OCISvcCtx
	hError   *C.OCIError
}

func NewDriverOracle() (d *DriverOracle, err int) {
	d = new(DriverOracle)
	d.hEnv = new(C.OCIEnv)
	res := C.OCIEnvCreate(&d.hEnv, C.OCI_DEFAULT, nil, nil, nil, nil, 0, nil)
	if C.OCI_SUCCESS != res && C.OCI_STILL_EXECUTING != res {
		fmt.Printf("Result1=%v\n", res)
		return nil, 1
	}

	pService := unsafe.Pointer(&d.hService)
	//    res = C.OCIHandleAlloc( unsafe.Pointer(env), (*unsafe.Pointer), C.OCI_HTYPE_SVCCTX, 0, nil);
	res = C.OCIHandleAlloc(unsafe.Pointer(d.hEnv), (*unsafe.Pointer)(pService), C.OCI_HTYPE_SVCCTX, 0, nil)
	if C.OCI_SUCCESS != res && C.OCI_STILL_EXECUTING != res {
		fmt.Printf("Result2=%v\n", res)
		return nil, 2
	}

	pError := unsafe.Pointer(&d.hError)
	res = C.OCIHandleAlloc(unsafe.Pointer(d.hEnv), (*unsafe.Pointer)(pError), C.OCI_HTYPE_ERROR, 0, nil)
	if C.OCI_SUCCESS != res && C.OCI_STILL_EXECUTING != res {
		fmt.Printf("Result3=%v\n", res)
		return nil, 3
	}
	return d, 0
}

func (d *DriverOracle) NewConnection(tns string, username string, password string) (c *ConnectionOracle, err int) {
	fmt.Printf("username=%v, password=%v, tns=%v\n", username, password, tns)
	fmt.Println("Connecting to Oracle...")

	c = new(ConnectionOracle)
	pServer := unsafe.Pointer(&c.hServer)
	res := C.OCIHandleAlloc(unsafe.Pointer(d.hEnv), (*unsafe.Pointer)(pServer), C.OCI_HTYPE_SERVER, 0, nil)
	if C.OCI_SUCCESS != res && C.OCI_STILL_EXECUTING != res {
		fmt.Printf("Resultado3=%v\n", res)
		return nil, 4
	}
	pTns := C.CString(tns)
	defer C.free(unsafe.Pointer(pTns))
	//    C.Printw (pTns);
	res = C.AuxOCIServerAttach(c.hServer, d.hError, pTns)
	if C.OCI_SUCCESS != res && C.OCI_STILL_EXECUTING != res {
		fmt.Printf("Resultado OCIServerAttach=%v\n", res)
		C.AuxOCIErrorGet(d.hError)
		return nil, 5
	}

	res = C.OCIAttrSet(unsafe.Pointer(d.hService), C.OCI_HTYPE_SVCCTX, unsafe.Pointer(c.hServer), 0, C.OCI_ATTR_SERVER, d.hError)
	if C.OCI_SUCCESS != res && C.OCI_STILL_EXECUTING != res {
		fmt.Printf("Resultado OCIAttrSet SERVICE-SERVER=%v\n", res)
		C.AuxOCIErrorGet(d.hError)
		return nil, 6
	}

	var test [1000]C.OraText
	res = C.OCIServerVersion(unsafe.Pointer(c.hServer), d.hError, &test[0], C.ub4(len(test)), C.OCI_HTYPE_SERVER)
	if C.OCI_SUCCESS != res && C.OCI_STILL_EXECUTING != res {
		fmt.Printf("Resultado OCIServerVersion=%v\n", res)
		C.AuxOCIErrorGet(d.hError)
		return nil, 7
	}
	fmt.Printf("Version=%s\n", test)

	pSession := unsafe.Pointer(&c.hSession)
	res = C.OCIHandleAlloc(unsafe.Pointer(d.hEnv), (*unsafe.Pointer)(pSession), C.OCI_HTYPE_SESSION, 0, nil)
	if C.OCI_SUCCESS != res && C.OCI_STILL_EXECUTING != res {
		fmt.Printf("Resultado hSession=%v\n", res)
		return nil, 8
	}
	pUsername := unsafe.Pointer(C.CString(username))
	defer C.free(pUsername)
	res = C.OCIAttrSet(unsafe.Pointer(c.hSession), C.OCI_HTYPE_SESSION, pUsername, C.ub4(len(username)), C.OCI_ATTR_USERNAME, d.hError)
	if C.OCI_SUCCESS != res && C.OCI_STILL_EXECUTING != res {
		fmt.Printf("Resultado OCIAttrSet USERNAME=%v\n", res)
		return nil, 9
	}
	cPassword := C.CString(password)
	//    pPassword := unsafe.Pointer(C.CString(password))
	pPassword := unsafe.Pointer(cPassword)
	defer C.free(pPassword)
	res = C.OCIAttrSet(unsafe.Pointer(c.hSession), C.OCI_HTYPE_SESSION, pPassword, C.ub4(len(password)), C.OCI_ATTR_PASSWORD, d.hError)
	if C.OCI_SUCCESS != res && C.OCI_STILL_EXECUTING != res {
		fmt.Printf("Resultado OCIAttrSet PASSWORD=%v\n", res)
		return nil, 10
	}
	res = C.OCISessionBegin(d.hService, d.hError, c.hSession, C.OCI_CRED_RDBMS, C.OCI_DEFAULT)
	if C.OCI_SUCCESS != res && C.OCI_STILL_EXECUTING != res {
		fmt.Printf("Resultado OCIAttrSet SessionBegin=%v\n", res)
		C.AuxOCIErrorGet(d.hError)
		return nil, 11
	}
	res = C.OCIAttrSet(unsafe.Pointer(d.hService), C.OCI_HTYPE_SVCCTX, unsafe.Pointer(c.hSession), 0, C.OCI_ATTR_SESSION, d.hError)
	if C.OCI_SUCCESS != res && C.OCI_STILL_EXECUTING != res {
		fmt.Printf("Resultado OCIAttrSet Service<-Session=%v\n", res)
		C.AuxOCIErrorGet(d.hError)
		return nil, 12
	}
	fmt.Printf("End NewConnection\n")

	c.driver = d
	return c, 0
}
