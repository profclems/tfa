// This package generates time-based one-time password
// all []byte in this program are treated as Big Endian
package totp

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type OTP struct {
	Name     string
	Code     string
	Key      []byte
	Password uint32
	Timer    int64
}

func toBytes(value int64) []byte {
	mask := int64(0xFF)
	shifts := [8]uint16{56, 48, 40, 32, 24, 16, 8, 0}
	result := make([]byte, 8)

	for i, shift := range shifts {
		result[i] = byte((value >> shift) & mask)
	}
	return result
}

func toUint32(bytes []byte) uint32 {
	return (uint32(bytes[0]) << 24) + (uint32(bytes[1]) << 16) +
		(uint32(bytes[2]) << 8) + uint32(bytes[3])
}

func oneTimePassword(key []byte, value []byte) (uint32, error) {
	// sign the value using HMAC-SHA1
	hmacSha1 := hmac.New(sha1.New, key)
	_, err := hmacSha1.Write(value)
	if err != nil {
		return 0, err
	}
	hash := hmacSha1.Sum(nil)

	// Using a subset of the generated hash.
	// Using the last nibble (half-byte) to choose the index to start from.
	// This number is always appropriate as it's maximum decimal 15, the hash will
	// have the maximum index 19 (20 bytes of SHA1) and we need 4 bytes.
	offset := hash[len(hash)-1] & 0x0F

	// get a 32-bit (4-byte) chunk from the hash starting at offset
	hashParts := hash[offset : offset+4]

	// ignore the most significant bit as per RFC 4226
	// http://tools.ietf.org/html/rfc4226
	hashParts[0] &= 0x7F

	number := toUint32(hashParts)

	// size to 6 digits
	// one million is the first number with 7 digits so the remainder
	// of the division will always return < 7 digits
	pwd := number % 1000000

	return pwd, nil
}

func (otp *OTP) Refresh() (err error) {
	// generate a one-time password using the time at 30-second intervals
	epochSeconds := time.Now().Unix()
	otp.Password, err = oneTimePassword(otp.Key, toBytes(epochSeconds/30))
	if err != nil {
		return err
	}

	otp.Timer = 30 - (epochSeconds % 30)

	// fmt.Printf("%06d (%d second(s) remaining)\n", pwd, secondsRemaining)
	return nil
}

func New(name, code string) (otp *OTP, err error) {
	code = strings.Replace(code, " ", "", -1)
	code = strings.ToUpper(code)
	c, err := strconv.Unquote(code)
	if err == nil {
		code = c
	}

	otp = &OTP{
		Name: name,
		Code: code,
	}

	// decode the key
	// TODO: base32 encoding expects A-Z, 2-7 and =. Validate code before decoding
	otp.Key, err = base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(otp.Code)
	if err != nil {
		return nil, fmt.Errorf("base32 decoding failed: %w", err)
	}

	err = otp.Refresh()
	if err != nil {
		return nil, err
	}

	return otp, nil
}
