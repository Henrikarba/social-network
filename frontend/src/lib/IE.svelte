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
	import NewPost from './NewPost.svelte'
	import Groups from './Groups.svelte'
	import ViewGroup from './ViewGroup.svelte'
	import NewGroup from './NewGroup.svelte'

	// Icons
	import FaRegWindowClose from 'svelte-icons/fa/FaRegWindowClose.svelte'
	import MdFullscreen from 'svelte-icons/md/MdFullscreen.svelte'
	import MdFullscreenExit from 'svelte-icons/md/MdFullscreenExit.svelte'

	// Stores
	import { postsStore, groupPostsStore } from '../stores/post'
	import { currentUser, currentUserGroups } from '../stores/user'
	import { groupStore, getGroups, groupInfo } from '../stores/groups'

	// ADS
	const ads = new URL('../assets/ads.png', import.meta.url).href
	const fb = new URL('../assets/fakebook.png', import.meta.url).href

	// utils
	import { getProfile } from '../utils'

	export let ieUrl
	export let z

	let left = 200
	let top = 100
	let full = false
	let moving = false

	let route = 'posts'

	async function groups() {
		const groups = await getGroups()
		groupStore.set(groups)
		route = 'groups'
	}

	let groupid
	let group

	// taskbar
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

	let ie
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
				groupid = parseInt(event.detail)
				group = await groupInfo(groupid)
				route = 'groups/' + groupid
				break
		}
	}

	function onPostClick(event) {
		groupPostFilter = event.detail
		route = 'group_posts'
	}

	function onSinglePost(event) {
		viewpostID = event.detail
		route = 'posts/' + event.detail
	}

	let groupPostFilter
	let viewpostID
	$: posts = $postsStore
	$: group_posts =
		groupPostFilter == 0 ? $groupPostsStore : $groupPostsStore.filter((item) => item.group.id == groupPostFilter)
	$: console.log(group_posts)
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
<div
	class="{full
		? 'w-full h-screen'
		: 'w-[1280px] h-[720px]'} border-2 rounded absolute border-b-4 border-zinc-500 select-none"
	style={full
		? 'left: 0px; top: 0px; z-index: ' + z + ';'
		: 'left: ' + left + 'px; top: ' + top + 'px; z-index: ' + z + ';'}
	in:scale|global={{ duration: 500, start: 0.5 }}
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
	<div class="{'full ? h-[88%] : h-[635px]'} overflow-y-scroll overflow-x-hidden bg-slate-100 p-6" bind:this={ie}>
		<img src={fb} class="h-20" alt="fakebook" />
		<nav class="flex items-center justify-center gap-4 font-bold text-4xl border-b-4 border-slate-900">
			<button on:click={() => (route = 'posts')}>[POSTS]</button>
			<button on:click={() => (route = 'group_posts')}>[GROUP POSTS]</button>
			<button on:click={() => groups()}>[GROUPS]</button>
			<button on:click={() => (route = 'profile')}>[PROFILE]</button>
		</nav>
		{#if route == 'posts' || route == 'group_posts'}
			<div class="flex justify-center">
				<button class="btn btn-ghost mt-4" on:click={() => (route = 'post/new')}>Create New PoSt</button>
			</div>
			{#if route == 'group_posts'}
				<div class="flex flex-col w-80">
					<h2>Filter by groups</h2>
					<select name="group" bind:value={groupPostFilter}>
						<option value="0">All</option>
						{#each $currentUserGroups as group}
							<option value={group.id}>{group.title}</option>
						{/each}
					</select>
				</div>
			{/if}
		{/if}
		{#if route == 'posts'}
			{#each posts as post, index (post.post_id)}
				<Post {post} on:user={onClick} />
			{/each}
		{:else if route == 'posts/' + viewpostID}
			{#each posts.filter((item) => item.post_id == viewpostID) as post, index (post.post_id)}
				<Post {post} on:user={onClick} />
			{/each}
		{:else if route == 'group_posts'}
			{#each group_posts as gPost, index (gPost.post_id)}
				<Post post={gPost} on:group={onClick} on:user={onClick} />
			{/each}
		{:else if route == 'profile'}
			<Profile on:user={onClick} />
		{:else if route == 'user/' + id}
			<ViewProfile {profile} on:user={onClick} on:singlePost={onSinglePost} />
		{:else if route == 'post/new'}
			<NewPost on:regular_post={() => (route = 'posts')} on:group_post={() => (route = 'group_posts')} />
		{:else if route == 'groups'}
			<div class="mt-4 flex justify-center w-full">
				<button on:click={() => (route = 'groups/new')} class="btn btn-primary">Create New Group</button>
			</div>
			<Groups on:group={onClick} />
		{:else if route == 'groups/new'}
			<NewGroup on:group={onClick} />
		{:else if route == 'groups/' + groupid}
			<ViewGroup {group} on:user={onClick} on:post={onPostClick} />
		{/if}
		<div class="flex justify-center flex-col items-center">
			<div class="mt-10">ADVERTISEMENT/SPONSORED CONTENT:</div>
			<img src={ads} alt="sponsored content" />
		</div>
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
