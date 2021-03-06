// Package utils is used for the common functions for dealing with
// conversion to and from hex, bytes, and strings, formatting time.
package utils

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
	"time"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/jpillora/backoff"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/sha3"
	null "gopkg.in/guregu/null.v3"
)

const (
	// HumanTimeFormat is the predefined layout for use in Time.Format and time.Parse
	HumanTimeFormat = "2006-01-02 15:04:05 MST"
	// EVMWordByteLen the length of an EVM Word Byte
	EVMWordByteLen = 32
	// EVMWordHexLen the length of an EVM Word Hex
	EVMWordHexLen = EVMWordByteLen * 2
)

var weiPerEth = big.NewInt(1e18)

// ZeroAddress is an empty address, otherwise in Ethereum as
// 0x0000000000000000000000000000000000000000
var ZeroAddress = common.Address{}

// WithoutZeroAddresses returns a list of addresses excluding the zero address.
func WithoutZeroAddresses(addresses []common.Address) []common.Address {
	var withoutZeros []common.Address
	for _, address := range addresses {
		if address != ZeroAddress {
			withoutZeros = append(withoutZeros, address)
		}
	}
	return withoutZeros
}

// HexToUint64 converts a given hex string to 64-bit unsigned integer.
func HexToUint64(hex string) (uint64, error) {
	return strconv.ParseUint(RemoveHexPrefix(hex), 16, 64)
}

// Uint64ToHex converts the given uint64 value to a hex-value string.
func Uint64ToHex(i uint64) string {
	return fmt.Sprintf("0x%x", i)
}

// EncodeTxToHex converts the given Ethereum Transaction type and
// returns its hex-value string.
func EncodeTxToHex(tx *types.Transaction) (string, error) {
	rlp := new(bytes.Buffer)
	if err := tx.EncodeRLP(rlp); err != nil {
		return "", err
	}
	return hexutil.Encode(rlp.Bytes()), nil
}

// ISO8601UTC formats given time to ISO8601.
func ISO8601UTC(t time.Time) string {
	return t.UTC().Format(time.RFC3339)
}

// NullISO8601UTC returns formatted time if valid, empty string otherwise.
func NullISO8601UTC(t null.Time) string {
	if t.Valid {
		return ISO8601UTC(t.Time)
	}
	return ""
}

// FormatJSON applies indent to format a JSON response.
func FormatJSON(v interface{}) ([]byte, error) {
	return json.MarshalIndent(v, "", "  ")
}

// NewBytes32Length holds the length of bytes needed for Bytes32ID.
const NewBytes32Length = 32

// NewBytes32ID returns a randomly generated UUID that conforms to
// Ethereum bytes32.
func NewBytes32ID() string {
	return strings.Replace(uuid.Must(uuid.NewV4()).String(), "-", "", -1)
}

// RemoveHexPrefix removes the prefix (0x) of a given hex string.
func RemoveHexPrefix(str string) string {
	if HasHexPrefix(str) {
		return str[2:]
	}
	return str
}

// HasHexPrefix returns true if the string starts with 0x.
func HasHexPrefix(str string) bool {
	return len(str) >= 2 && str[0] == '0' && str[1] == 'x'
}

// DecodeEthereumTx takes an RLP hex encoded Ethereum transaction and
// returns a Transaction struct with all the fields accessible.
func DecodeEthereumTx(hex string) (types.Transaction, error) {
	var tx types.Transaction
	b, err := hexutil.Decode(hex)
	if err != nil {
		return tx, err
	}
	return tx, rlp.DecodeBytes(b, &tx)
}

// IsEmptyAddress checks that the address is empty, synonymous with the zero
// account/address. No logs can come from this address, as there is no contract
// present there.
//
// See https://stackoverflow.com/questions/48219716/what-is-address0-in-solidity
// for the more info on the zero address.
func IsEmptyAddress(addr common.Address) bool {
	return addr == ZeroAddress
}

// StringToHex converts a standard string to a hex encoded string.
func StringToHex(in string) string {
	return AddHexPrefix(hex.EncodeToString([]byte(in)))
}

// AddHexPrefix adds the previx (0x) to a given hex string.
func AddHexPrefix(str string) string {
	if len(str) < 2 || len(str) > 1 && strings.ToLower(str[0:2]) != "0x" {
		str = "0x" + str
	}
	return str
}

// ToFilterQueryFor returns a struct that encapsulates desired arguments used to filter
// event logs.
func ToFilterQueryFor(fromBlock *big.Int, addresses []common.Address) ethereum.FilterQuery {
	return ethereum.FilterQuery{
		FromBlock: fromBlock,
		Addresses: WithoutZeroAddresses(addresses),
	}
}

// ToFilterArg filters logs with the given FilterQuery
// https://github.com/ethereum/go-ethereum/blob/762f3a48a00da02fe58063cb6ce8dc2d08821f15/ethclient/ethclient.go#L363
func ToFilterArg(q ethereum.FilterQuery) interface{} {
	arg := map[string]interface{}{
		"fromBlock": toBlockNumArg(q.FromBlock),
		"toBlock":   toBlockNumArg(q.ToBlock),
		"address":   q.Addresses,
		"topics":    q.Topics,
	}
	if q.FromBlock == nil {
		arg["fromBlock"] = "0x0"
	}
	return arg
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	return hexutil.EncodeBig(number)
}

// Sleeper interface is used for tasks that need to be done on some
// interval, excluding Cron, like reconnecting.
type Sleeper interface {
	Reset()
	Sleep()
	After() time.Duration
	Duration() time.Duration
}

// BackoffSleeper is a counter to assist with reattempts.
type BackoffSleeper struct {
	*backoff.Backoff
}

// NewBackoffSleeper returns a BackoffSleeper that is configured to
// sleep for 1 second minimum, and 10 seconds maximum.
func NewBackoffSleeper() BackoffSleeper {
	return BackoffSleeper{&backoff.Backoff{
		Min: 1 * time.Second,
		Max: 10 * time.Second,
	}}
}

// Sleep waits for the given duration, incrementing the back off.
func (bs BackoffSleeper) Sleep() {
	time.Sleep(bs.Backoff.Duration())
}

// After returns the duration for the next stop, and increments the backoff.
func (bs BackoffSleeper) After() time.Duration {
	return bs.Backoff.Duration()
}

// Duration returns the current duration value.
func (bs BackoffSleeper) Duration() time.Duration {
	return bs.ForAttempt(bs.Attempt())
}

// ConstantSleeper is to assist with reattempts with
// the same sleep duration.
type ConstantSleeper struct {
	interval time.Duration
}

// NewConstantSleeper returns a ConstantSleeper that is configured to
// sleep for a constant duration based on the input.
func NewConstantSleeper(d time.Duration) ConstantSleeper {
	return ConstantSleeper{interval: d}
}

// Reset is a no op since sleep time is constant.
func (cs ConstantSleeper) Reset() {}

// Sleep waits for the given duration before reattempting.
func (cs ConstantSleeper) Sleep() {
	time.Sleep(cs.interval)
}

// After returns the duration.
func (cs ConstantSleeper) After() time.Duration {
	return cs.interval
}

// Duration returns the duration value.
func (cs ConstantSleeper) Duration() time.Duration {
	return cs.interval
}

// MaxUint64 finds the maximum value of a list of uint64s.
func MaxUint64(uints ...uint64) uint64 {
	var max uint64
	for _, n := range uints {
		if n > max {
			max = n
		}
	}
	return max
}

// MaxInt finds the maximum value of a list of ints.
func MaxInt(ints ...int) int {
	var max int
	for _, n := range ints {
		if n > max {
			max = n
		}
	}
	return max
}

// ConcatBytes appends a bunch of byte arrays into a single byte array
func ConcatBytes(bufs ...[]byte) ([]byte, error) {
	buffer := bytes.NewBuffer([]byte{})
	for _, b := range bufs {
		_, err := buffer.Write(b)
		if err != nil {
			return nil, err
		}
	}
	return buffer.Bytes(), nil
}

// EVMWordUint64 returns a uint64 as an EVM word byte array.
func EVMWordUint64(val uint64) []byte {
	word := make([]byte, EVMWordByteLen)
	binary.BigEndian.PutUint64(word[EVMWordByteLen-8:], val)
	return word
}

// EVMWordSignedBigInt returns a big.Int as an EVM word byte array, with
// support for a signed representation. Returns error on overflow.
func EVMWordSignedBigInt(val *big.Int) ([]byte, error) {
	bytes := val.Bytes()
	if val.BitLen() > (8*EVMWordByteLen - 1) {
		return nil, fmt.Errorf("Overflow saving signed big.Int to EVM word: %v", val)
	}
	if val.Sign() == -1 {
		twosComplement := new(big.Int).Add(val, MaxUint256)
		bytes = new(big.Int).Add(twosComplement, big.NewInt(1)).Bytes()
	}
	return common.LeftPadBytes(bytes, EVMWordByteLen), nil
}

// EVMWordBigInt returns a big.Int as an EVM word byte array, with support for
// a signed representation. Returns error on overflow.
func EVMWordBigInt(val *big.Int) ([]byte, error) {
	if val.Sign() == -1 {
		return nil, errors.New("Uint256 cannot be negative")
	}
	bytes := val.Bytes()
	if len(bytes) > EVMWordByteLen {
		return nil, fmt.Errorf("Overflow saving big.Int to EVM word: %v", val)
	}
	return common.LeftPadBytes(bytes, EVMWordByteLen), nil
}

// CoerceInterfaceMapToStringMap converts map[interface{}]interface{} (interface maps) to
// map[string]interface{} (string maps) and []interface{} with interface maps to string maps.
// Relevant when serializing between CBOR and JSON.
func CoerceInterfaceMapToStringMap(in interface{}) (interface{}, error) {
	switch typed := in.(type) {
	case map[string]interface{}:
		for k, v := range typed {
			coerced, err := CoerceInterfaceMapToStringMap(v)
			if err != nil {
				return nil, err
			}
			typed[k] = coerced
		}
		return typed, nil
	case map[interface{}]interface{}:
		m := map[string]interface{}{}
		for k, v := range typed {
			coercedKey, ok := k.(string)
			if !ok {
				return nil, fmt.Errorf("Unable to coerce key %T %v to a string", k, k)
			}
			coerced, err := CoerceInterfaceMapToStringMap(v)
			if err != nil {
				return nil, err
			}
			m[coercedKey] = coerced
		}
		return m, nil
	case []interface{}:
		r := make([]interface{}, len(typed))
		for i, v := range typed {
			coerced, err := CoerceInterfaceMapToStringMap(v)
			if err != nil {
				return nil, err
			}
			r[i] = coerced
		}
		return r, nil
	default:
		return in, nil
	}
}

// ParseUintHex parses an unsigned integer out of a hex string.
func ParseUintHex(hex string) (*big.Int, error) {
	amount, ok := new(big.Int).SetString(hex, 0)
	if !ok {
		return amount, fmt.Errorf("unable to decode hex to integer: %s", hex)
	}
	return amount, nil
}

// HashPassword wraps around bcrypt.GenerateFromPassword for a friendlier API.
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// CheckPasswordHash wraps around bcrypt.CompareHashAndPassword for a friendlier API.
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// FileExists returns true if a file at the passed string exists.
func FileExists(name string) bool {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return false
	}
	return true
}

// Keccak256 is a simplified interface for the legacy SHA3 implementation that
// Ethereum uses.
func Keccak256(in []byte) ([]byte, error) {
	hash := sha3.NewLegacyKeccak256()
	_, err := hash.Write(in)
	return hash.Sum(nil), err
}

// StripBearer removes the 'Bearer: ' prefix from the HTTP Authorization header.
func StripBearer(authorizationStr string) string {
	return strings.TrimPrefix(strings.TrimSpace(authorizationStr), "Bearer ")
}

// IsQuoted checks if the first and last characters are either " or '.
func IsQuoted(input []byte) bool {
	return len(input) >= 2 &&
		((input[0] == '"' && input[len(input)-1] == '"') ||
			(input[0] == '\'' && input[len(input)-1] == '\''))
}

// RemoveQuotes removes the first and last character if they are both either
// " or ', otherwise it is a noop.
func RemoveQuotes(input []byte) []byte {
	if IsQuoted(input) {
		return input[1 : len(input)-1]
	}
	return input
}

// "Constants" used by EVM words
var (
	maxUint257 = &big.Int{}
	// MaxUint256 represents the largest number represented by an EVM word
	MaxUint256 = &big.Int{}
	// MaxInt256 represents the largest number represented by an EVM word using
	// signed encoding.
	MaxInt256 = &big.Int{}
	// MinInt256 represents the smallest number represented by an EVM word using
	// signed encoding.
	MinInt256 = &big.Int{}
)

func init() {
	maxUint257 = new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)
	MaxUint256 = new(big.Int).Sub(maxUint257, big.NewInt(1))
	MaxInt256 = new(big.Int).Div(MaxUint256, big.NewInt(2))
	MinInt256 = new(big.Int).Neg(MaxInt256)
}
