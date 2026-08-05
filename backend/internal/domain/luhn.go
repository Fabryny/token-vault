package domain

// IsValidLuhn aplica o checksum mod 10 (Hans Peter Luhn, IBM, 1954).
func IsValidLuhn(pan string) bool {
	if len(pan) < 13 || len(pan) > 19 {
		return false
	}

	sum := 0
	double := false

	for i := len(pan) - 1; i >= 0; i-- {
		c := pan[i]
		if c < '0' || c > '9' {
			return false
		}

		d := int(c - '0')
		if double {
			d *= 2
			if d > 9 {
				d -= 9
			}
		}
		sum += d
		double = !double
	}

	return sum%10 == 0
}
