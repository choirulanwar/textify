package helper

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"mime/multipart"
	"os"
	"os/exec"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/bwmarrin/snowflake"
	"github.com/choirulanwar/textify/backend/models"
	"github.com/choirulanwar/textify/backend/pkg/pagination"

	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/hashicorp/go-uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func GenerateBaseSnowId(num int, n *snowflake.Node) string {
	if n == nil {
		localIp, err := GetLocalIpToInt()
		if err != nil {
			return ""
		}
		node, err := snowflake.NewNode(int64(localIp) % 1023)
		n = node
	}
	id := n.Generate()
	switch num {
	case 2:
		return id.Base2()
	case 32:
		return id.Base32()
	case 36:
		return id.Base36()
	case 58:
		return id.Base58()
	case 64:
		return id.Base64()
	default:
		return gconv.String(id.Int64())
	}
}

func GenerateUuid(size int) string {
	str, err := uuid.GenerateUUID()
	if err != nil {
		return ""
	}
	return gstr.SubStr(str, 0, size)
}

func GeneratePasswordHash(password string, salt string) string {
	md5 := md5.New()
	io.WriteString(md5, password)
	s := sha256.New()
	io.WriteString(s, password+salt)
	str := fmt.Sprintf("%x", s.Sum(nil))
	return str
}

func IsPathExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func MakeMultiDir(filePath string) error {
	if !IsPathExist(filePath) {
		err := os.MkdirAll(filePath, os.ModePerm)
		if err != nil {
			return err
		}
		return err
	}
	return nil
}

func MakeFileOrPath(path string) (*os.File, error) {
	pathArr := strings.Split(path, "/")
	pathUrl := strings.Join(pathArr[:len(pathArr)-1], "/")
	if err := MakeMultiDir(pathUrl); err != nil {
		return nil, err
	}
	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return file, nil
}

func WriteContentToFile(file *multipart.FileHeader, filePath string) error {
	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()
	open, err := file.Open()
	if err != nil {
		return err
	}
	defer open.Close()
	fileBytes, err := ioutil.ReadAll(open)
	if err != nil {
		return err
	}
	if _, err := f.Write(fileBytes); err != nil {
		return err
	}
	return nil
}

func MakeTimeFormatDir(rootPath, pathName, timeFormat string) (string, error) {
	filePath := "upload/"
	if pathName != "" {
		filePath += pathName + "/"
	}
	filePath += time.Now().Format(timeFormat) + "/"
	if err := MakeMultiDir(rootPath + filePath); err != nil {
		return "", err
	}
	return filePath, nil
}

func Cmd(commandName string, params []string) (string, error) {
	cmd := exec.Command(commandName, params...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return "", err
	}
	err = cmd.Wait()
	return out.String(), err
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func GetFileMD5(path string) (MD5 string, err error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	md5Hash := md5.New()
	if _, err := io.Copy(md5Hash, f); err != nil {
		return "", err
	}
	return hex.EncodeToString(md5Hash.Sum(nil)), nil
}

func ArrayChunk[T any](arr []T, size int) [][]T {
	var chunks [][]T
	for i := 0; i < len(arr); i += size {
		end := i + size
		if end > len(arr) {
			end = len(arr)
		}
		chunks = append(chunks, arr[i:end])
	}
	return chunks
}

func GetRuntimeUserHomeDir() (runtimeDir string) {
	switch runtime.GOOS {
	case "windows":
		runtimeDir = os.Getenv("APPDATA")
	case "darwin":
		runtimeDir, _ = os.UserHomeDir()
	default:
		runtimeDir, _ = os.UserHomeDir()
	}
	return
}

func GetStructColumnName(s interface{}, _type int) ([]string, error) {
	v := reflect.ValueOf(s)
	if v.Kind() != reflect.Struct {
		return []string{}, fmt.Errorf("interface is not a struct")
	}
	t := v.Type()
	var fields []string
	for i := 0; i < v.NumField(); i++ {
		var field string
		if _type == 1 {
			field = t.Field(i).Tag.Get("json")
			if field == "" {
				tagSetting := schema.ParseTagSetting(t.Field(i).Tag.Get("gorm"), ";")
				field = tagSetting["COLUMN"]
			}
		} else {
			field = t.Field(i).Name
		}
		fields = append(fields, field)
	}
	return fields, nil
}

func GetSettingInfo(db *gorm.DB) (*models.Setting, error) {
	var settingInfo models.Setting
	err := db.Where("id = ?", 1).First(&settingInfo).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("setting not found")
		}
		return nil, errors.New("internal server error")
	}
	return &settingInfo, nil
}

func ReverseArray(arr []int) []int {
	reversed := make([]int, len(arr))
	for i := len(arr) - 1; i >= 0; i-- {
		reversed[len(arr)-1-i] = arr[i]
	}
	return reversed
}

func ReverseSlice(slice interface{}) {
	value := reflect.ValueOf(slice)

	for i, j := 0, value.Len()-1; i < j; i, j = i+1, j-1 {
		temp := reflect.ValueOf(value.Index(i).Interface())
		value.Index(i).Set(value.Index(j))
		value.Index(j).Set(temp)
	}
}

type FilterFunc func(db *gorm.DB) *gorm.DB

func Paginate(value interface{}, pagination *pagination.Pagination, db *gorm.DB, filters ...FilterFunc) func(db *gorm.DB) *gorm.DB {
	var totalRows int64
	query := db.Model(value)

	for _, filter := range filters {
		query = filter(query)
	}

	query.Count(&totalRows)

	pagination.TotalRows = totalRows
	totalPages := int(math.Ceil(float64(totalRows) / float64(pagination.Limit)))
	pagination.TotalPages = totalPages

	return func(db *gorm.DB) *gorm.DB {
		for _, filter := range filters {
			db = filter(db)
		}

		return db.Offset(pagination.GetOffset()).Limit(pagination.GetLimit()).Order(pagination.GetSort())
	}
}
