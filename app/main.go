/**
 * @file app/main.go
 * @description Provides the compilable tokenizer starter entry point before any challenge stage is implemented.
 * @depends_on none
 * @depended_on_by your_program.sh
 */

package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Fprintln(os.Stderr, "Tokenizer starter: no challenge stage has been implemented yet.")
	os.Exit(1)
}
