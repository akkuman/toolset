package utils

import (
	"archive/zip"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FileData 文件数据信息
type FileData struct {
	Name string
	Body []byte
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789+/")

func RandStringRunes(n int64) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func RandInt63(from int64, to int64) int64 {
	if to <= from || to <= 0 {
		to = 100
		from = 10
	}
	return rand.Int63n(to-from) + from
}

func PathExist(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func GetExecutableDir() string {
	ex, err := os.Executable()
	if err != nil {
		return ""
	}
	exPath := filepath.Dir(ex)
	return exPath
}

// XorEncryptDecrypt 异或加解密
func XorEncryptDecrypt(input, key []byte) (output []byte) {
	for i := 0; i < len(input); i++ {
		output = append(output, byte(input[i]^key[i%len(key)]))
	}
	return output
}

// CopyFile 拷贝文件
func CopyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()
	_, err = io.Copy(destination, source)
	if err != nil {
		return err
	}
	err = destination.Sync()
	return err
}

// CopyDir 拷贝一个目录到另一个目录
func CopyDir(src, dst string) error {
	err := filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		relPath := strings.Replace(path, src, "", 1)
		if relPath == "" {
			return nil
		}
		dstPath := filepath.Join(dst, relPath)
		if info.IsDir() {
			return os.Mkdir(dstPath, 0755)
		}
		return CopyFile(path, dstPath)
	})
	return err
}

// CustomEncryptData 私有的文件加密逻辑
func CustomEncryptData(data []byte) []byte {
	keyLen := RandInt63(10, 100)
	key := RandStringRunes(keyLen)
	buf := new(bytes.Buffer)
	encryptShellcode := XorEncryptDecrypt(data, []byte(key))
	binary.Write(buf, binary.LittleEndian, uint64(keyLen))
	binary.Write(buf, binary.LittleEndian, []byte(key))
	binary.Write(buf, binary.LittleEndian, uint64(len(encryptShellcode)))
	binary.Write(buf, binary.LittleEndian, []byte(encryptShellcode))
	return buf.Bytes()
}

// ZipData 压缩数据
func ZipData(files []FileData) (zipPath string, err error) {
	f, err := ioutil.TempFile(os.TempDir(), "tmp-*.zip")
	if err != nil {
		return "", err
	}
	defer f.Close()
	zipWriter := zip.NewWriter(f)
	defer zipWriter.Close()
	for _, file := range files {
		zipFile, err := zipWriter.Create(file.Name)
		if err != nil {
			return "", err
		}
		_, err = zipFile.Write(file.Body)
		if err != nil {
			return "", err
		}
	}
	return f.Name(), nil
}
