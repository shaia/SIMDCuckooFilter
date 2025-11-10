//go:build amd64
// +build amd64

#include "textflag.h"

// Note: Scalar prime64 constants are defined in xxhash_amd64.s
// This file only defines AVX2-specific vector constants

// AVX2 constants - 4x 64-bit values packed
DATA avx2_prime64_1<>+0(SB)/8, $11400714785074694791
DATA avx2_prime64_1<>+8(SB)/8, $11400714785074694791
DATA avx2_prime64_1<>+16(SB)/8, $11400714785074694791
DATA avx2_prime64_1<>+24(SB)/8, $11400714785074694791
GLOBL avx2_prime64_1<>(SB), RODATA, $32

DATA avx2_prime64_2<>+0(SB)/8, $14029467366897019727
DATA avx2_prime64_2<>+8(SB)/8, $14029467366897019727
DATA avx2_prime64_2<>+16(SB)/8, $14029467366897019727
DATA avx2_prime64_2<>+24(SB)/8, $14029467366897019727
GLOBL avx2_prime64_2<>(SB), RODATA, $32

DATA avx2_prime64_3<>+0(SB)/8, $1609587929392839161
DATA avx2_prime64_3<>+8(SB)/8, $1609587929392839161
DATA avx2_prime64_3<>+16(SB)/8, $1609587929392839161
DATA avx2_prime64_3<>+24(SB)/8, $1609587929392839161
GLOBL avx2_prime64_3<>(SB), RODATA, $32

DATA avx2_prime64_4<>+0(SB)/8, $9650029242287828579
DATA avx2_prime64_4<>+8(SB)/8, $9650029242287828579
DATA avx2_prime64_4<>+16(SB)/8, $9650029242287828579
DATA avx2_prime64_4<>+24(SB)/8, $9650029242287828579
GLOBL avx2_prime64_4<>(SB), RODATA, $32

DATA avx2_prime64_5<>+0(SB)/8, $2870177450012600261
DATA avx2_prime64_5<>+8(SB)/8, $2870177450012600261
DATA avx2_prime64_5<>+16(SB)/8, $2870177450012600261
DATA avx2_prime64_5<>+24(SB)/8, $2870177450012600261
GLOBL avx2_prime64_5<>(SB), RODATA, $32

// Shift constants for rotation
DATA avx2_shift_31<>+0(SB)/8, $31
DATA avx2_shift_31<>+8(SB)/8, $31
DATA avx2_shift_31<>+16(SB)/8, $31
DATA avx2_shift_31<>+24(SB)/8, $31
GLOBL avx2_shift_31<>(SB), RODATA, $32

DATA avx2_shift_33<>+0(SB)/8, $33
DATA avx2_shift_33<>+8(SB)/8, $33
DATA avx2_shift_33<>+16(SB)/8, $33
DATA avx2_shift_33<>+24(SB)/8, $33
GLOBL avx2_shift_33<>(SB), RODATA, $32

DATA avx2_shift_27<>+0(SB)/8, $27
DATA avx2_shift_27<>+8(SB)/8, $27
DATA avx2_shift_27<>+16(SB)/8, $27
DATA avx2_shift_27<>+24(SB)/8, $27
GLOBL avx2_shift_27<>(SB), RODATA, $32

DATA avx2_shift_37<>+0(SB)/8, $37
DATA avx2_shift_37<>+8(SB)/8, $37
DATA avx2_shift_37<>+16(SB)/8, $37
DATA avx2_shift_37<>+24(SB)/8, $37
GLOBL avx2_shift_37<>(SB), RODATA, $32

// processBatchXXHashAVX2 computes XXHash64 for multiple items in batch using AVX2
// Processes 4 items in parallel using 256-bit SIMD registers
// func processBatchXXHashAVX2(items [][]byte, results []HashResult, fingerprintBits, numBuckets uint)
TEXT ·processBatchXXHashAVX2(SB), NOSPLIT, $136-72
    // Load arguments
    MOVQ items_base+0(FP), DI        // DI = items slice base
    MOVQ items_len+8(FP), SI         // SI = number of items
    MOVQ results_base+32(FP), R10    // R10 = results slice base
    MOVQ fingerprintBits+56(FP), R11 // R11 = fingerprintBits
    MOVQ numBuckets+64(FP), R12      // R12 = numBuckets

    // Save non-volatile state to stack
    MOVQ R11, 0(SP)   // fingerprintBits
    MOVQ R12, 8(SP)   // numBuckets

    // Calculate masks
    MOVQ R12, R13
    DECQ R13          // R13 = numBuckets - 1
    MOVQ R13, 16(SP)  // Save numBuckets mask

    MOVQ $1, R14
    MOVQ R11, CX
    SHLQ CL, R14
    DECQ R14          // R14 = fingerprint mask
    MOVQ R14, 24(SP)  // Save fingerprint mask

    XORQ AX, AX       // AX = item index

    // Check if we can process 4 items in parallel
    MOVQ SI, R15
    SUBQ $3, R15
    JL   scalar_loop  // If items < 4, use scalar

// Process 4 items in parallel using AVX2
simd_loop:
    CMPQ AX, R15
    JG   scalar_loop

    // Load constants into YMM registers
    VMOVDQU avx2_prime64_5<>(SB), Y0   // Y0 = prime64_5 (4x)

    // Load 4 item lengths and initialize hashes
    // Stack layout: items[AX..AX+3] pointers and lengths
    MOVQ AX, BX
    IMULQ $24, BX     // BX = offset to items[AX]

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

    // Load item 2
    ADDQ $24, BX
    MOVQ (DI)(BX*1), R8
    MOVQ 8(DI)(BX*1), R9
    MOVQ R8, 64(SP)   // item2 data ptr
    MOVQ R9, 72(SP)   // item2 length

    // Load item 3
    ADDQ $24, BX
    MOVQ (DI)(BX*1), R8
    MOVQ 8(DI)(BX*1), R9
    MOVQ R8, 80(SP)   // item3 data ptr
    MOVQ R9, 88(SP)   // item3 length

    // Initialize hash vector: prime64_5 + length for each item
    // Load lengths individually and construct YMM register
    MOVQ 40(SP), R8   // len0
    MOVQ 56(SP), R9   // len1
    VPINSRQ $0, R8, X2, X2
    VPINSRQ $1, R9, X2, X2

    MOVQ 72(SP), R8   // len2
    MOVQ 88(SP), R9   // len3
    VPINSRQ $0, R8, X3, X3
    VPINSRQ $1, R9, X3, X3

    VINSERTI128 $1, X3, Y2, Y2    // Y2 = 4x lengths
    VPADDQ Y2, Y0, Y1              // Y1 = prime64_5 + len for each item

    // Find minimum length for aligned processing
    MOVQ 40(SP), R8
    MOVQ 56(SP), R9
    CMPQ R9, R8
    CMOVQLT R9, R8
    MOVQ 72(SP), R9
    CMPQ R9, R8
    CMOVQLT R9, R8
    MOVQ 88(SP), R9
    CMPQ R9, R8
    CMOVQLT R9, R8    // R8 = min length

    MOVQ R8, 96(SP)   // Save min length

    // Process 8-byte chunks for all 4 items in parallel
    XORQ CX, CX       // CX = offset within data

simd_chunk_loop:
    CMPQ CX, R8
    JGE  simd_remainder

    MOVQ R8, R9
    SUBQ CX, R9
    CMPQ R9, $8
    JL   simd_remainder

    // Load 4x 8-byte values from each item
    MOVQ 32(SP), BX
    MOVQ (BX)(CX*1), R9
    VPINSRQ $0, R9, X2, X2

    MOVQ 48(SP), BX
    MOVQ (BX)(CX*1), R9
    VPINSRQ $1, R9, X2, X2

    MOVQ 64(SP), BX
    MOVQ (BX)(CX*1), R9
    VPINSRQ $0, R9, X3, X3

    MOVQ 80(SP), BX
    MOVQ (BX)(CX*1), R9
    VPINSRQ $1, R9, X3, X3

    VINSERTI128 $1, X3, Y2, Y2    // Y2 = 4x uint64 values

    // k *= prime64_2
    // AVX2 doesn't have 64-bit multiply, so we extract to scalar, multiply, and reinsert
    VMOVDQU avx2_prime64_2<>(SB), Y3

    // Extract each 64-bit value, multiply, and reinsert
    VEXTRACTI128 $0, Y2, X4
    VPEXTRQ $0, X4, R9
    IMULQ prime64_2<>(SB), R9
    VPINSRQ $0, R9, X4, X4

    VPEXTRQ $1, X4, R9
    IMULQ prime64_2<>(SB), R9
    VPINSRQ $1, R9, X4, X4

    VEXTRACTI128 $1, Y2, X5
    VPEXTRQ $0, X5, R9
    IMULQ prime64_2<>(SB), R9
    VPINSRQ $0, R9, X5, X5

    VPEXTRQ $1, X5, R9
    IMULQ prime64_2<>(SB), R9
    VPINSRQ $1, R9, X5, X5

    VINSERTI128 $1, X5, Y2, Y2

    // k = rotl64(k, 31) - requires scalar extraction for now
    // For true SIMD, we'd need to extract, rotate, and reinsert
    // Simplified: process scalar for complex operations
    VEXTRACTI128 $0, Y2, X3
    VPEXTRQ $0, X3, R9
    ROLQ $31, R9
    VPINSRQ $0, R9, X3, X3

    VPEXTRQ $1, X3, R9
    ROLQ $31, R9
    VPINSRQ $1, R9, X3, X3

    VEXTRACTI128 $1, Y2, X4
    VPEXTRQ $0, X4, R9
    ROLQ $31, R9
    VPINSRQ $0, R9, X4, X4

    VPEXTRQ $1, X4, R9
    ROLQ $31, R9
    VPINSRQ $1, R9, X4, X4

    VINSERTI128 $1, X4, Y2, Y2

    // k *= prime64_1
    VMOVDQU avx2_prime64_1<>(SB), Y3

    // Extract each 64-bit value, multiply, and reinsert
    VEXTRACTI128 $0, Y2, X4
    VPEXTRQ $0, X4, R9
    IMULQ prime64_1<>(SB), R9
    VPINSRQ $0, R9, X4, X4

    VPEXTRQ $1, X4, R9
    IMULQ prime64_1<>(SB), R9
    VPINSRQ $1, R9, X4, X4

    VEXTRACTI128 $1, Y2, X5
    VPEXTRQ $0, X5, R9
    IMULQ prime64_1<>(SB), R9
    VPINSRQ $0, R9, X5, X5

    VPEXTRQ $1, X5, R9
    IMULQ prime64_1<>(SB), R9
    VPINSRQ $1, R9, X5, X5

    VINSERTI128 $1, X5, Y2, Y2

    // hash ^= k
    VPXOR Y2, Y1, Y1

    // hash = rotl64(hash, 27) - extract and rotate
    VEXTRACTI128 $0, Y1, X3
    VPEXTRQ $0, X3, R9
    ROLQ $27, R9
    VPINSRQ $0, R9, X3, X3

    VPEXTRQ $1, X3, R9
    ROLQ $27, R9
    VPINSRQ $1, R9, X3, X3

    VEXTRACTI128 $1, Y1, X4
    VPEXTRQ $0, X4, R9
    ROLQ $27, R9
    VPINSRQ $0, R9, X4, X4

    VPEXTRQ $1, X4, R9
    ROLQ $27, R9
    VPINSRQ $1, R9, X4, X4

    VINSERTI128 $1, X4, Y2, Y2

    // hash *= prime64_1
    // Extract each 64-bit value, multiply, and reinsert
    VEXTRACTI128 $0, Y2, X4
    VPEXTRQ $0, X4, R9
    IMULQ prime64_1<>(SB), R9
    VPINSRQ $0, R9, X4, X4

    VPEXTRQ $1, X4, R9
    IMULQ prime64_1<>(SB), R9
    VPINSRQ $1, R9, X4, X4

    VEXTRACTI128 $1, Y2, X5
    VPEXTRQ $0, X5, R9
    IMULQ prime64_1<>(SB), R9
    VPINSRQ $0, R9, X5, X5

    VPEXTRQ $1, X5, R9
    IMULQ prime64_1<>(SB), R9
    VPINSRQ $1, R9, X5, X5

    VINSERTI128 $1, X5, Y2, Y2

    // hash + prime64_4
    VMOVDQU avx2_prime64_4<>(SB), Y3
    VPADDQ Y3, Y2, Y1

    ADDQ $8, CX
    JMP  simd_chunk_loop

simd_remainder:
    // Process remaining bytes for each item individually (fallback to scalar)
    // Extract hashes and continue scalar processing
    VEXTRACTI128 $0, Y1, X2
    VPEXTRQ $0, X2, R8
    MOVQ R8, 104(SP)  // hash0
    VPEXTRQ $1, X2, R8
    MOVQ R8, 112(SP)  // hash1

    VEXTRACTI128 $1, Y1, X2
    VPEXTRQ $0, X2, R8
    MOVQ R8, 120(SP)  // hash2
    VPEXTRQ $1, X2, R8
    MOVQ R8, 128(SP)  // hash3

    // Process each item's remainder and finalize individually
    MOVQ AX, BX       // Save item index
    XORQ DX, DX       // DX = 0..3 for each of the 4 items

simd_finalize_loop:
    CMPQ DX, $4
    JGE  simd_finalize_done

    // Get item data ptr, length, and hash
    MOVQ DX, CX
    SHLQ $4, CX       // CX = DX * 16
    ADDQ $32, CX
    MOVQ (SP)(CX*1), R8     // data ptr
    MOVQ 8(SP)(CX*1), R9    // length

    MOVQ DX, CX
    SHLQ $3, CX
    ADDQ $104, CX
    MOVQ (SP)(CX*1), BP     // hash value

    // Process from min_length to actual length
    MOVQ 96(SP), CX         // CX = min_length (already processed)

simd_item_remainder:
    CMPQ CX, R9
    JGE  simd_item_finalize

    MOVBQZX (R8)(CX*1), R15
    MOVQ prime64_5<>(SB), R14
    IMULQ R14, R15
    XORQ R15, BP

    MOVQ BP, R15
    ROLQ $11, R15
    MOVQ R15, BP

    MOVQ prime64_1<>(SB), R15
    IMULQ R15, BP

    INCQ CX
    JMP  simd_item_remainder

simd_item_finalize:
    // Avalanche
    MOVQ BP, R15
    SHRQ $33, R15
    XORQ R15, BP

    MOVQ prime64_2<>(SB), R15
    IMULQ R15, BP

    MOVQ BP, R15
    SHRQ $29, R15
    XORQ R15, BP

    MOVQ prime64_3<>(SB), R15
    IMULQ R15, BP

    MOVQ BP, R15
    SHRQ $32, R15
    XORQ R15, BP

    // Extract fingerprint
    MOVQ 24(SP), R14        // fingerprint mask
    MOVQ BP, R15
    ANDQ R14, R15
    MOVB R15, CL
    TESTB CL, CL
    JNZ  simd_fp_ok
    MOVB $1, CL
simd_fp_ok:

    // Calculate i1
    MOVQ 16(SP), R13        // numBuckets mask
    MOVQ BP, R8
    ANDQ R13, R8            // i1

    // Calculate i2: Compute XXHash64 of fingerprint byte in CL, store in R9
    MOVBQZX CL, R9                  // R9 = fp
    MOVQ    prime64_5<>(SB), R11   // R11 = prime64_5
    ADDQ    $1, R11                // seed = prime64_5 + 1 (length=1)
    ADDQ    R9, R11                // hash = seed + fp
    MOVQ    R11, R9                // R9 = hash
    // Avalanche (as in XXHash64 for <=8 bytes)
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
    XORQ    R8, R9
    ANDQ    R13, R9            // i2

    // Store result
    MOVQ BX, R15            // Original item index
    ADDQ DX, R15            // + current offset (0..3)
    IMULQ $24, R15

    MOVQ R8, (R10)(R15*1)
    MOVQ R9, 8(R10)(R15*1)
    MOVB CL, 16(R10)(R15*1)

    INCQ DX
    JMP  simd_finalize_loop

simd_finalize_done:
    ADDQ $4, AX
    JMP  simd_loop

// Scalar processing for remaining items (< 4)
scalar_loop:
    CMPQ AX, SI
    JGE  done

    // Process single item using scalar code
    MOVQ AX, BX
    IMULQ $24, BX
    MOVQ (DI)(BX*1), R8     // data ptr
    MOVQ 8(DI)(BX*1), CX    // length

    MOVQ prime64_5<>(SB), BP
    ADDQ CX, BP

scalar_chunk_loop:
    CMPQ CX, $8
    JL   scalar_final_bytes

    MOVQ (R8), R9
    MOVQ R9, R15
    IMULQ prime64_2<>(SB), R15

    ROLQ $31, R15
    IMULQ prime64_1<>(SB), R15

    XORQ R15, BP
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
    VZEROUPPER
    RET
