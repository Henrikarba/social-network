<script>
	// Svelte
	import { slide } from 'svelte/transition'
	import { quintIn, quintInOut } from 'svelte/easing'
	import { flip } from 'svelte/animate'

	// Icons
	import FaRegQuestionCircle from 'svelte-icons/fa/FaRegQuestionCircle.svelte'

	// Stores
	import { currentUserGroups, currentUserFollowers } from '../stores/user'

	let post_target = 'regular_post'
	let privacy = 'public'

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

	// $: followers = $currentUserFollowers.filter((item) => item.status == 'accepted')
	let followers = [
		{
			'id': 1,
			'first_name': 'Linus',
			'last_name': 'Torvalds',
			'privacy': 0,
			'status': 'accepted',
		},
		{
			'id': 2,
			'first_name': 'Blue',
			'last_name': 'Torvalds',
			'privacy': 0,
			'status': 'accepted',
		},
		{
			'id': 3,
			'first_name': 'Blue',
			'last_name': 'Torvalds',
			'privacy': 0,
			'status': 'accepted',
		},
		{
			'id': 4,
			'first_name': 'Blue',
			'last_name': 'Torvalds',
			'privacy': 0,
			'status': 'accepted',
		},
		{
			'id': 5,
			'first_name': 'Blue',
			'last_name': 'Greek',
			'privacy': 0,
			'status': 'accepted',
		},
	]

	let localFollowers = followers
	let postFollowers = []
</script>

<main data-theme="dracula" class="bg-base-100 mt-4 rounded h-[100vh] p-10 flex gap-4">
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
						data-tip="Everyone can see public posts. Private posts can be seen only by your followers. Privatus Maximus lets you pick specific followers who can see your post"
					>
						<div class="w-8 text-warning mt-2 ml-2"><FaRegQuestionCircle /></div>
					</div>
				</div>

				<label class="text-info" for="privacy">Public</label>
				<input bind:group={privacy} type="radio" name="privacy" class="radio radio-error" value="public" />
				<label class="text-info ml-4" for="privacy">Private</label>
				<input bind:group={privacy} type="radio" name="privacy" class="radio radio-error" value="private" />
				<label class="text-info ml-4" for="privacy">Privatus Maximus</label>
				<input bind:group={privacy} type="radio" name="privacy" class="radio radio-error" value="followers_only" />
			</div>
			{#if privacy == 'followers_only'}
				<div
					in:slide|global={{ delay: 300, duration: 300, easing: quintIn }}
					out:slide|global={{ delay: 100, duration: 300, easing: quintInOut }}
				>
					<div class=" rounded w-full flex items-center flex-wrap gap-2">
						<ul>
							{#each [localFollowers, postFollowers] as follower, index}
								<h2 class="text-info text-4xl mt-4">{index === 0 ? 'Your followers' : 'Post followers'}</h2>
								{#if index == 0}
									<p class="text-accent">Click on follower to add them to post followers</p>
								{/if}
								<ul class="h-10 flex gap-2">
									{#each follower as foll, index (foll.id)}
										<!-- svelte-ignore a11y-click-events-have-key-events -->
										<!-- svelte-ignore a11y-no-noninteractive-element-interactions -->
										<li
											class="badge mt-2 cursor-pointer"
											animate:flip={{ duration: 200 }}
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
				<select class="select select-info w-full max-w-xs mt-4 text-white">
					{#each $currentUserGroups as group}
						<option value={group.id}>{group.title}</option>
					{/each}
				</select>
			</div>
		{/if}
	</div>
</main>
