<script>
	// Icons
	import FaRegWindowClose from 'svelte-icons/fa/FaRegWindowClose.svelte'
	import { get_current_component } from 'svelte/internal'
	const THISComponent = get_current_component()
	import { createEventDispatcher } from 'svelte'
	const dispatch = createEventDispatcher()
	import { scale } from 'svelte/transition'

	export let type
	export let id
	export let z
	$: console.log(z)

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
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
<div
	class="w-[1280px] h-[720px]'} border-2 rounded absolute border-b-4 border-zinc-500 select-none {z}"
	style="left: {left}px; top: {top}px; z-index: {z};"
	in:scale|global={{ duration: 500, start: 0.5 }}
	on:click={() => dispatch('last', 'chat')}
>
	<div
		on:mousedown={onMouseDown}
		class=" bg-gradient-to-r from-sky-500 to-blue-700 h-10 border-b-2 border-blue-950 flex justify-end items-center draggable"
	>
		<div class="h-6 w-6 mr-2 text-black cursor-pointer" on:click={destroySelf}><FaRegWindowClose /></div>
	</div>
</div>
<svelte:window on:mouseup={onMouseUp} on:mousemove={onMouseMove} />

<style>
	.draggable {
		cursor: move;
	}
</style>
