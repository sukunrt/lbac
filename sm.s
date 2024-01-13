	.section	__TEXT,__text,regular,pure_instructions
	.globl _eval
	.p2align	2
_eval:
	.cfi_startproc
	mov	x0, #32
	str	x0, [sp, -16]!
	mov	x1, #31
	ldr	x0, [sp], 16
	mul	x0, x1, x0
	str	x0, [sp, -16]!
	mov	x0, #38
	str	x0, [sp, -16]!
	mov	x1, #39
	ldr	x0, [sp], 16
	sdiv	x0, x0, x1
	str	x0, [sp, -16]!
	ldr	x1, [sp], 16
	ldr	x0, [sp], 16
	add	x0, x1, x0
	str	x0, [sp, -16]!
	ldr	x0, [sp], 16
	ret
	.cfi_endproc
.subsections_via_symbols
