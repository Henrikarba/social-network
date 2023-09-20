<script>
  import { fly, fade } from "svelte/transition";
  import { socket, createWebSocket } from "../ws";
  import { LoginPage } from "../stores/signin";
  const xplogo = new URL("../assets/xplogo2.png", import.meta.url).href;

  let email;
  let password;
  let first_name;
  let last_name;
  let birthdate;
  let avatar;
  let nickname;
  let about;
  let registering = false;
  let msg = "";
  let selectedImage = null;
  async function register() {
    registering = true;
    msg = "........registering....";
    console.log(
      email,
      password,
      first_name,
      last_name,
      birthdate,
      avatar,
      nickname,
      about
    );
    const creds = {
      email: email,
      password: password,
      first_name: first_name,
      last_name: last_name,
      birthdate: birthdate,
      avatar: avatar,
      nickname: nickname,
      about: about,
    };
    try {
      const response = await fetch("http://localhost:80/register", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(creds),
        credentials: "include", // Important
      });
      if (!response.ok) {
        msg = "an error occured when trying to display error";
        throw new Error("Network response was not ok");
      }
      createWebSocket();
    } catch (error) {
      msg = "an error occured when trying to display error";
      console.error("Error:", error);
    }
  }
  function handleImageChange(event) {
    const file = event.target.files[0];
    if (file) {
      selectedImage = file;
    }
  }
</script>

<main>
  <div
    class=" h-3/4 w-4/12 select-none"
    in:fade|global={{ duration: 500 }}
    out:fly|global={{ duration: 300, x: 800 }}
  >
    <div
      class="h-10 bg-gradient-to-t from-blue-500 to-blue-700 flex justify-center items-center border-2 rounded border-blue-950"
    >
      <h2 class="text-white font-extrabold">Register to Windows</h2>
    </div>
    <div
      class="h-30 bg-gray-500 logo flex justify-center border-b-2 border-stone-700"
    >
      <img src={xplogo} class="h-20" alt="XP LOGO" />
    </div>
    <div
      class="h-100 p-5 bg-gray-500 flex flex-col items-center justify-center gap-5 border-b-2 border-gray-500 rounded-b-lg"
    >
      <div class="flex items-center gap-2">
        <label for="first_name" class="w-24 text-right">First Name</label>
        <input
          required
          bind:value={first_name}
          type="first_name"
          placeholder="Type here"
          class="input w-full max-w-xs h-8 bg-white focus:outline-none"
        />
        <label for="last_name" class="w-24 text-right">Last Name</label>
        <input
          required
          bind:value={last_name}
          type="last_name"
          placeholder="Type here"
          class="input w-full max-w-xs h-8 bg-white focus:outline-none"
        />
      </div>
      <div class="flex items-center gap-2">
        <label for="email" class="w-24 text-right">Email</label>
        <input
          required
          bind:value={email}
          type="email"
          placeholder="Type here"
          class="input w-full max-w-xs h-8 bg-white focus:outline-none"
        />
      </div>
      <div class="flex items-center gap-2">
        <label for="password" class="w-24 text-right">Password</label>
        <input
          required
          bind:value={password}
          type="password"
          placeholder="Type here"
          class="input w-full max-w-xs h-8 bg-white focus:outline-none"
        />
      </div>
      <div class="flex items-center gap-2">
        <label for="birthdate" class="w-24 text-right">Date of Birth</label>
        <input
          required
          bind:value={birthdate}
          type="date"
          placeholder="Type here"
          class="input w-full max-w-xs h-8 bg-white focus:outline-none"
        />
      </div>
      <div class="flex items-center gap-2">
        <label for="avatar" class="w-24 text-right">Avatar {"(optional)"}</label
        >
        <input
          bind:value={avatar}
          on:change={handleImageChange}
          type="file"
          accept="image/*"
          class="w-full max-w-xs h-8 focus:outline-none"
        />
      </div>
      <div class="flex items-center gap-2">
        <label for="nickname" class="w-24 text-right"
          >Nickname {"(optional)"}</label
        >
        <input
          bind:value={nickname}
          type="nickname"
          placeholder="Type here"
          class="input w-full max-w-xs h-8 focus:outline-none"
        />
      </div>
      <div class="flex items-center gap-2">
        <label for="about" class="w-24 text-right">About {"(optional)"}</label>
        <input
          bind:value={about}
          type="about"
          placeholder="Type here"
          class="input w-full max-w-xs h-8 bg-white focus:outline-none"
        />
      </div>
      <div class="flex flex-col items-center">
        <div class="flex flex-row">
          <button class="btn mr-3" on:click={register}>REGISTER</button>
          <button
            class="btn"
            on:click={() => {
              LoginPage.set("login");
            }}>TO LOGIN</button
          >
        </div>
        {#if registering}
          <h2 class="text-xl text-red-700 font-extrabold">{msg}</h2>
        {/if}
      </div>
    </div>
  </div>
</main>

<style lang="scss">
  main {
    background-color: #004e98;
    z-index: 99999999999999;
    height: 100vh;
    display: flex;
    align-items: center;
    justify-content: center;

    .logo {
      background: linear-gradient(to right, #6286e1, #6286e1, #9fbbf6, #7899e9);
    }
  }
</style>
