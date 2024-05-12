package blokchain

import "regexp"

func IsValidEthereumAddress(address string) bool {
	var validEnvName = regexp.MustCompile(`^0x[0-9A-Fa-f]{40}$`)
	return validEnvName.MatchString(address)
}
