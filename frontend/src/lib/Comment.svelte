<script>
	import { slide } from 'svelte/transition'
	import { formatTime } from '../utils'
	export let comment
</script>

<div class="border-2 flex flex-col w-96" in:slide={{ duration: 500, axis: 'y' }}>
	<h2 class="">{comment.content}</h2>
	{#if comment.image_url && comment.image_url != ''}
		<img src="http://localhost:80/images/{comment.image_url}" alt="" />
	{/if}
	<p class="opacity-50 text-right">{formatTime(comment.created_at)}</p>
</div>
