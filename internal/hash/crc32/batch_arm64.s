//go:build arm64
// +build arm64

#include "textflag.h"

// func batchCRC32Hardware(items [][]byte, results []uint32)
// Processes CRC32C checksums using ARMv8 hardware CRC32C instructions
// The ARMv8 CRC32C instructions compute the Castagnoli polynomial,
// which is exactly what we need for hash/crc32.Castagnoli.
//
// items: ptr(0), len(8), cap(16) = 24 bytes
// results: ptr(24), len(32), cap(40) = 24 bytes
// Total frame size: 24 + 24 = 48 bytes
TEXT ·batchCRC32Hardware(SB), NOSPLIT, $0-48
    // Load arguments
    MOVD items_base+0(FP), R6      // R6 = items slice base pointer
    MOVD items_len+8(FP), R7       // R7 = items length (number of items)
    MOVD results_base+24(FP), R8   // R8 = results slice base pointer

    // Check if we have any items
    CBZ R7, done

process_loop:
    // Load current item ([]byte structure: ptr, len, cap)
    MOVD 0(R6), R9     // R9 = item.ptr (data pointer)
    MOVD 8(R6), R10    // R10 = item.len (data length)

    // Initialize CRC32 to 0xFFFFFFFF (standard CRC32 initialization)
    MOVD $0xFFFFFFFF, R11
    MOVW R11, R12      // R12 = CRC accumulator (32-bit)

    // Check if data length is 0
    CBZ R10, empty_item

    // Process 8 bytes at a time using CRC32CX instruction
loop8:
    CMP $8, R10
    BLT loop4

    // Load 8 bytes and process with CRC32CX (CRC32C for 64-bit)
    MOVD (R9), R13
    CRC32CX R13, R12   // R12 = CRC32C(R12, R13) - 64-bit CRC32C
    ADD $8, R9
    SUB $8, R10
    B loop8

loop4:
    // Process 4 bytes at a time using CRC32CW instruction
    CMP $4, R10
    BLT loop2

    // Load 4 bytes and process with CRC32CW (CRC32C for 32-bit)
    MOVW (R9), R13
    CRC32CW R13, R12   // R12 = CRC32C(R12, R13) - 32-bit CRC32C
    ADD $4, R9
    SUB $4, R10
    B loop4

loop2:
    // Process 2 bytes at a time using CRC32CH instruction
    CMP $2, R10
    BLT loop1

    // Load 2 bytes and process with CRC32CH (CRC32C for 16-bit)
    MOVH (R9), R13
    CRC32CH R13, R12   // R12 = CRC32C(R12, R13) - 16-bit CRC32C
    ADD $2, R9
    SUB $2, R10
    B loop2

loop1:
    // Process remaining bytes one at a time using CRC32CB instruction
    CBZ R10, finalize

    // Load 1 byte and process with CRC32CB (CRC32C for 8-bit)
    MOVB (R9), R13
    CRC32CB R13, R12   // R12 = CRC32C(R12, R13) - 8-bit CRC32C
    ADD $1, R9
    SUB $1, R10
    B loop1

empty_item:
finalize:
    // Finalize CRC32: invert all bits (standard CRC32 finalization)
    MVN R12, R12       // R12 = ~R12 (bitwise NOT)

    // Store result (32-bit value)
    MOVW R12, (R8)

    // Move to next item and result
    ADD $24, R6        // sizeof([]byte) = 24 bytes (ptr, len, cap)
    ADD $4, R8         // sizeof(uint32) = 4 bytes
    SUB $1, R7
    CBNZ R7, process_loop

done:
    RET
