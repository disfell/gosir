package common

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims JWT 声明
type JWTClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// JWTManager JWT 管理器
type JWTManager struct {
	secretKey  string
	expiration time.Duration
}

// NewJWTManager 创建 JWT 管理器
func NewJWTManager(secretKey string, expiration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:  secretKey,
		expiration: expiration,
	}
}

// GenerateToken 生成 JWT token
func (m *JWTManager) GenerateToken(userID string) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}

// ValidateToken 验证 JWT token
func (m *JWTManager) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("无效的签名方法: %v", token.Header["alg"])
		}
		return []byte(m.secretKey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("无效的 token")
	}

	return claims, nil
}

// RefreshToken 刷新 token
func (m *JWTManager) RefreshToken(tokenString string) (string, error) {
	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	// 检查是否快过期（例如剩余时间少于 1 小时）
	if time.Until(claims.ExpiresAt.Time) > time.Hour {
		return "", errors.New("token 未到刷新时间")
	}

	return m.GenerateToken(claims.UserID)
}

// 全局 JWT 管理器实例
var jwtManager *JWTManager

// InitJWT 初始化全局 JWT 管理器
func InitJWT(secretKey string, expiryHours int) {
	jwtManager = NewJWTManager(secretKey, time.Duration(expiryHours)*time.Hour)
}

// GetJWTManager 获取全局 JWT 管理器实例
func GetJWTManager() *JWTManager {
	return jwtManager
}

// GenerateToken 使用全局 JWT 管理器生成 token
func GenerateToken(userID string) (string, error) {
	if jwtManager == nil {
		return "", errors.New("JWT 管理器未初始化")
	}
	return jwtManager.GenerateToken(userID)
}
