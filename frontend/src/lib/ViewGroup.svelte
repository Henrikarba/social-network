<script>
	import { createEventDispatcher } from 'svelte'
	const dispatch = createEventDispatcher()

	export let group

	console.log(group)
	let owner = group.members[0]
	let members = group.members.slice(1)
</script>

<div class="border-4 border-red-400 mt-4 p-4 rounded">
	<h2 class="text-center text-4xl font-extrabold">{group.title}</h2>
	<p class="text-center">{group.description}</p>
	<p class="text-center">
		Created by <span
			class="font-extrabold text-primary hover:cursor-pointer"
			on:click={() => dispatch('user', owner.id)}>{owner.first_name} {owner.last_name}</span
		>
	</p>
</div>
