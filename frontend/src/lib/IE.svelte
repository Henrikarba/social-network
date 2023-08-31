<script>
	// Svelte
	import { get_current_component } from 'svelte/internal'
	const THISComponent = get_current_component()
	import { createEventDispatcher } from 'svelte'
	const dispatch = createEventDispatcher()
	import { scale } from 'svelte/transition'

	// Components
	import Profile from './Profile.svelte'
	import Post from './Post.svelte'
	import ViewProfile from './ViewProfile.svelte'

	// Icons
	import FaRegWindowClose from 'svelte-icons/fa/FaRegWindowClose.svelte'
	import MdFullscreen from 'svelte-icons/md/MdFullscreen.svelte'
	import MdFullscreenExit from 'svelte-icons/md/MdFullscreenExit.svelte'

	// Stores
	import { postsStore, groupPostsStore } from '../stores/post'
	import { currentUser } from '../stores/user'

	// utils
	import { getProfile } from '../utils'

	export let ieUrl
	export let z

	let left = 400
	let top = 100
	let full = false
	let moving = false

	let route = 'posts'

	function onMouseDown() {
		dispatch('last', 'ie')
		moving = true
	}

	function onMouseMove(e) {
		if (full) return
		if (moving) {
			left += e.movementX
			top += e.movementY
		}
	}

	function onMouseUp() {
		moving = false
	}

	function destroySelf() {
		dispatch('close')
		THISComponent.$destroy()
	}

	let profile
	let id
	async function onClick(event) {
		id = event.detail
		switch (event.type) {
			case 'user':
				if (id == $currentUser.id) {
					route = 'profile'
					return
				}
				profile = await getProfile(id)
				route = 'user/' + id
				break
			case 'group':
				console.log(id)
				break
		}
	}

	$: posts = $postsStore
	$: group_posts = $groupPostsStore
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
<div
	class="{full
		? 'w-full h-screen'
		: 'w-[1280px] h-[720px]'} {z} border-2 rounded absolute border-b-4 border-zinc-500 select-none"
	style={full ? 'left: 0px; top: 0px;' : 'left: ' + left + 'px; top: ' + top + 'px;'}
	in:scale|global={{ duration: 500, start: 0.5 }}
	on:click={() => dispatch('last', 'ie')}
>
	<!-- svelte-ignore a11y-no-static-element-interactions -->
	<div
		on:mousedown={onMouseDown}
		class=" bg-gradient-to-r from-sky-500 to-blue-700 h-10 border-b-2 border-blue-950 flex justify-between items-center draggable"
	>
		<div class=" ml-4 font-bold select-none text-slate-800 text-lg flex items-center gap-2">
			<!-- svelte-ignore a11y-missing-attribute -->
			<img class="h-6 w-6" src={ieUrl} />
			<h2>INTERNET EXPLORER</h2>
		</div>
		<!-- svelte-ignore a11y-click-events-have-key-events -->
		<div class="flex">
			<!-- svelte-ignore a11y-click-events-have-key-events -->
			<div class="h-6 w-6 mr-2 text-black cursor-pointer" on:click={() => (full = !full)}>
				{#if !full}
					<MdFullscreen />
				{:else}
					<MdFullscreenExit />
				{/if}
			</div>
			<div class="h-6 w-6 mr-2 text-black cursor-pointer" on:click={destroySelf}>
				<FaRegWindowClose />
			</div>
		</div>
	</div>
	<div class="h-10 bg-gradient-to-r from-zinc-400 to-zinc-600 border-b-2 border-black flex items-center">
		<div class="w-[90%] mx-auto bg-slate-200 rounded select-none">
			<p class="ml-4">http://localhost:5000/{route}</p>
		</div>
	</div>
	<div class="{'full ? h-[88%] : h-[635px]'} overflow-y-scroll overflow-x-hidden bg-slate-100 p-6">
		<nav class="flex items-center justify-center gap-4 font-bold text-4xl border-b-4 border-slate-900">
			<button on:click={() => (route = 'posts')}>[POSTS]</button>
			<button on:click={() => (route = 'group_posts')}>[GROUP POSTS]</button>
			<button on:click={() => (route = 'groups')}>[GROUPS]</button>
			<button on:click={() => (route = 'profile')}>[PROFILE]</button>
		</nav>
		{#if route == 'posts'}
			{#each posts as post, index (post.post_id)}
				<Post {post} on:user={onClick} />
			{/each}
		{:else if route == 'group_posts'}
			{#each group_posts as gPost, index (gPost.post_id)}
				<Post post={gPost} on:group={onClick} on:user={onClick} />
			{/each}
		{:else if route == 'profile'}
			<Profile on:user={onClick} />
		{:else if route == 'user/' + id}
			<ViewProfile {profile} on:user={onClick} />
		{/if}
	</div>
</div>
<svelte:window on:mouseup={onMouseUp} on:mousemove={onMouseMove} />

<style>
	div {
		left: 23rem;
		top: 8rem;
	}

	.draggable {
		cursor: move;
	}
</style>
