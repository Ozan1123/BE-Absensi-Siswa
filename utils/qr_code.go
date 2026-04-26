package utils

import "github.com/skip2/go-qrcode"

func GenerateQRCode(token string) ([]byte,error){
    return qrcode.Encode(
        token,
        qrcode.Medium,
        256,
    )
}