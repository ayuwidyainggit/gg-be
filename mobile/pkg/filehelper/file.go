package filehelper

import (
	"bufio"
	"fmt"
	"io"
	"mime/multipart"
)

func MultipartToByte(file *multipart.FileHeader) (byteFile []byte, err error) {
	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	defer src.Close()
	size := file.Size

	// Read the file into a byte slice
	bs := make([]byte, size)
	_, err = bufio.NewReader(src).Read(bs)
	if err != nil && err != io.EOF {
		fmt.Println(err)
		return nil, err
	}
	return bs, nil
}
