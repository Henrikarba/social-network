<script>
	//
	import { socket } from '../ws'
	import { currentChat } from '../stores/chat'
	import { currentUser } from '../stores/user'
	import { formatTime } from '../utils'
	import Emoji from './Emoji.svelte'
	import { messagesStore } from '../stores/chat'
	// Icons
	import FaRegWindowClose from 'svelte-icons/fa/FaRegWindowClose.svelte'
	import MdInsertEmoticon from 'svelte-icons/md/MdInsertEmoticon.svelte'
	import MdDoNotDisturbAlt from 'svelte-icons/md/MdDoNotDisturbAlt.svelte'

	// Svelte
	import { get_current_component } from 'svelte/internal'
	const THISComponent = get_current_component()
	import { afterUpdate, createEventDispatcher, onMount } from 'svelte'
	const dispatch = createEventDispatcher()
	import { scale } from 'svelte/transition'

	export let type
	export let id
	export let z
	function fetchChat() {
		const data = {
			action: 'get_chat',
			data: {
				sender_id: parseInt(id),
			},
		}
		socket.send(JSON.stringify(data))
	}

	let input
	let msg
	let disabled = true
	$: {
		disabled = !(msg && msg.length > 0 && msg.trim() != '')
	}

	function newMessage() {
		const data = {
			action: 'new_message',
			data: {
				recipient_id: parseInt(id),
				created_at: new Date(),
				content: input.value,
			},
		}
		socket.send(JSON.stringify(data))
		$currentChat.messages = [...$currentChat.messages, data.data]
		input.value = ''
		disabled = true
	}

	onMount(() => {
		fetchChat()
	})

	let title
	$: if ($currentChat) {
		if ($currentChat.partner.id == id && type == 'regular') {
			title = `${$currentChat.partner.first_name} ${$currentChat.partner.last_name}`
		}
		$currentChat.messages.sort((a, b) => {
			const dateA = new Date(a.created_at).getTime()
			const dateB = new Date(b.created_at).getTime()
			return dateA - dateB
		})
	}
	$: console.log($currentChat)
	$: {
		if ($messagesStore && $messagesStore.some((item) => item.sender_id == id)) {
			const messageToAdd = $messagesStore.find((item) => item.sender_id == id)
			if (messageToAdd) {
				$currentChat.messages = [...$currentChat.messages, messageToAdd]
			}
			$messagesStore = $messagesStore.filter((item) => item.sender_id != id)
		}
	}

	// moving window
	let left = 300
	let top = 20
	let moving = false

	function onMouseDown() {
		dispatch('last', 'chat')
		moving = true
	}

	function onMouseMove(e) {
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

	let container
	const scrollToBottom = () => {
		if (container) {
			container.scrollTop = container.scrollHeight
		}
	}
	afterUpdate(scrollToBottom)
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
<div
	class="bg-slate-700 w-[1260px] h-[700px] border-2 rounded absolute border-b-4 border-zinc-500 select-none {z}"
	style="left: {left}px; top: {top}px; z-index: {z};"
	in:scale|global={{ duration: 500, start: 0.5 }}
	on:click={() => dispatch('last', 'chat')}
>
	<div
		on:mousedown={onMouseDown}
		class=" bg-gradient-to-r from-sky-500 to-blue-700 h-10 border-b-2 border-blue-950 flex justify-between items-center draggable"
	>
		<h2 class="ml-4 font-bold select-none text-slate-800 text-lg">Chat with: {title}</h2>
		<div class="h-6 w-6 mr-2 text-black cursor-pointer" on:click={destroySelf}><FaRegWindowClose /></div>
	</div>
	<div data-theme="dracula" class="flex w-full">
		<div class="w-[1260px]">
			<div class="bg-base-200 h-[700px] overflow-y-scroll overflow-x-hidden px-6 py-4" bind:this={container}>
				{#if $currentChat}
					{#each $currentChat.messages as chat}
						{#if chat.sender_id == id}
							<div class="chat chat-start">
								<div class="chat-image avatar">
									<div class="w-10 rounded-full">
										<img src="http://localhost:80/images/{$currentChat.partner.avatar}" />
									</div>
								</div>
								<div class="chat-header text-accent font-bold mb-2">
									{$currentChat.partner.first_name}
									{$currentChat.partner.last_name}
									<time class="text-xs opacity-50">{formatTime(chat.created_at)}</time>
								</div>
								<div class="chat-bubble bg-info text-base-200 font-semibold">{chat.content}</div>
							</div>
						{:else}
							<div class="chat chat-end">
								<div class="chat-image avatar">
									<div class="w-10 rounded-full">
										<img src="http://localhost:80/images/{$currentUser.avatar}" />
									</div>
								</div>
								<div class="chat-header text-accent font-bold mb-2">
									{$currentUser.first_name}
									{$currentUser.last_name}
									<time class="text-xs opacity-50">{formatTime(chat.created_at)}</time>
								</div>
								<div class="chat-bubble bg-secondary text-base-200 font-semibold flex-wrap">
									{chat.content}
								</div>
							</div>
						{/if}
					{/each}
				{/if}
			</div>

			<div class="flex sticky bottom-0">
				<form on:submit|preventDefault={newMessage} class="w-full flex">
					<input
						type="text"
						placeholder="Type here"
						class="input text-white input-accent w-full rounded-none focus:outline-none"
						bind:this={input}
						bind:value={msg}
					/>
					<details class="dropdown dropdown-right dropdown-end">
						<summary class="btn rounded-none"><MdInsertEmoticon /> </summary>
						<div class="w-[260px] h-[740px] ml-16 dropdown-content">
							<Emoji {input} />
						</div>
					</details>

					<button {disabled} class="btn rounded-none w-20 btn-success disabled:btn-error">
						{#if disabled}
							<div class="h-8">
								<MdDoNotDisturbAlt />
							</div>
						{:else}
							SEND
						{/if}
					</button>
				</form>
			</div>
		</div>
	</div>
</div>
<svelte:window on:mouseup={onMouseUp} on:mousemove={onMouseMove} />

<style>
	.draggable {
		cursor: move;
	}
</style>
