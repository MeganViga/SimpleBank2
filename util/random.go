package util

import (
	"math/rand"
	"strings"
	"time"
)
func RandInt(min,max int)int{
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max - min + 1) + min
}

func RandString(lenght int) string{
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	var str strings.Builder
	for i :=0;i < lenght;i++{
		rand.Seed(time.Now().UnixNano())
		c := letterBytes[rand.Intn((len(letterBytes)))]
		str.WriteByte(c)
	}
	return str.String()
}

func RandomOwner()string{
	return RandString(6)
}

func RandomMoney()int{
	return RandInt(0,1000)
}

func RandomCurrency()string{
	currency := []string{"INR","USD","SGP"}
	rand.Seed(time.Now().UnixNano())
	return currency[rand.Intn(len(currency))]
}
