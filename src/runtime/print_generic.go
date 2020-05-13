// +build !nintendoswitch

package runtime

//go:nobounds
func printstring(s string) {
	for i := 0; i < len(s); i++ {
		putchar(s[i])
	}
}
