package service_test

import (
	"log"
	"testing"

	"github.com/Vasily-van-Zaam/ushortener/internal/service"
)

func TestUtil(t *testing.T) {
	res := service.Convert62(1161)
	log.Println(res)
}
