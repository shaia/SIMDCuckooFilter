//go:build amd64
// +build amd64

#include "textflag.h"

// func containsAVX2(data []uint16, fp uint16) bool
TEXT ·containsAVX2(SB), NOSPLIT, $0-33
	MOVD    data_base+0(FP), SI    // SI = data pointer
	MOVD    data_len+8(FP), BX     // BX = length
	MOVW    fp+24(FP), AX          // AX = fingerprint (uint16)

	// Broadcast fp to YMM0
	MOVQ         AX, X0
	VPBROADCASTW X0, Y0

	// Check length to jump to appropriate handler
	// Check common sizes first (4, 8)
	CMPQ    BX, $4
	JE      len4
	CMPQ    BX, $8
	JE      len8
	CMPQ    BX, $16
	JE      len16
	CMPQ    BX, $32
	JE      len32
	
	// Fallback for other sizes (loop)
	JMP     generic_loop

len32:
	// 32 items = 64 bytes = 2 YMM registers
	// Use memory operand to save load instruction
	VPCMPEQW (SI), Y0, Y1     // Compare first 16 items
	VPMOVMSKB Y1, CX        // Extract mask
	TESTL   CX, CX          // Check if any set
	JNZ     found

	VPCMPEQW 32(SI), Y0, Y1   // Compare next 16 items
	VPMOVMSKB Y1, CX
	TESTL   CX, CX
	JNZ     found
	
	JMP     not_found

len16:
	// 16 items = 32 bytes = 1 YMM register
	VPCMPEQW (SI), Y0, Y1
	VPMOVMSKB Y1, CX
	TESTL   CX, CX
	JNZ     found
	JMP     not_found

len8:
	// 8 items = 16 bytes = 1 XMM register
	// Use X0 (low part of Y0) which has broadcasted fp
	VPCMPEQW (SI), X0, X1
	VPMOVMSKB X1, CX
	TESTL   CX, CX
	JNZ     found
	JMP     not_found

len4:
	// 4 items = 8 bytes. 
	// Use VMOVQ to load 8 bytes into XMM register (zero extends to 128 bits)
	// This avoids the potential penalty of partial register access or complex masking
	VMOVQ   (SI), X1
	VPCMPEQW X0, X1, X1
	VPMOVMSKB X1, CX
	// Mask out high 8 bits (corresponding to the zero-extended part)
	// VMOVQ zero extends, so the high bytes of X1 are 0.
	// X0 has broadcasted fp. 
	// If fp != 0, then 0 == fp is false, so high bits of result are 0.
	// If fp == 0, then 0 == fp is true, so high bits of result are 1.
	// We need to mask them out regardless to be safe, especially if searching for 0.
	ANDL    $0xFF, CX
	TESTL   CX, CX
	JNZ     found
	JMP     not_found

generic_loop:
	// Process 16 items at a time
	CMPQ    BX, $16
	JL      scalar_loop
	
	VPCMPEQW (SI), Y0, Y1
	VPMOVMSKB Y1, CX
	TESTL   CX, CX
	JNZ     found
	
	ADDQ    $32, SI    // Advance ptr by 32 bytes (16 items)
	SUBQ    $16, BX    // Decrement count
	JMP     generic_loop

scalar_loop:
	TESTQ   BX, BX
	JZ      not_found
	CMPW    (SI), AX
	JE      found
	ADDQ    $2, SI
	DECQ    BX
	JMP     scalar_loop

found:
	VZEROUPPER
	MOVB    $1, ret+32(FP)
	RET

not_found:
	VZEROUPPER
	MOVB    $0, ret+32(FP)
	RET
