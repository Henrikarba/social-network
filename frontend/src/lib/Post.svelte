<script>
	// Svelte
	import { createEventDispatcher } from 'svelte'
	const dispatch = createEventDispatcher()
	// Icons
	import FaRegArrowAltCircleDown from 'svelte-icons/fa/FaRegArrowAltCircleDown.svelte'
	import FaArrowUp from 'svelte-icons/fa/FaArrowUp.svelte'

	import Comment from './Comment.svelte'
	import Emoji from './Emoji.svelte'
	let selectedInput

	import { slide, fly } from 'svelte/transition'
	import { circIn } from 'svelte/easing'
	import { formatTime, isValidFileType } from '../utils'
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
	let files
	async function handleNewComment(formData) {
		try {
			const response = await fetch(`http://localhost:80/new/comment`, {
				method: 'POST',
				body: formData,
				credentials: 'include',
			})
			if (!response.ok) {
				const errorMessage = await response.text()
				throw new Error(`Request failed: ${errorMessage}`)
			}
			return await response.json()
		} catch (error) {
			throw error
		}
	}

	async function submitComment() {
		if (commentContent.trim() == '') {
			error = true
		}
		let post_target = 'regular_post'
		if (post?.group) {
			post_target = 'group_post'
		}

		const formData = new FormData()

		formData.append('content', commentContent)
		formData.append('post_target', post_target)
		formData.append('post_id', post.post_id)
		if (post_target == 'group_post') {
			formData.append('group_id', post.group.id)
		}
		if (files && files[0] && isValidFileType(files[0].type)) {
			formData.append('image', files[0])
		}

		const result = await handleNewComment(formData)
		if (!post?.comments) post.comments = []
		post.comments = [...post.comments, result.comment]
		commentContent = ''
	}
</script>

<div
	in:slide|global={{ delay: 500, duration: 500, axis: 'y' }}
	out:fly|global={{ duration: 500, x: 500 }}
	class="flex flex-col border-2 border-black mt-4 rounded p-4 items-center bg-slate-300"
>
	<article class="prose">
		<h2 class="font-bold text-4xl text-center">{post.title}</h2>
		<p class="text-center">
			<!-- svelte-ignore a11y-click-events-have-key-events -->
			Posted by
			<!-- svelte-ignore a11y-no-static-element-interactions -->
			<span on:click={() => dispatch('user', post.created_by.id)} class="cursor-pointer font-bold text-violet-800"
				>{post.created_by.first_name}
				{post.created_by.last_name}</span
			>
			on {formatTime(post.created_at)}
			{#if post.group}
				<!-- svelte-ignore a11y-no-static-element-interactions -->
				in
				<!-- svelte-ignore a11y-click-events-have-key-events -->
				<span on:click={() => dispatch('group', post.group.id)} class="cursor-pointer text-orange-600 font-extrabold">
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
			<!-- svelte-ignore a11y-click-events-have-key-events -->
			<!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
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
			<div transition:slide={{ axis: 'y' }} class="flex">
				<div>
					<div class="flex flex-col gap-1">
						<input
							bind:this={selectedInput}
							type="text"
							placeholder={error ? 'try to do better' : 'Type here'}
							class="input w-96 {error ? 'border-2 border-red-700 w-[90%] ml-6' : ''}"
							bind:value={commentContent}
						/>
						<input
							bind:files
							type="file"
							class="file-input file-input-accent w-full max-w-xs text-primary"
							accept=".png, .jpg, .jpeg, .gif"
						/>
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
				<div class="w-full h-72 mt-4">
					<Emoji input={selectedInput} />
				</div>
			</div>
		{/if}
	</div>
</div>
