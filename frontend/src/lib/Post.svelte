<script>
	import Comment from './Comment.svelte'

	import { slide, fly } from 'svelte/transition'
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
			Posted by <span class="font-bold text-violet-800"
				>{post.created_by.first_name}
				{post.created_by.last_name}</span
			>
			on {formatTime(post.created_at)}
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
		<h2 class="cursor-pointer" on:click={toggleComments}>{post?.comments ? post.comments.length : ''} comments</h2>

		{#if showComments}
			<div transition:slide={{ axis: 'y' }}>
				<input type="text" placeholder="Type here" class="input w-full max-w-xs" />
				<button class="btn">Submit</button>
				{#if post?.comments}
					{#each post.comments as comment, index (comment.ID)}
						<Comment {comment} />
					{/each}
				{:else}
					<h2>Be the first one to comment!</h2>
				{/if}
			</div>
		{/if}
	</div>
</div>
