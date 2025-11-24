//go:build amd64
// +build amd64

#include "textflag.h"

// Constants for XXHash - stored as data
// Note: These constants are shared with batch_avx2_amd64.s, so they cannot be file-private
DATA prime64_1+0(SB)/8, $11400714785074694791
DATA prime64_2+0(SB)/8, $14029467366897019727
DATA prime64_3+0(SB)/8, $1609587929392839161
DATA prime64_4+0(SB)/8, $9650029242287828579
DATA prime64_5+0(SB)/8, $2870177450012600261
GLOBL prime64_1(SB), RODATA|NOPTR, $8
GLOBL prime64_2(SB), RODATA|NOPTR, $8
GLOBL prime64_3(SB), RODATA|NOPTR, $8
GLOBL prime64_4(SB), RODATA|NOPTR, $8
GLOBL prime64_5(SB), RODATA|NOPTR, $8

// hash64XXHashInternal computes XXHash64 for a single item
// func hash64XXHashInternal(data []byte) uint64
TEXT ·hash64XXHashInternal(SB), NOSPLIT, $0-32
    // Load arguments
    MOVQ data_base+0(FP), BX   // BX = data pointer
    MOVQ data_len+8(FP), CX    // CX = data length
    
    // Initialize hash = prime64_5 + len
    MOVQ prime64_5(SB), AX   // AX = hash = prime64_5
    ADDQ CX, AX                // hash += len
    
    // Check if length >= 8
    CMPQ CX, $8
    JL   short_data
    
    // Process 8-byte chunks
chunk_loop:
    CMPQ CX, $8
    JL   final_bytes
    
    // Load 8 bytes (little-endian) - MOVQ handles this
    MOVQ (BX), DX
    
    // Process block: k = DX
    // k *= prime64_2
    MOVQ DX, SI
    IMULQ prime64_2(SB), SI
    // k = rotl64(k, 31)
    MOVQ SI, DX
    SHLQ $31, DX
    SHRQ $33, SI
    ORQ  DX, SI
    // k *= prime64_1
    IMULQ prime64_1(SB), SI
    // hash ^= k
    XORQ SI, AX
    // hash = rotl64(hash, 27) * prime64_1 + prime64_4
    MOVQ AX, SI
    SHLQ $27, SI
    SHRQ $37, AX
    ORQ  SI, AX
    IMULQ prime64_1(SB), AX
    MOVQ prime64_4(SB), SI
    ADDQ SI, AX
    
    ADDQ $8, BX
    SUBQ $8, CX
    JMP  chunk_loop
    
final_bytes:
    // Process remaining bytes one by one
    TESTQ CX, CX
    JZ    hash_finalize
    
final_byte_loop:
    MOVBQZX (BX), SI
    IMULQ prime64_5(SB), SI  // Multiply byte by prime64_5 first
    XORQ SI, AX                // Then XOR with hash
    // Rotate left by 11
    MOVQ AX, SI
    SHLQ $11, SI
    SHRQ $53, AX
    ORQ  SI, AX
    IMULQ prime64_1(SB), AX
    INCQ BX
    DECQ CX
    JNZ  final_byte_loop
    
    JMP  hash_finalize
    
short_data:
    // For data < 8 bytes, just process final bytes
    TESTQ CX, CX
    JZ    hash_finalize
    
short_byte_loop:
    MOVBQZX (BX), SI
    IMULQ prime64_5(SB), SI  // Multiply byte by prime64_5 first
    XORQ SI, AX                // Then XOR with hash
    MOVQ AX, SI
    SHLQ $11, SI
    SHRQ $53, AX
    ORQ  SI, AX
    IMULQ prime64_1(SB), AX
    INCQ BX
    DECQ CX
    JNZ  short_byte_loop
    
hash_finalize:
    // Finalize: xxhash_mix64
    // hash ^= hash >> 33
    MOVQ AX, SI
    SHRQ $33, SI
    XORQ SI, AX
    // hash *= prime64_2
    MOVQ prime64_2(SB), SI
    IMULQ SI, AX
    // hash ^= hash >> 29
    MOVQ AX, SI
    SHRQ $29, SI
    XORQ SI, AX
    // hash *= prime64_3
    MOVQ prime64_3(SB), SI
    IMULQ SI, AX
    // hash ^= hash >> 32
    MOVQ AX, SI
    SHRQ $32, SI
    XORQ SI, AX
    
    // Return hash in AX
    MOVQ AX, ret+24(FP)
    RET

