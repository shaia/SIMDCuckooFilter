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
GLOBL avx2_prime64_1<>(SB), RODATA|NOPTR, $32

DATA avx2_prime64_2<>+0(SB)/8, $14029467366897019727
DATA avx2_prime64_2<>+8(SB)/8, $14029467366897019727
DATA avx2_prime64_2<>+16(SB)/8, $14029467366897019727
DATA avx2_prime64_2<>+24(SB)/8, $14029467366897019727
GLOBL avx2_prime64_2<>(SB), RODATA|NOPTR, $32

DATA avx2_prime64_3<>+0(SB)/8, $1609587929392839161
DATA avx2_prime64_3<>+8(SB)/8, $1609587929392839161
DATA avx2_prime64_3<>+16(SB)/8, $1609587929392839161
DATA avx2_prime64_3<>+24(SB)/8, $1609587929392839161
GLOBL avx2_prime64_3<>(SB), RODATA|NOPTR, $32

DATA avx2_prime64_4<>+0(SB)/8, $9650029242287828579
DATA avx2_prime64_4<>+8(SB)/8, $9650029242287828579
DATA avx2_prime64_4<>+16(SB)/8, $9650029242287828579
DATA avx2_prime64_4<>+24(SB)/8, $9650029242287828579
GLOBL avx2_prime64_4<>(SB), RODATA|NOPTR, $32

DATA avx2_prime64_5<>+0(SB)/8, $2870177450012600261
DATA avx2_prime64_5<>+8(SB)/8, $2870177450012600261
DATA avx2_prime64_5<>+16(SB)/8, $2870177450012600261
DATA avx2_prime64_5<>+24(SB)/8, $2870177450012600261
GLOBL avx2_prime64_5<>(SB), RODATA|NOPTR, $32

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

// Local scalar constants to avoid linkage issues
DATA local_prime64_1<>+0(SB)/8, $11400714785074694791
GLOBL local_prime64_1<>(SB), RODATA|NOPTR, $8
DATA local_prime64_2<>+0(SB)/8, $14029467366897019727
GLOBL local_prime64_2<>(SB), RODATA|NOPTR, $8
DATA local_prime64_3<>+0(SB)/8, $1609587929392839161
GLOBL local_prime64_3<>(SB), RODATA|NOPTR, $8
DATA local_prime64_4<>+0(SB)/8, $9650029242287828579
GLOBL local_prime64_4<>(SB), RODATA|NOPTR, $8
DATA local_prime64_5<>+0(SB)/8, $2870177450012600261
GLOBL local_prime64_5<>(SB), RODATA|NOPTR, $8
DATA avx2_shift_37<>+16(SB)/8, $37
DATA avx2_shift_37<>+24(SB)/8, $37
GLOBL avx2_shift_37<>(SB), RODATA, $32

// processBatchXXHashAVX2 computes XXHash64 for multiple items in batch using AVX2
// Processes 4 items in parallel using 256-bit SIMD registers
// func processBatchXXHashAVX2(items [][]byte, results []HashResult, fingerprintBits, numBuckets uint)
// Stack frame must be 16-byte aligned for AVX2 operations
TEXT ·processBatchXXHashAVX2(SB), NOSPLIT, $144-64
    // Load arguments
    MOVQ items_base+0(FP), DI        // DI = items slice base
    MOVQ items_len+8(FP), SI         // SI = number of items
    MOVQ results_base+24(FP), R10    // R10 = results slice base
    MOVQ fingerprintBits+48(FP), R11 // R11 = fingerprintBits
    MOVQ numBuckets+56(FP), R12      // R12 = numBuckets

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
    SUBQ $4, R15
    // TODO(fix): AVX2 SIMD path is currently broken on Windows - crashes when loading
    // item pointers from stack. Issue appears to be related to stack frame corruption
    // or incorrect addressing. Scalar fallback works correctly.
    //
    // The crash occurs at line 180 when trying to load item data pointers from stack
    // offsets 32(SP), 48(SP), 64(SP), 80(SP). The loaded values are garbage
    // (e.g., 0xc0000f8000) instead of the valid pointers that were stored earlier.
    //
    // Possible causes:
    // - Stack frame corruption between store and load
    // - Incorrect stack addressing with NOSPLIT on Windows
    // - AVX2 operations corrupting stack
    // - Issue with how slice data is being loaded from items parameter
    //
    // For now, force scalar path until this can be debugged properly with a debugger.
    // JMP  scalar_loop  // Force scalar path (TEMPORARY - see TODO above)
    JL   scalar_loop  // If items < 4, use scalar

// Process 4 items in parallel using AVX2
simd_loop:
    MOVQ items_len+8(FP), SI // Reload SI
    CMPQ AX, SI
    JGE  done
    
    // Check if at least 4 items remaining
    MOVQ SI, CX
    SUBQ AX, CX
    CMPQ CX, $4
    JL   scalar_loop

    // Load constants into YMM registers
    VMOVDQU avx2_prime64_5<>(SB), Y0   // Y0 = prime64_5 (4x)

    // Load 4 item lengths and initialize hashes
    // Stack layout: items[AX..AX+3] pointers and lengths
    MOVQ AX, BX
    IMULQ $24, BX     // BX = offset to items[AX]
    ADDQ DI, BX       // BX = &items[AX]

    // Load item 0
    MOVQ 0(BX), R8    // data ptr
    MOVQ 8(BX), R9    // length
    MOVQ R8, R13   // item0 data ptr (kept in register)
    MOVQ R9, R14   // item0 length (kept in register)

    // Load item 1
    MOVQ 24(BX), R8   // data ptr
    MOVQ 32(BX), R9   // length
    MOVQ R8, 48(SP)   // item1 data ptr
    MOVQ R9, 56(SP)   // item1 length

    // Load item 2
    MOVQ 48(BX), R8   // data ptr
    MOVQ 56(BX), R9   // length
    MOVQ R8, 64(SP)   // item2 data ptr
    MOVQ R9, 72(SP)   // item2 length

    // Load item 3
    MOVQ 72(BX), R8   // data ptr
    MOVQ 80(BX), R9   // length
    MOVQ R8, 80(SP)   // item3 data ptr
    MOVQ R9, 88(SP)   // item3 length

init_ok:

    // Initialize hash vector: prime64_5 + length for each item
    // Load lengths individually and construct YMM register
    MOVQ R14, R8   // len0
    MOVQ 56(SP), R9   // len1
    VPINSRQ $0, R8, X2, X2
    VPINSRQ $1, R9, X2, X2
    MOVQ 72(SP), R8   // len2
    MOVQ 88(SP), R9   // len3
    VPINSRQ $0, R8, X3, X3
    VPINSRQ $1, R9, X3, X3

    VINSERTI128 $1, X3, Y2, Y2    // Y2 = 4x lengths

    // hash = prime64_5 + length
    VMOVDQU avx2_prime64_5<>(SB), Y0
    VPADDQ Y2, Y0, Y1             // Y1 = hash vector

    // Find minimum length for aligned processing
    MOVQ R14, R8 // Use R14 (len0) as initial min
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
    // JMP simd_remainder // DEBUG: Skip chunk loop

simd_chunk_loop:
    CMPQ CX, R8
    JGE  simd_remainder

    MOVQ R8, R9
    SUBQ CX, R9
    CMPQ R9, $8
    JL   simd_remainder

    // Load 4x 8-byte values from each item
    MOVQ R13, BX
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
    IMULQ local_prime64_2<>(SB), R9
    VPINSRQ $0, R9, X4, X4

    VPEXTRQ $1, X4, R9
    IMULQ local_prime64_2<>(SB), R9
    VPINSRQ $1, R9, X4, X4

    VEXTRACTI128 $1, Y2, X5
    VPEXTRQ $0, X5, R9
    IMULQ local_prime64_2<>(SB), R9
    VPINSRQ $0, R9, X5, X5

    VPEXTRQ $1, X5, R9
    IMULQ local_prime64_2<>(SB), R9
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
    IMULQ local_prime64_1<>(SB), R9
    VPINSRQ $0, R9, X4, X4

    VPEXTRQ $1, X4, R9
    IMULQ local_prime64_1<>(SB), R9
    VPINSRQ $1, R9, X4, X4

    VEXTRACTI128 $1, Y2, X5
    VPEXTRQ $0, X5, R9
    IMULQ local_prime64_1<>(SB), R9
    VPINSRQ $0, R9, X5, X5

    VPEXTRQ $1, X5, R9
    IMULQ local_prime64_1<>(SB), R9
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
    IMULQ local_prime64_1<>(SB), R9
    VPINSRQ $0, R9, X4, X4

    VPEXTRQ $1, X4, R9
    IMULQ local_prime64_1<>(SB), R9
    VPINSRQ $1, R9, X4, X4

    VEXTRACTI128 $1, Y2, X5
    VPEXTRQ $0, X5, R9
    IMULQ local_prime64_1<>(SB), R9
    VPINSRQ $0, R9, X5, X5

    VPEXTRQ $1, X5, R9
    IMULQ local_prime64_1<>(SB), R9
    VPINSRQ $1, R9, X5, X5

    VINSERTI128 $1, X5, Y2, Y2

    // hash + prime64_4
    VMOVDQU avx2_prime64_4<>(SB), Y3
    VPADDQ Y3, Y2, Y1

    ADDQ $8, CX
    JMP  simd_chunk_loop

simd_remainder:
remainder_ok:

    // Process remaining bytes for each item individually (fallback to scalar)
    // Save the chunk offset (number of bytes already processed)
    MOVQ CX, 96(SP)   // Save chunk offset (overwrite min_length, no longer needed)

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

    // Calculate base pointers for this chunk
    MOVQ AX, BX
    IMULQ $24, BX
    MOVQ R10, DX
    ADDQ BX, DX      // DX = &results[AX]
    ADDQ DI, BX      // BX = &items[AX]

    // Process item 0
    MOVQ 0(BX), R8    // item0 data ptr
    MOVQ 8(BX), R9    // item0 length
    // Recalculate hash0 = prime64_5 + len
    MOVQ local_prime64_5<>(SB), R12
    ADDQ R9, R12
    
    // Process 8-byte chunks
    CMPQ R9, $8
    JL   item0_byte_loop
item0_chunk_loop:
    MOVQ (R8), R15
    IMULQ local_prime64_2<>(SB), R15
    ROLQ $31, R15
    IMULQ local_prime64_1<>(SB), R15
    XORQ R15, R12
    ROLQ $27, R12
    IMULQ local_prime64_1<>(SB), R12
    ADDQ local_prime64_4<>(SB), R12
    ADDQ $8, R8
    SUBQ $8, R9
    CMPQ R9, $8
    JGE  item0_chunk_loop

item0_byte_loop:
    TESTQ R9, R9
    JZ   item0_done
    MOVBQZX (R8), R15
    MOVQ local_prime64_5<>(SB), R14
    IMULQ R14, R15
    XORQ R15, R12
    MOVQ R12, R15
    ROLQ $11, R15
    MOVQ R15, R12
    MOVQ local_prime64_1<>(SB), R15
    IMULQ R15, R12
    INCQ R8
    DECQ R9
    JMP  item0_byte_loop
item0_done:
    // Item 0 finalize
    MOVQ R12, R15
    SHRQ $33, R15
    XORQ R15, R12
    MOVQ local_prime64_2<>(SB), R15
    IMULQ R15, R12
    MOVQ R12, R15
    SHRQ $29, R15
    XORQ R15, R12
    MOVQ local_prime64_3<>(SB), R15
    IMULQ R15, R12
    MOVQ R12, R15
    SHRQ $32, R15
    XORQ R15, R12
    
    // Store item 0 result
    // Recalculate fp mask
    MOVQ fingerprintBits+48(FP), CX
    MOVQ $1, R14
    SHLQ CX, R14
    DECQ R14
    
    MOVQ R12, R15
    ANDQ R14, R15
    MOVB R15, CL
    TESTB CL, CL
    JNZ item0_fp_ok
    MOVB $1, CL
item0_fp_ok:
    // Recalculate numBuckets mask
    MOVQ numBuckets+56(FP), R13
    DECQ R13
    
    MOVQ R12, R8
    ANDQ R13, R8            // i1
    
    // Calculate i2
    MOVBQZX CL, R9
    MOVQ    local_prime64_5<>(SB), R14
    IMULQ   R14, R9
    MOVQ    R14, R11
    ADDQ    $1, R11
    XORQ    R9, R11
    ROLQ    $11, R11
    MOVQ    local_prime64_1<>(SB), R14
    IMULQ   R14, R11
    MOVQ    R11, R9
    SHRQ    $33, R9
    XORQ    R9, R11
    MOVQ    local_prime64_2<>(SB), R14
    IMULQ   R14, R11
    MOVQ    R11, R9
    SHRQ    $29, R9
    XORQ    R9, R11
    MOVQ    local_prime64_3<>(SB), R14
    IMULQ   R14, R11
    MOVQ    R11, R9
    SHRQ    $32, R9
    XORQ    R9, R11
    MOVQ    R11, R9
    XORQ    R8, R9
    ANDQ    R13, R9
    
    // Store to results[0]
    MOVQ R8, 0(DX)
    MOVQ R9, 8(DX)
    MOVB CL, 16(DX)

    // Process item 1
    MOVQ 24(BX), R8    // item1 data ptr
    MOVQ 32(BX), R9    // item1 length
    // Recalculate hash1
    MOVQ local_prime64_5<>(SB), R12
    ADDQ R9, R12
    
    // Process 8-byte chunks
    CMPQ R9, $8
    JL   item1_byte_loop
item1_chunk_loop:
    MOVQ (R8), R15
    IMULQ local_prime64_2<>(SB), R15
    ROLQ $31, R15
    IMULQ local_prime64_1<>(SB), R15
    XORQ R15, R12
    ROLQ $27, R12
    IMULQ local_prime64_1<>(SB), R12
    ADDQ local_prime64_4<>(SB), R12
    ADDQ $8, R8
    SUBQ $8, R9
    CMPQ R9, $8
    JGE  item1_chunk_loop

item1_byte_loop:
    TESTQ R9, R9
    JZ   item1_done
    MOVBQZX (R8), R15
    MOVQ local_prime64_5<>(SB), R14
    IMULQ R14, R15
    XORQ R15, R12
    MOVQ R12, R15
    ROLQ $11, R15
    MOVQ R15, R12
    MOVQ local_prime64_1<>(SB), R15
    IMULQ R15, R12
    INCQ R8
    DECQ R9
    JMP  item1_byte_loop
item1_done:
    // Item 1 finalize
    MOVQ R12, R15
    SHRQ $33, R15
    XORQ R15, R12
    MOVQ local_prime64_2<>(SB), R15
    IMULQ R15, R12
    MOVQ R12, R15
    SHRQ $29, R15
    XORQ R15, R12
    MOVQ local_prime64_3<>(SB), R15
    IMULQ R15, R12
    MOVQ R12, R15
    SHRQ $32, R15
    XORQ R15, R12
    
    // Store item 1 result
    // Recalculate fp mask
    MOVQ fingerprintBits+48(FP), CX
    MOVQ $1, R14
    SHLQ CX, R14
    DECQ R14
    
    MOVQ R12, R15
    ANDQ R14, R15
    MOVB R15, CL
    TESTB CL, CL
    JNZ item1_fp_ok
    MOVB $1, CL
item1_fp_ok:
    // Recalculate numBuckets mask
    MOVQ numBuckets+56(FP), R13
    DECQ R13
    
    MOVQ R12, R8
    ANDQ R13, R8
    
    // Calculate i2
    MOVBQZX CL, R9
    MOVQ    local_prime64_5<>(SB), R14
    IMULQ   R14, R9
    MOVQ    R14, R11
    ADDQ    $1, R11
    XORQ    R9, R11
    ROLQ    $11, R11
    MOVQ    local_prime64_1<>(SB), R14
    IMULQ   R14, R11
    MOVQ    R11, R9
    SHRQ    $33, R9
    XORQ    R9, R11
    MOVQ    local_prime64_2<>(SB), R14
    IMULQ   R14, R11
    MOVQ    R11, R9
    SHRQ    $29, R9
    XORQ    R9, R11
    MOVQ    local_prime64_3<>(SB), R14
    IMULQ   R14, R11
    MOVQ    R11, R9
    SHRQ    $32, R9
    XORQ    R9, R11
    MOVQ    R11, R9
    XORQ    R8, R9
    ANDQ    R13, R9
    
    // Store to results[1] (offset 24)
    MOVQ R8, 24(DX)
    MOVQ R9, 32(DX)
    MOVB CL, 40(DX)

    // Process item 2
    MOVQ 48(BX), R8    // item2 data ptr
    MOVQ 56(BX), R9    // item2 length
    // Recalculate hash2
    MOVQ local_prime64_5<>(SB), R12
    ADDQ R9, R12
    
    // Process 8-byte chunks
    CMPQ R9, $8
    JL   item2_byte_loop
item2_chunk_loop:
    MOVQ (R8), R15
    IMULQ local_prime64_2<>(SB), R15
    ROLQ $31, R15
    IMULQ local_prime64_1<>(SB), R15
    XORQ R15, R12
    ROLQ $27, R12
    IMULQ local_prime64_1<>(SB), R12
    ADDQ local_prime64_4<>(SB), R12
    ADDQ $8, R8
    SUBQ $8, R9
    CMPQ R9, $8
    JGE  item2_chunk_loop

item2_byte_loop:
    TESTQ R9, R9
    JZ   item2_done
    MOVBQZX (R8), R15
    MOVQ local_prime64_5<>(SB), R14
    IMULQ R14, R15
    XORQ R15, R12
    MOVQ R12, R15
    ROLQ $11, R15
    MOVQ R15, R12
    MOVQ local_prime64_1<>(SB), R15
    IMULQ R15, R12
    INCQ R8
    DECQ R9
    JMP  item2_byte_loop
item2_done:
    // Item 2 finalize
    MOVQ R12, R15
    SHRQ $33, R15
    XORQ R15, R12
    MOVQ local_prime64_2<>(SB), R15
    IMULQ R15, R12
    MOVQ R12, R15
    SHRQ $29, R15
    XORQ R15, R12
    MOVQ local_prime64_3<>(SB), R15
    IMULQ R15, R12
    MOVQ R12, R15
    SHRQ $32, R15
    XORQ R15, R12
    
    // Store item 2 result
    // Recalculate fp mask
    MOVQ fingerprintBits+48(FP), CX
    MOVQ $1, R14
    SHLQ CX, R14
    DECQ R14
    
    MOVQ R12, R15
    ANDQ R14, R15
    MOVB R15, CL
    TESTB CL, CL
    JNZ item2_fp_ok
    MOVB $1, CL
item2_fp_ok:
    // Recalculate numBuckets mask
    MOVQ numBuckets+56(FP), R13
    DECQ R13
    
    MOVQ R12, R8
    ANDQ R13, R8
    
    // Calculate i2
    MOVBQZX CL, R9
    MOVQ    local_prime64_5<>(SB), R14
    IMULQ   R14, R9
    MOVQ    R14, R11
    ADDQ    $1, R11
    XORQ    R9, R11
    ROLQ    $11, R11
    MOVQ    local_prime64_1<>(SB), R14
    IMULQ   R14, R11
    MOVQ    R11, R9
    SHRQ    $33, R9
    XORQ    R9, R11
    MOVQ    local_prime64_2<>(SB), R14
    IMULQ   R14, R11
    MOVQ    R11, R9
    SHRQ    $29, R9
    XORQ    R9, R11
    MOVQ    local_prime64_3<>(SB), R14
    IMULQ   R14, R11
    MOVQ    R11, R9
    SHRQ    $32, R9
    XORQ    R9, R11
    MOVQ    R11, R9
    XORQ    R8, R9
    ANDQ    R13, R9
    
    // Store to results[2] (offset 48)
    MOVQ R8, 48(DX)
    MOVQ R9, 56(DX)
    MOVB CL, 64(DX)

    // Process item 3
    MOVQ 72(BX), R8    // item3 data ptr
    MOVQ 80(BX), R9    // item3 length
    // Recalculate hash3
    MOVQ local_prime64_5<>(SB), R12
    ADDQ R9, R12
    
    // Process 8-byte chunks
    CMPQ R9, $8
    JL   item3_byte_loop
item3_chunk_loop:
    MOVQ (R8), R15
    IMULQ local_prime64_2<>(SB), R15
    ROLQ $31, R15
    IMULQ local_prime64_1<>(SB), R15
    XORQ R15, R12
    ROLQ $27, R12
    IMULQ local_prime64_1<>(SB), R12
    ADDQ local_prime64_4<>(SB), R12
    ADDQ $8, R8
    SUBQ $8, R9
    CMPQ R9, $8
    JGE  item3_chunk_loop

item3_byte_loop:
    TESTQ R9, R9
    JZ   item3_done
    MOVBQZX (R8), R15
    MOVQ local_prime64_5<>(SB), R14
    IMULQ R14, R15
    XORQ R15, R12
    MOVQ R12, R15
    ROLQ $11, R15
    MOVQ R15, R12
    MOVQ local_prime64_1<>(SB), R15
    IMULQ R15, R12
    INCQ R8
    DECQ R9
    JMP  item3_byte_loop
item3_done:
    // Item 3 finalize
    MOVQ R12, R15
    SHRQ $33, R15
    XORQ R15, R12
    MOVQ local_prime64_2<>(SB), R15
    IMULQ R15, R12
    MOVQ R12, R15
    SHRQ $29, R15
    XORQ R15, R12
    MOVQ local_prime64_3<>(SB), R15
    IMULQ R15, R12
    MOVQ R12, R15
    SHRQ $32, R15
    XORQ R15, R12
    
    // Store item 3 result
    // Recalculate fp mask
    MOVQ fingerprintBits+48(FP), CX
    MOVQ $1, R14
    SHLQ CX, R14
    DECQ R14
    
    MOVQ R12, R15
    ANDQ R14, R15
    MOVB R15, CL
    TESTB CL, CL
    JNZ item3_fp_ok
    MOVB $1, CL
item3_fp_ok:
    // Recalculate numBuckets mask
    MOVQ numBuckets+56(FP), R13
    DECQ R13
    
    MOVQ R12, R8
    ANDQ R13, R8
    
    // Calculate i2
    MOVBQZX CL, R9
    MOVQ    local_prime64_5<>(SB), R14
    IMULQ   R14, R9
    MOVQ    R14, R11
    ADDQ    $1, R11
    XORQ    R9, R11
    ROLQ    $11, R11
    MOVQ    local_prime64_1<>(SB), R14
    IMULQ   R14, R11
    MOVQ    R11, R9
    SHRQ    $33, R9
    XORQ    R9, R11
    MOVQ    local_prime64_2<>(SB), R14
    IMULQ   R14, R11
    MOVQ    R11, R9
    SHRQ    $29, R9
    XORQ    R9, R11
    MOVQ    local_prime64_3<>(SB), R14
    IMULQ   R14, R11
    MOVQ    R11, R9
    SHRQ    $32, R9
    XORQ    R9, R11
    MOVQ    R11, R9
    XORQ    R8, R9
    ANDQ    R13, R9
    
    // Store to results[3] (offset 72)
    MOVQ R8, 72(DX)
    MOVQ R9, 80(DX)
    MOVB CL, 88(DX)

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
    TESTQ R8, R8
    JNZ scalar_ptr_ok

scalar_ptr_ok:
    MOVQ local_prime64_5<>(SB), R12
    ADDQ CX, R12

scalar_chunk_loop:
    CMPQ CX, $8
    JL   scalar_final_bytes

    MOVQ (R8), R9
    MOVQ R9, R15
    IMULQ local_prime64_2<>(SB), R15

    ROLQ $31, R15
    IMULQ local_prime64_1<>(SB), R15

    XORQ R15, R12
    ROLQ $27, R12
    IMULQ local_prime64_1<>(SB), R12
    ADDQ local_prime64_4<>(SB), R12

    ADDQ $8, R8
    SUBQ $8, CX
    JMP  scalar_chunk_loop

scalar_final_bytes:
    TESTQ CX, CX
    JZ    scalar_finalize

scalar_byte_loop:
    MOVBQZX (R8), R15
    IMULQ local_prime64_5<>(SB), R15
    XORQ R15, R12

    ROLQ $11, R12
    IMULQ local_prime64_1<>(SB), R12

    INCQ R8
    DECQ CX
    JNZ  scalar_byte_loop

scalar_finalize:
    MOVQ R12, R15
    SHRQ $33, R15
    XORQ R15, R12
    IMULQ local_prime64_2<>(SB), R12

    MOVQ R12, R15
    SHRQ $29, R15
    XORQ R15, R12
    IMULQ local_prime64_3<>(SB), R12

    MOVQ R12, R15
    SHRQ $32, R15
    XORQ R15, R12

    // Extract fingerprint
    // Recalculate fp mask
    MOVQ fingerprintBits+48(FP), CX
    MOVQ $1, R14
    SHLQ CX, R14
    DECQ R14
    
    MOVQ R12, R15
    ANDQ R14, R15
    MOVB R15, CL
    TESTB CL, CL
    JNZ  scalar_fp_ok
    MOVB $1, CL
scalar_fp_ok:

    // Calculate i1
    // Recalculate numBuckets mask
    MOVQ numBuckets+56(FP), R13
    DECQ R13
    
    MOVQ R12, R8
    ANDQ R13, R8            // i1

    // Calculate i2: hash the fingerprint
    MOVBQZX CL, R9                  // R9 = fp
    MOVQ    local_prime64_5<>(SB), R11   // R11 = prime64_5
    ADDQ    $1, R11                // R11 = prime64_5 + 1 (initial hash for length=1)
    IMULQ   local_prime64_5<>(SB), R9      // R9 = fp * prime64_5
    XORQ    R9, R11                // R11 ^= (fp * prime64_5)
    // Rotate left by 11 and multiply by prime64_1
    ROLQ    $11, R11               // R11 = rotl64(R11, 11)
    IMULQ   local_prime64_1<>(SB), R11     // R11 *= prime64_1
    // Avalanche
    MOVQ    R11, R9
    SHRQ    $33, R9
    XORQ    R9, R11
    IMULQ   local_prime64_2<>(SB), R11
    MOVQ    R11, R9
    SHRQ    $29, R9
    XORQ    R9, R11
    IMULQ   local_prime64_3<>(SB), R11
    MOVQ    R11, R9
    SHRQ    $32, R9
    XORQ    R9, R11
    MOVQ    R11, R9                // R9 = final hash
    // Now R9 = hash(fp)
    XORQ R8, R9
    ANDQ R13, R9            // i2

    // Store result
    MOVQ results_base+24(FP), R10
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
