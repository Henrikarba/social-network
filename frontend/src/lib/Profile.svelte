<script>
  import { slide } from "svelte/transition";
  // Svelte
  import { createEventDispatcher } from "svelte";
  const dispatch = createEventDispatcher();
  //
  import { formatDateTime } from "../utils.js";
  import { socket } from "../ws";
  import { postsStore } from "../stores/post";

  // Icons
  import FaBirthdayCake from "svelte-icons/fa/FaBirthdayCake.svelte";

  // Stores
  import {
    currentUser,
    currentUserFollowing,
    currentUserGroups,
    currentUserFollowers,
  } from "../stores/user";

  let user = $currentUser;
  let groups = $currentUserGroups;
  $: followers = $currentUserFollowers
    ? $currentUserFollowers.filter((item) => item.status != "pending")
    : [];
  $: following = $currentUserFollowing
    ? $currentUserFollowing.filter((item) => item.status != "pending")
    : [];

  // Birthday
  let birthday = formatDateTime(user.date_of_birth);

  const hackerman = new URL("../assets/robot.jpg", import.meta.url).href;

  function togglePrivacy() {
    user.privacy == 1 ? (user.privacy = 0) : (user.privacy = 1);
    const data = {
      action: "toggle_privacy",
    };
    socket.send(JSON.stringify(data));
  }
  $: console.log(followers, following);
  $: userPosts = $postsStore.filter((item) => item.user_id == $currentUser.id);
</script>

<div class="mt-10 flex flex-col items-center">
  <div class="flex">
    <h2 class="text-4xl">Welcome back {user.first_name} {user.last_name}!</h2>
    <img class="w-10" src="http://localhost:80/images/{user.avatar}" alt="" />
  </div>
  <div class="flex items-center">
    <div class="w-10 mr-4">
      <FaBirthdayCake />
    </div>
    <div class="text-xl">
      {birthday}
    </div>
  </div>
  <div
    class="p-4 mt-2 border-r-emerald-700 border-8 rounded-xl border-t-red-500 border-l-cyan-600 border-b-orange-400"
  >
    <h2>
      Congratz, your birthday is coming soon. Or later. Depends on when you were
      born
    </h2>
  </div>
  <div class="mt-2">
    {#if user?.about_me}
      <h2>About me:</h2>
      {user.about_me}
    {/if}
  </div>
  <div class="flex gap-2">
    <div class="mt-2 border-2 border-zinc-400 p-4">
      {#if groups && groups.length > 0}
        <h2 class="border-b-2">My groups</h2>
        {#each groups.filter((item) => item.status == "joined") as group, index (group.id)}
          <h2>{group.title}</h2>
        {/each}
      {:else}
        <h2>You aren't in any groups :(</h2>
      {/if}
    </div>
    <div class="mt-2 border-2 border-zinc-400 p-4">
      {#if followers && followers.length > 0}
        <h2 class="border-b-2 font-bold">My followers</h2>
        {#each followers as follower, index (follower.id)}
          <h2
            on:click={() => dispatch("user", follower.id)}
            class="text-orange-800 cursor-pointer font-extrabold"
          >
            {follower.first_name}
            {follower.last_name}
          </h2>
        {/each}
      {:else}
        <h2 class="border-b-2 font-bold">
          My impressive list of no one who follows me:
        </h2>
      {/if}
    </div>
    <div class="mt-2 border-2 border-zinc-400 p-4">
      {#if following}
        <h2 class="border-b-2 font-bold">I follow these interesting guys:</h2>
        {#each following as follower, index (follower.id)}
          <h2
            on:click={() => dispatch("user", follower.id)}
            class="text-orange-800 cursor-pointer font-extrabold"
          >
            {follower.first_name}
            {follower.last_name}
          </h2>
        {/each}
      {:else}
        <h2 class="border-b-2 font-bold">no one follows me</h2>
      {/if}
    </div>
  </div>
  {#if userPosts && userPosts.length > 0}
    <h2 class="font-bold text-xl">Posts created by me:</h2>
    {#each userPosts as post}
      <h2
        class="text-primary hover:cursor-pointer"
        on:click={() => dispatch("singlePost", post.post_id)}
      >
        {post.title}
      </h2>
    {/each}
  {/if}
  <div class="mt-6 flex flex-col justify-center items-center">
    <h2>USE VPN?! Toggle privacy mode</h2>
    <div class="flex items-center">
      <input
        on:click={() => {
          togglePrivacy();
        }}
        bind:checked={user.privacy}
        id="checked-checkbox"
        type="checkbox"
        value=""
        class="w-4 h-4 text-blue-600 bg-gray-100 border-gray-300 rounded focus:ring-blue-500 dark:focus:ring-blue-600 dark:ring-offset-gray-800 focus:ring-2 dark:bg-gray-700 dark:border-gray-600"
      />
    </div>
    {#if user.privacy == 0}
      <div class="flex items-center">
        <h2>
          you are currently not using vpn and tor and blockchain and gfx card to
          hide your identity
        </h2>
      </div>
    {:else}
      <div>
        <h2 class="text-4xl font-extrabold text-green-500">STATUS:</h2>
        <img class="w-6/12" src={hackerman} alt="you are private citizen!" />
      </div>
    {/if}
  </div>
</div>

<style>
</style>
