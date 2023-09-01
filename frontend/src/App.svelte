<script>
	import { fade } from 'svelte/transition'
	import { onMount } from 'svelte'
	// WS
	import { isAuthenticated, socket, createWebSocket } from './ws'

	// Components
	import Footer from './lib/Footer.svelte'
	import Shortcut from './lib/Shortcut.svelte'
	import IE from './lib/IE.svelte'
	import MSN from './lib/MSN.svelte'
	import Login from './lib/Login.svelte'

	const msnUrl = new URL('./assets/msn.png', import.meta.url).href
	const ieUrl = new URL('./assets/ie.png', import.meta.url).href

	let ieOpen = false
	let msnOpen = false
	function openIE() {
		ieOpen = !ieOpen
	}

	function openMSN() {
		msnOpen = !msnOpen
	}
	let last
	let zMax

	$: last = zMax
	function zindex(event) {
		zMax = event.detail
	}
	let loading = true
	let authenticated = false

	onMount(() => {
		createWebSocket()
		setTimeout(() => {
			loading = false
		}, 3000)
	})

	function bsod() {
		setTimeout(() => {
			loading = true
		}, 1500)
		setTimeout(() => {
			loading = false
		}, 3000)
	}

	$: authenticated = $isAuthenticated
</script>

{#if authenticated && !loading}
	<!-- svelte-ignore a11y-click-events-have-key-events -->
	<!-- svelte-ignore a11y-no-static-element-interactions -->
	<div on:click={() => (last = 'msn')}>
		<Shortcut imgurl={msnUrl} left={300} on:open={openMSN}>MSN</Shortcut>
	</div>

	<!-- svelte-ignore a11y-click-events-have-key-events -->
	<!-- svelte-ignore a11y-no-static-element-interactions -->
	<div on:click={() => (last = 'ie')}>
		<Shortcut imgurl={ieUrl} left={100} on:open={openIE}>Internet Explorer</Shortcut>
	</div>
	<Footer on:bsod={bsod} />

	{#if ieOpen}
		<div>
			<IE {ieUrl} on:close={openIE} on:last={zindex} z={last == 'ie' ? 'z-top' : 'z-low'} />
		</div>
	{/if}
	{#if msnOpen}
		<div>
			<MSN {msnUrl} on:close={openMSN} on:last={zindex} z={last == 'msn' ? 'z-top' : 'z-low'} />
		</div>
	{/if}
{/if}

{#if !loading && !authenticated}
	<Login />
{/if}
{#if loading}
	<div out:fade={{ duration: 300 }} class="loader" />
{/if}

<style>
	.loader {
		background-image: url('./assets/loading.gif');
		position: absolute;
		width: 100%;
		height: 100vh;
		z-index: 99999999;
		background-position: center center;
		background-size: 100%;
		background-attachment: fixed;
	}

	:global(.z-top) {
		z-index: 333;
	}
	:global(.z-low) {
		z-index: 125;
	}
</style>
