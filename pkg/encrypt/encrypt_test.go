package encrypt

import (
	"bufio"
	"os"
	"slices"
	"testing"
)

func TestMD5(t *testing.T) {
	var ss [16]byte
	ss[0] = '2'
	an := make([]byte, 16)
	an[0] = '2'
	if !slices.Equal(byte16ToBytes(ss), an) {
		t.Fatal("error", an, byte16ToBytes(ss))
	}
}

func TestEncryptMobile(t *testing.T) {
	mobile := "13800138000"
	encryptedMobile, err := EncMobile(mobile)
	if err != nil {
		t.Fatal(err)
	}
	decryptedMobile, err := DecMobile(encryptedMobile)
	if err != nil {
		t.Fatal(err)
	}
	if mobile != decryptedMobile {
		t.Fatalf("expected %s, but got %s", mobile, decryptedMobile)
	}
}
func TestFileMd5(t *testing.T) {
	data := []byte("5.5@#$")
	file, _ := os.OpenFile("a.txt", os.O_CREATE|os.O_RDWR, os.ModePerm)
	reader := bufio.NewReader(file)
	buf := make([]byte, 1024)
	num, err := reader.Read(buf)
	if err != nil {
		return
	}
	if Md5Sum(buf[:num]) != Md5Sum(data) {
		t.Fatal("md5 error", data, buf[:num])
	}
}
