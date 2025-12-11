package common

import (
	"testing"

	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/encoding/gbase64"
	"github.com/gogf/gf/v2/frame/g"

	"github.com/denghuo98/zzframe/web/zencrypt"
	"github.com/denghuo98/zzframe/zconsts"
)

// encryptPassword 将明文密码加密为前端传输格式
// 先 AES 加密，再 base64 编码
func encryptPassword(plainText string) string {
	encrypted, err := zencrypt.AesECBEncrypt([]byte(plainText), zconsts.RequestEncryptKey)
	if err != nil {
		panic(err)
	}
	return gbase64.EncodeToString(encrypted)
}

func TestVerifyPassword(t *testing.T) {
	s := NewCommonSite()
	salt := "abc123"
	plainPassword := "testPassword123"

	// 计算正确的密码 hash
	correctHash := gmd5.MustEncryptString(plainPassword + salt)

	tests := []struct {
		name      string
		input     string
		salt      string
		hash      string
		wantError bool
		errorMsg  string
	}{
		{
			name:      "正确的密码验证",
			input:     encryptPassword(plainPassword),
			salt:      salt,
			hash:      correctHash,
			wantError: false,
		},
		{
			name:      "错误的密码",
			input:     encryptPassword("wrongPassword"),
			salt:      salt,
			hash:      correctHash,
			wantError: true,
			errorMsg:  "用户密码错误",
		},
		{
			name:      "错误的盐值",
			input:     encryptPassword(plainPassword),
			salt:      "wrongSalt",
			hash:      correctHash,
			wantError: true,
			errorMsg:  "用户密码错误",
		},
		{
			name:      "无效的 base64 编码",
			input:     "!!!invalid-base64!!!",
			salt:      salt,
			hash:      correctHash,
			wantError: true,
		},
		{
			name:      "无效的 AES 加密数据",
			input:     gbase64.EncodeToString([]byte("invalid aes data")),
			salt:      salt,
			hash:      correctHash,
			wantError: true,
		},
		{
			name:      "空密码",
			input:     encryptPassword(""),
			salt:      salt,
			hash:      correctHash,
			wantError: true,
			errorMsg:  "用户密码错误",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := s.verifyPassword(tt.input, tt.salt, tt.hash)
			if tt.wantError {
				if err == nil {
					t.Errorf("期望返回错误，但返回了 nil")
				} else if tt.errorMsg != "" && err.Error() != tt.errorMsg {
					t.Errorf("错误信息不匹配，期望: %s，实际: %s", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("不期望返回错误，但返回了: %v", err)
				}
			}
		})
	}
}

func TestVerifyPasswordWithDifferentSalts(t *testing.T) {
	plainPassword := "123456"
	encryptedPassword := encryptPassword(plainPassword)
	g.Log().Infof(t.Context(), "encryptedPassword: %s", encryptedPassword)
}
