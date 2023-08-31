<script>
	// Svelte
	import { scale } from 'svelte/transition'
	import { get_current_component } from 'svelte/internal'
	import { createEventDispatcher } from 'svelte'
	// Icons
	import FaRegWindowClose from 'svelte-icons/fa/FaRegWindowClose.svelte'

	const THISComponent = get_current_component()
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

	function onMouseMove(e) {
		if (full) return
		if (moving) {
			left += e.movementX
			top += e.movementY
		}
	}

	const dispatch = createEventDispatcher()
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
	style={full ? 'left: 0px; top: 0px;' : 'left: ' + left + 'px; top: ' + top + 'px;'}
	in:scale|global={{ duration: 500, start: 0.5 }}
	on:click={() => dispatch('last', 'msn')}
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
	<div class="bg-slate-100 h-[96%] border-b-4 border-neutral-600">
		<h2>asd</h2>
	</div>
</div>
<svelte:window on:mouseup={onMouseUp} on:mousemove={onMouseMove} />

<style>
	.draggable {
		cursor: move;
	}
</style>
