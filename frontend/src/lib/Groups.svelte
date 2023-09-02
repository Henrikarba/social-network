<script>
	import { slide } from 'svelte/transition'
	// Svelte
	import { createEventDispatcher } from 'svelte'
	const dispatch = createEventDispatcher()
	import { groupStore } from '../stores/groups'
	import { formatTime } from '../utils'
	import { currentUserGroups, currentUser } from '../stores/user'

	$: groups = $groupStore
</script>

{#if groups}
	<div
		class="overflow-x-auto"
		in:slide|global={{ duration: 300, delay: 500, axis: 'y' }}
		out:slide|global={{ duration: 200, axis: 'x' }}
	>
		<table class="table">
			<!-- head -->
			<thead>
				<tr>
					<th />
					<th>Name</th>
					<th>Role</th>
					<th>Created At</th>
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
						{:else}
							<td>-</td>
						{/if}
						<td class="cursor-pointer font-extrabold">{formatTime(group.created_at)}</td>
					</tr>
				{/each}
			</tbody>
		</table>
	</div>
{/if}
