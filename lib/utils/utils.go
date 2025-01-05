package utils

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Id = uint64
type Uuid = string

func MakeUuid() Uuid {
	return strings.ReplaceAll(uuid.New().String(), "-", "")
}

// Basic time defined in milliseconds.
type Time = int64

// We consider standard time.Time as a Date, other mentions of Time are ms
// timestamps.
type Date = time.Time
type Dict = map[string]any

// NOTE:
// If we use `Unwrap(e error)` signature, somehow golang behaves this way:
// if you pass nil pointer to the function, inside the function it will become
// non-nil. We don't know why this happens yet. A solution to use pointer
// to actual struct, this will limit us from passing generic errors, but will
// resolve the problem for the time being.
func Unwrap(e *Error) {
	if e != nil {
		panic(e.Error())
	}
}

// Reference: https://stackoverflow.com/a/13295158/14748231
func TimeToDate(t Time) (Date, error) {
	return time.Unix(0, t*int64(time.Millisecond)), nil
}

func TimeNow() Time {
	return DateNow().UnixMilli()
}

func DateNow() Date {
	return time.Now()
}

type Code = int16
type Error struct {
	code    Code
	message string
}

func (e *Error) Error() string {
	text := fmt.Sprintf("[Error %d] %s", e.code, TranslateCode(e.code))
	if len(e.message) > 0 {
		text += ": " + e.message
	}
	return text
}

func (e *Error) Message() string {
	return e.message
}

func (e *Error) Code() Code {
	return e.code
}

func (e *Error) IsCode(anycode ...Code) bool {
	for _, code := range anycode {
		if e.code == code {
			return true
		}
	}
	return false
}

// Convert from one error to another using conversion map.
// If code is not found in the conversion map, default error is created and
// returned.
func (e *Error) Convert(conversion map[Code]Code) *Error {
	target, ok := conversion[e.code]
	if !ok {
		return CreateDefaultError("")
	}
	return CreateError(target, "")
}

func CreateError(code Code, message string) *Error {
	return &Error{code, message}
}

func CreateDefaultErrorFromBase(be error) *Error {
	return &Error{1, be.Error()}
}

func CreateDefaultError(message string) *Error {
	return CreateError(1, message)
}

type Locale = string
type TranslationKey = string

var translationMap = map[Locale]map[TranslationKey]string{}
var translationLocale string = "en"

const CodeOk Code = 0
const CodeError Code = 1

// Register a translation from a CSV file.
// CSV file structure:
// key(string),text(string)
//
// This function can be called many times, each new call the old matching
// entries will be overwritten.
//
// Text may contain placeholders in form of `%` to accept incoming value,
// which will always be converted to string.
//
// For list of locales refer to https://docs.godotengine.org/en/4.3/tutorials/i18n/locales.html
func LoadTranslationCsv(path string, locale Locale, delimiter rune) *Error {
	locale = strings.ToLower(locale)

	file, e := os.Open(path)
	if e != nil {
		return CreateDefaultError("")
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = delimiter
	records, e := reader.ReadAll()
	if e != nil {
		return CreateDefaultError("")
	}

	localeMap, ok := translationMap[locale]
	if !ok {
		localeMap = map[TranslationKey]string{}
		translationMap[locale] = localeMap
	}

	for i, record := range records {
		if len(record) != 2 {
			return CreateError(CodeError, "")
		}
		if i == 0 {
			continue
		}
		localeMap[strings.TrimSpace(record[0])] = strings.TrimSpace(record[1])
	}

	return nil
}

// Codes are translated using keys `CODE_%`, where `%` is the number.
func TranslateCode(code Code) string {
	return Translate(fmt.Sprintf("CODE_%d", code))
}

func Translate(key TranslationKey, args ...any) string {
	key = strings.ToUpper(key)
	localeMap, ok := translationMap[translationLocale]
	if !ok {
		return key
	}
	text, ok := localeMap[strings.ToUpper(key)]
	if !ok {
		return key
	}
	return text
}

func Sleep(duration Time) {
	time.Sleep(time.Duration(duration * 1000))
}

// Logging implementation.
func Log(message string) {
	// TODO: write to sink, which may be stderr
	fmt.Println(message)
}

// Order is not important.
// https://stackoverflow.com/a/37335777/14748231
func RemoveFromUnorderedSlice[T any](s []T, i int) []T {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func PowInt(x, y int) int {
	return int(math.Pow(float64(x), float64(y)))
}

func IsAlnum(s string) bool {
	match, _ := regexp.MatchString("^[a-zA-Z0-9]+$", s)
	return match
}

func IsFloat(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

func IsInt(s string) bool {
	_, err := strconv.ParseInt(s, 10, 64)
	return err == nil
}

func RemoveSpaces(s string) string {
	re := regexp.MustCompile(`\s+`)
	return re.ReplaceAllString(s, "")
}

func WrapString(s string, wrapper string) string {
	return fmt.Sprintf("%s%s%s", wrapper, s, wrapper)
}
