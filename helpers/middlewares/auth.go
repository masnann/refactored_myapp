package middlewares

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"io"
	"log"
	"myapp/config"
	"myapp/constants"
	"myapp/handler"
	"myapp/helpers"
	"myapp/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var result models.Response
		// Extract the JWT token from the request header
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			result = helpers.ResponseJSON(false, constants.VALIDATION_ERROR_CODE, "Missing authorization header", nil)
			return c.JSON(http.StatusBadRequest, result)
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate the JWT token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.JWTSecret), nil
		})
		if err != nil {
			result = helpers.ResponseJSON(false, constants.UNAUTHORIZED_CODE, err.Error(), nil)
			return c.JSON(http.StatusUnauthorized, result)
		}

		// Extract claims and create a User struct
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := claims["userID"].(float64)
			userRole := claims["role"].(string)
			userEmail := claims["email"].(string)

			user := models.CurrentUserModels{
				ID:    int64(userID),
				Role:  userRole,
				Email: userEmail,
			}

			// Set the user struct in the context
			c.Set("user", user)
		} else {
			result = helpers.ResponseJSON(false, constants.UNAUTHORIZED_CODE, "Invalid token", nil)
			return c.JSON(http.StatusUnauthorized, result)
		}

		return next(c)
	}
}

// PermissionMiddleware to check user permissions
// func PermissionMiddleware(handler handler.Handler, permissionGroup, permissionName string, next echo.HandlerFunc) echo.HandlerFunc {
// 	return func(c echo.Context) error {
// 		var result models.Response
// 		currentUser, ok := c.Get("user").(models.CurrentUserModels)
// 		if !ok {
// 			result = helpers.ResponseJSON(false, constants.UNAUTHORIZED_CODE, "Failed to get user from context", nil)
// 			return c.JSON(http.StatusUnauthorized, result)
// 		}

// 		// Check permissions for Merchant Admin role
// 		if currentUser.Role == constants.MerchantAdminRole {
// 			hasPermission, err := handler.RolePermissionService.RolePermissionIsMerchantAdminHavePermission(currentUser.ID, permissionGroup, permissionName)
// 			if err != nil {
// 				log.Printf("Error checking permissions for Merchant Admin ID %d, Email %s: %v", currentUser.ID, currentUser.Email, err)
// 				result = helpers.ResponseJSON(false, constants.SYSTEM_ERROR_CODE, "Failed to check Merchant Admin permissions", nil)
// 				return c.JSON(http.StatusInternalServerError, result)
// 			}
// 			if hasPermission {
// 				return next(c)
// 			}
// 		} else {
// 			// Check permissions for other roles
// 			hasPermission, err := handler.RolePermissionService.RolePermissionIsRoleHavePermission(currentUser.ID, permissionGroup, permissionName)
// 			if err != nil {
// 				log.Printf("Error checking permissions for user ID %d, Email %s: %v", currentUser.ID, currentUser.Email, err)
// 				result = helpers.ResponseJSON(false, constants.SYSTEM_ERROR_CODE, "Failed to check user permissions", nil)
// 				return c.JSON(http.StatusInternalServerError, result)
// 			}
// 			if hasPermission {
// 				return next(c)
// 			}
// 		}

// 		// If no permissions are found, return 403 Forbidden
// 		log.Printf("Access denied for user ID %d, Email %s. Permission Group: %s, Permission Name: %s", currentUser.ID, currentUser.Email, permissionGroup, permissionName)
// 		result = helpers.ResponseJSON(false, constants.FORBIDDEN_CODE, constants.AccessDenied, nil)
// 		return c.JSON(http.StatusForbidden, result)
// 	}
// }

// SuperAdminMiddleware checks if the current user is a SuperAdmin.
func SuperAdminMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var result models.Response
		currentUser, ok := c.Get("user").(models.CurrentUserModels)
		if !ok {
			result = helpers.ResponseJSON(false, constants.UNAUTHORIZED_CODE, "Failed to get user from context", nil)
			return c.JSON(http.StatusUnauthorized, result)
		}
		if currentUser.Role != "Customer" {
			log.Printf("Access denied for Email %s", currentUser.Email)
			result = helpers.ResponseJSON(false, constants.FORBIDDEN_CODE, constants.AccessDenied, nil)
			return c.JSON(http.StatusForbidden, result)
		}
		// Call next handler if user is SuperAdmin
		return next(c)
	}
}

func SignatureValidatorMiddleware(handler handler.Handler) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			var result models.Response

			// Baca dan simpan body ke dalam buffer
			bodyBytes, err := io.ReadAll(ctx.Request().Body)
			if err != nil {
				result = helpers.ResponseJSON(false, constants.INTERNAL_SERVER_ERROR, "Failed to read request body", nil)
				return ctx.JSON(http.StatusInternalServerError, result)
			}

			// Reset body agar bisa dibaca ulang di handler
			ctx.Request().Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			// Validasi timestamp
			timestampStr := ctx.Request().Header.Get("X-TIMESTAMP")
			if timestampStr == "" {
				result = helpers.ResponseJSON(false, constants.BAD_REQUEST_CODE, "Timestamp is required", nil)
				return ctx.JSON(http.StatusBadRequest, result)
			}

			// Parsing timestamp dengan format ISO 8601
			timestamp, err := time.Parse(time.RFC3339, timestampStr)
			if err != nil {
				log.Printf("Invalid timestamp format: %v", err)
				result = helpers.ResponseJSON(false, constants.BAD_REQUEST_CODE, "Invalid timestamp format", nil)
				return ctx.JSON(http.StatusBadRequest, result)
			}

			// Validasi kedaluwarsa (contoh: 5 menit = 300 detik)
			now := time.Now()
			if timestamp.After(now) || now.Sub(timestamp) > 5*time.Minute {
				log.Printf("Timestamp expired or in the future: now=%v, timestamp=%v", now, timestamp)
				result = helpers.ResponseJSON(false, constants.BAD_REQUEST_CODE, "Request timestamp expired or invalid", nil)
				return ctx.JSON(http.StatusBadRequest, result)
			}

			// Validasi header lainnya
			signature := ctx.Request().Header.Get("X-SIGNATURE")
			if signature == "" {
				result = helpers.ResponseJSON(false, constants.BAD_REQUEST_CODE, "Signature is required", nil)
				return ctx.JSON(http.StatusBadRequest, result)
			}

			partnerIDStr := ctx.Request().Header.Get("X-PARTNERID")
			if partnerIDStr == "" {
				result = helpers.ResponseJSON(false, constants.BAD_REQUEST_CODE, "Partner ID is required", nil)
				return ctx.JSON(http.StatusBadRequest, result)
			}

			partnerID, err := strconv.ParseInt(partnerIDStr, 10, 64)
			if err != nil {
				result = helpers.ResponseJSON(false, constants.BAD_REQUEST_CODE, "Invalid Partner ID format", nil)
				return ctx.JSON(http.StatusBadRequest, result)
			}

			secretKey := ctx.Request().Header.Get("X-SECRETKEY")
			if secretKey == "" {
				result = helpers.ResponseJSON(false, constants.BAD_REQUEST_CODE, "Secret Key is required", nil)
				return ctx.JSON(http.StatusBadRequest, result)
			}

			// Ambil partner dan validasi public key
			partner, err := handler.PartnerService.PartnerFindByPartnerIDandSecretKey(partnerID, secretKey)
			if err != nil {
				result = helpers.ResponseJSON(false, constants.INTERNAL_SERVER_ERROR, "Failed to load partner", nil)
				return ctx.JSON(http.StatusInternalServerError, result)
			}

			decodedPublicKey, err := handler.Utils.DecryptEncryptedText(partner.PublicKey)
			if err != nil {
				result = helpers.ResponseJSON(false, constants.INTERNAL_SERVER_ERROR, "Failed to decrypt public key", nil)
				return ctx.JSON(http.StatusInternalServerError, result)
			}

			block, _ := pem.Decode([]byte(decodedPublicKey))
			if block == nil {
				result = helpers.ResponseJSON(false, constants.BAD_REQUEST_CODE, "Invalid public key format", nil)
				return ctx.JSON(http.StatusBadRequest, result)
			}

			publicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
			if err != nil {
				result = helpers.ResponseJSON(false, constants.BAD_REQUEST_CODE, "Invalid public key format", nil)
				return ctx.JSON(http.StatusBadRequest, result)
			}

			rsaPublicKey, ok := publicKey.(*rsa.PublicKey)
			if !ok {
				result = helpers.ResponseJSON(false, constants.BAD_REQUEST_CODE, "Public key is not of type *rsa.PublicKey", nil)
				return ctx.JSON(http.StatusBadRequest, result)
			}

			// Verifikasi tanda tangan
			valid, err := helpers.VerifySignature(string(bodyBytes), signature, rsaPublicKey)
			if err != nil || !valid {
				result = helpers.ResponseJSON(false, constants.FORBIDDEN_CODE, "Invalid signature", nil)
				return ctx.JSON(http.StatusForbidden, result)
			}

			// Tambahkan partner ke context untuk digunakan di handler
			ctx.Set("partner", partner)

			// Lanjutkan ke handler berikutnya
			return next(ctx)
		}
	}
}
