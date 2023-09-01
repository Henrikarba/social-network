<script>
	const milf = new URL('../assets/milfs.webp', import.meta.url).href

	const positions = [
		{ top: '0', left: '0' },
		{ left: '600px', right: '0' },
		{ top: '200px', left: '0' },
		{ top: '0', right: '0' },
		{ top: '200px', left: '100px' },
	]
	let currentPosition = 0

	function getRandomPositionIndex() {
		return Math.floor(Math.random() * positions.length)
	}

	function nextPosition() {
		currentPosition = (currentPosition + 1) % positions.length
	}
</script>

<img
	src={milf}
	alt="milf"
	class="w-40 absolute image-transition"
	style="top: {positions[currentPosition].top}; left: {positions[currentPosition].left}; right: {positions[
		currentPosition
	].right}; bottom: {positions[currentPosition].bottom}"
	on:mouseover={nextPosition}
/>

<style>
	.image-transition {
		transition: top 0.2s, left 0.2s, right 0.2s, bottom 0.2s;
	}
</style>
