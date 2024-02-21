	.text
	.file	"sm.c"
	.globl	eval                            # -- Begin function eval
	.p2align	4, 0x90
	.type	eval,@function
eval:
	pushq	%rbp
	movq	%rsp, %rbp
	
	pushq	$1
	pushq	$1
	popq	%rdi
	popq	%rax
	addq	%rdi, %rax
	pushq	%rax
	pushq	$2
	pushq	$2
	popq	%rdi
	popq	%rax
	imulq	%rdi, %rax
	pushq	%rax
	pushq	$3
	popq	%rdi
	popq	%rdx
	movq	$1, %rbx
	movq	$1, %rax
JL0:
	cmp	$0, %rdi
	je	JL1
	imulq	%rdx, %rax
	subq	%rbx, %rdi
	jmp	JL0
JL1:
	pushq	%rax
	popq	%rdi
	popq	%rax
	imulq	%rdi, %rax
	pushq	%rax
	popq	%rax
	movq	%rbp, %rsp
	popq	%rbp
	retq
.Lfunc_end0:
	.size	eval, .Lfunc_end0-eval
					# -- End function
	.ident	"clang version 16.0.6"
	.section	".note.GNU-stack","",@progbits
	.addrsig
	
