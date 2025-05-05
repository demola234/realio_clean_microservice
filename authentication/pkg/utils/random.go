package utils

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomString generates a random string of length n
func RandomString(n int) string {
	var letter = []rune("abcdefghijklmnopqrstuvwxyz")

	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

func RandomOtp() string {
	return fmt.Sprintf("%06v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
}

// RandomInt generates a random integer between min and max
func RandomInt(min, max int) int {
	return rand.Intn(max-min) + min
}

// RandomOwner generates a random owner name
func RandomOwner() string {
	return RandomString(6)
}

func RandomRole() string {
	roles := []string{"admin", "user", "buyer", "seller"}
	return roles[rand.Intn(len(roles))]
}

func RandomProfilePicture() string {
	return fmt.Sprintf("https://picsum.photos/200/300?random=%d", rand.Intn(1000))
}

func RandomBio() string {
	return fmt.Sprintf("Bio %d", rand.Intn(1000))
}

func RandomProvider() string {
	providers := []string{"google", "facebook", "twitter", "github", "local"}
	return providers[rand.Intn(len(providers))]
}

// RandomMoney generates a random amount of money
func RandomMoney() int64 {
	return int64(RandomInt(0, 1000))
}

func RandomEmail() string {
	return fmt.Sprintf("%s@demola.dev", RandomString(6))
}

func SplitStrings(s string) []string {
	var r []string
	for _, v := range s {
		r = append(r, string(v))
	}
	return r
}

func RandomPhoneNumber() string {
	// Nigerian Phone Number
	// should start with 0 and have 11 digits
	// e.g 08012345678
	randomInt := RandomInt(10000000, 99999999)
	phoneNumber := strconv.Itoa(randomInt)

	return fmt.Sprintf("080%s", phoneNumber)
}

// enums -> email, google, facebook


// roles -> admin, user, buyer, seller
