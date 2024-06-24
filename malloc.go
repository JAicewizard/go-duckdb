//go:build !windows

package duckdb

/*
#include <stdlib.h>
*/
import "C"

type mallocT = C.ulong
