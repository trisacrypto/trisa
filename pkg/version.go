/*
Package pkg describes the TRISA Go package (reference implementation).
*/
package pkg

import "fmt"

// Version component constants for the current build.
const (
	VersionMajor         = 0
	VersionMinor         = 5
	VersionPatch         = 0
	VersionReleaseLevel  = "v1beta1"
	VersionReleaseNumber = 10
)

// Version returns the semantic version for the current build.
func Version() string {
	var versionCore string
	if VersionPatch > 0 {
		versionCore = fmt.Sprintf("%d.%d.%d", VersionMajor, VersionMinor, VersionPatch)
	} else {
		versionCore = fmt.Sprintf("%d.%d", VersionMajor, VersionMinor)
	}

	if VersionReleaseLevel != "" {
		if VersionReleaseNumber > 0 {
			return fmt.Sprintf("%s (%s revision %d)", versionCore, VersionReleaseLevel, VersionReleaseNumber)
		}
		return fmt.Sprintf("%s (%s)", versionCore, VersionReleaseLevel)
	}
	return versionCore
}
