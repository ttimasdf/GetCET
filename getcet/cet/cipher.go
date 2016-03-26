package cet

/*
#cgo LDFLAGS: -L/usr/local/opt/openssl/lib -lcrypto
#cgo CPPFLAGS: -I/usr/local/opt/openssl/include
#include <stdlib.h>
#include <openssl/des.h>
*/
import "C"

import (
	"unsafe"
	"fmt"
)

const (
	DECRYPT = 0
	ENCRYPT = 1
)

type CETCipher struct {
	TicketKey  string
	RequestKey string
}

func NewCETCipher(tKey string, rKey string) *CETCipher {
	return &CETCipher{TicketKey:tKey, RequestKey:rKey}
}

func (cc *CETCipher) DecryptTicket(etData []byte) ([]byte, error) {
	output, err := ProcessData(etData, cc.TicketKey, DECRYPT)
	return output, err
}

func (cc *CETCipher) EncryptRequest(rqData []byte) ([]byte, error) {
	output, err := ProcessData(rqData, cc.RequestKey, ENCRYPT)
	return output, err
}

func ProcessData(input []byte, key string, isEnc int) ([]byte, error) {
	if keyLen := len(key); keyLen != 8 {
		return nil, fmt.Errorf("cet: The key of DES must be 8 length, got %d instead.", keyLen)
	}

	cIsEnc := C.int(isEnc)
	inputSize := C.long(len(input))
	n := C.int(0)

	pInData := (*C.uchar)(unsafe.Pointer(&input[0]))
	pKeySchedule := &C.DES_key_schedule{}

	pOutData := C.calloc(C.size_t(inputSize), 1)
	pKey := (*C.DES_cblock)(unsafe.Pointer(C.CString(key)))

	defer func() {
		C.free(pOutData)
		C.free(unsafe.Pointer(pKey))
	}()

	C.DES_set_odd_parity(pKey)
	C.DES_set_key_checked((*C.const_DES_cblock)(pKey), pKeySchedule)
	C.DES_cfb64_encrypt(pInData, (*C.uchar)(pOutData), inputSize, pKeySchedule, pKey, &n, cIsEnc)

	output := C.GoBytes(pOutData, C.int(inputSize))
	if len(output) == 0 {
		return nil, fmt.Errorf("cet: Process Data failed.")
	}
	return output, nil
}

