package common

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTClaims JWT 声明
type JWTClaims struct {
	UserID string `json:"user_id"`
	JTI    string `json:"jti"` // JWT ID，用于标识唯一 token
	jwt.RegisteredClaims
}

// TokenBlacklistEntry 黑名单条目
type TokenBlacklistEntry struct {
	JTI         string    // JWT ID
	ExpiredAt   time.Time // token 原始过期时间
	Blacklisted time.Time // 加入黑名单的时间
}

// JWTManager JWT 管理器
type JWTManager struct {
	secretKey  string
	expiration time.Duration
	issuer     string    // 签发者
	blacklist  *sync.Map // token 黑名单 (本地缓存)
}

// NewJWTManager 创建 JWT 管理器
func NewJWTManager(secretKey string, expiration time.Duration) *JWTManager {
	return &JWTManager{
		secretKey:  secretKey,
		expiration: expiration,
		issuer:     "gosir",
		blacklist:  &sync.Map{},
	}
}

// GenerateToken 生成 JWT token
func (m *JWTManager) GenerateToken(userID string) (string, error) {
	now := time.Now()
	claims := JWTClaims{
		UserID: userID,
		JTI:    uuid.New().String(), // 生成唯一的 JTI
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(m.expiration)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    m.issuer,
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

	// 检查 token 是否在黑名单中
	if m.IsTokenBlacklisted(claims.JTI) {
		return nil, errors.New("token 已失效")
	}

	return claims, nil
}

// RefreshToken 刷新 token
func (m *JWTManager) RefreshToken(tokenString string) (string, error) {
	claims, err := m.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}

	// 检查 token 是否在黑名单中
	if m.IsTokenBlacklisted(claims.JTI) {
		return "", errors.New("token 已失效")
	}

	// 直接刷新，不限制时间
	return m.GenerateToken(claims.UserID)
}

// AddToBlacklist 将 token 加入黑名单
func (m *JWTManager) AddToBlacklist(jti string, expiredAt time.Time) {
	entry := TokenBlacklistEntry{
		JTI:         jti,
		ExpiredAt:   expiredAt,
		Blacklisted: time.Now(),
	}
	m.blacklist.Store(jti, entry)
}

// IsTokenBlacklisted 检查 token 是否在黑名单中
func (m *JWTManager) IsTokenBlacklisted(jti string) bool {
	if value, ok := m.blacklist.Load(jti); ok {
		entry := value.(TokenBlacklistEntry)

		// 如果 token 已经过期，从黑名单中删除
		if time.Now().After(entry.ExpiredAt) {
			m.blacklist.Delete(jti)
			return false
		}

		return true
	}
	return false
}

// CleanupExpiredBlacklist 清理黑名单中已过期的 token
func (m *JWTManager) CleanupExpiredBlacklist() {
	m.blacklist.Range(func(key, value interface{}) bool {
		entry := value.(TokenBlacklistEntry)
		if time.Now().After(entry.ExpiredAt) {
			m.blacklist.Delete(key)
		}
		return true
	})
}

// GetBlacklistSize 获取黑名单大小
func (m *JWTManager) GetBlacklistSize() int {
	size := 0
	m.blacklist.Range(func(key, value interface{}) bool {
		size++
		return true
	})
	return size
}

// 全局 JWT 管理器实例
var jwtManager *JWTManager

// InitJWT 初始化全局 JWT 管理器
func InitJWT(secretKey string, expiryHours int) {
	expiration := time.Duration(expiryHours) * time.Hour
	jwtManager = NewJWTManager(secretKey, expiration)
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
