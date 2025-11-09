//go:build amd64
// +build amd64

#include "textflag.h"

// func batchCRC32SIMD(items [][]byte, results []uint32)
// Processes CRC32C checksums using SSE4.2 CRC32C instruction
// items: ptr(0), len(8), cap(16) = 24 bytes
// results: ptr(24), len(32), cap(40) = 12 bytes
// Total frame size: 24 + 12 = 36 bytes (round up to 40 for alignment)
TEXT ·batchCRC32SIMD(SB), NOSPLIT, $0-48
    MOVQ items_base+0(FP), SI      // SI = items slice base
    MOVQ items_len+8(FP), CX       // CX = items length
    MOVQ results_base+24(FP), DI   // DI = results slice base

    // Process all items sequentially
    TESTQ CX, CX
    JZ done

process_loop:
    // Load item ptr and len
    MOVQ 0(SI), AX     // item.ptr
    MOVQ 8(SI), BX     // item.len

    // Initialize CRC32 to 0xFFFFFFFF
    MOVL $0xFFFFFFFF, R8

    // Check if length is 0
    TESTQ BX, BX
    JZ empty_item

    // Save registers we'll modify
    MOVQ AX, R9        // data ptr
    MOVQ BX, R10       // data len

    // Process 8 bytes at a time
loop8:
    CMPQ R10, $8
    JL loop1

    MOVQ 0(R9), R11
    CRC32Q R11, R8
    ADDQ $8, R9
    SUBQ $8, R10
    JMP loop8

loop1:
    TESTQ R10, R10
    JZ finalize

    MOVBQZX 0(R9), R11
    CRC32B R11, R8
    INCQ R9
    DECQ R10
    JMP loop1

empty_item:
finalize:
    // Finalize CRC32 (invert)
    NOTL R8

    // Store result
    MOVL R8, 0(DI)

    // Move to next item
    ADDQ $24, SI       // sizeof([]byte) = 24
    ADDQ $4, DI        // sizeof(uint32) = 4
    DECQ CX
    JNZ process_loop

done:
    RET
