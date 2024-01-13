	.text
	.file	"sm.c"
	.globl	eval                            // -- Begin function eval
	.p2align	2
	.type	eval,@function
eval:                                   // @eval
	.cfi_startproc
// %bb.0:
	mov	w0, #42
	ret
.Lfunc_end0:
	.size	eval, .Lfunc_end0-eval
	.cfi_endproc
                                        // -- End function
	.ident	"Apple clang version 14.0.0 (clang-1400.0.29.202)"
	.section	".note.GNU-stack","",@progbits
	.addrsig
