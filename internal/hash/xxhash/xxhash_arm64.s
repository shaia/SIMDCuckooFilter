//go:build arm64
// +build arm64

#include "textflag.h"

// Constants for XXHash
DATA prime64_1<>+0(SB)/8, $11400714785074694791
DATA prime64_2<>+0(SB)/8, $14029467366897019727
DATA prime64_3<>+0(SB)/8, $1609587929392839161
DATA prime64_4<>+0(SB)/8, $9650029242287828579
DATA prime64_5<>+0(SB)/8, $2870177450012600261
GLOBL prime64_1<>(SB), RODATA, $8
GLOBL prime64_2<>(SB), RODATA, $8
GLOBL prime64_3<>(SB), RODATA, $8
GLOBL prime64_4<>(SB), RODATA, $8
GLOBL prime64_5<>(SB), RODATA, $8

// hash64XXHashInternal computes XXHash64 for a single item
// func hash64XXHashInternal(data []byte) uint64
TEXT ·hash64XXHashInternal(SB), NOSPLIT, $0-32
    // Load arguments
    MOVD data_base+0(FP), R0   // R0 = data pointer
    MOVD data_len+8(FP), R1    // R1 = data length

    // Load prime64_5
    MOVD $prime64_5<>(SB), R2
    MOVD (R2), R3              // R3 = prime64_5

    // Initialize hash = prime64_5 + len
    ADD  R1, R3, R4            // R4 = hash = prime64_5 + len

    // Check if length >= 8
    CMP  $8, R1
    BLT  short_data

    // Load constants
    MOVD $prime64_1<>(SB), R10
    MOVD (R10), R10            // R10 = prime64_1
    MOVD $prime64_2<>(SB), R11
    MOVD (R11), R11            // R11 = prime64_2
    MOVD $prime64_4<>(SB), R12
    MOVD (R12), R12            // R12 = prime64_4

    // Process 8-byte chunks
chunk_loop:
    CMP  $8, R1
    BLT  final_bytes

    // Load 8 bytes (little-endian)
    MOVD (R0), R5

    // k *= prime64_2
    MUL  R11, R5, R6

    // k = rotl64(k, 31)
    ROR  $33, R6, R6           // ARM64: ROR by 33 = ROL by 31

    // k *= prime64_1
    MUL  R10, R6

    // hash ^= k
    EOR  R6, R4

    // hash = rotl64(hash, 27)
    ROR  $37, R4, R4           // ARM64: ROR by 37 = ROL by 27

    // hash *= prime64_1
    MUL  R10, R4

    // hash += prime64_4
    ADD  R12, R4

    // Advance pointer and decrement length
    ADD  $8, R0
    SUB  $8, R1
    B    chunk_loop

final_bytes:
    // Check if there are remaining bytes
    CBZ  R1, hash_finalize

    // Load prime64_5 and prime64_1 for byte processing
    MOVD $prime64_5<>(SB), R13
    MOVD (R13), R13            // R13 = prime64_5

final_byte_loop:
    // Load one byte
    MOVBU (R0), R5

    // hash ^= byte * prime64_5
    MUL  R13, R5
    EOR  R5, R4

    // hash = rotl64(hash, 11)
    ROR  $53, R4, R4           // ARM64: ROR by 53 = ROL by 11

    // hash *= prime64_1
    MUL  R10, R4

    // Advance pointer and decrement length
    ADD  $1, R0
    SUB  $1, R1
    CBNZ R1, final_byte_loop

    B    hash_finalize

short_data:
    // For data < 8 bytes, process byte by byte
    CBZ  R1, hash_finalize

    // Load constants
    MOVD $prime64_1<>(SB), R10
    MOVD (R10), R10
    MOVD $prime64_5<>(SB), R13
    MOVD (R13), R13

short_byte_loop:
    MOVBU (R0), R5
    MUL  R13, R5
    EOR  R5, R4
    ROR  $53, R4, R4
    MUL  R10, R4
    ADD  $1, R0
    SUB  $1, R1
    CBNZ R1, short_byte_loop

hash_finalize:
    // Avalanche mixing
    // Load remaining constants
    MOVD $prime64_2<>(SB), R11
    MOVD (R11), R11
    MOVD $prime64_3<>(SB), R12
    MOVD (R12), R12

    // hash ^= hash >> 33
    LSR  $33, R4, R5
    EOR  R5, R4

    // hash *= prime64_2
    MUL  R11, R4

    // hash ^= hash >> 29
    LSR  $29, R4, R5
    EOR  R5, R4

    // hash *= prime64_3
    MUL  R12, R4

    // hash ^= hash >> 32
    LSR  $32, R4, R5
    EOR  R5, R4

    // Return hash in R4
    MOVD R4, ret+24(FP)
    RET
