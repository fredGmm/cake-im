package path

import (
	"log"
	"testing"
)

func TestGetConfigPath(t *testing.T) {

	configPath, err := GetConfigPath("")

	if err != nil {
		t.Error(err)
	}
	log.Println(configPath)
	log.Printf("default success")

	configPath2,err2 := GetConfigPath("E:\\code\\golang\\src\\imChong\\config\\Config.ini")
	if err2 != nil {
		t.Error(err2)
	}
	log.Println(configPath2)
	log.Printf("success")

	configPath3,err3 := GetConfigPath("E:\\code\\golang\\src\\imChong\\configxxx\\Config.ini")
	if err3 == nil {
		t.Error("异常情况未捕捉")
	}
	log.Println(configPath3)
}
