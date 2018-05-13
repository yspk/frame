package sensitive_test

import (
	"fmt"
	"coding.net/baoquan2017/candy-backend/src/common/sensitive"
)

func ExampleFilter() {
	filter := sensitive.New()
	filter.LoadWordDict("dict/dict.txt")
	filter.AddWord("长者")

	fmt.Println(filter.Filter("我为长者续一秒"))
	fmt.Println(filter.Replace("我为长者续一秒", 42))
	fmt.Println(filter.FindIn("我为长者续一秒"))
	fmt.Println(filter.FindIn("我为长 者续一秒"))
}
