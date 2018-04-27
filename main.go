package main

import "fmt"
import "github.com/go-lession/function"

import "io"
import "encoding/json"
import "strings"
import "time"
import (
	"os"
	"reflect"
)

func main1() {
	value := `[
		{name: "Hari", address: "Halberstratur stratse"}
	]`
	fmt.Println(function.Add(2, 3))
	dec := json.NewDecoder(strings.NewReader(value))
	t, err := dec.Token()
	t, err = dec.Token()
	if err != nil {
		fmt.Println("Error occurred", err)
	}
	fmt.Println("Token is", t)
	time.Now().Weekday()
}

func main2() {
	file, err := os.Open("./glide.yaml")
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("file does not exist")
			os.Exit(1)
		}
		fmt.Println("Some other error")
		os.Exit(2)
	}
	defer file.Close()

	for {
		data := make([]byte, 20)
		n, err := file.Read(data)
		if err != nil && err == io.EOF {
			fmt.Println("Reached EOF")
			break
		}
		fmt.Println("No of bytes read", n)
		fmt.Println("Data read", string(data))
	}
	//os.Stat()
}

func main3() {
	ms := &MyStringerImpl{"ok"}
	fmt.Printf("%s\n", ms)

	//var msi interface{} = ms
	var ss fmt.Stringer = ms
	fmt.Printf("%s", ss)

	ms2 := &MyStringerImpl{"ok"}
	var msval OnlyStringer = ms2
	fmt.Println(msval)
}

func main() {
	val := []string{"1", "2"}
	fmt.Println(reflect.TypeOf(val))
	fmt.Println(reflect.TypeOf(val).Kind())
	fmt.Println(reflect.ValueOf(val).IsValid())
}

type OnlyStringer interface {
	String() string
}

type MyStringer interface {
	String() string
	Bling()
}

type MyStringerImpl struct {
	val string
}

func (ms *MyStringerImpl) String() string {
	return fmt.Sprintf("I am %s", ms.val)
}

func (ms *MyStringerImpl) Bling() {
	fmt.Println("Bling")
}

type MyStringerImpl2 struct {
	val string
}

func (ms *MyStringerImpl2) String() string {
	return fmt.Sprintf("2 I am %s", ms.val)
}