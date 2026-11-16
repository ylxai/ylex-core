// Copyright 2015 by the Authors
// This file is part of the go-core library.
//
// The go-core library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-core library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-core library. If not, see <http://www.gnu.org/licenses/>.

package common

import (
	"bytes"
	"database/sql/driver"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"reflect"
	"strconv"
	"strings"

	"github.com/ylxai/ylex-core/v2/common/hexutil"
)

// Lengths of hashes and addresses in bytes.
const (
	// HashLength is the expected length of the hash
	HashLength = 32
	// AddressLength is the expected length of the address
	AddressLength = 22
)

// DefaultNetworkID is the network id whose address prefix is used when
// deriving and validating addresses. YLEX uses the "ylx" prefix.
var DefaultNetworkID = YlexNetworkID

var (
	hashT    = reflect.TypeOf(Hash{})
	addressT = reflect.TypeOf(Address{})
)

var invalidPrefix = errors.New("Invalid network id prefix in address")
var invalidChecksum = errors.New("Invalid checksum in address")

// Hash represents the 32 byte SHA3 hash of arbitrary data.
type Hash [HashLength]byte

// BytesToHash sets b to hash.
// If b is larger than len(h), b will be cropped from the left.
func BytesToHash(b []byte) Hash {
	var h Hash
	h.SetBytes(b)
	return h
}

// BigToHash sets byte representation of b to hash.
// If b is larger than len(h), b will be cropped from the left.
func BigToHash(b *big.Int) Hash { return BytesToHash(b.Bytes()) }

// HexToHash sets byte representation of s to hash.
// If b is larger than len(h), b will be cropped from the left.
func HexToHash(s string) Hash { return BytesToHash(FromHex(s)) }

// Bytes gets the byte representation of the underlying hash.
func (h Hash) Bytes() []byte { return h[:] }

// Big converts a hash to a big integer.
func (h Hash) Big() *big.Int { return new(big.Int).SetBytes(h[:]) }

// Hex converts a hash to a hex string.
func (h Hash) Hex() string { return hexutil.Encode(h[:]) }

// TerminalString implements log.TerminalStringer, formatting a string for console
// output during logging.
func (h Hash) TerminalString() string {
	return fmt.Sprintf("%x…%x", h[:3], h[29:])
}

// String implements the stringer interface and is used also by the logger when
// doing full logging into a file.
func (h Hash) String() string {
	return h.Hex()
}

// Format implements fmt.Formatter.
// Hash supports the %v, %s, %v, %x, %X and %d format verbs.
func (h Hash) Format(s fmt.State, c rune) {
	hexb := make([]byte, 2+len(h)*2)
	copy(hexb, "0x")
	hex.Encode(hexb[2:], h[:])

	switch c {
	case 'x', 'X':
		if !s.Flag('#') {
			hexb = hexb[2:]
		}
		if c == 'X' {
			hexb = bytes.ToUpper(hexb)
		}
		fallthrough
	case 'v', 's':
		s.Write(hexb)
	case 'q':
		q := []byte{'"'}
		s.Write(q)
		s.Write(hexb)
		s.Write(q)
	case 'd':
		fmt.Fprint(s, ([len(h)]byte)(h))
	default:
		fmt.Fprintf(s, "%%!%c(hash=%x)", c, h)
	}
}

// UnmarshalText parses a hash in hex syntax.
func (h *Hash) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("Hash", input, h[:])
}

// UnmarshalJSON parses a hash in hex syntax.
func (h *Hash) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(hashT, input, h[:])
}

// MarshalText returns the hex representation of h.
func (h Hash) MarshalText() ([]byte, error) {
	return hexutil.Bytes(h[:]).MarshalText()
}

// SetBytes sets the hash to the value of b.
// If b is larger than len(h), b will be cropped from the left.
func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}

	copy(h[HashLength-len(b):], b)
}

// Generate implements testing/quick.Generator.
func (h Hash) Generate(rand *rand.Rand, size int) reflect.Value {
	m := rand.Intn(len(h))
	for i := len(h) - 1; i > m; i-- {
		h[i] = byte(rand.Uint32())
	}
	return reflect.ValueOf(h)
}

// Scan implements Scanner for database/sql.
func (h *Hash) Scan(src interface{}) error {
	srcB, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("can't scan %T into Hash", src)
	}
	if len(srcB) != HashLength {
		return fmt.Errorf("can't scan []byte of len %d into Hash, want %d", len(srcB), HashLength)
	}
	copy(h[:], srcB)
	return nil
}

// Value implements valuer for database/sql.
func (h Hash) Value() (driver.Value, error) {
	return h[:], nil
}

// ImplementsGraphQLType returns true if Hash implements the specified GraphQL type.
func (Hash) ImplementsGraphQLType(name string) bool { return name == "Bytes32" }

// UnmarshalGraphQL unmarshals the provided GraphQL query data.
func (h *Hash) UnmarshalGraphQL(input interface{}) error {
	var err error
	switch input := input.(type) {
	case string:
		err = h.UnmarshalText([]byte(input))
	default:
		err = fmt.Errorf("unexpected type %T for Hash", input)
	}
	return err
}

// UnprefixedHash allows marshaling a Hash without 0x prefix.
type UnprefixedHash Hash

// UnmarshalText decodes the hash from hex. The 0x prefix is optional.
func (h *UnprefixedHash) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedUnprefixedText("UnprefixedHash", input, h[:])
}

// MarshalText encodes the hash as hex.
func (h UnprefixedHash) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(h[:])), nil
}

/////////// Address

// Address represents the 22 byte address of an Core account.
type Address [AddressLength]byte

// PrecompiledContracts contains the default set of pre-compiled Core
// contracts used in the release.
var (
	Addr0 = BytesToAddress([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0})
	Addr1 = BytesToAddress([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1})
	Addr2 = BytesToAddress([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2})
	Addr3 = BytesToAddress([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 3})
	Addr4 = BytesToAddress([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4})
	Addr5 = BytesToAddress([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5})
	Addr6 = BytesToAddress([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 6})
	Addr7 = BytesToAddress([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 7})
	Addr8 = BytesToAddress([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 8})
	Addr9 = BytesToAddress([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 9})

	DevAddress = BytesToAddress([]byte{0xcb, 0x03, 0xa5, 0xfd, 0x22, 0xb9, 0xbe, 0xe8, 0xb8, 0xab, 0x87, 0x7c, 0x86, 0xe0, 0xa2, 0xc2, 0x17, 0x65, 0xe1, 0xd5, 0xbf, 0xc5}) // cb03a5fd22b9bee8b8ab877c86e0a2c21765e1d5bfc5
)

func isPrecompiledAddress(addr Address) bool {
	switch addr {
	case Addr0, Addr1, Addr2, Addr3, Addr4, Addr5, Addr6, Addr7, Addr8, Addr9:
		return true
	}
	return false
}

func isDevAddress(addr Address) bool {
	switch addr {
	case DevAddress:
		return true
	}
	return false
}

func CalculateChecksum(address, prefix []byte) string {
	addrString := strings.ToUpper(Bytes2Hex(append(address, prefix...)) + "00")
	// Replace each letter in the string with two digits, thereby expanding the string, where A = 10, B = 11, ..., F = 15
	mods := ""
	for _, c := range addrString {
		// Get character code point value
		i := int(c)
		// Check if c is characters A-F (codepoint 65 - 90)
		if i > 64 && i < 91 {
			// A=10, B=11 etc...
			i -= 55
			// Add int as string to mod string
			mods += strconv.Itoa(i)
		} else {
			mods += string(c)
		}
	}

	// Create bignum from mod string and perform module
	bigVal, success := new(big.Int).SetString(mods, 10)
	if !success {
		return ""
	}

	val97 := new(big.Int).SetInt64(97)
	val98 := new(big.Int).SetInt64(98)

	remainder := new(big.Int).Mod(bigVal, val97)

	checksum := new(big.Int).Sub(val98, remainder)

	resInt := checksum.Int64()
	if resInt < 10 {
		return "0" + strconv.Itoa(int(resInt))
	}
	return strconv.Itoa(int(resInt))
}

func VerifyAddress(addr Address) error {
	if isPrecompiledAddress(addr) || isDevAddress(addr) {
		return nil
	}
	prefixLen := len(DefaultNetworkID.Bytes())
	if prefixLen == 0 || len(addr) < prefixLen+1 {
		return invalidPrefix
	}
	if !bytes.Equal(addr[:prefixLen], DefaultNetworkID.Bytes()) {
		return invalidPrefix
	} else if Bytes2Hex(addr[prefixLen:prefixLen+1]) != CalculateChecksum(addr[prefixLen+1:], addr[:prefixLen]) {
		return invalidChecksum
	}
	return nil
}

// BytesToAddress returns Address with value b.
// If b is larger than len(h), b will be cropped from the left.
func BytesToAddress(b []byte) Address {
	var a Address
	a.SetBytes(b)
	return a
}

// BigToAddress returns Address with byte values of b.
// If b is larger than len(h), b will be cropped from the left.
func BigToAddress(b *big.Int) Address { return BytesToAddress(b.Bytes()) }

// HexToAddress returns Address with byte values of s.
// If s is larger than len(h), s will be cropped from the left.
func HexToAddress(s string) (Address, error) {
	hexS, err := Hex2BytesWithError(s)
	if err != nil {
		return Address{}, err
	}
	addr := BytesToAddress(hexS)

	if err := VerifyAddress(addr); err != nil {
		return addr, err
	}
	return addr, nil
}

// IsHexAddress verifies whether a string can represent a valid hex-encoded
// Core address or not.
func IsHexAddress(s string) bool {
	if has0xPrefix(s) {
		s = s[2:]
	}
	return len(s) == 2*AddressLength && isHex(s)
}

// Bytes gets the string representation of the underlying address.
func (a Address) Bytes() []byte {
	return a[:]
}

// YlexGenesisCoinbase is the protocol address credited with the genesis
// coinbase on YLEX.
//
// It is derived the same way every other address is: a network prefix, a
// checksum and a payload. The payload is all zeroes, the conventional
// "unspendable" marker.
//
// It is built in init() because it depends on networkIdPrefixes, and Go
// initialises package level variables in source order.
var YlexGenesisCoinbase Address

func init() {
	EnsureNetworkPrefixes()
	YlexGenesisCoinbase = buildYlexAddress(make([]byte, AddressLength-len(networkIdPrefixes[Ylex])-1))
}

// EnsureNetworkPrefixes makes sure the network prefix table is populated.
//
// It is idempotent and safe to call from any package. It exists because a
// package level variable in another package can be initialised before this
// package's init function has run, and reading the table before then would
// see an empty map.
func EnsureNetworkPrefixes() {
	if len(networkIdPrefixes) == 0 {
		networkIdPrefixes = map[NetworkID][]byte{
			Mainnet: {0xcb},
			Devin:   {0xab},
			Ylex:    {'y', 'l', 'x'},
		}
	}
}

// buildYlexAddress assembles an address for the Ylex network from a payload.
func buildYlexAddress(payload []byte) Address {
	prefix := networkIdPrefixes[Ylex]
	raw := make([]byte, 0, AddressLength)
	raw = append(raw, prefix...)
	raw = append(raw, Hex2Bytes(CalculateChecksum(payload, prefix))...)
	raw = append(raw, payload...)
	return BytesToAddress(raw)
}

// TrimAddressPayload trims a hash to the number of bytes an address payload
// holds for the current network.
//
// An address is prefix || checksum || payload. The prefix length depends on
// the network: the upstream networks use a one byte prefix, while Ylex uses
// the three ASCII bytes "ylx". The payload takes whatever is left of the 22
// byte address.
func TrimAddressPayload(hash []byte) []byte {
	prefixLen := len(DefaultNetworkID.Bytes())
	payloadLen := AddressLength - prefixLen - 1 // minus prefix and checksum
	if payloadLen < 0 {
		payloadLen = 0
	}
	if len(hash) > payloadLen {
		return hash[len(hash)-payloadLen:]
	}
	return hash
}

// Hash converts an address to a hash by left-padding it with zeros.
func (a Address) Hash() Hash { return BytesToHash(a[:]) }

// Hex returns hex string representation of the address.
func (a Address) Hex() string {
	return hex.EncodeToString(a[:])
}

// String implements fmt.Stringer.
func (a Address) String() string {
	return a.Hex()
}

// Format implements fmt.Formatter, forcing the byte slice to be formatted as is,
// without going through the stringer interface used for logging.
func (a Address) Format(s fmt.State, c rune) {
	switch c {
	case 'v', 's':
		s.Write([]byte(a.Hex()))
	case 'q':
		q := []byte{'"'}
		s.Write(q)
		s.Write([]byte(a.Hex()))
		s.Write(q)
	case 'x', 'X':
		// %x disables the checksum.
		hex := a.Hex()
		s.Write([]byte(hex))
	case 'd':
		fmt.Fprint(s, ([len(a)]byte)(a))
	default:
		fmt.Fprintf(s, "%%!%c(address=%x)", c, a)
	}
}

// SetBytes sets the address to the value of b.
// If b is larger than len(a), b will be cropped from the left.
func (a *Address) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

// MarshalText returns the hex representation of a.
func (a Address) MarshalText() ([]byte, error) {
	return []byte(a.Hex()), nil
}

// UnmarshalText parses a hash in hex syntax.
func (a *Address) UnmarshalText(input []byte) error {
	addr, err := HexToAddress(Bytes2Hex(input))
	if err != nil {
		return err
	}
	copy(a[:], addr[:])
	return nil
}

// UnmarshalJSON parses a hash in hex syntax.
func (a *Address) UnmarshalJSON(input []byte) error {
	input = []byte(strings.Replace(string(input), `"`, "", -1))
	addr, err := HexToAddress(string(input))
	if err != nil {
		return err
	}
	copy(a[:], addr[:])
	return nil
}

// Scan implements Scanner for database/sql.
func (a *Address) Scan(src interface{}) error {
	srcB, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("can't scan %T into Address", src)
	}
	if len(srcB) != AddressLength {
		return fmt.Errorf("can't scan []byte of len %d into Address, want %d", len(srcB), AddressLength)
	}
	copy(a[:], srcB)
	return nil
}

// Value implements valuer for database/sql.
func (a Address) Value() (driver.Value, error) {
	return a[:], nil
}

// ImplementsGraphQLType returns true if Hash implements the specified GraphQL type.
func (a Address) ImplementsGraphQLType(name string) bool { return name == "Address" }

// UnmarshalGraphQL unmarshals the provided GraphQL query data.
func (a *Address) UnmarshalGraphQL(input interface{}) error {
	var err error
	switch input := input.(type) {
	case string:
		*a, err = HexToAddress(input)
	default:
		err = fmt.Errorf("unexpected type %T for Address", input)
	}
	return err
}

// UnprefixedAddress allows marshaling an Address without 0x prefix.
type UnprefixedAddress Address

// UnmarshalText decodes the address from hex. The 0x prefix is optional.
func (a *UnprefixedAddress) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedUnprefixedText("UnprefixedAddress", input, a[:])
}

// MarshalText encodes the address as hex.
func (a UnprefixedAddress) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(a[:])), nil
}

type NetworkID uint64

const (
	Mainnet NetworkID = 1
	Devin   NetworkID = 3
	Ylex    NetworkID = 4
)

// The typed aliases below keep the constants usable as NetworkID values in
// contexts where the untyped constant would be inferred as int.
const (
	YlexNetworkID NetworkID = Ylex
)

var networkIds = [...]string{
	"cb",       // 1 - Mainnet
	"wrong id", // 2 - wrong
	"ab",       // 3 - Devin
	"ylx",      // 4 - Ylex
}

// networkIdPrefixes holds the raw byte prefix for each network id. The
// prefix is what makes an address recognisable at a glance and what keeps
// addresses from one network from being valid on another.
//
// Ylex uses the three ASCII bytes "ylx" rather than a hex pair, because
// the prefix is displayed verbatim in the address string. The upstream
// networks use two hex characters, so their addresses stay 22 bytes with a
// one byte prefix; Ylex therefore uses a three byte prefix and a 22 byte
// address with a shorter payload. The checksum is computed over the prefix
// and the payload exactly as it is for the other networks.
var networkIdPrefixes = map[NetworkID][]byte{
	Mainnet: {0xcb},
	Devin:   {0xab},
	Ylex:    {'y', 'l', 'x'},
}

func (network NetworkID) String() string {
	if network > 3 {
		return "ce" // everything else except mainnet, devin and ylex - private nets
	}
	return networkIds[network-1]
}

// Bytes returns the raw address prefix for this network id.
//
// Every network on this client shares the three byte "ylx" prefix: mainnet,
// the test network and private networks all derive and validate addresses the
// same way. Only the chain id differs between them, which is what keeps a
// transaction from one network valid on another. Sharing the prefix means a
// wallet, an explorer or a block explorer URL works unchanged on every
// network the client can join, and it avoids the failure mode where a node
// silently switches prefix depending on which network id it was started with.
//
// The networkIdPrefixes table is kept so the prefix can be looked up by name
// in tests and documentation, but address derivation does not depend on the
// running network id any more.
func (network NetworkID) Bytes() []byte {
	return networkIdPrefixes[Ylex]
}

func (network NetworkID) Int64() int64 {
	return int64(network)
}
