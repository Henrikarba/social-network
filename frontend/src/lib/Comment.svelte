<script>
	import { formatTime } from '../utils'
	export let comment
</script>

<div class="border-2 flex flex-col w-96">
	<h2 class="">{comment.content}</h2>
	<p class="opacity-50 text-right">{formatTime(comment.created_at)}</p>
</div>
