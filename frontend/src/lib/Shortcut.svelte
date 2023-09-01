<script>
	import { createEventDispatcher } from 'svelte'
	const dispatch = createEventDispatcher()

	export let imgurl

	export let left = 100
	export let top = 200
</script>

<!-- svelte-ignore a11y-click-events-have-key-events -->
<!-- svelte-ignore a11y-no-static-element-interactions -->
<div
	class="w-10 h-10 text-white text-sm draggable"
	style="left: {left}px; top: {top}px;"
	on:click={() => dispatch('open')}
>
	<!-- svelte-ignore a11y-missing-attribute -->
	<img src={imgurl} />
	<p class="w-full text-center"><slot /></p>
</div>

<style lang="scss">
	.draggable {
		user-select: none;
		cursor: pointer;
		position: absolute;
	}
</style>
