<script>
	// Svelte
	import { onMount } from 'svelte'
	import { createEventDispatcher } from 'svelte'
	const dispatch = createEventDispatcher()

	import { getGroups } from '../utils'
	import { currentUserGroups, currentUser } from '../stores/user'
	let groups

	$: console.log($currentUserGroups)
	onMount(async () => {
		groups = await getGroups()
	})
</script>

{#if groups}
	<div class="overflow-x-auto">
		<table class="table">
			<!-- head -->
			<thead>
				<tr>
					<th />
					<th>Name</th>
					<th>Role</th>
				</tr>
			</thead>
			<tbody>
				{#each groups as group, index (group.id)}
					<tr class="hover">
						<th>{index + 1}</th>
						<td class="cursor-pointer font-extrabold" on:click={() => dispatch('group', group.id)}>{group.title}</td>
						{#if group.creator_id == $currentUser.id}
							<td class="text-primary uppercase font-extrabold text-lg">Admin</td>
						{:else if $currentUserGroups.some((userGroup) => userGroup.id === group.id && userGroup.status == 'joined')}
							<td class="text-cyan-700 uppercase font-extrabold text-lg">Member</td>
						{:else if $currentUserGroups.some((userGroup) => userGroup.id === group.id && userGroup.status == 'requested')}
							<td class="text-orange-800 uppercase font-extrabold text-lg">Request pending...</td>
						{/if}
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
{/if}
