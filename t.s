	.text
	.file	"sm.c"
	.globl	eval                            # -- Begin function eval
	.p2align	4, 0x90
	.type	eval,@function
eval:
	pushq	%rbp
	movq	%rsp, %rbp
	
	pushq	$5
	pushq	$5
	pushq	-8(%rbp)
	pushq	-16(%rbp)
	popq	%rdi
	popq	%rax
	subq	%rdi, %rax
	pushq	%rax
	popq	%rax
	cmpq	$0, %rax
	je	JL0
	pushq	$2
	pushq	$2
	popq	%rdi
	popq	%rax
	addq 	%rdi, %rax
	pushq	%rax
	popq	%rax
	movq	%rbp, %rdi
	movq	$16, %rbx
	subq	%rbx, %rdi
	movq	%rdi, %rsp
	pushq	%rax
	jmp	JL1
JL0:
	pushq	$2
	pushq	$3
	popq	%rdi
	popq	%rax
	addq 	%rdi, %rax
	pushq	%rax
	popq	%rax
	movq	%rbp, %rdi
	movq	$24, %rbx
	subq	%rbx, %rdi
	movq	%rdi, %rsp
	pushq	%rax
JL1:
	popq	%rax
	movq	%rbp, %rdi
	movq	$0, %rbx
	subq	%rbx, %rdi
	movq	%rdi, %rsp
	pushq	%rax
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
	
