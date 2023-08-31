<script>
	// Svelte
	import { get_current_component } from 'svelte/internal'
	const THISComponent = get_current_component()
	import { createEventDispatcher } from 'svelte'
	import { scale } from 'svelte/transition'
	// Icons
	import FaRegWindowClose from 'svelte-icons/fa/FaRegWindowClose.svelte'
	import MdFullscreen from 'svelte-icons/md/MdFullscreen.svelte'
	import MdFullscreenExit from 'svelte-icons/md/MdFullscreenExit.svelte'

	export let ieUrl
	export let z
	let left = 400
	let top = 100
	let full = false
	let moving = false

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

	const dispatch = createEventDispatcher()
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
			<p class="ml-4">http://localhost:5000/posts</p>
		</div>
	</div>
	<div class="{'full ? h-full : h-[635px]'} overflow-y-auto bg-slate-100">
		<h2>hello</h2>
		<h2>hello</h2>
		<h2>hello</h2>
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
