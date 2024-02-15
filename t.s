	.text
	.file	"sm.c"
	.globl	eval                            # -- Begin function eval
	.p2align	4, 0x90
	.type	eval,@function
eval:
pushq $1
pushq $1
	popq	%rdi
	popq	%rax
	addq 	%rdi, %rax
	pushq %rax
pushq $1
pushq $10
	popq	%rdi
	popq	%rax
	cqto
	imulq	%rdi, %rax
	pushq	%rax
	popq	%rdi
	popq	%rax
	cqto
	imulq	%rdi, %rax
	pushq	%rax
	popq %rax
	retq
.Lfunc_end0:
	.size	eval, .Lfunc_end0-eval
					# -- End function
	.ident	"clang version 16.0.6"
	.section	".note.GNU-stack","",@progbits
	.addrsig
	
