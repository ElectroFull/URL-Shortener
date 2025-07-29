package helpers

const alphabet = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func Base62Encode(number int) string {
	var encoding []byte
	for number > 0 {
		encoding = append(encoding, alphabet[number%62])
		number /= 62
	}
	for i, j := 0, len(encoding)-1; i < j; i, j = i+1, j-1 {
		encoding[i], encoding[j] = encoding[j], encoding[i]
	}
	return string(encoding)
}
