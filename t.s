	.text
	.file	"sm.c"
	.globl	eval                            # -- Begin function eval
	.p2align	4, 0x90
	.type	eval,@function
eval:
pushq $2
pushq $3
	popq	%rdi
	movq	$-1, %rax
	imulq	%rax, %rdi
	pushq	%rdi
	popq	%rdi
	popq	%rax
	addq 	%rdi, %rax
	pushq %rax
	popq	%rdi
	movq	$-1, %rax
	imulq	%rax, %rdi
	pushq	%rdi
	popq %rax
	retq
.Lfunc_end0:
	.size	eval, .Lfunc_end0-eval
					# -- End function
	.ident	"clang version 16.0.6"
	.section	".note.GNU-stack","",@progbits
	.addrsig
	
