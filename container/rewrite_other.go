//go:build !amd64

package container

// _RewriteWide is the portable stand-in for the amd64 kernels. There is no
// wider step on this architecture, so the word path does the whole tail.
//
// Revisions:
//   - 2026-10-13 17:25: initial creation
func _RewriteWide(dst []byte, src []byte, span *_CaseSpan) int {
	return 0
}
