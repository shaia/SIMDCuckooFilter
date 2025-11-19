//go:build arm64
// +build arm64

#include "textflag.h"

// func bucketLookupNEON(fingerprints []byte, target byte) bool
// Uses ARM NEON SIMD instructions to search for a byte in a slice
// Processes 16 bytes at a time using NEON vector registers
TEXT ·bucketLookupNEON(SB), NOSPLIT, $0-33
	// Load parameters
	MOVD fingerprints_base+0(FP), R0  // R0 = pointer to fingerprints
	MOVD fingerprints_len+8(FP), R1   // R1 = length of fingerprints
	MOVBU target+24(FP), R2            // R2 = target byte

	// Check for empty slice
	CMP $0, R1
	BEQ notfound

	// Duplicate target byte across 128-bit vector
	// V0 will contain 16 copies of the target byte
	VDUP R2, V0.B16

	// Calculate number of 16-byte chunks
	LSR $4, R1, R3                     // R3 = len / 16
	CMP $0, R3
	BEQ handle_remaining               // If less than 16 bytes, handle remainder

loop_16:
	// Load 16 bytes from fingerprints into V1
	VLD1 (R0), [V1.B16]

	// Compare all 16 bytes with target
	// CMEQ sets each byte to 0xFF if equal, 0x00 if not equal
	VCMEQ V0.B16, V1.B16, V2.B16

	// Check if any bytes matched by moving vector to general registers
	// Extract lanes and OR them together
	VMOV V2.D[0], R4
	VMOV V2.D[1], R5
	ORR R4, R5, R6
	CMP $0, R6
	BNE found                          // If any bit is set, we found a match

	// Move to next 16-byte chunk
	ADD $16, R0
	SUB $1, R3
	CMP $0, R3
	BNE loop_16

handle_remaining:
	// Handle remaining bytes (less than 16)
	AND $15, R1, R5                    // R5 = len % 16
	CMP $0, R5
	BEQ notfound                       // No remaining bytes

scalar_loop:
	// Process remaining bytes one at a time
	MOVBU (R0), R6
	CMP R2, R6
	BEQ found

	ADD $1, R0
	SUB $1, R5
	CMP $0, R5
	BNE scalar_loop

notfound:
	MOVD $0, R0
	MOVB R0, ret+32(FP)
	RET

found:
	MOVD $1, R0
	MOVB R0, ret+32(FP)
	RET
