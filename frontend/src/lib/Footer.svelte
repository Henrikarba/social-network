<script>
	import { createEventDispatcher } from 'svelte'
	import { onMount } from 'svelte'
	import FaWindows from 'svelte-icons/fa/FaWindows.svelte'
	import FaAngleUp from 'svelte-icons/fa/FaAngleUp.svelte'
	let time = new Date()
	const dispatch = createEventDispatcher()

	onMount(() => {
		const interval = setInterval(() => {
			time = new Date()
		}, 1000)

		return () => {
			clearInterval(interval)
		}
	})

	let bsod = false

	$: currentTime = `${time.getHours()}:${time.getMinutes().toString().padStart(2, '0')}`
</script>

{#if !bsod}
	<div
		class="select-none shadow fixed bottom-0 left-0 h-10 w-full bg-gradient-to-r from-sky-500 to-blue-700 border-t-2 border-blue-900 flex z-[9999999999]"
	>
		<button
			class="h-full w-32 px-2 bg-gradient-to-r from-green-700 to-green-800 rounded-r-lg"
			on:click={() => {
				bsod = true
				dispatch('bsod')
			}}
		>
			<div class="flex items-center justify-center gap-2">
				<div class="h-6 w-6 text-slate-200"><FaWindows /></div>
				<p class="text-slate-300 uppercase font-bold">Start</p>
			</div>
		</button>

		<div class="flex justify-between w-full">
			<div>hello</div>
			<div class="flex justify-center items-center mr-6 text-gray-200">
				<div class="w-6 h-6 mr-2 cursor-pointer"><FaAngleUp /></div>
				<div>{currentTime}</div>
			</div>
		</div>
	</div>
{/if}
{#if bsod}
	<div class="bsod" />
{/if}

<style lang="scss">
	.shadow {
		box-shadow: 0 -10px 100px rgba(8, 112, 184, 0.7);
	}

	.bsod {
		background-image: url('../assets/bsod.png');
		position: absolute;
		width: 100%;
		height: 100vh;
		z-index: 99999999;
		background-position: center center;
		background-size: 100%;
		background-attachment: fixed;
	}
</style>
