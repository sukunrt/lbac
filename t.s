	.text
	.file	"sm.c"
	.globl	eval                            # -- Begin function eval
	.p2align	4, 0x90
	.type	eval,@function
eval:
	pushq	%rbp
	movq	%rsp, %rbp
	
	pushq	$100
	pushq	$15
JL0:
	pushq	-8(%rbp)
	pushq	-16(%rbp)
	popq	%rdi
	popq	%rax
	subq	%rdi, %rax
	pushq	%rax
	popq	%rax
	cmpq	$0, %rax
	je	JL1
	pushq	-8(%rbp)
	pushq	$2
	popq	%rdi
	popq	%rax
	subq	%rdi, %rax
	pushq	%rax
	popq	%rax
	movq	%rax, -8(%rbp)
	jmp	JL0
JL1:
	pushq	-8(%rbp)
	popq	%rax
	movq %rbp, %rsp
	popq %rbp
	retq
.Lfunc_end0:
	.size	eval, .Lfunc_end0-eval
					# -- End function
	.ident	"clang version 16.0.6"
	.section	".note.GNU-stack","",@progbits
	.addrsig
	
