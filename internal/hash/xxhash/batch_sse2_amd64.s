//go:build amd64
// +build amd64

#include "textflag.h"

// SSE2 constants - 2x 64-bit values packed
DATA sse2_prime64_1<>+0(SB)/8, $11400714785074694791
DATA sse2_prime64_1<>+8(SB)/8, $11400714785074694791
GLOBL sse2_prime64_1<>(SB), RODATA, $16

DATA sse2_prime64_2<>+0(SB)/8, $14029467366897019727
DATA sse2_prime64_2<>+8(SB)/8, $14029467366897019727
GLOBL sse2_prime64_2<>(SB), RODATA, $16

DATA sse2_prime64_3<>+0(SB)/8, $1609587929392839161
DATA sse2_prime64_3<>+8(SB)/8, $1609587929392839161
GLOBL sse2_prime64_3<>(SB), RODATA, $16

DATA sse2_prime64_4<>+0(SB)/8, $9650029242287828579
DATA sse2_prime64_4<>+8(SB)/8, $9650029242287828579
GLOBL sse2_prime64_4<>(SB), RODATA, $16

DATA sse2_prime64_5<>+0(SB)/8, $2870177450012600261
DATA sse2_prime64_5<>+8(SB)/8, $2870177450012600261
GLOBL sse2_prime64_5<>(SB), RODATA, $16

// processBatchXXHashSSE2 computes XXHash64 for multiple items in batch using SSE2
// Processes 2 items in parallel using 128-bit SIMD registers
// func processBatchXXHashSSE2(items [][]byte, results []HashResult, fingerprintBits, numBuckets uint)
TEXT ·processBatchXXHashSSE2(SB), NOSPLIT, $96-72
    // Load arguments
    MOVQ items_base+0(FP), DI
    MOVQ items_len+8(FP), SI
    MOVQ results_base+32(FP), R10
    MOVQ fingerprintBits+56(FP), R11
    MOVQ numBuckets+64(FP), R12

    // Save to stack
    MOVQ R11, 0(SP)   // fingerprintBits
    MOVQ R12, 8(SP)   // numBuckets

    // Calculate masks
    MOVQ R12, R13
    DECQ R13
    MOVQ R13, 16(SP)  // numBuckets mask

    MOVQ $1, R14
    MOVQ R11, CX
    SHLQ CL, R14
    DECQ R14
    MOVQ R14, 24(SP)  // fingerprint mask

    XORQ AX, AX       // item index

    // Check if we can process 2 items in parallel
    MOVQ SI, R15
    SUBQ $1, R15
    JL   scalar_loop

// Process 2 items in parallel using SSE2
simd_loop:
    CMPQ AX, R15
    JG   scalar_loop

    // Load 2 items
    MOVQ AX, BX
    IMULQ $24, BX

    // Load item 0
    MOVQ (DI)(BX*1), R8
    MOVQ 8(DI)(BX*1), R9
    MOVQ R8, 32(SP)   // item0 data ptr
    MOVQ R9, 40(SP)   // item0 length

    // Load item 1
    ADDQ $24, BX
    MOVQ (DI)(BX*1), R8
    MOVQ 8(DI)(BX*1), R9
    MOVQ R8, 48(SP)   // item1 data ptr
    MOVQ R9, 56(SP)   // item1 length

    // Initialize hash vector
    MOVDQU sse2_prime64_5<>(SB), X0
    MOVQ 40(SP), R8
    MOVQ 56(SP), R9
    MOVQ R8, X1
    PINSRQ $1, R9, X1
    PADDQ X1, X0      // X0 = hash values

    // Find minimum length
    CMPQ R9, R8
    CMOVQLT R9, R8
    MOVQ R8, 64(SP)   // min length

    XORQ CX, CX       // offset within data

simd_chunk_loop:
    CMPQ CX, R8
    JGE  simd_remainder

    MOVQ R8, R9
    SUBQ CX, R9
    CMPQ R9, $8
    JL   simd_remainder

    // Load 2x 8-byte values
    MOVQ 32(SP), BX
    MOVQ (BX)(CX*1), R9
    MOVQ R9, X1
    MOVQ 48(SP), BX
    MOVQ (BX)(CX*1), R9
    PINSRQ $1, R9, X1

    // k *= prime64_2 (emulated, SSE2 doesn't have PMULLQ)
    PEXTRQ $0, X1, R8
    PEXTRQ $1, X1, R9
    IMULQ prime64_2<>(SB), R8
    IMULQ prime64_2<>(SB), R9
    MOVQ R8, X1
    PINSRQ $1, R9, X1

    // k = rotl64(k, 31)
    PEXTRQ $0, X1, R8
    PEXTRQ $1, X1, R9
    ROLQ $31, R8
    ROLQ $31, R9
    MOVQ R8, X1
    PINSRQ $1, R9, X1

    // k *= prime64_1
    IMULQ prime64_1<>(SB), R8
    IMULQ prime64_1<>(SB), R9
    MOVQ R8, X1
    PINSRQ $1, R9, X1

    // hash ^= k
    PXOR X1, X0

    // hash = rotl64(hash, 27)
    PEXTRQ $0, X0, R8
    PEXTRQ $1, X0, R9
    ROLQ $27, R8
    ROLQ $27, R9

    // hash *= prime64_1 + prime64_4
    IMULQ prime64_1<>(SB), R8
    IMULQ prime64_1<>(SB), R9
    ADDQ prime64_4<>(SB), R8
    ADDQ prime64_4<>(SB), R9
    MOVQ R8, X0
    PINSRQ $1, R9, X0

    ADDQ $8, CX
    JMP  simd_chunk_loop

simd_remainder:
    // Extract hashes
    PEXTRQ $0, X0, R8
    MOVQ R8, 72(SP)   // hash0
    PEXTRQ $1, X0, R8
    MOVQ R8, 80(SP)   // hash1

    // Process each item's remainder
    MOVQ AX, BX       // Save item index
    XORQ DX, DX       // DX = 0..1

simd_finalize_loop:
    CMPQ DX, $2
    JGE  simd_finalize_done

    // Get item data ptr, length, and hash
    MOVQ DX, CX
    SHLQ $4, CX
    ADDQ $32, CX
    MOVQ (SP)(CX*1), R8     // data ptr
    MOVQ 8(SP)(CX*1), R9    // length

    MOVQ DX, CX
    SHLQ $3, CX
    ADDQ $72, CX
    MOVQ (SP)(CX*1), BP     // hash value

    // Process from min_length to actual length
    MOVQ 64(SP), CX

simd_item_remainder:
    CMPQ CX, R9
    JGE  simd_item_finalize

    MOVBQZX (R8)(CX*1), R15
    IMULQ prime64_5<>(SB), R15
    XORQ R15, BP
    ROLQ $11, BP
    IMULQ prime64_1<>(SB), BP

    INCQ CX
    JMP  simd_item_remainder

simd_item_finalize:
    // Avalanche
    MOVQ BP, R15
    SHRQ $33, R15
    XORQ R15, BP
    IMULQ prime64_2<>(SB), BP

    MOVQ BP, R15
    SHRQ $29, R15
    XORQ R15, BP
    IMULQ prime64_3<>(SB), BP

    MOVQ BP, R15
    SHRQ $32, R15
    XORQ R15, BP

    // Extract fingerprint
    MOVQ 24(SP), R14
    MOVQ BP, R15
    ANDQ R14, R15
    MOVB R15, CL
    TESTB CL, CL
    JNZ  simd_fp_ok
    MOVB $1, CL
simd_fp_ok:

    // Calculate i1
    MOVQ 16(SP), R13
    MOVQ BP, R8
    ANDQ R13, R8            // i1

    // Calculate i2: hash the fingerprint
    MOVBQZX CL, R9                  // R9 = fp
    MOVQ    prime64_5<>(SB), R11   // R11 = prime64_5
    ADDQ    $1, R11                // seed = prime64_5 + 1 (length=1)
    ADDQ    R9, R11                // hash = seed + fp
    MOVQ    R11, R9                // R9 = hash
    // Avalanche
    MOVQ    R9, R12
    SHRQ    $33, R12
    XORQ    R12, R9
    IMULQ   prime64_2<>(SB), R9
    MOVQ    R9, R12
    SHRQ    $29, R12
    XORQ    R12, R9
    IMULQ   prime64_3<>(SB), R9
    MOVQ    R9, R12
    SHRQ    $32, R12
    XORQ    R12, R9
    // Now R9 = hash(fp)
    XORQ R8, R9
    ANDQ R13, R9            // i2

    // Store result
    MOVQ BX, R15
    ADDQ DX, R15
    IMULQ $24, R15

    MOVQ R8, (R10)(R15*1)
    MOVQ R9, 8(R10)(R15*1)
    MOVB CL, 16(R10)(R15*1)

    INCQ DX
    JMP  simd_finalize_loop

simd_finalize_done:
    ADDQ $2, AX
    JMP  simd_loop

// Scalar processing for remaining items
scalar_loop:
    CMPQ AX, SI
    JGE  done

    MOVQ AX, BX
    IMULQ $24, BX
    MOVQ (DI)(BX*1), R8
    MOVQ 8(DI)(BX*1), CX

    MOVQ prime64_5<>(SB), BP
    ADDQ CX, BP

scalar_chunk_loop:
    CMPQ CX, $8
    JL   scalar_final_bytes

    MOVQ (R8), R9
    IMULQ prime64_2<>(SB), R9
    ROLQ $31, R9
    IMULQ prime64_1<>(SB), R9

    XORQ R9, BP
    ROLQ $27, BP
    IMULQ prime64_1<>(SB), BP
    ADDQ prime64_4<>(SB), BP

    ADDQ $8, R8
    SUBQ $8, CX
    JMP  scalar_chunk_loop

scalar_final_bytes:
    TESTQ CX, CX
    JZ    scalar_finalize

scalar_byte_loop:
    MOVBQZX (R8), R15
    IMULQ prime64_5<>(SB), R15
    XORQ R15, BP
    ROLQ $11, BP
    IMULQ prime64_1<>(SB), BP

    INCQ R8
    DECQ CX
    JNZ  scalar_byte_loop

scalar_finalize:
    MOVQ BP, R15
    SHRQ $33, R15
    XORQ R15, BP
    IMULQ prime64_2<>(SB), BP

    MOVQ BP, R15
    SHRQ $29, R15
    XORQ R15, BP
    IMULQ prime64_3<>(SB), BP

    MOVQ BP, R15
    SHRQ $32, R15
    XORQ R15, BP

    // Extract fingerprint
    MOVQ 24(SP), R14
    MOVQ BP, R15
    ANDQ R14, R15
    MOVB R15, CL
    TESTB CL, CL
    JNZ  scalar_fp_ok
    MOVB $1, CL
scalar_fp_ok:

    // Calculate i1
    MOVQ 16(SP), R13
    MOVQ BP, R8
    ANDQ R13, R8            // i1

    // Calculate i2: hash the fingerprint
    MOVBQZX CL, R9                  // R9 = fp
    MOVQ    prime64_5<>(SB), R11   // R11 = prime64_5
    ADDQ    $1, R11                // seed = prime64_5 + 1 (length=1)
    ADDQ    R9, R11                // hash = seed + fp
    MOVQ    R11, R9                // R9 = hash
    // Avalanche
    MOVQ    R9, R12
    SHRQ    $33, R12
    XORQ    R12, R9
    IMULQ   prime64_2<>(SB), R9
    MOVQ    R9, R12
    SHRQ    $29, R12
    XORQ    R12, R9
    IMULQ   prime64_3<>(SB), R9
    MOVQ    R9, R12
    SHRQ    $32, R12
    XORQ    R12, R9
    // Now R9 = hash(fp)
    XORQ R8, R9
    ANDQ R13, R9            // i2

    // Store result
    MOVQ AX, R15
    IMULQ $24, R15
    MOVQ R8, (R10)(R15*1)
    MOVQ R9, 8(R10)(R15*1)
    MOVB CL, 16(R10)(R15*1)

    INCQ AX
    JMP  scalar_loop

done:
    RET
