package utils

func WipeData(data []byte) {
	for i := range data {
		data[i] = 0
	}
}
