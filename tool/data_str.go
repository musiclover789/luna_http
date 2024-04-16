package tool

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"strconv"
	"time"
)

// CalculateMD5 计算字符串的 MD5 值
func CalculateMD5(input string) string {
	// 将字符串转换为字节数组
	data := []byte(input)

	// 计算 MD5 值
	hash := md5.Sum(data)

	// 将 MD5 值转换为十六进制字符串
	md5Str := hex.EncodeToString(hash[:])

	return md5Str
}

func GenerateIncrementalID() string {
	return CalculateMD5(generateRandomString())
}

func generateRandomString() string {
	// 获取当前时间的毫秒数
	millis := time.Now().UnixNano() / 1e6

	// 生成 10 位随机数
	rand.Seed(time.Now().UnixNano())
	randomNum := rand.Intn(10000000000)

	// 将毫秒数和随机数拼接成字符串
	result := strconv.FormatInt(millis, 10) + strconv.Itoa(randomNum)

	return result
}
