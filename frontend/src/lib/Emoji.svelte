<script>
	import { onMount } from 'svelte'
	import 'emoji-picker-element'
	import insertText from 'insert-text-at-cursor'
	export let input

	onMount(async () => {
		const emojiPicker = document.querySelector('emoji-picker')
		emojiPicker.addEventListener('emoji-click', insertEmoji)
	})

	function insertEmoji(e) {
		insertText(input, e.detail.unicode)
	}
</script>

<emoji-picker />

<style>
	emoji-picker {
		width: 100%;
		height: 100%;
	}
</style>
