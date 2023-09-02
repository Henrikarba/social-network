<script>
	// Svelte
	import { slide } from 'svelte/transition'
	import { quintIn, quintInOut } from 'svelte/easing'
	import { flip } from 'svelte/animate'
	import { createEventDispatcher } from 'svelte'
	import { onMount } from 'svelte'
	const dispatch = createEventDispatcher()

	// Icons
	import FaRegQuestionCircle from 'svelte-icons/fa/FaRegQuestionCircle.svelte'

	// Emoji
	import Emoji from './Emoji.svelte'

	let selectedInput
	function handleCurrentInput(event) {
		selectedInput = event.target
	}

	// Stores
	import { currentUserGroups, currentUserFollowers, currentUser } from '../stores/user'
	import { postsStore, groupPostsStore } from '../stores/post'

	// Utils
	import { isValidFileType } from '../utils'

	let post_target = 'regular_post'
	let privacy = 'public'

	//
	let disabled = true
	let title
	let content = ''
	let files
	let selectedGroup
	let buttonToolTipError

	onMount(() => title.focus())

	function handlePostFollowers(followerId) {
		const followerIndex = localFollowers.findIndex((f) => f.id === followerId)
		if (followerIndex > -1) {
			const [follower] = localFollowers.splice(followerIndex, 1)
			localFollowers = [...localFollowers]
			postFollowers = [...postFollowers, follower]
		} else {
			const postFollowerIndex = postFollowers.findIndex((f) => f.id === followerId)
			if (postFollowerIndex > -1) {
				const [follower] = postFollowers.splice(postFollowerIndex, 1)
				postFollowers = [...postFollowers]
				localFollowers = [...localFollowers, follower]
			}
		}
	}

	async function submitForm() {
		if (!disabled) {
			const formData = new FormData()
			formData.append('title', title.value)
			formData.append('content', content)
			formData.append('privacy', privacy)
			formData.append('post_target', post_target)
			if (files && files[0] && isValidFileType(files[0].type)) {
				formData.append('image', files[0])
			}
			if (privacy === '2' && postFollowers.length > 0) {
				const followerIds = postFollowers.map((follower) => follower.id)
				formData.append('followers', JSON.stringify(followerIds))
			}
			if (post_target == 'group_post') {
				formData.append('groupid', selectedGroup)
			}

			const result = await handleNewPost(formData)
			result.post.created_by = $currentUser
			result.post.created_at = new Date()
			result.post.group = $currentUserGroups.find((item) => item.id == parseInt(selectedGroup))
			post_target == 'regular_post'
				? ($postsStore = [result.post, ...$postsStore])
				: ($groupPostsStore = [result.post, ...$groupPostsStore])
		}

		dispatch(post_target)
	}

	async function handleNewPost(formData) {
		try {
			const response = await fetch(`http://localhost:80/new/post`, {
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

	$: followers = $currentUserFollowers ? $currentUserFollowers.filter((item) => item.status == 'accepted') : []
	$: localFollowers = followers

	let postFollowers = []
	// Validate
	$: {
		const hasTitleAndContent = title && content
		const hasPostFollowers = postFollowers.length > 0
		buttonToolTipError = ''
		if (privacy === '2') {
			disabled = !hasTitleAndContent || !hasPostFollowers
			buttonToolTipError = !hasPostFollowers
				? 'Add some followers to your post or choose another privacy setting!'
				: 'Add some content to your post.'
		} else {
			disabled = !hasTitleAndContent
			buttonToolTipError = 'Add some content to your post.'
		}
		if (post_target == 'group_post' && isNaN(selectedGroup)) {
			disabled = !hasTitleAndContent
			buttonToolTipError = 'Pick a group or make regular post!'
		}
		if (files && files[0]) {
			if (!isValidFileType(files[0].type)) {
				disabled = true
				buttonToolTipError = 'Supported formats: .png, .jpg, .jpeg, .gif'
			}
		}
	}
</script>

<main
	data-theme="dracula"
	class="bg-base-100 mt-4 rounded-3xl h-[100vh] p-10 flex gap-4"
	in:slide|global={{ delay: 500, duration: 200, axis: 'y' }}
	out:slide|global={{ duration: 200, axis: 'x' }}
>
	<div class="w-1/3 flex flex-col gap-4">
		<div>
			<div class="text-accent text-4xl">General/Group post</div>
			<label class="text-info" for="type">General</label>
			<input bind:group={post_target} type="radio" name="type" class="radio radio-primary" value="regular_post" />
			{#if $currentUserGroups && $currentUserGroups.length > 0}
				<label class="text-info ml-4" for="type">Group</label>
				<input bind:group={post_target} type="radio" name="type" class="radio radio-primary" value="group_post" />
			{/if}
		</div>
		<!--  -->
		{#if post_target == 'regular_post'}
			<div
				in:slide|global={{ delay: 300, duration: 300, easing: quintIn }}
				out:slide|global={{ delay: 100, duration: 300, easing: quintInOut }}
			>
				<div class="flex">
					<div class="text-accent text-4xl">Privacy setting</div>
					<div
						class="tooltip"
						data-tip={followers.length > 0
							? 'Everyone can see public posts. Private posts can be seen only by your followers. Privatus Maximus lets you pick specific followers who can see your post'
							: 'You can only make public posts as you have no followers'}
					>
						<div class="w-8 text-warning mt-2 ml-2"><FaRegQuestionCircle /></div>
					</div>
				</div>

				<label class="text-info" for="privacy">Public</label>
				<input bind:group={privacy} type="radio" name="privacy" class="radio radio-error" value="public" />
				{#if followers && followers.length > 0}
					<label class="text-info ml-4" for="privacy">Private</label>
					<input bind:group={privacy} type="radio" name="privacy" class="radio radio-error" value="private" />
					<label class="text-info ml-4" for="privacy">Privatus Maximus</label>
					<input bind:group={privacy} type="radio" name="privacy" class="radio radio-error" value="followers_only" />
				{/if}
			</div>
			{#if privacy == 'followers_only'}
				<div
					in:slide|global={{ delay: 300, duration: 300, easing: quintIn }}
					out:slide|global={{ delay: 100, duration: 300, easing: quintInOut }}
				>
					<div class="rounded w-full flex items-center flex-wrap gap-2">
						<ul>
							{#each [localFollowers, postFollowers] as follower, index}
								<h2 class="text-info text-4xl mt-4">{index === 0 ? 'Your followers' : 'Post followers'}</h2>

								<p class="text-accent">
									{index === 0 ? 'Click on follower to add them to post followers' : 'Click to remove'}
								</p>

								<ul class="h-10 flex gap-2">
									{#each follower as foll, index (foll.id)}
										<!-- svelte-ignore a11y-click-events-have-key-events -->
										<!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
										<li
											class="badge py-6 mt-2 cursor-pointer"
											animate:flip={{ duration: 500 }}
											on:click={() => handlePostFollowers(foll.id)}
										>
											{foll.first_name}
											{foll.last_name}
										</li>
									{/each}
								</ul>
							{/each}
						</ul>
					</div>
				</div>
			{/if}
		{:else if post_target == 'group_post'}
			<div
				in:slide|global={{ delay: 300, duration: 300, easing: quintIn }}
				out:slide|global={{ delay: 100, duration: 300, easing: quintInOut }}
			>
				<div class="text-accent text-4xl">Select group</div>
				<select class="select select-info w-full max-w-xs mt-4 text-white" bind:value={selectedGroup}>
					{#each $currentUserGroups as group}
						<option value={group.id}>{group.title}</option>
					{/each}
				</select>
			</div>
		{/if}
		<div class="w-full h-96 mt-4">
			<Emoji input={selectedInput} />
		</div>
	</div>
	<div class="flex flex-col gap-2 w-2/3">
		<label for="title" class="text-xl text-info uppercase">Title</label>
		<input
			on:focus={handleCurrentInput}
			type="text"
			bind:this={title}
			placeholder="Title of your post"
			class="input w-full max-w focus:border-accent focus:outline-none border-2 border-primary text-accent"
		/>
		<label for="content" class="text-xl text-info uppercase">Content</label>
		<textarea
			on:focus={handleCurrentInput}
			bind:value={content}
			placeholder="Content of your post"
			class="textarea textarea-bordered textarea-lg w-full resize-none h-60 focus:border-accent focus:outline-none border-2 border-primary text-accent"
		/>
		<div class="flex justify-between">
			<input
				bind:files
				type="file"
				class="file-input file-input-accent w-full max-w-xs text-primary"
				accept=".png, .jpg, .jpeg, .gif"
			/>
			<div class={disabled ? 'tooltip tooltip-hover tooltip-left' : ''} data-tip={buttonToolTipError}>
				<button on:click={submitForm} class="btn btn-success font-bold" {disabled}>Submit</button>
			</div>
		</div>
	</div>
</main>
