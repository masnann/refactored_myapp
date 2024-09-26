package middlewares

import (
	"log"
	"myapp/config"
	"myapp/constants"
	"myapp/helpers"
	"myapp/models"
	"net/http"
	"strings"

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
