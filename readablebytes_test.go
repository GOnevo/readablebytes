package readablebytes

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func ExampleHumanBinarySize() {
	fmt.Println(HumanBinarySize(1024))
	fmt.Println(HumanBinarySize(1024 * 1024))
	fmt.Println(HumanBinarySize(1048576))
	fmt.Println(HumanBinarySize(2 * MiB))
	fmt.Println(HumanBinarySize(3.42 * GiB))
	fmt.Println(HumanBinarySize(5.372 * TiB))
	fmt.Println(HumanBinarySize(2.22 * PiB))
}

func ExampleHumanSize() {
	fmt.Println(HumanSize(1000))
	fmt.Println(HumanSize(1024))
	fmt.Println(HumanSize(1000000))
	fmt.Println(HumanSize(1048576))
	fmt.Println(HumanSize(2 * MB))
	fmt.Println(HumanSize(3.42 * GB))
	fmt.Println(HumanSize(5.372 * TB))
	fmt.Println(HumanSize(2.22 * PB))
}

func ExampleFromHumanString() {
	fmt.Println(FromHumanString("32"))
	fmt.Println(FromHumanString("32b"))
	fmt.Println(FromHumanString("32B"))
	fmt.Println(FromHumanString("32k"))
	fmt.Println(FromHumanString("32K"))
	fmt.Println(FromHumanString("32kb"))
	fmt.Println(FromHumanString("32Kb"))
	fmt.Println(FromHumanString("32Mb"))
	fmt.Println(FromHumanString("32Gb"))
	fmt.Println(FromHumanString("32Tb"))
	fmt.Println(FromHumanString("32Pb"))
}

func TestHumanBinarySize(t *testing.T) {
	assertEquals(t, "1KiB", HumanBinarySize(1024))
	assertEquals(t, "1MiB", HumanBinarySize(1024*1024))
	assertEquals(t, "1MiB", HumanBinarySize(1048576))
	assertEquals(t, "2MiB", HumanBinarySize(2*MiB))
	assertEquals(t, "3.42GiB", HumanBinarySize(3.42*GiB))
	assertEquals(t, "5.372TiB", HumanBinarySize(5.372*TiB))
	assertEquals(t, "2.22PiB", HumanBinarySize(2.22*PiB))
	assertEquals(t, "1.049e+06YiB", HumanBinarySize(KiB*KiB*KiB*KiB*KiB*PiB))
}

func TestHumanDecimalSize(t *testing.T) {
	assertEquals(t, "1kB", HumanDecimalSize(1000))
	assertEquals(t, "1.024kB", HumanDecimalSize(1024))
	assertEquals(t, "1MB", HumanDecimalSize(1000000))
	assertEquals(t, "1.049MB", HumanDecimalSize(1048576))
	assertEquals(t, "2MB", HumanDecimalSize(2*MB))
	assertEquals(t, "3.42GB", HumanDecimalSize(3.42*GB))
	assertEquals(t, "5.372TB", HumanDecimalSize(5.372*TB))
	assertEquals(t, "2.22PB", HumanDecimalSize(2.22*PB))
	assertEquals(t, "1e+04YB", HumanDecimalSize(float64(10000000000000*PB)))
}

func TestHumanSize(t *testing.T) {
	assertEquals(t, "1kB", HumanSize(1000))
	assertEquals(t, "1.024kB", HumanSize(1024))
	assertEquals(t, "1MB", HumanSize(1000000))
	assertEquals(t, "1.049MB", HumanSize(1048576))
	assertEquals(t, "2MB", HumanSize(2*MB))
	assertEquals(t, "3.42GB", HumanSize(3.42*GB))
	assertEquals(t, "5.372TB", HumanSize(5.372*TB))
	assertEquals(t, "2.22PB", HumanSize(2.22*PB))
	assertEquals(t, "1e+04YB", HumanSize(float64(10000000000000*PB)))
}

func TestFromHumanSize(t *testing.T) {
	assertSuccessEquals(t, 32, FromHumanString, "32")
	assertSuccessEquals(t, 32, FromHumanString, "32b")
	assertSuccessEquals(t, 32, FromHumanString, "32B")
	assertSuccessEquals(t, 32*KB, FromHumanString, "32k")
	assertSuccessEquals(t, 32*KB, FromHumanString, "32K")
	assertSuccessEquals(t, 32*KB, FromHumanString, "32kb")
	assertSuccessEquals(t, 32*KB, FromHumanString, "32Kb")
	assertSuccessEquals(t, 32*MB, FromHumanString, "32Mb")
	assertSuccessEquals(t, 32*GB, FromHumanString, "32Gb")
	assertSuccessEquals(t, 32*TB, FromHumanString, "32Tb")
	assertSuccessEquals(t, 32*PB, FromHumanString, "32Pb")

	assertSuccessEquals(t, 32.5*KB, FromHumanString, "32.5kB")
	assertSuccessEquals(t, 32.5*KB, FromHumanString, "32.5 kB")
	assertSuccessEquals(t, 32, FromHumanString, "32.5 B")

	assertError(t, FromHumanString, "")
	assertError(t, FromHumanString, "hello")
	assertError(t, FromHumanString, "-32")
	assertError(t, FromHumanString, ".3kB")
	assertError(t, FromHumanString, " 32 ")
	assertError(t, FromHumanString, "32m b")
	assertError(t, FromHumanString, "32bm")
}

func TestFromHumanString(t *testing.T) {
	assertSuccessEquals(t, 32, FromHumanString, "32")
	assertSuccessEquals(t, 32, FromHumanString, "32b")
	assertSuccessEquals(t, 32, FromHumanString, "32B")
	assertSuccessEquals(t, 32*KiB, FromHumanString, "32ki")
	assertSuccessEquals(t, 32*KiB, FromHumanString, "32Ki")
	assertSuccessEquals(t, 32*KiB, FromHumanString, "32kib")
	assertSuccessEquals(t, 32*KiB, FromHumanString, "32Kib")
	assertSuccessEquals(t, 32*KiB, FromHumanString, "32Kib")
	assertSuccessEquals(t, 32*KiB, FromHumanString, "32KIB")
	assertSuccessEquals(t, 32*MiB, FromHumanString, "32Mib")
	assertSuccessEquals(t, 32*GiB, FromHumanString, "32Gib")
	assertSuccessEquals(t, 32*TiB, FromHumanString, "32Tib")
	assertSuccessEquals(t, 32*PiB, FromHumanString, "32Pib")
	assertSuccessEquals(t, 32*PiB, FromHumanString, "32PiB")
	assertSuccessEquals(t, 32*PiB, FromHumanString, "32Pi")

	assertSuccessEquals(t, 32, FromHumanString, "32.3")
	tmp := 32.3 * MiB
	assertSuccessEquals(t, int64(tmp), FromHumanString, "32.3 mib")

	assertError(t, FromHumanString, "")
	assertError(t, FromHumanString, "hello")
	assertError(t, FromHumanString, "-32")
	assertError(t, FromHumanString, " 32 ")
	assertError(t, FromHumanString, "32m b")
	assertError(t, FromHumanString, "32bm")
}

func assertEquals(t *testing.T, expected, actual interface{}) {
	if expected != actual {
		t.Errorf("Expected '%v' but got '%v'", expected, actual)
	}
}

// func that maps to the parse function signatures as testing abstraction
type parseFn func(string) (int64, error)

// Define 'String()' for pretty-print
func (fn parseFn) String() string {
	fnName := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
	return fnName[strings.LastIndex(fnName, ".")+1:]
}

func assertSuccessEquals(t *testing.T, expected int64, fn parseFn, arg string) {
	res, err := fn(arg)
	if err != nil || res != expected {
		t.Errorf("%s(\"%s\") -> expected '%d' but got '%d' with error '%v'", fn.String(), arg, expected, res, err)
	}
}

func assertError(t *testing.T, fn parseFn, arg string) {
	res, err := fn(arg)
	if err == nil && res != -1 {
		t.Errorf("%s(\"%s\") -> expected error but got '%d'", fn.String(), arg, res)
	}
}
