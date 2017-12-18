package blockchain

import (
	"fmt"
	"log"

	"github.com/blang/semver"
)

const BlockchainVersion = "0.0.1" // TODO autogenerate

func Version() semver.Version {
	v, _ := semver.Make(BlockchainVersion)
	return v
}

func IsCompatibleWithCurrent(challenge string) bool {
	return IsCompatible(BlockchainVersion, challenge)
}

func IsCompatible(version string, challenge string) bool {
	v, err := semver.Make(version)
	if err != nil {
		log.Println(fmt.Sprintf("Error parsing semantic version: %s with error %s", v, err.Error()))
		return false
	}

	cv, err := semver.Make(challenge)
	if err != nil {
		log.Println(fmt.Sprintf("Error parsing semantic version: %s with error %s", challenge, err.Error()))
		return false
	}

	if v.Major != cv.Major {
		return false
	}

	result := v.Compare(cv)

	if result >= 0 {
		return true
	}
	return false
}
