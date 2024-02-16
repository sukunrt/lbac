	.text
	.file	"sm.c"
	.globl	eval                            # -- Begin function eval
	.p2align	4, 0x90
	.type	eval,@function
eval:
	pushq	%rbp
	movq	%rsp, %rbp
	
pushq	$5
pushq	$15
pushq	$2
pushq	$15
pushq	-32(%rbp)
pushq	-16(%rbp)
pushq	$10
popq	%rdi
popq	%rax
	subq	%rdi, %rax
pushq	%rax
popq	%rdi
popq	%rax
	addq 	%rdi, %rax
pushq	%rax
pushq	-24(%rbp)
popq	%rdi
popq	%rax
cqto
	idivq	%rdi
pushq	%rax
pushq	$1
popq	%rdi
popq	%rax
	addq 	%rdi, %rax
pushq	%rax
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
	
