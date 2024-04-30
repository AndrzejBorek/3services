package utils

import (
	"testing"
)

func BenchmarkRandStringBytesMaskImpr(b *testing.B) {
	var gen = createNewCustomRandomGenerator(1)
	input := 9
	for i := 0; i < b.N; i++ {
		_ = randStringBytesMaskImpr(input, gen)
	}
}

func TestRandStringBytesMaskImpr(t *testing.T) {
	for input := 1; input < 10; input++ {
		for i := 0; i < 10000; i++ {
			src1 := createNewCustomRandomGenerator(int64(i))
			var str1 = randStringBytesMaskImpr(input, src1)

			src2 := createNewCustomRandomGenerator(int64(i))
			var str2 = randStringBytesMaskImpr(input, src2)

			if str1 != str2 {
				t.Errorf("str1 != str2 (%s != %s)", str1, str2)
			}
		}
	}
}

func BenchmarkCreateRandomJson(b *testing.B) {
	for j := 1; j < 100000; j++ {
		for i := 0; i < b.N; i++ {
			_ = createRandomJson(int64(j))
		}
	}
}
