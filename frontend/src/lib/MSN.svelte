<script>
	// Svelte
	import { scale } from 'svelte/transition'
	import { get_current_component } from 'svelte/internal'
	const THISComponent = get_current_component()

	import { createEventDispatcher } from 'svelte'
	const dispatch = createEventDispatcher()

	// Icons
	import FaRegWindowClose from 'svelte-icons/fa/FaRegWindowClose.svelte'
	import IoIosPeople from 'svelte-icons/io/IoIosPeople.svelte'

	import { currentUserGroups, currentUser, chatStore } from '../stores/user'
	import { messagesStore } from '../stores/chat'
	$: joinedGroups = $currentUserGroups ? $currentUserGroups.filter((group) => group.status == 'joined') : []
	$: chats = $chatStore ? $chatStore.filter((chat) => chat.sender.id != $currentUser.id) : []
	$: console.log($chatStore)
	$: console.log(joinedGroups)

	export let msnUrl
	export let z

	let full = false

	let moving = false
	let left = 1000
	let top = 10
	function onMouseDown() {
		dispatch('last', 'msn')
		moving = true
	}
	$: console.log($messagesStore)
	function onMouseMove(e) {
		if (full) return
		if (moving) {
			left += e.movementX
			top += e.movementY
		}
	}

	function destroySelf() {
		dispatch('close')
		THISComponent.$destroy()
	}

	function onMouseUp() {
		moving = false
	}
</script>

<!-- svelte-ignore a11y-no-static-element-interactions -->
<!-- svelte-ignore a11y-click-events-have-key-events -->
<div
	class="{full
		? 'w-full h-screen'
		: 'w-[400px] h-[900px]'} {z} border-2 rounded absolute border-b-4 border-zinc-500 select-none"
	style={full ? 'left: 0px; top: 0px;' : 'left: ' + left + 'px; top: ' + top + 'px; z-index: ' + z + ';'}
	in:scale|global={{ duration: 500, start: 0.5 }}
>
	<!-- svelte-ignore a11y-no-static-element-interactions -->
	<div
		on:mousedown={onMouseDown}
		class=" bg-gradient-to-r from-sky-500 to-blue-700 h-10 border-b-2 border-blue-950 flex justify-between items-center draggable"
	>
		<div class=" ml-4 font-bold select-none text-slate-800 text-lg flex items-center gap-2">
			<!-- svelte-ignore a11y-missing-attribute -->
			<img class="h-6 w-6" src={msnUrl} />
			<h2>MESSENGER</h2>
		</div>
		<!-- svelte-ignore a11y-click-events-have-key-events -->
		<div class="flex">
			<!-- svelte-ignore a11y-no-static-element-interactions -->
			<div class="h-6 w-6 mr-2 text-black cursor-pointer" on:click={destroySelf}>
				<FaRegWindowClose />
			</div>
		</div>
	</div>
	<div class="bg-slate-100 h-[96%] border-b-4 border-neutral-600 flex flex-col px-2">
		<div class="h-32 border-b-2 p-4">
			<div class="flex flex-col items-center font-bold justify-center text-xl">
				<img
					src="http://localhost:80/images/{$currentUser.avatar}"
					class="w-20 h-20"
					alt="{$currentUser.first_name} {$currentUser.last_name}"
				/>
				<h2 class="text-center">{$currentUser.first_name} {$currentUser.last_name}</h2>
			</div>
		</div>
		<div class="mt-4 pb-2 border-b-2">
			<h2 class="font-bold text-xl flex flex-col">Regular Chats</h2>
			{#if chats && chats.length > 0}
				{#each chats as chat}
					<div
						class="flex items-center cursor-pointer"
						on:click|stopPropagation={() => dispatch('chat', { type: 'regular', id: chat.sender.id })}
					>
						<div class="w-8 text-primary">
							<IoIosPeople />
						</div>
						<h2 class="ml-1 font-bold hover:text-orange-500">{chat.sender.first_name} {chat.sender.last_name}</h2>
						{#if $messagesStore && $messagesStore.length > 0 && $messagesStore.some((item) => item.sender_id == chat.sender.id)}
							<h2 class="ml-4 text-red-600 font-bold">NEW MESSAGES</h2>
						{/if}
					</div>
				{/each}
			{:else}
				<h2>No chat history. Find an user on fakebook and initiate chat.</h2>
			{/if}
		</div>
		<div class="flex flex-col">
			<h2 class="font-bold text-xl flex flex-col">Group chats</h2>
			{#if joinedGroups && joinedGroups.length > 0}
				{#each joinedGroups as group}
					<div
						class="flex items-center cursor-pointer"
						on:click|stopPropagation={() =>
							dispatch('chat', { type: 'group', id: parseInt(group.id), name: group.title })}
					>
						<div class="w-8 text-primary">
							<IoIosPeople />
						</div>
						<h2 class="ml-1 font-bold hover:text-orange-500">{group.title}</h2>
					</div>
				{/each}
			{:else}
				<h2>You have not joined any groups</h2>
			{/if}
		</div>
	</div>
</div>
<svelte:window on:mouseup={onMouseUp} on:mousemove={onMouseMove} />

<style>
	.draggable {
		cursor: move;
	}
</style>
