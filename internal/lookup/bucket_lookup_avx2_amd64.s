// +build amd64

#include "textflag.h"

// func bucketLookupAVX2(fingerprints []byte, target byte) bool
TEXT ·bucketLookupAVX2(SB), NOSPLIT, $0-40
    MOVQ    fingerprints_base+0(FP), AX  // AX = pointer to fingerprints
    MOVQ    fingerprints_len+8(FP), CX   // CX = length
    MOVBQZX target+24(FP), DX            // DX = target fingerprint

    // Broadcast target to all bytes of YMM0
    MOVD    DX, X0
    VPBROADCASTB X0, Y0

    // Process 32 bytes at a time
    XORQ    BX, BX  // BX = offset

loop:
    CMPQ    CX, $32
    JL      remainder

    // Load 32 fingerprints
    VMOVDQU (AX)(BX*1), Y1

    // Compare with target
    VPCMPEQB Y0, Y1, Y2

    // Check if any match
    VPMOVMSKB Y2, DI
    TESTL   DI, DI
    JNZ     found

    ADDQ    $32, BX
    SUBQ    $32, CX
    JMP     loop

remainder:
    // Handle remaining bytes
    TESTQ   CX, CX
    JZ      notfound

rem_loop:
    MOVBQZX (AX)(BX*1), DI
    CMPB    DI, DX
    JE      found
    INCQ    BX
    DECQ    CX
    JNZ     rem_loop

notfound:
    VZEROUPPER
    MOVB    $0, ret+32(FP)
    RET

found:
    VZEROUPPER
    MOVB    $1, ret+32(FP)
    RET
