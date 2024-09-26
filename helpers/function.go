package helpers

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
	"myapp/models"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

func ResponseJSON(success bool, code, message string, result interface{}) models.Response {
	response := models.Response{
		StatusCode:       code,
		Success:          success,
		Message:          message,
		ResponseDateTime: time.Now(),
		Result:           result,
	}

	return response
}

func TimeStampNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func ReplaceSQL(old, searchPattern string) string {
	tmpCount := strings.Count(old, searchPattern)
	for m := 1; m <= tmpCount; m++ {
		old = strings.Replace(old, searchPattern, "$"+strconv.Itoa(m), 1)
	}
	return old
}

func IsValidDate(date string) bool {
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}

func IsDateRangeValid(dateStart, dateEnd string) bool {
	start, _ := time.Parse("2006-01-02", dateStart)
	end, _ := time.Parse("2006-01-02", dateEnd)
	return !end.Before(start)
}

var validate *validator.Validate

func init() {
	validate = validator.New()
	err := validate.RegisterValidation("noSpace", func(fl validator.FieldLevel) bool {
		password := fl.Field().String()
		return !regexp.MustCompile(`\s`).MatchString(password)
	})
	if err != nil {
		log.Fatalf("error registering validation: %v", err)
	}
}

func ValidateStruct(ctx echo.Context, s interface{}) error {
	// Bind the request body to the struct
	if err := ctx.Bind(s); err != nil {
		return fmt.Errorf("Invalid request body")
	}

	// Validate the struct
	err := validate.Struct(s)
	if err != nil {
		var customErrors []string

		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			case "required":
				customErrors = append(customErrors, fmt.Sprintf("Field '%s' is required", err.Field()))
			case "min":
				customErrors = append(customErrors, fmt.Sprintf("Field '%s' must be at least %s characters long", err.Field(), err.Param()))
			case "email":
				customErrors = append(customErrors, fmt.Sprintf("Field '%s' must be a valid email address", err.Field()))
			case "noSpace":
				customErrors = append(customErrors, fmt.Sprintf("Field '%s' cannot contain spaces", err.Field()))
			case "alphanum":
				customErrors = append(customErrors, fmt.Sprintf("Field '%s' must be alphanumeric", err.Field()))
			case "max":
				customErrors = append(customErrors, fmt.Sprintf("Field '%s' cannot be longer than %s characters", err.Field(), err.Param()))
			default:
				customErrors = append(customErrors, fmt.Sprintf("Field '%s' validation failed with tag '%s'", err.Field(), err.Tag()))
			}
		}

		return fmt.Errorf("Validation error: %s", strings.Join(customErrors, "; "))
	}
	return nil
}

func ContainsStringInSlice(slice []string, str string) bool {
	for _, item := range slice {
		if item == str {
			return true
		}
	}
	return false
}

func GetCurrentUser(ctx echo.Context) (models.CurrentUserModels, error) {
	currentUser, ok := ctx.Get("user").(models.CurrentUserModels)
	if !ok {
		return currentUser, errors.New("failed to get user from context")
	}

	return currentUser, nil
}

func AsymetricSignature(input, key string) string {

	var rawkey []byte

	fmt.Println("STRING KEY : ", key)

	privPem, errs := pem.Decode([]byte(key))
	if privPem == nil {
		fmt.Println(string(errs))
		fmt.Println(privPem)
		return "error"
	}

	if privPem.Type != "PRIVATE KEY" {
		fmt.Printf("RSA private key is of the wrong type :%s\n", privPem.Type)
	}

	key_for_sign, err := parsePrivateKey(privPem.Bytes)

	// key_for_sign, err := parsePrivateKey([]byte(key))

	if err != nil {
		fmt.Println("ERROR PARSING RSA KEY : ", err.Error())

		return ""
	}

	h := sha256.New()
	h.Write([]byte(input))
	d := h.Sum(nil)

	fmt.Println("request raw", input)
	fmt.Println("request hashed : ", d)

	rawkey, err = rsa.SignPKCS1v15(rand.Reader, key_for_sign, crypto.SHA256, d)

	if err != nil {
		fmt.Println("ERROR ENCRYPT SIGNATURE : ", err.Error())
		return ""
	}
	base64key := base64.StdEncoding.EncodeToString(rawkey)

	return base64key
}

func parsePrivateKey(pemBytes []byte) (*rsa.PrivateKey, error) {
	privatekey, err := x509.ParsePKCS8PrivateKey(pemBytes)

	privatekeyrsa, _ := privatekey.(*rsa.PrivateKey)

	return privatekeyrsa, err
}
