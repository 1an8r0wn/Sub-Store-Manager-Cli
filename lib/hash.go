package lib

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// CreateHash 创建随机哈希值
func CreateHash() string {
	// 创建一个新的哈希对象
	h := sha256.New()

	// 生成随机数据
	randData := make([]byte, 16)
	_, err := rand.Read(randData)
	if err != nil {
		PrintError("generating random hash failed:", err)
		return ""
	}

	// 将随机数据写入哈希对象
	h.Write(randData)

	// 获取哈希值
	hash := h.Sum(nil)

	// 将哈希值转换为十六进制字符串
	hashStr := hex.EncodeToString(hash)

	fmt.Println("Random Hash:", hashStr)
	return hashStr[:28] // 返回前28位
}
