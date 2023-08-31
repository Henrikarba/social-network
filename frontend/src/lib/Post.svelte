<script>
	// Svelte
	import { createEventDispatcher } from 'svelte'
	const dispatch = createEventDispatcher()
	// Icons
	import FaRegArrowAltCircleDown from 'svelte-icons/fa/FaRegArrowAltCircleDown.svelte'
	import FaArrowUp from 'svelte-icons/fa/FaArrowUp.svelte'

	import Comment from './Comment.svelte'

	import { slide, fly } from 'svelte/transition'
	import { circIn } from 'svelte/easing'
	import { formatTime } from '../utils'
	export let post
	let content = post.content

	function addBlockquotes(content) {
		return content.replace(/"([^"]*)"/g, '<blockquote>$1</blockquote>')
	}

	let contentWithBlockquotes = addBlockquotes(content)

	let showComments = false
	function toggleComments() {
		showComments = !showComments
	}

	let error = false
	let commentContent = ''
	function submitComment() {
		if (commentContent.trim() == '') {
			error = true
		}
	}
	$: console.log(post)
</script>

<div
	in:slide|global={{ delay: 500, duration: 500, axis: 'y' }}
	out:fly|global={{ duration: 500, x: 500 }}
	class="flex flex-col border-2 border-black mt-4 rounded p-4 items-center"
>
	<article class="prose">
		<h2 class="font-bold text-4xl text-center">{post.title}</h2>
		<p class="text-center">
			Posted by <span
				on:click={() => dispatch('user', post.created_by.id)}
				class="cursor-pointer font-bold text-violet-800"
				>{post.created_by.first_name}
				{post.created_by.last_name}</span
			>
			on {formatTime(post.created_at)}
			{#if post.group}
				in <span
					on:click={() => dispatch('group', post.group.id)}
					class="cursor-pointer text-orange-600 font-extrabold"
				>
					{post.group.title}
				</span>
			{/if}
		</p>

		<div class="container flex flex-col">
			<p>{@html contentWithBlockquotes}</p>
			{#if post.image_url != '' && post.image_url}
				<div class="flex justify-center w-vw items-center">
					<img class="w-96 rounded" src="http://localhost:80/images/{post.image_url}" alt="" />
				</div>
			{/if}
		</div>
	</article>
	<div class=" text-zinc-800 font-bold w-full text-center">
		<div class="flex">
			<h2 class="cursor-pointer" on:click={toggleComments}>{post?.comments ? post.comments.length : ''} comments</h2>
			<div class="w-6">
				{#if showComments}
					<FaRegArrowAltCircleDown />
				{:else}
					<FaArrowUp />
				{/if}
			</div>
		</div>
		{#if showComments}
			<div transition:slide={{ axis: 'y' }}>
				<div class="flex flex-col gap-1">
					<input
						type="text"
						placeholder={error ? 'try to do better' : 'Type here'}
						class="input w-96 {error ? 'border-2 border-red-700 w-[90%] ml-6' : ''}"
						bind:value={commentContent}
					/>
					<input type="file" class="file-input w-full max-w-xs" />
					<button on:click={submitComment} class="btn w-20">Submit</button>
				</div>
				{#if post?.comments}
					{#each post.comments as comment, index (comment.ID)}
						<Comment {comment} />
					{/each}
				{:else}
					<h2>Be the first one to comment!</h2>
				{/if}
				{#if error}
					<div transition:slide={{ duration: 1000, axis: 'x', easing: circIn }} class="alert alert-error w-96">
						<svg
							xmlns="http://www.w3.org/2000/svg"
							class="stroke-current shrink-0 h-6 w-6"
							fill="none"
							viewBox="0 0 24 24"
							><path
								stroke-linecap="round"
								stroke-linejoin="round"
								stroke-width="2"
								d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
							/></svg
						>
						<span class="text-red-600">Error! Task failed successfully.</span>
					</div>
				{/if}
			</div>
		{/if}
	</div>
</div>