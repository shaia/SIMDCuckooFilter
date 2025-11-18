//go:build arm64
// +build arm64

#include "textflag.h"

// func containsNEON(data []byte, fp byte) bool
// Optimized ARM64 implementation using unrolled scalar operations
TEXT ·containsNEON(SB), NOSPLIT, $0-33
	MOVD    data_base+0(FP), R0    // R0 = data pointer
	MOVD    data_len+8(FP), R1     // R1 = length
	MOVBU   fp+24(FP), R2          // R2 = fingerprint

	CMP     $64, R1
	BEQ     len64
	CMP     $32, R1
	BEQ     len32
	CMP     $16, R1
	BEQ     len16
	CMP     $8, R1
	BEQ     len8
	CMP     $4, R1
	BEQ     len4
	B       scalar

len64:
	// Manually unroll loop for 64 bytes
	MOVBU   0(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   1(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   2(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   3(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   4(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   5(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   6(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   7(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   8(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   9(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   10(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   11(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   12(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   13(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   14(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   15(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   16(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   17(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   18(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   19(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   20(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   21(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   22(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   23(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   24(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   25(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   26(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   27(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   28(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   29(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   30(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   31(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   32(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   33(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   34(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   35(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   36(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   37(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   38(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   39(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   40(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   41(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   42(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   43(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   44(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   45(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   46(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   47(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   48(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   49(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   50(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   51(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   52(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   53(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   54(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   55(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   56(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   57(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   58(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   59(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   60(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   61(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   62(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   63(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVD    $0, R0
	MOVB    R0, ret+32(FP)
	RET

len32:
	// Manually unroll loop for 32 bytes
	MOVBU   0(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   1(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   2(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   3(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   4(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   5(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   6(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   7(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   8(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   9(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   10(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   11(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   12(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   13(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   14(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   15(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   16(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   17(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   18(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   19(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   20(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   21(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   22(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   23(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   24(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   25(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   26(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   27(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   28(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   29(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   30(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   31(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVD    $0, R0
	MOVB    R0, ret+32(FP)
	RET

len16:
	// Manually unroll loop for 16 bytes
	MOVBU   0(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   1(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   2(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   3(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   4(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   5(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   6(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   7(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   8(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   9(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   10(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   11(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   12(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   13(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   14(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   15(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVD    $0, R0
	MOVB    R0, ret+32(FP)
	RET

len8:
	// Manually unroll loop for 8 bytes
	MOVBU   0(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   1(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   2(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   3(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   4(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   5(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   6(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   7(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVD    $0, R0
	MOVB    R0, ret+32(FP)
	RET

len4:
	MOVBU   0(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   1(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   2(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVBU   3(R0), R3
	CMP     R2, R3
	BEQ     found
	MOVD    $0, R0
	MOVB    R0, ret+32(FP)
	RET

found:
	MOVD    $1, R0
	MOVB    R0, ret+32(FP)
	RET

scalar:
	CBZ     R1, scalar_notfound

scalar_loop:
	MOVBU   (R0), R3
	CMP     R2, R3
	BEQ     found
	ADD     $1, R0
	SUB     $1, R1
	CBNZ    R1, scalar_loop

scalar_notfound:
	MOVD    $0, R0
	MOVB    R0, ret+32(FP)
	RET

// func isFullNEON(data []byte) bool
TEXT ·isFullNEON(SB), NOSPLIT, $0-25
	MOVD    data_base+0(FP), R0
	MOVD    data_len+8(FP), R1

	CMP     $64, R1
	BEQ     full_len64
	CMP     $32, R1
	BEQ     full_len32
	CMP     $16, R1
	BEQ     full_len16
	CMP     $8, R1
	BEQ     full_len8
	CMP     $4, R1
	BEQ     full_len4
	B       full_scalar

full_len64:
	MOVBU   0(R0), R2
	CBZ     R2, not_full
	MOVBU   1(R0), R2
	CBZ     R2, not_full
	MOVBU   2(R0), R2
	CBZ     R2, not_full
	MOVBU   3(R0), R2
	CBZ     R2, not_full
	MOVBU   4(R0), R2
	CBZ     R2, not_full
	MOVBU   5(R0), R2
	CBZ     R2, not_full
	MOVBU   6(R0), R2
	CBZ     R2, not_full
	MOVBU   7(R0), R2
	CBZ     R2, not_full
	MOVBU   8(R0), R2
	CBZ     R2, not_full
	MOVBU   9(R0), R2
	CBZ     R2, not_full
	MOVBU   10(R0), R2
	CBZ     R2, not_full
	MOVBU   11(R0), R2
	CBZ     R2, not_full
	MOVBU   12(R0), R2
	CBZ     R2, not_full
	MOVBU   13(R0), R2
	CBZ     R2, not_full
	MOVBU   14(R0), R2
	CBZ     R2, not_full
	MOVBU   15(R0), R2
	CBZ     R2, not_full
	MOVBU   16(R0), R2
	CBZ     R2, not_full
	MOVBU   17(R0), R2
	CBZ     R2, not_full
	MOVBU   18(R0), R2
	CBZ     R2, not_full
	MOVBU   19(R0), R2
	CBZ     R2, not_full
	MOVBU   20(R0), R2
	CBZ     R2, not_full
	MOVBU   21(R0), R2
	CBZ     R2, not_full
	MOVBU   22(R0), R2
	CBZ     R2, not_full
	MOVBU   23(R0), R2
	CBZ     R2, not_full
	MOVBU   24(R0), R2
	CBZ     R2, not_full
	MOVBU   25(R0), R2
	CBZ     R2, not_full
	MOVBU   26(R0), R2
	CBZ     R2, not_full
	MOVBU   27(R0), R2
	CBZ     R2, not_full
	MOVBU   28(R0), R2
	CBZ     R2, not_full
	MOVBU   29(R0), R2
	CBZ     R2, not_full
	MOVBU   30(R0), R2
	CBZ     R2, not_full
	MOVBU   31(R0), R2
	CBZ     R2, not_full
	MOVBU   32(R0), R2
	CBZ     R2, not_full
	MOVBU   33(R0), R2
	CBZ     R2, not_full
	MOVBU   34(R0), R2
	CBZ     R2, not_full
	MOVBU   35(R0), R2
	CBZ     R2, not_full
	MOVBU   36(R0), R2
	CBZ     R2, not_full
	MOVBU   37(R0), R2
	CBZ     R2, not_full
	MOVBU   38(R0), R2
	CBZ     R2, not_full
	MOVBU   39(R0), R2
	CBZ     R2, not_full
	MOVBU   40(R0), R2
	CBZ     R2, not_full
	MOVBU   41(R0), R2
	CBZ     R2, not_full
	MOVBU   42(R0), R2
	CBZ     R2, not_full
	MOVBU   43(R0), R2
	CBZ     R2, not_full
	MOVBU   44(R0), R2
	CBZ     R2, not_full
	MOVBU   45(R0), R2
	CBZ     R2, not_full
	MOVBU   46(R0), R2
	CBZ     R2, not_full
	MOVBU   47(R0), R2
	CBZ     R2, not_full
	MOVBU   48(R0), R2
	CBZ     R2, not_full
	MOVBU   49(R0), R2
	CBZ     R2, not_full
	MOVBU   50(R0), R2
	CBZ     R2, not_full
	MOVBU   51(R0), R2
	CBZ     R2, not_full
	MOVBU   52(R0), R2
	CBZ     R2, not_full
	MOVBU   53(R0), R2
	CBZ     R2, not_full
	MOVBU   54(R0), R2
	CBZ     R2, not_full
	MOVBU   55(R0), R2
	CBZ     R2, not_full
	MOVBU   56(R0), R2
	CBZ     R2, not_full
	MOVBU   57(R0), R2
	CBZ     R2, not_full
	MOVBU   58(R0), R2
	CBZ     R2, not_full
	MOVBU   59(R0), R2
	CBZ     R2, not_full
	MOVBU   60(R0), R2
	CBZ     R2, not_full
	MOVBU   61(R0), R2
	CBZ     R2, not_full
	MOVBU   62(R0), R2
	CBZ     R2, not_full
	MOVBU   63(R0), R2
	CBZ     R2, not_full
	MOVD    $1, R0
	MOVB    R0, ret+24(FP)
	RET

full_len32:
	MOVBU   0(R0), R2
	CBZ     R2, not_full
	MOVBU   1(R0), R2
	CBZ     R2, not_full
	MOVBU   2(R0), R2
	CBZ     R2, not_full
	MOVBU   3(R0), R2
	CBZ     R2, not_full
	MOVBU   4(R0), R2
	CBZ     R2, not_full
	MOVBU   5(R0), R2
	CBZ     R2, not_full
	MOVBU   6(R0), R2
	CBZ     R2, not_full
	MOVBU   7(R0), R2
	CBZ     R2, not_full
	MOVBU   8(R0), R2
	CBZ     R2, not_full
	MOVBU   9(R0), R2
	CBZ     R2, not_full
	MOVBU   10(R0), R2
	CBZ     R2, not_full
	MOVBU   11(R0), R2
	CBZ     R2, not_full
	MOVBU   12(R0), R2
	CBZ     R2, not_full
	MOVBU   13(R0), R2
	CBZ     R2, not_full
	MOVBU   14(R0), R2
	CBZ     R2, not_full
	MOVBU   15(R0), R2
	CBZ     R2, not_full
	MOVBU   16(R0), R2
	CBZ     R2, not_full
	MOVBU   17(R0), R2
	CBZ     R2, not_full
	MOVBU   18(R0), R2
	CBZ     R2, not_full
	MOVBU   19(R0), R2
	CBZ     R2, not_full
	MOVBU   20(R0), R2
	CBZ     R2, not_full
	MOVBU   21(R0), R2
	CBZ     R2, not_full
	MOVBU   22(R0), R2
	CBZ     R2, not_full
	MOVBU   23(R0), R2
	CBZ     R2, not_full
	MOVBU   24(R0), R2
	CBZ     R2, not_full
	MOVBU   25(R0), R2
	CBZ     R2, not_full
	MOVBU   26(R0), R2
	CBZ     R2, not_full
	MOVBU   27(R0), R2
	CBZ     R2, not_full
	MOVBU   28(R0), R2
	CBZ     R2, not_full
	MOVBU   29(R0), R2
	CBZ     R2, not_full
	MOVBU   30(R0), R2
	CBZ     R2, not_full
	MOVBU   31(R0), R2
	CBZ     R2, not_full
	MOVD    $1, R0
	MOVB    R0, ret+24(FP)
	RET

full_len16:
	MOVBU   0(R0), R2
	CBZ     R2, not_full
	MOVBU   1(R0), R2
	CBZ     R2, not_full
	MOVBU   2(R0), R2
	CBZ     R2, not_full
	MOVBU   3(R0), R2
	CBZ     R2, not_full
	MOVBU   4(R0), R2
	CBZ     R2, not_full
	MOVBU   5(R0), R2
	CBZ     R2, not_full
	MOVBU   6(R0), R2
	CBZ     R2, not_full
	MOVBU   7(R0), R2
	CBZ     R2, not_full
	MOVBU   8(R0), R2
	CBZ     R2, not_full
	MOVBU   9(R0), R2
	CBZ     R2, not_full
	MOVBU   10(R0), R2
	CBZ     R2, not_full
	MOVBU   11(R0), R2
	CBZ     R2, not_full
	MOVBU   12(R0), R2
	CBZ     R2, not_full
	MOVBU   13(R0), R2
	CBZ     R2, not_full
	MOVBU   14(R0), R2
	CBZ     R2, not_full
	MOVBU   15(R0), R2
	CBZ     R2, not_full
	MOVD    $1, R0
	MOVB    R0, ret+24(FP)
	RET

full_len8:
	MOVBU   0(R0), R2
	CBZ     R2, not_full
	MOVBU   1(R0), R2
	CBZ     R2, not_full
	MOVBU   2(R0), R2
	CBZ     R2, not_full
	MOVBU   3(R0), R2
	CBZ     R2, not_full
	MOVBU   4(R0), R2
	CBZ     R2, not_full
	MOVBU   5(R0), R2
	CBZ     R2, not_full
	MOVBU   6(R0), R2
	CBZ     R2, not_full
	MOVBU   7(R0), R2
	CBZ     R2, not_full
	MOVD    $1, R0
	MOVB    R0, ret+24(FP)
	RET

full_len4:
	MOVBU   0(R0), R2
	CBZ     R2, not_full
	MOVBU   1(R0), R2
	CBZ     R2, not_full
	MOVBU   2(R0), R2
	CBZ     R2, not_full
	MOVBU   3(R0), R2
	CBZ     R2, not_full
	MOVD    $1, R0
	MOVB    R0, ret+24(FP)
	RET

not_full:
	MOVD    $0, R0
	MOVB    R0, ret+24(FP)
	RET

full_scalar:
	CBZ     R1, full_scalar_empty

full_scalar_loop:
	MOVBU   (R0), R2
	CBZ     R2, not_full
	ADD     $1, R0
	SUB     $1, R1
	CBNZ    R1, full_scalar_loop

full_scalar_empty:
	MOVD    $1, R0
	MOVB    R0, ret+24(FP)
	RET

// func countNEON(data []byte) uint
TEXT ·countNEON(SB), NOSPLIT, $0-32
	MOVD    data_base+0(FP), R0
	MOVD    data_len+8(FP), R1
	MOVD    $0, R3                 // count = 0

	CMP     $64, R1
	BEQ     count_len64
	CMP     $32, R1
	BEQ     count_len32
	CMP     $16, R1
	BEQ     count_len16
	CMP     $8, R1
	BEQ     count_len8
	CMP     $4, R1
	BEQ     count_len4
	B       count_scalar

count_len64:
	MOVBU   0(R0), R2
	CBZ     R2, count_64_1
	ADD     $1, R3
count_64_1:
	MOVBU   1(R0), R2
	CBZ     R2, count_64_2
	ADD     $1, R3
count_64_2:
	MOVBU   2(R0), R2
	CBZ     R2, count_64_3
	ADD     $1, R3
count_64_3:
	MOVBU   3(R0), R2
	CBZ     R2, count_64_4
	ADD     $1, R3
count_64_4:
	MOVBU   4(R0), R2
	CBZ     R2, count_64_5
	ADD     $1, R3
count_64_5:
	MOVBU   5(R0), R2
	CBZ     R2, count_64_6
	ADD     $1, R3
count_64_6:
	MOVBU   6(R0), R2
	CBZ     R2, count_64_7
	ADD     $1, R3
count_64_7:
	MOVBU   7(R0), R2
	CBZ     R2, count_64_8
	ADD     $1, R3
count_64_8:
	MOVBU   8(R0), R2
	CBZ     R2, count_64_9
	ADD     $1, R3
count_64_9:
	MOVBU   9(R0), R2
	CBZ     R2, count_64_10
	ADD     $1, R3
count_64_10:
	MOVBU   10(R0), R2
	CBZ     R2, count_64_11
	ADD     $1, R3
count_64_11:
	MOVBU   11(R0), R2
	CBZ     R2, count_64_12
	ADD     $1, R3
count_64_12:
	MOVBU   12(R0), R2
	CBZ     R2, count_64_13
	ADD     $1, R3
count_64_13:
	MOVBU   13(R0), R2
	CBZ     R2, count_64_14
	ADD     $1, R3
count_64_14:
	MOVBU   14(R0), R2
	CBZ     R2, count_64_15
	ADD     $1, R3
count_64_15:
	MOVBU   15(R0), R2
	CBZ     R2, count_64_16
	ADD     $1, R3
count_64_16:
	MOVBU   16(R0), R2
	CBZ     R2, count_64_17
	ADD     $1, R3
count_64_17:
	MOVBU   17(R0), R2
	CBZ     R2, count_64_18
	ADD     $1, R3
count_64_18:
	MOVBU   18(R0), R2
	CBZ     R2, count_64_19
	ADD     $1, R3
count_64_19:
	MOVBU   19(R0), R2
	CBZ     R2, count_64_20
	ADD     $1, R3
count_64_20:
	MOVBU   20(R0), R2
	CBZ     R2, count_64_21
	ADD     $1, R3
count_64_21:
	MOVBU   21(R0), R2
	CBZ     R2, count_64_22
	ADD     $1, R3
count_64_22:
	MOVBU   22(R0), R2
	CBZ     R2, count_64_23
	ADD     $1, R3
count_64_23:
	MOVBU   23(R0), R2
	CBZ     R2, count_64_24
	ADD     $1, R3
count_64_24:
	MOVBU   24(R0), R2
	CBZ     R2, count_64_25
	ADD     $1, R3
count_64_25:
	MOVBU   25(R0), R2
	CBZ     R2, count_64_26
	ADD     $1, R3
count_64_26:
	MOVBU   26(R0), R2
	CBZ     R2, count_64_27
	ADD     $1, R3
count_64_27:
	MOVBU   27(R0), R2
	CBZ     R2, count_64_28
	ADD     $1, R3
count_64_28:
	MOVBU   28(R0), R2
	CBZ     R2, count_64_29
	ADD     $1, R3
count_64_29:
	MOVBU   29(R0), R2
	CBZ     R2, count_64_30
	ADD     $1, R3
count_64_30:
	MOVBU   30(R0), R2
	CBZ     R2, count_64_31
	ADD     $1, R3
count_64_31:
	MOVBU   31(R0), R2
	CBZ     R2, count_64_32
	ADD     $1, R3
count_64_32:
	MOVBU   32(R0), R2
	CBZ     R2, count_64_33
	ADD     $1, R3
count_64_33:
	MOVBU   33(R0), R2
	CBZ     R2, count_64_34
	ADD     $1, R3
count_64_34:
	MOVBU   34(R0), R2
	CBZ     R2, count_64_35
	ADD     $1, R3
count_64_35:
	MOVBU   35(R0), R2
	CBZ     R2, count_64_36
	ADD     $1, R3
count_64_36:
	MOVBU   36(R0), R2
	CBZ     R2, count_64_37
	ADD     $1, R3
count_64_37:
	MOVBU   37(R0), R2
	CBZ     R2, count_64_38
	ADD     $1, R3
count_64_38:
	MOVBU   38(R0), R2
	CBZ     R2, count_64_39
	ADD     $1, R3
count_64_39:
	MOVBU   39(R0), R2
	CBZ     R2, count_64_40
	ADD     $1, R3
count_64_40:
	MOVBU   40(R0), R2
	CBZ     R2, count_64_41
	ADD     $1, R3
count_64_41:
	MOVBU   41(R0), R2
	CBZ     R2, count_64_42
	ADD     $1, R3
count_64_42:
	MOVBU   42(R0), R2
	CBZ     R2, count_64_43
	ADD     $1, R3
count_64_43:
	MOVBU   43(R0), R2
	CBZ     R2, count_64_44
	ADD     $1, R3
count_64_44:
	MOVBU   44(R0), R2
	CBZ     R2, count_64_45
	ADD     $1, R3
count_64_45:
	MOVBU   45(R0), R2
	CBZ     R2, count_64_46
	ADD     $1, R3
count_64_46:
	MOVBU   46(R0), R2
	CBZ     R2, count_64_47
	ADD     $1, R3
count_64_47:
	MOVBU   47(R0), R2
	CBZ     R2, count_64_48
	ADD     $1, R3
count_64_48:
	MOVBU   48(R0), R2
	CBZ     R2, count_64_49
	ADD     $1, R3
count_64_49:
	MOVBU   49(R0), R2
	CBZ     R2, count_64_50
	ADD     $1, R3
count_64_50:
	MOVBU   50(R0), R2
	CBZ     R2, count_64_51
	ADD     $1, R3
count_64_51:
	MOVBU   51(R0), R2
	CBZ     R2, count_64_52
	ADD     $1, R3
count_64_52:
	MOVBU   52(R0), R2
	CBZ     R2, count_64_53
	ADD     $1, R3
count_64_53:
	MOVBU   53(R0), R2
	CBZ     R2, count_64_54
	ADD     $1, R3
count_64_54:
	MOVBU   54(R0), R2
	CBZ     R2, count_64_55
	ADD     $1, R3
count_64_55:
	MOVBU   55(R0), R2
	CBZ     R2, count_64_56
	ADD     $1, R3
count_64_56:
	MOVBU   56(R0), R2
	CBZ     R2, count_64_57
	ADD     $1, R3
count_64_57:
	MOVBU   57(R0), R2
	CBZ     R2, count_64_58
	ADD     $1, R3
count_64_58:
	MOVBU   58(R0), R2
	CBZ     R2, count_64_59
	ADD     $1, R3
count_64_59:
	MOVBU   59(R0), R2
	CBZ     R2, count_64_60
	ADD     $1, R3
count_64_60:
	MOVBU   60(R0), R2
	CBZ     R2, count_64_61
	ADD     $1, R3
count_64_61:
	MOVBU   61(R0), R2
	CBZ     R2, count_64_62
	ADD     $1, R3
count_64_62:
	MOVBU   62(R0), R2
	CBZ     R2, count_64_63
	ADD     $1, R3
count_64_63:
	MOVBU   63(R0), R2
	CBZ     R2, count_64_done
	ADD     $1, R3
count_64_done:
	MOVD    R3, ret+24(FP)
	RET

count_len32:
	MOVBU   0(R0), R2
	CBZ     R2, count_32_1
	ADD     $1, R3
count_32_1:
	MOVBU   1(R0), R2
	CBZ     R2, count_32_2
	ADD     $1, R3
count_32_2:
	MOVBU   2(R0), R2
	CBZ     R2, count_32_3
	ADD     $1, R3
count_32_3:
	MOVBU   3(R0), R2
	CBZ     R2, count_32_4
	ADD     $1, R3
count_32_4:
	MOVBU   4(R0), R2
	CBZ     R2, count_32_5
	ADD     $1, R3
count_32_5:
	MOVBU   5(R0), R2
	CBZ     R2, count_32_6
	ADD     $1, R3
count_32_6:
	MOVBU   6(R0), R2
	CBZ     R2, count_32_7
	ADD     $1, R3
count_32_7:
	MOVBU   7(R0), R2
	CBZ     R2, count_32_8
	ADD     $1, R3
count_32_8:
	MOVBU   8(R0), R2
	CBZ     R2, count_32_9
	ADD     $1, R3
count_32_9:
	MOVBU   9(R0), R2
	CBZ     R2, count_32_10
	ADD     $1, R3
count_32_10:
	MOVBU   10(R0), R2
	CBZ     R2, count_32_11
	ADD     $1, R3
count_32_11:
	MOVBU   11(R0), R2
	CBZ     R2, count_32_12
	ADD     $1, R3
count_32_12:
	MOVBU   12(R0), R2
	CBZ     R2, count_32_13
	ADD     $1, R3
count_32_13:
	MOVBU   13(R0), R2
	CBZ     R2, count_32_14
	ADD     $1, R3
count_32_14:
	MOVBU   14(R0), R2
	CBZ     R2, count_32_15
	ADD     $1, R3
count_32_15:
	MOVBU   15(R0), R2
	CBZ     R2, count_32_16
	ADD     $1, R3
count_32_16:
	MOVBU   16(R0), R2
	CBZ     R2, count_32_17
	ADD     $1, R3
count_32_17:
	MOVBU   17(R0), R2
	CBZ     R2, count_32_18
	ADD     $1, R3
count_32_18:
	MOVBU   18(R0), R2
	CBZ     R2, count_32_19
	ADD     $1, R3
count_32_19:
	MOVBU   19(R0), R2
	CBZ     R2, count_32_20
	ADD     $1, R3
count_32_20:
	MOVBU   20(R0), R2
	CBZ     R2, count_32_21
	ADD     $1, R3
count_32_21:
	MOVBU   21(R0), R2
	CBZ     R2, count_32_22
	ADD     $1, R3
count_32_22:
	MOVBU   22(R0), R2
	CBZ     R2, count_32_23
	ADD     $1, R3
count_32_23:
	MOVBU   23(R0), R2
	CBZ     R2, count_32_24
	ADD     $1, R3
count_32_24:
	MOVBU   24(R0), R2
	CBZ     R2, count_32_25
	ADD     $1, R3
count_32_25:
	MOVBU   25(R0), R2
	CBZ     R2, count_32_26
	ADD     $1, R3
count_32_26:
	MOVBU   26(R0), R2
	CBZ     R2, count_32_27
	ADD     $1, R3
count_32_27:
	MOVBU   27(R0), R2
	CBZ     R2, count_32_28
	ADD     $1, R3
count_32_28:
	MOVBU   28(R0), R2
	CBZ     R2, count_32_29
	ADD     $1, R3
count_32_29:
	MOVBU   29(R0), R2
	CBZ     R2, count_32_30
	ADD     $1, R3
count_32_30:
	MOVBU   30(R0), R2
	CBZ     R2, count_32_31
	ADD     $1, R3
count_32_31:
	MOVBU   31(R0), R2
	CBZ     R2, count_32_done
	ADD     $1, R3
count_32_done:
	MOVD    R3, ret+24(FP)
	RET

count_len16:
	MOVBU   0(R0), R2
	CBZ     R2, count_16_1
	ADD     $1, R3
count_16_1:
	MOVBU   1(R0), R2
	CBZ     R2, count_16_2
	ADD     $1, R3
count_16_2:
	MOVBU   2(R0), R2
	CBZ     R2, count_16_3
	ADD     $1, R3
count_16_3:
	MOVBU   3(R0), R2
	CBZ     R2, count_16_4
	ADD     $1, R3
count_16_4:
	MOVBU   4(R0), R2
	CBZ     R2, count_16_5
	ADD     $1, R3
count_16_5:
	MOVBU   5(R0), R2
	CBZ     R2, count_16_6
	ADD     $1, R3
count_16_6:
	MOVBU   6(R0), R2
	CBZ     R2, count_16_7
	ADD     $1, R3
count_16_7:
	MOVBU   7(R0), R2
	CBZ     R2, count_16_8
	ADD     $1, R3
count_16_8:
	MOVBU   8(R0), R2
	CBZ     R2, count_16_9
	ADD     $1, R3
count_16_9:
	MOVBU   9(R0), R2
	CBZ     R2, count_16_10
	ADD     $1, R3
count_16_10:
	MOVBU   10(R0), R2
	CBZ     R2, count_16_11
	ADD     $1, R3
count_16_11:
	MOVBU   11(R0), R2
	CBZ     R2, count_16_12
	ADD     $1, R3
count_16_12:
	MOVBU   12(R0), R2
	CBZ     R2, count_16_13
	ADD     $1, R3
count_16_13:
	MOVBU   13(R0), R2
	CBZ     R2, count_16_14
	ADD     $1, R3
count_16_14:
	MOVBU   14(R0), R2
	CBZ     R2, count_16_15
	ADD     $1, R3
count_16_15:
	MOVBU   15(R0), R2
	CBZ     R2, count_16_done
	ADD     $1, R3
count_16_done:
	MOVD    R3, ret+24(FP)
	RET

count_len8:
	MOVBU   0(R0), R2
	CBZ     R2, count_8_1
	ADD     $1, R3
count_8_1:
	MOVBU   1(R0), R2
	CBZ     R2, count_8_2
	ADD     $1, R3
count_8_2:
	MOVBU   2(R0), R2
	CBZ     R2, count_8_3
	ADD     $1, R3
count_8_3:
	MOVBU   3(R0), R2
	CBZ     R2, count_8_4
	ADD     $1, R3
count_8_4:
	MOVBU   4(R0), R2
	CBZ     R2, count_8_5
	ADD     $1, R3
count_8_5:
	MOVBU   5(R0), R2
	CBZ     R2, count_8_6
	ADD     $1, R3
count_8_6:
	MOVBU   6(R0), R2
	CBZ     R2, count_8_7
	ADD     $1, R3
count_8_7:
	MOVBU   7(R0), R2
	CBZ     R2, count_8_done
	ADD     $1, R3
count_8_done:
	MOVD    R3, ret+24(FP)
	RET

count_len4:
	MOVBU   0(R0), R2
	CBZ     R2, count_4_1
	ADD     $1, R3
count_4_1:
	MOVBU   1(R0), R2
	CBZ     R2, count_4_2
	ADD     $1, R3
count_4_2:
	MOVBU   2(R0), R2
	CBZ     R2, count_4_3
	ADD     $1, R3
count_4_3:
	MOVBU   3(R0), R2
	CBZ     R2, count_4_done
	ADD     $1, R3
count_4_done:
	MOVD    R3, ret+24(FP)
	RET

count_scalar:
	CBZ     R1, count_scalar_done

count_scalar_loop:
	MOVBU   (R0), R2
	CBZ     R2, count_scalar_skip
	ADD     $1, R3

count_scalar_skip:
	ADD     $1, R0
	SUB     $1, R1
	CBNZ    R1, count_scalar_loop

count_scalar_done:
	MOVD    R3, ret+24(FP)
	RET

// func findFirstZeroNEON(data []byte) uint
TEXT ·findFirstZeroNEON(SB), NOSPLIT, $0-32
	MOVD    data_base+0(FP), R0
	MOVD    data_len+8(FP), R1

	CMP     $64, R1
	BEQ     find_len64
	CMP     $32, R1
	BEQ     find_len32
	CMP     $16, R1
	BEQ     find_len16
	CMP     $8, R1
	BEQ     find_len8
	CMP     $4, R1
	BEQ     find_len4
	B       find_scalar

find_len64:
	MOVBU   0(R0), R2
	CBZ     R2, find_64_0
	MOVBU   1(R0), R2
	CBZ     R2, find_64_1
	MOVBU   2(R0), R2
	CBZ     R2, find_64_2
	MOVBU   3(R0), R2
	CBZ     R2, find_64_3
	MOVBU   4(R0), R2
	CBZ     R2, find_64_4
	MOVBU   5(R0), R2
	CBZ     R2, find_64_5
	MOVBU   6(R0), R2
	CBZ     R2, find_64_6
	MOVBU   7(R0), R2
	CBZ     R2, find_64_7
	MOVBU   8(R0), R2
	CBZ     R2, find_64_8
	MOVBU   9(R0), R2
	CBZ     R2, find_64_9
	MOVBU   10(R0), R2
	CBZ     R2, find_64_10
	MOVBU   11(R0), R2
	CBZ     R2, find_64_11
	MOVBU   12(R0), R2
	CBZ     R2, find_64_12
	MOVBU   13(R0), R2
	CBZ     R2, find_64_13
	MOVBU   14(R0), R2
	CBZ     R2, find_64_14
	MOVBU   15(R0), R2
	CBZ     R2, find_64_15
	MOVBU   16(R0), R2
	CBZ     R2, find_64_16
	MOVBU   17(R0), R2
	CBZ     R2, find_64_17
	MOVBU   18(R0), R2
	CBZ     R2, find_64_18
	MOVBU   19(R0), R2
	CBZ     R2, find_64_19
	MOVBU   20(R0), R2
	CBZ     R2, find_64_20
	MOVBU   21(R0), R2
	CBZ     R2, find_64_21
	MOVBU   22(R0), R2
	CBZ     R2, find_64_22
	MOVBU   23(R0), R2
	CBZ     R2, find_64_23
	MOVBU   24(R0), R2
	CBZ     R2, find_64_24
	MOVBU   25(R0), R2
	CBZ     R2, find_64_25
	MOVBU   26(R0), R2
	CBZ     R2, find_64_26
	MOVBU   27(R0), R2
	CBZ     R2, find_64_27
	MOVBU   28(R0), R2
	CBZ     R2, find_64_28
	MOVBU   29(R0), R2
	CBZ     R2, find_64_29
	MOVBU   30(R0), R2
	CBZ     R2, find_64_30
	MOVBU   31(R0), R2
	CBZ     R2, find_64_31
	MOVBU   32(R0), R2
	CBZ     R2, find_64_32
	MOVBU   33(R0), R2
	CBZ     R2, find_64_33
	MOVBU   34(R0), R2
	CBZ     R2, find_64_34
	MOVBU   35(R0), R2
	CBZ     R2, find_64_35
	MOVBU   36(R0), R2
	CBZ     R2, find_64_36
	MOVBU   37(R0), R2
	CBZ     R2, find_64_37
	MOVBU   38(R0), R2
	CBZ     R2, find_64_38
	MOVBU   39(R0), R2
	CBZ     R2, find_64_39
	MOVBU   40(R0), R2
	CBZ     R2, find_64_40
	MOVBU   41(R0), R2
	CBZ     R2, find_64_41
	MOVBU   42(R0), R2
	CBZ     R2, find_64_42
	MOVBU   43(R0), R2
	CBZ     R2, find_64_43
	MOVBU   44(R0), R2
	CBZ     R2, find_64_44
	MOVBU   45(R0), R2
	CBZ     R2, find_64_45
	MOVBU   46(R0), R2
	CBZ     R2, find_64_46
	MOVBU   47(R0), R2
	CBZ     R2, find_64_47
	MOVBU   48(R0), R2
	CBZ     R2, find_64_48
	MOVBU   49(R0), R2
	CBZ     R2, find_64_49
	MOVBU   50(R0), R2
	CBZ     R2, find_64_50
	MOVBU   51(R0), R2
	CBZ     R2, find_64_51
	MOVBU   52(R0), R2
	CBZ     R2, find_64_52
	MOVBU   53(R0), R2
	CBZ     R2, find_64_53
	MOVBU   54(R0), R2
	CBZ     R2, find_64_54
	MOVBU   55(R0), R2
	CBZ     R2, find_64_55
	MOVBU   56(R0), R2
	CBZ     R2, find_64_56
	MOVBU   57(R0), R2
	CBZ     R2, find_64_57
	MOVBU   58(R0), R2
	CBZ     R2, find_64_58
	MOVBU   59(R0), R2
	CBZ     R2, find_64_59
	MOVBU   60(R0), R2
	CBZ     R2, find_64_60
	MOVBU   61(R0), R2
	CBZ     R2, find_64_61
	MOVBU   62(R0), R2
	CBZ     R2, find_64_62
	MOVBU   63(R0), R2
	CBZ     R2, find_64_63
	MOVD    $64, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_0:
	MOVD    $0, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_1:
	MOVD    $1, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_2:
	MOVD    $2, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_3:
	MOVD    $3, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_4:
	MOVD    $4, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_5:
	MOVD    $5, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_6:
	MOVD    $6, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_7:
	MOVD    $7, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_8:
	MOVD    $8, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_9:
	MOVD    $9, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_10:
	MOVD    $10, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_11:
	MOVD    $11, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_12:
	MOVD    $12, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_13:
	MOVD    $13, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_14:
	MOVD    $14, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_15:
	MOVD    $15, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_16:
	MOVD    $16, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_17:
	MOVD    $17, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_18:
	MOVD    $18, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_19:
	MOVD    $19, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_20:
	MOVD    $20, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_21:
	MOVD    $21, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_22:
	MOVD    $22, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_23:
	MOVD    $23, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_24:
	MOVD    $24, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_25:
	MOVD    $25, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_26:
	MOVD    $26, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_27:
	MOVD    $27, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_28:
	MOVD    $28, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_29:
	MOVD    $29, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_30:
	MOVD    $30, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_31:
	MOVD    $31, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_32:
	MOVD    $32, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_33:
	MOVD    $33, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_34:
	MOVD    $34, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_35:
	MOVD    $35, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_36:
	MOVD    $36, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_37:
	MOVD    $37, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_38:
	MOVD    $38, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_39:
	MOVD    $39, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_40:
	MOVD    $40, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_41:
	MOVD    $41, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_42:
	MOVD    $42, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_43:
	MOVD    $43, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_44:
	MOVD    $44, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_45:
	MOVD    $45, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_46:
	MOVD    $46, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_47:
	MOVD    $47, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_48:
	MOVD    $48, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_49:
	MOVD    $49, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_50:
	MOVD    $50, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_51:
	MOVD    $51, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_52:
	MOVD    $52, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_53:
	MOVD    $53, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_54:
	MOVD    $54, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_55:
	MOVD    $55, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_56:
	MOVD    $56, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_57:
	MOVD    $57, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_58:
	MOVD    $58, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_59:
	MOVD    $59, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_60:
	MOVD    $60, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_61:
	MOVD    $61, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_62:
	MOVD    $62, R0
	MOVD    R0, ret+24(FP)
	RET
find_64_63:
	MOVD    $63, R0
	MOVD    R0, ret+24(FP)
	RET

find_len32:
	MOVBU   0(R0), R2
	CBZ     R2, find_32_0
	MOVBU   1(R0), R2
	CBZ     R2, find_32_1
	MOVBU   2(R0), R2
	CBZ     R2, find_32_2
	MOVBU   3(R0), R2
	CBZ     R2, find_32_3
	MOVBU   4(R0), R2
	CBZ     R2, find_32_4
	MOVBU   5(R0), R2
	CBZ     R2, find_32_5
	MOVBU   6(R0), R2
	CBZ     R2, find_32_6
	MOVBU   7(R0), R2
	CBZ     R2, find_32_7
	MOVBU   8(R0), R2
	CBZ     R2, find_32_8
	MOVBU   9(R0), R2
	CBZ     R2, find_32_9
	MOVBU   10(R0), R2
	CBZ     R2, find_32_10
	MOVBU   11(R0), R2
	CBZ     R2, find_32_11
	MOVBU   12(R0), R2
	CBZ     R2, find_32_12
	MOVBU   13(R0), R2
	CBZ     R2, find_32_13
	MOVBU   14(R0), R2
	CBZ     R2, find_32_14
	MOVBU   15(R0), R2
	CBZ     R2, find_32_15
	MOVBU   16(R0), R2
	CBZ     R2, find_32_16
	MOVBU   17(R0), R2
	CBZ     R2, find_32_17
	MOVBU   18(R0), R2
	CBZ     R2, find_32_18
	MOVBU   19(R0), R2
	CBZ     R2, find_32_19
	MOVBU   20(R0), R2
	CBZ     R2, find_32_20
	MOVBU   21(R0), R2
	CBZ     R2, find_32_21
	MOVBU   22(R0), R2
	CBZ     R2, find_32_22
	MOVBU   23(R0), R2
	CBZ     R2, find_32_23
	MOVBU   24(R0), R2
	CBZ     R2, find_32_24
	MOVBU   25(R0), R2
	CBZ     R2, find_32_25
	MOVBU   26(R0), R2
	CBZ     R2, find_32_26
	MOVBU   27(R0), R2
	CBZ     R2, find_32_27
	MOVBU   28(R0), R2
	CBZ     R2, find_32_28
	MOVBU   29(R0), R2
	CBZ     R2, find_32_29
	MOVBU   30(R0), R2
	CBZ     R2, find_32_30
	MOVBU   31(R0), R2
	CBZ     R2, find_32_31
	MOVD    R1, R0
	MOVD    R0, ret+24(FP)
	RET

find_32_0:
	MOVD    $0, R0
	MOVD    R0, ret+24(FP)
	RET
find_32_1:
	MOVD    $1, R0
	MOVD    R0, ret+24(FP)
	RET
find_32_2:
	MOVD    $2, R0
	MOVD    R0, ret+24(FP)
	RET
find_32_3:
	MOVD    $3, R0
	MOVD    R0, ret+24(FP)
	RET
find_32_4:
	MOVD    $4, R0
	MOVD    R0, ret+24(FP)
	RET
find_32_5:
	MOVD    $5, R0
	MOVD    R0, ret+24(FP)
	RET
find_32_6:
	MOVD    $6, R0
	MOVD    R0, ret+24(FP)
	RET
find_32_7:
	MOVD    $7, R0
	MOVD    R0, ret+24(FP)
	RET
find_32_8:
	MOVD    $8, R0
	MOVD    R0, ret+24(FP)
	RET
find_32_9:
	MOVD    $9, R0
	MOVD    R0, ret+24(FP)
	RET
find_32_10:
	MOVD    $10, R0
	MOVD    R0, ret+24(FP)
	RET
find_32_11:
	MOVD    $11, R0
	MOVD    R0, ret+24(FP)
	RET
find_32_12:
	MOVD    $12, R0
	MOVD    R0, ret+24(FP)
	RET
find_32_13:
	MOVD    $13, R0
	MOVD    R0, ret+24(FP)
	RET
find_32_14:
	MOVD    $14, R0
	MOVD    R0, ret+24(FP)
	RET
find_32_15:
	MOVD    $15, R0
	MOVD    R0, ret+24(FP)
	RET
find_32_16:
	MOVD    $16, R0
	MOVD    R0, ret+24(FP)
	RET
find_32_17:
	MOVD    $17, R0
	MOVD    R0, ret+24(FP)
	RET
find_32_18:
	MOVD    $18, R0
	MOVD    R0, ret+24(FP)
	RET
find_32_19:
	MOVD    $19, R0
	MOVD    R0, ret+24(FP)
	RET
find_32_20:
	MOVD    $20, R0
	MOVD    R0, ret+24(FP)
	RET
find_32_21:
	MOVD    $21, R0
	MOVD    R0, ret+24(FP)
	RET
find_32_22:
	MOVD    $22, R0
	MOVD    R0, ret+24(FP)
	RET
find_32_23:
	MOVD    $23, R0
	MOVD    R0, ret+24(FP)
	RET
find_32_24:
	MOVD    $24, R0
	MOVD    R0, ret+24(FP)
	RET
find_32_25:
	MOVD    $25, R0
	MOVD    R0, ret+24(FP)
	RET
find_32_26:
	MOVD    $26, R0
	MOVD    R0, ret+24(FP)
	RET
find_32_27:
	MOVD    $27, R0
	MOVD    R0, ret+24(FP)
	RET
find_32_28:
	MOVD    $28, R0
	MOVD    R0, ret+24(FP)
	RET
find_32_29:
	MOVD    $29, R0
	MOVD    R0, ret+24(FP)
	RET
find_32_30:
	MOVD    $30, R0
	MOVD    R0, ret+24(FP)
	RET
find_32_31:
	MOVD    $31, R0
	MOVD    R0, ret+24(FP)
	RET

find_len16:
	MOVBU   0(R0), R2
	CBZ     R2, find_16_0
	MOVBU   1(R0), R2
	CBZ     R2, find_16_1
	MOVBU   2(R0), R2
	CBZ     R2, find_16_2
	MOVBU   3(R0), R2
	CBZ     R2, find_16_3
	MOVBU   4(R0), R2
	CBZ     R2, find_16_4
	MOVBU   5(R0), R2
	CBZ     R2, find_16_5
	MOVBU   6(R0), R2
	CBZ     R2, find_16_6
	MOVBU   7(R0), R2
	CBZ     R2, find_16_7
	MOVBU   8(R0), R2
	CBZ     R2, find_16_8
	MOVBU   9(R0), R2
	CBZ     R2, find_16_9
	MOVBU   10(R0), R2
	CBZ     R2, find_16_10
	MOVBU   11(R0), R2
	CBZ     R2, find_16_11
	MOVBU   12(R0), R2
	CBZ     R2, find_16_12
	MOVBU   13(R0), R2
	CBZ     R2, find_16_13
	MOVBU   14(R0), R2
	CBZ     R2, find_16_14
	MOVBU   15(R0), R2
	CBZ     R2, find_16_15
	MOVD    R1, R0
	MOVD    R0, ret+24(FP)
	RET

find_16_0:
	MOVD    $0, R0
	MOVD    R0, ret+24(FP)
	RET
find_16_1:
	MOVD    $1, R0
	MOVD    R0, ret+24(FP)
	RET
find_16_2:
	MOVD    $2, R0
	MOVD    R0, ret+24(FP)
	RET
find_16_3:
	MOVD    $3, R0
	MOVD    R0, ret+24(FP)
	RET
find_16_4:
	MOVD    $4, R0
	MOVD    R0, ret+24(FP)
	RET
find_16_5:
	MOVD    $5, R0
	MOVD    R0, ret+24(FP)
	RET
find_16_6:
	MOVD    $6, R0
	MOVD    R0, ret+24(FP)
	RET
find_16_7:
	MOVD    $7, R0
	MOVD    R0, ret+24(FP)
	RET
find_16_8:
	MOVD    $8, R0
	MOVD    R0, ret+24(FP)
	RET
find_16_9:
	MOVD    $9, R0
	MOVD    R0, ret+24(FP)
	RET
find_16_10:
	MOVD    $10, R0
	MOVD    R0, ret+24(FP)
	RET
find_16_11:
	MOVD    $11, R0
	MOVD    R0, ret+24(FP)
	RET
find_16_12:
	MOVD    $12, R0
	MOVD    R0, ret+24(FP)
	RET
find_16_13:
	MOVD    $13, R0
	MOVD    R0, ret+24(FP)
	RET
find_16_14:
	MOVD    $14, R0
	MOVD    R0, ret+24(FP)
	RET
find_16_15:
	MOVD    $15, R0
	MOVD    R0, ret+24(FP)
	RET

find_len8:
	MOVBU   0(R0), R2
	CBZ     R2, find_8_0
	MOVBU   1(R0), R2
	CBZ     R2, find_8_1
	MOVBU   2(R0), R2
	CBZ     R2, find_8_2
	MOVBU   3(R0), R2
	CBZ     R2, find_8_3
	MOVBU   4(R0), R2
	CBZ     R2, find_8_4
	MOVBU   5(R0), R2
	CBZ     R2, find_8_5
	MOVBU   6(R0), R2
	CBZ     R2, find_8_6
	MOVBU   7(R0), R2
	CBZ     R2, find_8_7
	MOVD    R1, R0
	MOVD    R0, ret+24(FP)
	RET

find_8_0:
	MOVD    $0, R0
	MOVD    R0, ret+24(FP)
	RET
find_8_1:
	MOVD    $1, R0
	MOVD    R0, ret+24(FP)
	RET
find_8_2:
	MOVD    $2, R0
	MOVD    R0, ret+24(FP)
	RET
find_8_3:
	MOVD    $3, R0
	MOVD    R0, ret+24(FP)
	RET
find_8_4:
	MOVD    $4, R0
	MOVD    R0, ret+24(FP)
	RET
find_8_5:
	MOVD    $5, R0
	MOVD    R0, ret+24(FP)
	RET
find_8_6:
	MOVD    $6, R0
	MOVD    R0, ret+24(FP)
	RET
find_8_7:
	MOVD    $7, R0
	MOVD    R0, ret+24(FP)
	RET

find_len4:
	MOVBU   0(R0), R2
	CBZ     R2, find_4_0
	MOVBU   1(R0), R2
	CBZ     R2, find_4_1
	MOVBU   2(R0), R2
	CBZ     R2, find_4_2
	MOVBU   3(R0), R2
	CBZ     R2, find_4_3
	MOVD    R1, R0
	MOVD    R0, ret+24(FP)
	RET

find_4_0:
	MOVD    $0, R0
	MOVD    R0, ret+24(FP)
	RET
find_4_1:
	MOVD    $1, R0
	MOVD    R0, ret+24(FP)
	RET
find_4_2:
	MOVD    $2, R0
	MOVD    R0, ret+24(FP)
	RET
find_4_3:
	MOVD    $3, R0
	MOVD    R0, ret+24(FP)
	RET

find_scalar:
	MOVD    $0, R3
	CBZ     R1, find_scalar_done

find_scalar_loop:
	MOVBU   (R0), R2
	CBZ     R2, find_scalar_found
	ADD     $1, R3
	ADD     $1, R0
	SUB     $1, R1
	CBNZ    R1, find_scalar_loop

	MOVD    data_len+8(FP), R3

find_scalar_found:
find_scalar_done:
	MOVD    R3, ret+24(FP)
	RET
