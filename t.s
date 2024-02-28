	.text
	.file	"sm.c"
	.globl	eval                            # -- Begin function eval
	.p2align	4, 0x90
	.type	eval,@function
eval:
	pushq	%rbp
	movq	%rsp, %rbp
	
	jmp	endfib
fib:
	pushq	%rbp
	movq	%rsp, %rbp
	pushq	$0
	pushq	16(%rbp)
	pushq	$2
	popq	%rdi
	popq	%rax
	pushq	$0
	cmpq	%rdi, %rax
	jge	JL0
	popq	%rax
	pushq	$1
JL0:
	popq	%rax
	cmpq	$0, %rax
	je	JL1
	pushq	16(%rbp)
	popq	%rax
	movq	%rax, -8(%rbp)
	jmp	JL2
JL1:
	pushq	16(%rbp)
	pushq	$1
	popq	%rdi
	popq	%rax
	subq	%rdi, %rax
	pushq	%rax
	pushq	16(%rbp)
	pushq	$2
	popq	%rdi
	popq	%rax
	subq	%rdi, %rax
	pushq	%rax
	pushq	-16(%rbp)
	callq	fib
	popq	%rdx
	pushq	%rax
	pushq	-24(%rbp)
	callq	fib
	popq	%rdx
	pushq	%rax
	pushq	-32(%rbp)
	pushq	-40(%rbp)
	popq	%rdi
	popq	%rax
	addq	%rdi, %rax
	pushq	%rax
	popq	%rax
	movq	%rax, -8(%rbp)
JL2:
	pushq	-8(%rbp)
	popq	%rax
	movq	%rbp, %rsp
	popq	%rbp
	retq
endfib:
	pushq	$6
	callq	fib
	popq	%rdx
	pushq	%rax
	pushq	-8(%rbp)
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
	
